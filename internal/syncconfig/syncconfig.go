// Package syncconfig parses the hand-edited sync.yml registry that lists each
// extension and the non-binary metadata the sync tool needs (GitHub / GHCR
// coords, image tag to pull, license, tags, etc.).
package syncconfig

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the parsed sync.yml.
type Config struct {
	Extensions []Extension `yaml:"extensions"`
}

// Extension is one entry in sync.yml.
type Extension struct {
	ID           string   `yaml:"id"`
	Image        string   `yaml:"image"`
	Label        string   `yaml:"label"`
	Description  string   `yaml:"description"`
	Icon         string   `yaml:"icon"`
	Maintainer   string   `yaml:"maintainer"`
	License      string   `yaml:"license"`
	GitHub       GitRef   `yaml:"gitHub"`
	GHCR         GHCRRef  `yaml:"ghcr"`
	Homepage     string   `yaml:"homepage"`
	Installation string   `yaml:"installation"`
	Changelog    string   `yaml:"changelog"`
	ReleaseDate  string   `yaml:"releaseDate"`
	Tags         []string `yaml:"tags"`
}

type GitRef struct {
	Owner      string `yaml:"owner"`
	Repository string `yaml:"repository"`
}

type GHCRRef struct {
	Owner      string `yaml:"owner"`
	Repository string `yaml:"repository"`
	Package    string `yaml:"package"`
}

// Load reads and parses a sync.yml file from disk.
func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read sync config: %w", err)
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parse sync config: %w", err)
	}
	return &c, nil
}

// Lookup returns the extension entry with the given ID, or nil if absent.
func (c *Config) Lookup(id string) *Extension {
	for i := range c.Extensions {
		if c.Extensions[i].ID == id {
			return &c.Extensions[i]
		}
	}
	return nil
}
