package dbmodel

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTargetTypeMarshal_GoldenRabbitMQQueue(t *testing.T) {
	golden, err := os.ReadFile(filepath.Join("..", "..", "targetTypes",
		"com.steadybit.extension_rabbitmq.queue", "description.yml"))
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	normalized := strings.TrimRight(string(golden), "\n") + "\n"
	icon := extractIcon(t, string(golden))

	got, err := Marshal(TargetType{
		ID:          "com.steadybit.extension_rabbitmq.queue",
		LabelOne:    "RabbitMQ Queue",
		LabelOther:  "RabbitMQ Queues",
		Icon:        icon,
		Extension:   "com.steadybit.extension_rabbitmq",
		ReleaseDate: "2025-10-28",
		Tags:        []string{"RabbitMQ", "Message Queue"},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if string(got) != normalized {
		t.Errorf("marshal output does not match golden\n--- got ---\n%s--- want ---\n%s", got, normalized)
	}
}
