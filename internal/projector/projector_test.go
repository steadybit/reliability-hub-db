package projector

import (
	"encoding/base64"
	"testing"

	"github.com/steadybit/reliability-hub-db/internal/describe"
	"github.com/steadybit/reliability-hub-db/internal/syncconfig"
)

func TestAction(t *testing.T) {
	svg := `<svg><path d="M1 1"/></svg>`
	a := describe.ActionDescription{
		ID:          "com.steadybit.extension_rabbitmq.node.check",
		Label:       "Check Nodes",
		Description: "Monitor RabbitMQ node events.",
		Icon:        "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg)),
		Kind:        "check",
		TargetSelection: &describe.TargetSelectionDef{
			TargetType: "com.steadybit.extension_rabbitmq.node",
		},
	}
	e := syncconfig.Extension{
		ID:          "com.steadybit.extension_rabbitmq",
		ReleaseDate: "2025-10-28",
		Tags:        []string{"RabbitMQ", "Message Queue"},
	}
	got, err := Action(a, e)
	if err != nil {
		t.Fatalf("project action: %v", err)
	}
	if got.ID != a.ID {
		t.Errorf("ID: %q", got.ID)
	}
	if got.Icon != svg {
		t.Errorf("icon decoded: %q", got.Icon)
	}
	if got.TargetType != "com.steadybit.extension_rabbitmq.node" {
		t.Errorf("targetType: %q", got.TargetType)
	}
	if got.Extension != e.ID {
		t.Errorf("extension: %q", got.Extension)
	}
	if got.ReleaseDate != e.ReleaseDate {
		t.Errorf("releaseDate: %q", got.ReleaseDate)
	}
	if len(got.Tags) != 2 || got.Tags[0] != "RabbitMQ" {
		t.Errorf("tags: %v", got.Tags)
	}
}

func TestAction_NoTargetSelection(t *testing.T) {
	a := describe.ActionDescription{
		ID:    "com.example.action",
		Label: "X",
		Icon:  "<svg/>",
		Kind:  "attack",
	}
	got, err := Action(a, syncconfig.Extension{ID: "com.example"})
	if err != nil {
		t.Fatalf("project: %v", err)
	}
	if got.TargetType != "" {
		t.Errorf("expected empty targetType, got %q", got.TargetType)
	}
}

func TestTargetType(t *testing.T) {
	tt := describe.TargetDescription{
		ID:    "com.steadybit.extension_rabbitmq.queue",
		Label: describe.PluralLabel{One: "RabbitMQ Queue", Other: "RabbitMQ Queues"},
		Icon:  "<svg/>",
	}
	e := syncconfig.Extension{
		ID:          "com.steadybit.extension_rabbitmq",
		ReleaseDate: "2025-10-28",
		Tags:        []string{"RabbitMQ"},
	}
	got, err := TargetType(tt, e)
	if err != nil {
		t.Fatalf("project target: %v", err)
	}
	if got.LabelOne != "RabbitMQ Queue" || got.LabelOther != "RabbitMQ Queues" {
		t.Errorf("labels: one=%q other=%q", got.LabelOne, got.LabelOther)
	}
}

func TestAdvice(t *testing.T) {
	a := describe.AdviceDescription{
		ID:          "com.example.advice.foo",
		Label:       "Foo",
		Description: "Bar.",
		Icon:        "<svg/>",
		TargetTypes: []string{"com.example.target.x"},
	}
	e := syncconfig.Extension{
		ID:          "com.example",
		ReleaseDate: "2024-01-01",
		Tags:        []string{"Example"},
	}
	got, err := Advice(a, e)
	if err != nil {
		t.Fatalf("project advice: %v", err)
	}
	if got.TargetTypes[0] != "com.example.target.x" {
		t.Errorf("targetTypes: %v", got.TargetTypes)
	}
}

func TestExtension(t *testing.T) {
	e := syncconfig.Extension{
		ID:          "com.steadybit.extension_rabbitmq",
		Label:       "RabbitMQ",
		Description: "An extension.",
		Icon:        "<svg/>",
		Maintainer:  "com.steadybit",
		License:     "MIT",
		GitHub:      syncconfig.GitRef{Owner: "steadybit", Repository: "extension-rabbitmq"},
		GHCR:        syncconfig.GHCRRef{Owner: "steadybit", Repository: "extension-rabbitmq", Package: "extension-rabbitmq"},
		Homepage:    "https://hub.steadybit.com/extension/com.steadybit.extension_rabbitmq",
		ReleaseDate: "2025-10-28",
		Tags:        []string{"RabbitMQ"},
	}
	got, err := Extension(e)
	if err != nil {
		t.Fatalf("project extension: %v", err)
	}
	if got.GitHub.Owner != "steadybit" {
		t.Errorf("github: %+v", got.GitHub)
	}
	if got.GHCR.Package != "extension-rabbitmq" {
		t.Errorf("ghcr: %+v", got.GHCR)
	}
}
