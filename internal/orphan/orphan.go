// Package orphan computes per-extension namespace diffs between what an
// extension currently emits via --describe and what already exists on disk.
//
// An "orphan" is a folder under actions/<id>, targetTypes/<id>, or
// advice/<id> whose ID prefix matches an extension we just synced, but which
// is NOT in the extension's --describe output. The sync tool reports these so
// a human can review the deletion before it is merged.
package orphan

import (
	"os"
	"path/filepath"
	"strings"
)

// Scan returns the list of folders under <root>/<kindDir>/ whose ID starts
// with extensionPrefix and that are NOT in keep.
//
//	root            — absolute path to the reliability-hub-db checkout
//	kindDir         — "actions", "targetTypes", or "advice"
//	extensionPrefix — the extension ID, used as a prefix match
//	keep            — set of IDs the extension currently emits
//
// The returned slice is sorted lexicographically.
func Scan(root, kindDir, extensionPrefix string, keep map[string]struct{}) ([]string, error) {
	dir := filepath.Join(root, kindDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var orphans []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		id := e.Name()
		if !belongsToExtension(id, extensionPrefix) {
			continue
		}
		if _, kept := keep[id]; !kept {
			orphans = append(orphans, id)
		}
	}
	return orphans, nil
}

// belongsToExtension returns true if id is owned by the extension with the
// given prefix. An ID is "owned" if it equals the prefix or begins with
// prefix + "." (so com.steadybit.extension_rabbitmq doesn't match
// com.steadybit.extension_rabbitmq_extra).
func belongsToExtension(id, prefix string) bool {
	if id == prefix {
		return true
	}
	return strings.HasPrefix(id, prefix+".")
}
