package dbmodel

import "gopkg.in/yaml.v3"

// Extension is one entry under extensions/<id>/description.yml.
type Extension struct {
	ID           string
	Label        string
	Description  string
	Icon         string
	Maintainer   string
	License      string
	GitHub       GitRef
	GHCR         GHCRRef
	Homepage     string
	Installation string
	Changelog    string
	ReleaseDate  string
	Tags         []string
}

type GitRef struct {
	Owner      string
	Repository string
}

type GHCRRef struct {
	Owner      string
	Repository string
	Package    string
}

func mappingPairs(pairs ...string) *yaml.Node {
	n := &yaml.Node{Kind: yaml.MappingNode}
	for i := 0; i < len(pairs); i += 2 {
		n.Content = append(n.Content, scalar(pairs[i]), scalar(pairs[i+1]))
	}
	return n
}

func (e Extension) MarshalYAML() (any, error) {
	n := &yaml.Node{Kind: yaml.MappingNode}
	n.Content = append(n.Content,
		scalar("id"), scalar(e.ID),
		scalar("label"), scalar(e.Label),
		scalar("description"), scalar(e.Description),
		scalar("icon"), literalScalar(e.Icon),
		scalar("maintainer"), scalar(e.Maintainer),
		scalar("license"), scalar(e.License),
		scalar("gitHub"), mappingPairs("owner", e.GitHub.Owner, "repository", e.GitHub.Repository),
		scalar("ghcr"), mappingPairs("owner", e.GHCR.Owner, "repository", e.GHCR.Repository, "package", e.GHCR.Package),
	)
	if e.Homepage != "" {
		n.Content = append(n.Content, scalar("homepage"), scalar(e.Homepage))
	}
	if e.Installation != "" {
		n.Content = append(n.Content, scalar("installation"), scalar(e.Installation))
	}
	if e.Changelog != "" {
		n.Content = append(n.Content, scalar("changelog"), scalar(e.Changelog))
	}
	if e.ReleaseDate != "" {
		n.Content = append(n.Content, scalar("releaseDate"), singleQuotedScalar(e.ReleaseDate))
	}
	if len(e.Tags) > 0 {
		n.Content = append(n.Content, scalar("tags"), tagsSequence(e.Tags))
	}
	return n, nil
}
