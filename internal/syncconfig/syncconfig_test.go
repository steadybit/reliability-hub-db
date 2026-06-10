package syncconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Roundtrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sync.yml")
	contents := `extensions:
  - id: com.steadybit.extension_rabbitmq
    image: ghcr.io/steadybit/extension-rabbitmq:main
    label: RabbitMQ
    description: Steadybit extension for RabbitMQ.
    maintainer: com.steadybit
    license: MIT
    gitHub:
      owner: steadybit
      repository: extension-rabbitmq
    ghcr:
      owner: steadybit
      repository: extension-rabbitmq
      package: extension-rabbitmq
    homepage: https://hub.steadybit.com/extension/com.steadybit.extension_rabbitmq
    installation: https://github.com/steadybit/extension-rabbitmq#installation
    changelog: https://github.com/steadybit/extension-rabbitmq/blob/main/CHANGELOG.md
    releaseDate: '2025-10-28'
    tags:
      - RabbitMQ
      - Message Queue
`
	if err := os.WriteFile(path, []byte(contents), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	c, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(c.Extensions) != 1 {
		t.Fatalf("want 1 extension, got %d", len(c.Extensions))
	}
	e := c.Lookup("com.steadybit.extension_rabbitmq")
	if e == nil {
		t.Fatal("lookup returned nil")
	}
	if e.Image != "ghcr.io/steadybit/extension-rabbitmq:main" {
		t.Errorf("image: %q", e.Image)
	}
	if e.GitHub.Owner != "steadybit" || e.GitHub.Repository != "extension-rabbitmq" {
		t.Errorf("gitHub: %+v", e.GitHub)
	}
	if e.GHCR.Package != "extension-rabbitmq" {
		t.Errorf("ghcr.package: %q", e.GHCR.Package)
	}
	if e.ReleaseDate != "2025-10-28" {
		t.Errorf("releaseDate: %q", e.ReleaseDate)
	}
	if got, want := e.Tags, []string{"RabbitMQ", "Message Queue"}; len(got) != len(want) {
		t.Errorf("tags: %v", got)
	}

	if c.Lookup("does-not-exist") != nil {
		t.Error("expected nil lookup for missing extension")
	}
}
