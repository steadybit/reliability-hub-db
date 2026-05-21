// Package source acquires --describe JSON for an extension. v1 supports
// reading from a local file (offline mode); docker-based image execution can
// be added later for the daily drift workflow.
package source

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// FromFile reads --describe JSON from a local file. Use for tests and for
// manual sync runs where the operator has already captured a dump.
func FromFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read describe file: %w", err)
	}
	return b, nil
}

// FromImage runs `docker run --rm <image> --describe` and returns the JSON
// payload from stdout. The extension binary writes structured logs to stderr,
// which docker exposes separately; we read only stdout.
func FromImage(image string) ([]byte, error) {
	cmd := exec.Command("docker", "run", "--rm",
		"-e", "STEADYBIT_EXTENSION_MANAGEMENT_ENDPOINTS_JSON=[]",
		image, "--describe")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("docker run %s --describe: %w (stderr: %s)", image, err, stderr.String())
	}
	return stdout.Bytes(), nil
}
