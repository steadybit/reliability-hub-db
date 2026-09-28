// Command bootstrap is a one-shot migration that scrapes every
// extensions/<id>/description.yml in the repo root into a seed sync.yml.
//
// Run once when introducing the sync tool; after that, sync.yml is the source
// of truth for extension-level metadata and can be hand-edited.
//
// Usage:
//
//	bootstrap --root <dir> --out sync.yml [--default-image-tag main]
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type existingExtension struct {
	ID           string   `yaml:"id"`
	Label        string   `yaml:"label"`
	Description  string   `yaml:"description"`
	Icon         string   `yaml:"icon"`
	Maintainer   string   `yaml:"maintainer"`
	License      string   `yaml:"license"`
	GitHub       gitRef   `yaml:"gitHub"`
	GHCR         ghcrRef  `yaml:"ghcr"`
	Homepage     string   `yaml:"homepage"`
	Installation string   `yaml:"installation"`
	Changelog    string   `yaml:"changelog"`
	ReleaseDate  string   `yaml:"releaseDate"`
	Tags         []string `yaml:"tags"`
}

type gitRef struct {
	Owner      string `yaml:"owner"`
	Repository string `yaml:"repository"`
}

type ghcrRef struct {
	Owner      string `yaml:"owner"`
	Repository string `yaml:"repository"`
	Package    string `yaml:"package"`
}

// syncEntry is the shape we want to emit in sync.yml. Field order here drives
// the output YAML key order.
type syncEntry struct {
	ID           string   `yaml:"id"`
	Image        string   `yaml:"image,omitempty"`
	Label        string   `yaml:"label,omitempty"`
	Description  string   `yaml:"description,omitempty"`
	Maintainer   string   `yaml:"maintainer,omitempty"`
	License      string   `yaml:"license,omitempty"`
	GitHub       gitRef   `yaml:"gitHub"`
	GHCR         ghcrRef  `yaml:"ghcr"`
	Homepage     string   `yaml:"homepage,omitempty"`
	Installation string   `yaml:"installation,omitempty"`
	Changelog    string   `yaml:"changelog,omitempty"`
	ReleaseDate  string   `yaml:"releaseDate,omitempty"`
	Tags         []string `yaml:"tags,omitempty"`
	Icon         string   `yaml:"icon,omitempty"`
}

type syncConfig struct {
	Extensions []syncEntry `yaml:"extensions"`
}

func main() {
	root := flag.String("root", ".", "path to the reliability-hub-db checkout")
	out := flag.String("out", "sync.yml", "output path for the seed sync.yml")
	tag := flag.String("default-image-tag", "main", "default GHCR image tag to inject")
	flag.Parse()

	extensionsDir := filepath.Join(*root, "extensions")
	entries, err := os.ReadDir(extensionsDir)
	if err != nil {
		fatal("read %s: %v", extensionsDir, err)
	}

	var cfg syncConfig
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		path := filepath.Join(extensionsDir, e.Name(), "description.yml")
		b, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", path, err)
			continue
		}
		var ex existingExtension
		if err := yaml.Unmarshal(b, &ex); err != nil {
			fmt.Fprintf(os.Stderr, "skip %s: %v\n", path, err)
			continue
		}
		image := buildImageRef(ex.GHCR, *tag)
		cfg.Extensions = append(cfg.Extensions, syncEntry{
			ID:           ex.ID,
			Image:        image,
			Label:        ex.Label,
			Description:  ex.Description,
			Maintainer:   ex.Maintainer,
			License:      ex.License,
			GitHub:       ex.GitHub,
			GHCR:         ex.GHCR,
			Homepage:     ex.Homepage,
			Installation: ex.Installation,
			Changelog:    ex.Changelog,
			ReleaseDate:  ex.ReleaseDate,
			Tags:         ex.Tags,
			Icon:         ex.Icon,
		})
	}
	sort.Slice(cfg.Extensions, func(i, j int) bool {
		return cfg.Extensions[i].ID < cfg.Extensions[j].ID
	})

	f, err := os.Create(*out)
	if err != nil {
		fatal("create %s: %v", *out, err)
	}
	defer f.Close()
	enc := yaml.NewEncoder(f)
	enc.SetIndent(2)
	if err := enc.Encode(cfg); err != nil {
		fatal("encode: %v", err)
	}
	if err := enc.Close(); err != nil {
		fatal("close encoder: %v", err)
	}
	fmt.Printf("wrote %d extensions to %s\n", len(cfg.Extensions), *out)
}

func buildImageRef(g ghcrRef, tag string) string {
	if g.Owner == "" || g.Package == "" {
		return ""
	}
	return fmt.Sprintf("ghcr.io/%s/%s:%s", g.Owner, strings.TrimSpace(g.Package), tag)
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "bootstrap: "+format+"\n", args...)
	os.Exit(1)
}
