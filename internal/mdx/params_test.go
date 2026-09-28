package mdx

import (
	"strings"
	"testing"

	"github.com/steadybit/reliability-hub-db/internal/describe"
)

func boolPtr(b bool) *bool       { return &b }
func strPtr(s string) *string    { return &s }

func TestFormatParameterTable_SingleRequired(t *testing.T) {
	got := FormatParameterTable([]describe.Parameter{
		{
			Name:         "duration",
			Label:        "Duration",
			Description:  "How long should traffic been blocked?",
			Required:     boolPtr(true),
			DefaultValue: strPtr("60s"),
		},
	})
	want := `| Parameter | Description                           | Default |
|-----------|---------------------------------------|---------|
| Duration  | How long should traffic been blocked? | 60s     |
`
	if got != want {
		t.Errorf("table mismatch\n--- got ---\n%s--- want ---\n%s", got, want)
	}
}

func TestFormatParameterTable_OptionalNoDefault(t *testing.T) {
	got := FormatParameterTable([]describe.Parameter{
		{
			Name:        "expectedChanges",
			Label:       "Expected Changes",
			Description: "Which node-level changes to watch for.",
			// Required omitted => optional
		},
	})
	if !strings.Contains(got, "(optional) Which node-level changes to watch for.") {
		t.Errorf("missing (optional) prefix:\n%s", got)
	}
	if !strings.Contains(got, "| Expected Changes |") {
		t.Errorf("missing label cell:\n%s", got)
	}
}

func TestFormatParameterTable_EscapesPipes(t *testing.T) {
	got := FormatParameterTable([]describe.Parameter{
		{
			Name:        "x",
			Label:       "X",
			Description: "a | b",
			Required:    boolPtr(true),
		},
	})
	if !strings.Contains(got, `a \| b`) {
		t.Errorf("pipe not escaped:\n%s", got)
	}
}

func TestReplaceParameters_RoundTrip(t *testing.T) {
	in := `# Introduction

Block traffic.

# Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| OldName   | Stale       | 1s      |
`
	out, res, err := ReplaceParameters(in, []describe.Parameter{
		{
			Name:         "duration",
			Label:        "Duration",
			Description:  "How long should traffic been blocked?",
			Required:     boolPtr(true),
			DefaultValue: strPtr("60s"),
		},
	})
	if err != nil {
		t.Fatalf("replace: %v", err)
	}
	if res != ReplaceUpdated {
		t.Errorf("result: %s", res)
	}
	if strings.Contains(out, "OldName") {
		t.Errorf("old content not replaced:\n%s", out)
	}
	if !strings.Contains(out, "| Duration  | How long should traffic been blocked? | 60s     |") {
		t.Errorf("new table not present:\n%s", out)
	}
	if !strings.Contains(out, "# Introduction") || !strings.Contains(out, "Block traffic.") {
		t.Errorf("prose clobbered:\n%s", out)
	}
}

func TestReplaceParameters_PreservesContentAfterTable(t *testing.T) {
	in := `# Parameters

| Parameter | Description | Default |
|-----------|-------------|---------|
| X         | Y           | 1s      |

Some trailing prose.
`
	out, res, err := ReplaceParameters(in, []describe.Parameter{
		{Name: "a", Label: "A", Description: "B", Required: boolPtr(true), DefaultValue: strPtr("2s")},
	})
	if err != nil {
		t.Fatalf("replace: %v", err)
	}
	if res != ReplaceUpdated {
		t.Errorf("result: %s", res)
	}
	if !strings.Contains(out, "Some trailing prose.") {
		t.Errorf("trailing prose lost:\n%s", out)
	}
}

func TestReplaceParameters_NoHeading(t *testing.T) {
	in := `# Introduction

Just text. No params heading at all.
`
	out, res, err := ReplaceParameters(in, []describe.Parameter{
		{Name: "x", Label: "X", Description: "Y", Required: boolPtr(true)},
	})
	if err != nil {
		t.Fatalf("replace: %v", err)
	}
	if res != ReplaceNoHeading {
		t.Errorf("result: %s", res)
	}
	if out != in {
		t.Errorf("expected unchanged content; got:\n%s", out)
	}
}

func TestReplaceParameters_HeadingButNoTable(t *testing.T) {
	in := `# Parameters

There's no table here, just prose.
`
	out, res, err := ReplaceParameters(in, []describe.Parameter{
		{Name: "x", Label: "X", Description: "Y", Required: boolPtr(true)},
	})
	if err != nil {
		t.Fatalf("replace: %v", err)
	}
	if res != ReplaceNoTable {
		t.Errorf("result: %s", res)
	}
	if out != in {
		t.Errorf("expected unchanged content; got:\n%s", out)
	}
}

func TestReplaceParameters_Idempotent(t *testing.T) {
	params := []describe.Parameter{
		{Name: "duration", Label: "Duration", Description: "How long should traffic been blocked?", Required: boolPtr(true), DefaultValue: strPtr("60s")},
	}
	in := `# Parameters

| Parameter | Description                           | Default |
|-----------|---------------------------------------|---------|
| Duration  | How long should traffic been blocked? | 60s     |
`
	out, res, err := ReplaceParameters(in, params)
	if err != nil {
		t.Fatalf("first: %v", err)
	}
	if res != ReplaceUnchanged {
		t.Errorf("expected ReplaceUnchanged on identical input, got %s", res)
	}
	if out != in {
		t.Errorf("idempotent run changed content")
	}
}
