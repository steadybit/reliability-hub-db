package orphan

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{
		"actions/com.steadybit.extension_rabbitmq.node.check",
		"actions/com.steadybit.extension_rabbitmq.queue.check-backlog",
		"actions/com.steadybit.extension_rabbitmq.queue.publish-fixed-amount",
		"actions/com.steadybit.extension_aws.az.blackhole", // belongs to another extension
		"actions/com.steadybit.extension_rabbitmq_extra.foo", // prefix collision: must NOT be considered rabbitmq
	} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}

	keep := map[string]struct{}{
		"com.steadybit.extension_rabbitmq.node.check":             {},
		"com.steadybit.extension_rabbitmq.queue.check-backlog":    {},
		// rabbitmq.queue.publish-fixed-amount is intentionally NOT in keep => orphan
	}

	got, err := Scan(root, "actions", "com.steadybit.extension_rabbitmq", keep)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	want := []string{"com.steadybit.extension_rabbitmq.queue.publish-fixed-amount"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestScan_NonexistentDir(t *testing.T) {
	got, err := Scan(t.TempDir(), "advice", "com.example", nil)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestBelongsToExtension(t *testing.T) {
	prefix := "com.steadybit.extension_rabbitmq"
	cases := map[string]bool{
		"com.steadybit.extension_rabbitmq":                true,
		"com.steadybit.extension_rabbitmq.node.check":     true,
		"com.steadybit.extension_rabbitmq.queue":          true,
		"com.steadybit.extension_rabbitmq_extra.foo":      false,
		"com.steadybit.extension_aws":                     false,
		"com.steadybit.extension_aws.ec2-instance.state":  false,
	}
	for id, want := range cases {
		if got := belongsToExtension(id, prefix); got != want {
			t.Errorf("%s: got %v, want %v", id, got, want)
		}
	}
}
