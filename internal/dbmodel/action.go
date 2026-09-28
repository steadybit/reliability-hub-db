package dbmodel

import "gopkg.in/yaml.v3"

// Action is one entry under actions/<id>/description.yml.
type Action struct {
	ID          string
	Label       string
	Description string
	Icon        string
	Kind        string
	TargetType  string
	Extension   string
	ReleaseDate string
	Tags        []string
}

func (a Action) MarshalYAML() (any, error) {
	n := &yaml.Node{Kind: yaml.MappingNode}
	n.Content = append(n.Content,
		scalar("id"), scalar(a.ID),
		scalar("label"), scalar(a.Label),
		scalar("description"), scalar(a.Description),
		scalar("icon"), literalScalar(a.Icon),
		scalar("kind"), scalar(a.Kind),
	)
	if a.TargetType != "" {
		n.Content = append(n.Content, scalar("targetType"), scalar(a.TargetType))
	}
	n.Content = append(n.Content, scalar("extension"), scalar(a.Extension))
	if a.ReleaseDate != "" {
		n.Content = append(n.Content, scalar("releaseDate"), singleQuotedScalar(a.ReleaseDate))
	}
	if len(a.Tags) > 0 {
		n.Content = append(n.Content, scalar("tags"), tagsSequence(a.Tags))
	}
	return n, nil
}
