package dbmodel

import "gopkg.in/yaml.v3"

// Advice is one entry under advice/<id>/description.yml.
type Advice struct {
	ID          string
	Label       string
	Description string
	Icon        string
	TargetTypes []string
	Tags        []string
	Extension   string
	ReleaseDate string
}

func (a Advice) MarshalYAML() (any, error) {
	n := &yaml.Node{Kind: yaml.MappingNode}
	n.Content = append(n.Content,
		scalar("id"), scalar(a.ID),
		scalar("label"), scalar(a.Label),
		scalar("description"), scalar(a.Description),
		scalar("icon"), literalScalar(a.Icon),
	)
	if len(a.TargetTypes) > 0 {
		n.Content = append(n.Content, scalar("targetTypes"), targetTypesSequence(a.TargetTypes))
	}
	if len(a.Tags) > 0 {
		n.Content = append(n.Content, scalar("tags"), tagsSequence(a.Tags))
	}
	n.Content = append(n.Content, scalar("extension"), scalar(a.Extension))
	if a.ReleaseDate != "" {
		n.Content = append(n.Content, scalar("releaseDate"), singleQuotedScalar(a.ReleaseDate))
	}
	return n, nil
}
