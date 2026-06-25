// Package describe parses the JSON produced by an extension binary invoked
// with --describe. The output is a route-keyed map: "/" is the index listing
// which sub-paths hold actions / discoveries / target descriptions / advice,
// and each other key holds the full description for that route.
package describe

import (
	"encoding/json"
	"fmt"
)

// Output is the parsed --describe payload, dispatched by kind via the "/" index.
type Output struct {
	Actions     []ActionDescription
	TargetTypes []TargetDescription
	Advice      []AdviceDescription
}

// IndexBlock is the top-level "/" entry of the --describe output.
type IndexBlock struct {
	Actions          []RouteRef `json:"actions"`
	Discoveries      []RouteRef `json:"discoveries"`
	TargetTypes      []RouteRef `json:"targetTypes"`
	TargetAttributes []RouteRef `json:"targetAttributes"`
	EventListeners   []RouteRef `json:"eventListeners"`
	Advice           []RouteRef `json:"advice"`
}

// RouteRef is one entry in the index's per-kind array.
type RouteRef struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// ActionDescription is the per-action entry under its route in --describe output.
// Only the fields the sync tool needs are typed; unknown fields are ignored
// (Go's encoding/json default).
type ActionDescription struct {
	ID              string              `json:"id"`
	Label           string              `json:"label"`
	Description     string              `json:"description"`
	Icon            string              `json:"icon"`
	Kind            string              `json:"kind"`
	Parameters      []Parameter         `json:"parameters"`
	TargetSelection *TargetSelectionDef `json:"targetSelection,omitempty"`
}

// Parameter describes one action input parameter.
type Parameter struct {
	Name         string            `json:"name"`
	Label        string            `json:"label"`
	Description  string            `json:"description"`
	Type         string            `json:"type"`
	Required     *bool             `json:"required,omitempty"`
	DefaultValue *string           `json:"defaultValue,omitempty"`
	Options      []ParameterOption `json:"options,omitempty"`
}

// ParameterOption is one entry in a parameter's enum options.
type ParameterOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// TargetSelectionDef declares which target type an action operates on.
type TargetSelectionDef struct {
	TargetType string `json:"targetType"`
}

// TargetDescription is the per-target-type entry.
type TargetDescription struct {
	ID    string      `json:"id"`
	Label PluralLabel `json:"label"`
	Icon  string      `json:"icon"`
}

// PluralLabel is the singular/plural display label for a target type.
type PluralLabel struct {
	One   string `json:"one"`
	Other string `json:"other"`
}

// AdviceDescription is the per-advice entry. Schema follows advice-kit.
type AdviceDescription struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	TargetTypes []string `json:"targetTypes"`
}

// Parse parses a --describe JSON payload into a structured Output. Routes
// listed under "/".discoveries (and a few other non-stored sections) are
// silently skipped — the sync tool doesn't persist discovery metadata.
func Parse(payload []byte) (*Output, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(payload, &raw); err != nil {
		return nil, fmt.Errorf("parse --describe: %w", err)
	}
	indexRaw, ok := raw["/"]
	if !ok {
		return nil, fmt.Errorf("--describe output missing \"/\" index")
	}
	var index IndexBlock
	if err := json.Unmarshal(indexRaw, &index); err != nil {
		return nil, fmt.Errorf("parse --describe index: %w", err)
	}

	out := &Output{}
	for _, r := range index.Actions {
		body, ok := raw[r.Path]
		if !ok {
			return nil, fmt.Errorf("--describe action route %q missing body", r.Path)
		}
		var a ActionDescription
		if err := json.Unmarshal(body, &a); err != nil {
			return nil, fmt.Errorf("parse action %q: %w", r.Path, err)
		}
		out.Actions = append(out.Actions, a)
	}
	for _, r := range index.TargetTypes {
		body, ok := raw[r.Path]
		if !ok {
			return nil, fmt.Errorf("--describe targetType route %q missing body", r.Path)
		}
		var t TargetDescription
		if err := json.Unmarshal(body, &t); err != nil {
			return nil, fmt.Errorf("parse targetType %q: %w", r.Path, err)
		}
		out.TargetTypes = append(out.TargetTypes, t)
	}
	for _, r := range index.Advice {
		body, ok := raw[r.Path]
		if !ok {
			return nil, fmt.Errorf("--describe advice route %q missing body", r.Path)
		}
		var a AdviceDescription
		if err := json.Unmarshal(body, &a); err != nil {
			return nil, fmt.Errorf("parse advice %q: %w", r.Path, err)
		}
		out.Advice = append(out.Advice, a)
	}
	return out, nil
}

// RouteCount returns the total number of routes the index claims to expose.
// Callers use this as a sanity check: an extension binary returning fewer
// routes than expected (e.g. an empty index from a broken build) should be
// rejected before the sync tool deletes orphaned entities.
func (o *Output) RouteCount() int {
	return len(o.Actions) + len(o.TargetTypes) + len(o.Advice)
}
