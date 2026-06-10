// Package projector turns describe.* entities into dbmodel.* structs,
// combining them with the human-curated metadata from sync.yml.
package projector

import (
	"fmt"

	"github.com/steadybit/reliability-hub-db/internal/dbmodel"
	"github.com/steadybit/reliability-hub-db/internal/describe"
	"github.com/steadybit/reliability-hub-db/internal/icon"
	"github.com/steadybit/reliability-hub-db/internal/syncconfig"
)

// Action projects a describe.ActionDescription plus its extension's sync
// config entry into a dbmodel.Action.
func Action(a describe.ActionDescription, e syncconfig.Extension) (dbmodel.Action, error) {
	svg, err := icon.Decode(a.Icon)
	if err != nil {
		return dbmodel.Action{}, fmt.Errorf("decode action icon %q: %w", a.ID, err)
	}
	var targetType string
	if a.TargetSelection != nil {
		targetType = a.TargetSelection.TargetType
	}
	return dbmodel.Action{
		ID:          a.ID,
		Label:       a.Label,
		Description: a.Description,
		Icon:        svg,
		Kind:        a.Kind,
		TargetType:  targetType,
		Extension:   e.ID,
		ReleaseDate: e.ReleaseDate,
		Tags:        copyStrings(e.Tags),
	}, nil
}

// TargetType projects a describe.TargetDescription plus its extension's sync
// config entry into a dbmodel.TargetType.
func TargetType(t describe.TargetDescription, e syncconfig.Extension) (dbmodel.TargetType, error) {
	svg, err := icon.Decode(t.Icon)
	if err != nil {
		return dbmodel.TargetType{}, fmt.Errorf("decode target icon %q: %w", t.ID, err)
	}
	return dbmodel.TargetType{
		ID:          t.ID,
		LabelOne:    t.Label.One,
		LabelOther: t.Label.Other,
		Icon:        svg,
		Extension:   e.ID,
		ReleaseDate: e.ReleaseDate,
		Tags:        copyStrings(e.Tags),
	}, nil
}

// Advice projects a describe.AdviceDescription plus its extension's sync
// config entry into a dbmodel.Advice.
func Advice(a describe.AdviceDescription, e syncconfig.Extension) (dbmodel.Advice, error) {
	svg, err := icon.Decode(a.Icon)
	if err != nil {
		return dbmodel.Advice{}, fmt.Errorf("decode advice icon %q: %w", a.ID, err)
	}
	return dbmodel.Advice{
		ID:          a.ID,
		Label:       a.Label,
		Description: a.Description,
		Icon:        svg,
		TargetTypes: copyStrings(a.TargetTypes),
		Tags:        copyStrings(e.Tags),
		Extension:   e.ID,
		ReleaseDate: e.ReleaseDate,
	}, nil
}

// Extension projects a sync.yml entry into a dbmodel.Extension. The extension
// binary doesn't emit extension-level metadata, so everything comes from sync
// config.
func Extension(e syncconfig.Extension) (dbmodel.Extension, error) {
	svg, err := icon.Decode(e.Icon)
	if err != nil {
		return dbmodel.Extension{}, fmt.Errorf("decode extension icon %q: %w", e.ID, err)
	}
	return dbmodel.Extension{
		ID:          e.ID,
		Label:       e.Label,
		Description: e.Description,
		Icon:        svg,
		Maintainer:  e.Maintainer,
		License:     e.License,
		GitHub: dbmodel.GitRef{
			Owner:      e.GitHub.Owner,
			Repository: e.GitHub.Repository,
		},
		GHCR: dbmodel.GHCRRef{
			Owner:      e.GHCR.Owner,
			Repository: e.GHCR.Repository,
			Package:    e.GHCR.Package,
		},
		Homepage:     e.Homepage,
		Installation: e.Installation,
		Changelog:    e.Changelog,
		ReleaseDate:  e.ReleaseDate,
		Tags:         copyStrings(e.Tags),
	}, nil
}

func copyStrings(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}
