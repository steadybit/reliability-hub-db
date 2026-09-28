package dbmodel

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func scalar(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Value: value}
}

func literalScalar(value string) *yaml.Node {
	// yaml.v3 picks the chomp indicator from trailing newlines: 0 → "|-" (strip),
	// 1 → "|" (clip), 2+ → "|+" (keep). Existing description.yml files use "|",
	// so we ensure exactly one trailing newline.
	for len(value) > 0 && value[len(value)-1] == '\n' {
		value = value[:len(value)-1]
	}
	return &yaml.Node{Kind: yaml.ScalarNode, Value: value + "\n", Style: yaml.LiteralStyle}
}

func singleQuotedScalar(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Value: value, Style: yaml.SingleQuotedStyle}
}

func tagsSequence(tags []string) *yaml.Node {
	seq := &yaml.Node{Kind: yaml.SequenceNode}
	for _, t := range tags {
		seq.Content = append(seq.Content, scalar(t))
	}
	return seq
}

func targetTypesSequence(ts []string) *yaml.Node {
	seq := &yaml.Node{Kind: yaml.SequenceNode}
	for _, t := range ts {
		seq.Content = append(seq.Content, scalar(t))
	}
	return seq
}

// Marshal encodes a value to YAML using yaml.v3 with a 2-space indent and a
// leading "---" document marker, matching the existing file convention in
// reliability-hub-db.
func Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("---\n")
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
