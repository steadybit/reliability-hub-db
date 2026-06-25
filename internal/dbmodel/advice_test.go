package dbmodel

import (
	"strings"
	"testing"
)

// TestAdviceMarshal_CanonicalForm validates the canonical Advice output. The
// existing advice description.yml files store the icon as a single-line plain
// scalar; the sync tool normalizes to a `|` block literal (matching the
// convention used by actions / target types / extensions). First-sync runs
// will produce this one-time normalization diff.
func TestAdviceMarshal_CanonicalForm(t *testing.T) {
	icon := `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="M1 1L23 23"/></svg>`

	got, err := Marshal(Advice{
		ID:          "com.steadybit.extension_kubernetes.advice.k8s-cpu-request",
		Label:       "Requesting Reasonable CPU Resources",
		Description: "Validates that your Kubernetes resources request reasonable CPU.",
		Icon:        icon,
		TargetTypes: []string{
			"com.steadybit.extension_kubernetes.kubernetes-daemonset",
			"com.steadybit.extension_kubernetes.kubernetes-deployment",
			"com.steadybit.extension_kubernetes.kubernetes-statefulset",
		},
		Tags:        []string{"Kubernetes", "Requests", "CPU"},
		Extension:   "com.steadybit.extension_kubernetes",
		ReleaseDate: "2024-02-01",
	})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	want := `---
id: com.steadybit.extension_kubernetes.advice.k8s-cpu-request
label: Requesting Reasonable CPU Resources
description: Validates that your Kubernetes resources request reasonable CPU.
icon: |
  <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="M1 1L23 23"/></svg>
targetTypes:
  - com.steadybit.extension_kubernetes.kubernetes-daemonset
  - com.steadybit.extension_kubernetes.kubernetes-deployment
  - com.steadybit.extension_kubernetes.kubernetes-statefulset
tags:
  - Kubernetes
  - Requests
  - CPU
extension: com.steadybit.extension_kubernetes
releaseDate: '2024-02-01'
`

	if string(got) != want {
		t.Errorf("marshal output mismatch\n--- got ---\n%s--- want ---\n%s\n--- diff hint ---\ngot len=%d want len=%d", got, want, len(got), len(want))
		if strings.HasPrefix(string(got), want) || strings.HasPrefix(want, string(got)) {
			t.Logf("one is a prefix of the other (trailing-newline issue?)")
		}
	}
}
