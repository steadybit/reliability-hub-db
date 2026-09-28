package describe

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse_Rabbitmq(t *testing.T) {
	payload, err := os.ReadFile(filepath.Join("testdata", "rabbitmq.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	out, err := Parse(payload)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(out.Actions) == 0 {
		t.Errorf("want >=1 action, got 0")
	}
	if len(out.TargetTypes) == 0 {
		t.Errorf("want >=1 target type, got 0")
	}
	// rabbitmq has no advice today; skip assertion.

	var found bool
	for _, a := range out.Actions {
		if a.ID == "com.steadybit.extension_rabbitmq.node.check" {
			if a.Kind != "check" {
				t.Errorf("node.check kind: %q", a.Kind)
			}
			if a.TargetSelection == nil || a.TargetSelection.TargetType != "com.steadybit.extension_rabbitmq.node" {
				t.Errorf("node.check targetSelection: %+v", a.TargetSelection)
			}
			if len(a.Parameters) == 0 {
				t.Errorf("node.check has no parameters")
			}
			found = true
			break
		}
	}
	if !found {
		t.Errorf("did not find com.steadybit.extension_rabbitmq.node.check among %d actions", len(out.Actions))
	}

	var queueFound bool
	for _, tt := range out.TargetTypes {
		if tt.ID == "com.steadybit.extension_rabbitmq.queue" {
			if tt.Label.One != "RabbitMQ Queue" || tt.Label.Other != "RabbitMQ Queues" {
				t.Errorf("queue target type label: %+v", tt.Label)
			}
			queueFound = true
			break
		}
	}
	if !queueFound {
		t.Errorf("did not find com.steadybit.extension_rabbitmq.queue target type")
	}

	if out.RouteCount() < 5 {
		t.Errorf("route count = %d, want >=5", out.RouteCount())
	}
}

func TestParse_MissingIndex(t *testing.T) {
	_, err := Parse([]byte(`{"foo": {}}`))
	if err == nil {
		t.Fatal("expected error for missing \"/\" index")
	}
}
