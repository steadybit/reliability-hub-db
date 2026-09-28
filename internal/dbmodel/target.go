package dbmodel

import "gopkg.in/yaml.v3"

// TargetType is one entry under targetTypes/<id>/description.yml.
type TargetType struct {
	ID          string
	LabelOne    string
	LabelOther  string
	Icon        string
	Extension   string
	ReleaseDate string
	Tags        []string
}

func (t TargetType) MarshalYAML() (any, error) {
	labelNode := &yaml.Node{Kind: yaml.MappingNode}
	labelNode.Content = append(labelNode.Content,
		scalar("one"), scalar(t.LabelOne),
		scalar("other"), scalar(t.LabelOther),
	)

	n := &yaml.Node{Kind: yaml.MappingNode}
	n.Content = append(n.Content,
		scalar("id"), scalar(t.ID),
		scalar("label"), labelNode,
		scalar("icon"), literalScalar(t.Icon),
		scalar("extension"), scalar(t.Extension),
	)
	if t.ReleaseDate != "" {
		n.Content = append(n.Content, scalar("releaseDate"), singleQuotedScalar(t.ReleaseDate))
	}
	if len(t.Tags) > 0 {
		n.Content = append(n.Content, scalar("tags"), tagsSequence(t.Tags))
	}
	return n, nil
}
