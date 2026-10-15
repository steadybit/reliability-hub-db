// Package syncer is the high-level orchestrator: it takes a parsed describe
// output plus a sync-config entry and produces all the file mutations needed
// to bring the on-disk reliability-hub-db into sync.
package syncer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/steadybit/reliability-hub-db/internal/dbmodel"
	"github.com/steadybit/reliability-hub-db/internal/describe"
	"github.com/steadybit/reliability-hub-db/internal/mdx"
	"github.com/steadybit/reliability-hub-db/internal/orphan"
	"github.com/steadybit/reliability-hub-db/internal/projector"
	"github.com/steadybit/reliability-hub-db/internal/syncconfig"
)

// Options controls a sync run.
type Options struct {
	Root        string // path to the reliability-hub-db checkout
	DryRun      bool   // if true, no files are written or deleted
	KeepOrphans bool   // if true, orphans are reported but not deleted
	MinRoutes   int    // sanity guard: abort if --describe has fewer routes than this
}

// Report summarizes the file mutations performed (or that would have been
// performed under --dry-run).
type Report struct {
	ExtensionID string

	ActionsWritten     []string
	TargetTypesWritten []string
	AdviceWritten      []string
	ExtensionWritten   bool

	MDXUpdated       []string
	MDXUnchanged     []string
	MDXNoHeading     []string
	MDXNoTable       []string
	MDXMissingFile   []string

	OrphanActions     []string
	OrphanTargetTypes []string
	OrphanAdvice      []string
	OrphansDeleted    bool
}

// ErrSanityGuard is returned when --describe contains fewer routes than the
// configured minimum. This usually means the binary is broken or the docker
// image is empty; the sync is aborted so we don't mass-delete entities.
var ErrSanityGuard = errors.New("sanity guard tripped: too few routes in --describe output")

// Sync runs one extension's sync pass.
func Sync(out *describe.Output, ext syncconfig.Extension, opts Options) (*Report, error) {
	if out.RouteCount() < opts.MinRoutes {
		return nil, fmt.Errorf("%w: got %d, want >= %d", ErrSanityGuard, out.RouteCount(), opts.MinRoutes)
	}

	r := &Report{ExtensionID: ext.ID, OrphansDeleted: !opts.KeepOrphans && !opts.DryRun}

	// 1. Actions
	keepActions := make(map[string]struct{}, len(out.Actions))
	for _, a := range out.Actions {
		action, err := projector.Action(a, ext)
		if err != nil {
			return nil, fmt.Errorf("project action %q: %w", a.ID, err)
		}
		if err := writeYAML(opts, filepath.Join(opts.Root, "actions", action.ID, "description.yml"), action); err != nil {
			return nil, err
		}
		r.ActionsWritten = append(r.ActionsWritten, action.ID)
		keepActions[action.ID] = struct{}{}

		// MDX parameter table — update only if file + heading + existing table all present.
		if err := updateActionMDX(opts, action.ID, a.Parameters, r); err != nil {
			return nil, err
		}
	}

	// 2. Target types
	keepTargets := make(map[string]struct{}, len(out.TargetTypes))
	for _, t := range out.TargetTypes {
		tt, err := projector.TargetType(t, ext)
		if err != nil {
			return nil, fmt.Errorf("project target %q: %w", t.ID, err)
		}
		if err := writeYAML(opts, filepath.Join(opts.Root, "targetTypes", tt.ID, "description.yml"), tt); err != nil {
			return nil, err
		}
		r.TargetTypesWritten = append(r.TargetTypesWritten, tt.ID)
		keepTargets[tt.ID] = struct{}{}
	}

	// 3. Advice
	keepAdvice := make(map[string]struct{}, len(out.Advice))
	for _, a := range out.Advice {
		ad, err := projector.Advice(a, ext)
		if err != nil {
			return nil, fmt.Errorf("project advice %q: %w", a.ID, err)
		}
		if err := writeYAML(opts, filepath.Join(opts.Root, "advice", ad.ID, "description.yml"), ad); err != nil {
			return nil, err
		}
		r.AdviceWritten = append(r.AdviceWritten, ad.ID)
		keepAdvice[ad.ID] = struct{}{}
	}

	// 4. Extension-level description.yml — driven entirely by sync config.
	e, err := projector.Extension(ext)
	if err != nil {
		return nil, fmt.Errorf("project extension %q: %w", ext.ID, err)
	}
	if err := writeYAML(opts, filepath.Join(opts.Root, "extensions", e.ID, "description.yml"), e); err != nil {
		return nil, err
	}
	r.ExtensionWritten = true

	// 5. Orphans (per-extension namespace).
	if r.OrphanActions, err = orphan.Scan(opts.Root, "actions", ext.ID, keepActions); err != nil {
		return nil, fmt.Errorf("scan action orphans: %w", err)
	}
	if r.OrphanTargetTypes, err = orphan.Scan(opts.Root, "targetTypes", ext.ID, keepTargets); err != nil {
		return nil, fmt.Errorf("scan target orphans: %w", err)
	}
	if r.OrphanAdvice, err = orphan.Scan(opts.Root, "advice", ext.ID, keepAdvice); err != nil {
		return nil, fmt.Errorf("scan advice orphans: %w", err)
	}

	if !opts.KeepOrphans && !opts.DryRun {
		for _, kindOrphans := range [][2]any{
			{"actions", r.OrphanActions},
			{"targetTypes", r.OrphanTargetTypes},
			{"advice", r.OrphanAdvice},
		} {
			kind := kindOrphans[0].(string)
			ids := kindOrphans[1].([]string)
			for _, id := range ids {
				dir := filepath.Join(opts.Root, kind, id)
				if err := os.RemoveAll(dir); err != nil {
					return nil, fmt.Errorf("remove orphan %s: %w", dir, err)
				}
			}
		}
	}

	sortReportSlices(r)
	return r, nil
}

func updateActionMDX(opts Options, actionID string, params []describe.Parameter, r *Report) error {
	path := filepath.Join(opts.Root, "actions", actionID, "summary.mdx")
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			r.MDXMissingFile = append(r.MDXMissingFile, actionID)
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}
	updated, res, err := mdx.ReplaceParameters(string(b), params)
	if err != nil {
		return fmt.Errorf("replace params in %s: %w", path, err)
	}
	switch res {
	case mdx.ReplaceUpdated:
		r.MDXUpdated = append(r.MDXUpdated, actionID)
		if !opts.DryRun {
			if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
				return fmt.Errorf("write %s: %w", path, err)
			}
		}
	case mdx.ReplaceUnchanged:
		r.MDXUnchanged = append(r.MDXUnchanged, actionID)
	case mdx.ReplaceNoHeading:
		r.MDXNoHeading = append(r.MDXNoHeading, actionID)
	case mdx.ReplaceNoTable:
		r.MDXNoTable = append(r.MDXNoTable, actionID)
	}
	return nil
}

func writeYAML(opts Options, path string, v any) error {
	b, err := dbmodel.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	if opts.DryRun {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}
	// Idempotency: skip write if content already matches.
	if existing, err := os.ReadFile(path); err == nil {
		if string(existing) == string(b) {
			return nil
		}
	}
	if err := os.WriteFile(path, b, 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func sortReportSlices(r *Report) {
	for _, s := range [][]string{
		r.ActionsWritten, r.TargetTypesWritten, r.AdviceWritten,
		r.MDXUpdated, r.MDXUnchanged, r.MDXNoHeading, r.MDXNoTable, r.MDXMissingFile,
		r.OrphanActions, r.OrphanTargetTypes, r.OrphanAdvice,
	} {
		sort.Strings(s)
	}
}
