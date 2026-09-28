package dbmodel

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestActionMarshal_GoldenAWSBlackholeZone validates byte-identical output
// for an action with a multi-line indented SVG icon and no releaseDate.
func TestActionMarshal_GoldenAWSBlackholeZone(t *testing.T) {
	golden, err := os.ReadFile(filepath.Join("..", "..", "actions",
		"com.steadybit.extension_aws.az.blackhole", "description.yml"))
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}

	icon := extractIcon(t, string(golden))

	got, err := Marshal(Action{
		ID:          "com.steadybit.extension_aws.az.blackhole",
		Label:       "Blackhole Zone",
		Description: "Simulates an outage of an entire availability zone",
		Icon:        icon,
		Kind:        "attack",
		TargetType:  "com.steadybit.extension_aws.zone",
		Extension:   "com.steadybit.extension_aws",
		Tags:        []string{"AWS", "Cloud", "Network"},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if string(got) != string(golden) {
		t.Errorf("marshal output does not match golden\n--- got ---\n%s--- want ---\n%s", got, golden)
	}
}

// TestActionMarshal_GoldenRabbitMQQueueCheckBacklog validates byte-identical
// output for an action with a single-line SVG icon, releaseDate, and tags.
// Uses a frozen snapshot under testdata/ so the test is decoupled from the
// live DB file (which the sync tool itself may rewrite). The snapshot has
// a trailing blank line; we normalize to a single trailing newline.
func TestActionMarshal_GoldenRabbitMQQueueCheckBacklog(t *testing.T) {
	golden, err := os.ReadFile(filepath.Join("testdata", "rabbitmq-queue-check-backlog.original.yml"))
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	normalized := strings.TrimRight(string(golden), "\n") + "\n"
	icon := extractIcon(t, string(golden))

	got, err := Marshal(Action{
		ID:          "com.steadybit.extension_rabbitmq.queue.check-backlog",
		Label:       "Check Queue Backlog",
		Description: "Monitor the total message count (backlog) in a RabbitMQ queue and fail the experiment if it exceeds a threshold",
		Icon:        icon,
		Kind:        "check",
		TargetType:  "com.steadybit.extension_rabbitmq.queue",
		Extension:   "com.steadybit.extension_rabbitmq",
		ReleaseDate: "2025-10-28",
		Tags:        []string{"Message Queue", "RabbitMQ"},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if string(got) != normalized {
		t.Errorf("marshal output does not match (normalized) golden\n--- got ---\n%s--- want ---\n%s", got, normalized)
	}
}

// extractIcon pulls the SVG content out of an existing description.yml so the
// test fixture doesn't need to inline it. It strips the `  ` indent prefix
// from each line of the YAML block-literal and returns the result without
// a trailing newline (the icon string the caller supplies has no terminal LF).
func extractIcon(t *testing.T, content string) string {
	t.Helper()
	const start = "icon: |\n"
	startIdx := strings.Index(content, start)
	if startIdx < 0 {
		t.Fatalf("icon block start not found")
	}
	startIdx += len(start)
	rest := content[startIdx:]
	end := 0
	for i := 0; i < len(rest); {
		nl := strings.Index(rest[i:], "\n")
		if nl < 0 {
			end = len(rest)
			break
		}
		nl += i
		line := rest[i:nl]
		if len(line) < 2 || line[0] != ' ' || line[1] != ' ' {
			end = i
			break
		}
		i = nl + 1
	}
	block := rest[:end]
	var out strings.Builder
	for i := 0; i < len(block); {
		nl := strings.Index(block[i:], "\n")
		if nl < 0 {
			out.WriteString(block[i:])
			break
		}
		nl += i
		out.WriteString(block[i+2 : nl])
		out.WriteByte('\n')
		i = nl + 1
	}
	s := out.String()
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	return s
}
