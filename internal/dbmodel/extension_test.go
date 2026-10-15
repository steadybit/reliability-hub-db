package dbmodel

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtensionMarshal_GoldenRabbitMQ(t *testing.T) {
	golden, err := os.ReadFile(filepath.Join("..", "..", "extensions",
		"com.steadybit.extension_rabbitmq", "description.yml"))
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	normalized := strings.TrimRight(string(golden), "\n") + "\n"
	icon := extractIcon(t, string(golden))

	got, err := Marshal(Extension{
		ID:           "com.steadybit.extension_rabbitmq",
		Label:        "RabbitMQ",
		Description:  "A Steadybit extension with various actions and check about RabbitMQ.",
		Icon:         icon,
		Maintainer:   "com.steadybit",
		License:      "MIT",
		GitHub:       GitRef{Owner: "steadybit", Repository: "extension-rabbitmq"},
		GHCR:         GHCRRef{Owner: "steadybit", Repository: "extension-rabbitmq", Package: "extension-rabbitmq"},
		Homepage:     "https://hub.steadybit.com/extension/com.steadybit.extension_rabbitmq",
		Installation: "https://github.com/steadybit/extension-rabbitmq#installation",
		Changelog:    "https://github.com/steadybit/extension-rabbitmq/blob/main/CHANGELOG.md",
		ReleaseDate:  "2025-10-28",
		Tags:         []string{"RabbitMQ", "Message Queue"},
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	if string(got) != normalized {
		t.Errorf("marshal output does not match golden\n--- got ---\n%s--- want ---\n%s", got, normalized)
	}
}
