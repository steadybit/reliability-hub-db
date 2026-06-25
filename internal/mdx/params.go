// Package mdx regenerates the `# Parameters` Markdown table inside each
// action's summary.mdx, leaving everything else (Introduction, Use Cases,
// YouTube embeds, links) untouched.
package mdx

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/steadybit/reliability-hub-db/internal/describe"
)

// ReplaceResult tells the caller what happened. The sync tool uses it to
// classify per-file outcomes in the PR body.
type ReplaceResult int

const (
	// ReplaceUpdated means an existing parameter table was rewritten.
	ReplaceUpdated ReplaceResult = iota
	// ReplaceUnchanged means the existing table already matched the desired
	// output; the file was not rewritten.
	ReplaceUnchanged
	// ReplaceNoHeading means no `# Parameters` heading was found; the file is
	// left alone.
	ReplaceNoHeading
	// ReplaceNoTable means a heading was found but no Markdown table beneath
	// it; the file is left alone.
	ReplaceNoTable
)

func (r ReplaceResult) String() string {
	switch r {
	case ReplaceUpdated:
		return "updated"
	case ReplaceUnchanged:
		return "unchanged"
	case ReplaceNoHeading:
		return "no-parameters-heading"
	case ReplaceNoTable:
		return "no-parameters-table"
	}
	return "unknown"
}

var (
	// Match a Parameters heading. Allow any heading depth (# / ## / ###) since
	// existing files all use `# Parameters` but we want to be forgiving.
	headingRE = regexp.MustCompile(`(?m)^#+\s+Parameters\s*$`)
	// Match any heading line (for finding the next section after Parameters).
	anyHeadingRE = regexp.MustCompile(`(?m)^#+\s+\S`)
)

// ReplaceParameters rewrites the parameter table inside `content` to match
// `params`. The table is the first run of consecutive `|`-prefixed lines
// beneath the `# Parameters` heading. Returns the new content (or unchanged
// content) and a classification.
func ReplaceParameters(content string, params []describe.Parameter) (string, ReplaceResult, error) {
	headingLocs := headingRE.FindAllStringIndex(content, -1)
	if len(headingLocs) == 0 {
		return content, ReplaceNoHeading, nil
	}
	if len(headingLocs) > 1 {
		return content, ReplaceNoHeading, fmt.Errorf("multiple Parameters headings found")
	}
	// Find the line right after the heading line.
	headingEnd := headingLocs[0][1]
	// Locate where the next section starts (next `#+ ` heading after this one)
	// or EOF.
	sectionEnd := len(content)
	if next := anyHeadingRE.FindStringIndex(content[headingEnd:]); next != nil {
		sectionEnd = headingEnd + next[0]
	}

	// Within content[headingEnd:sectionEnd], find the run of `|`-prefixed lines.
	section := content[headingEnd:sectionEnd]
	tableStart, tableEnd, ok := findTable(section)
	if !ok {
		return content, ReplaceNoTable, nil
	}

	newTable := FormatParameterTable(params)
	// findTable returns tableEnd just past the trailing newline of the last
	// table row, and FormatParameterTable emits a trailing newline as well —
	// so splice as-is to preserve the surrounding line structure.
	rebuilt := content[:headingEnd] + section[:tableStart] + newTable + section[tableEnd:] + content[sectionEnd:]
	if rebuilt == content {
		return content, ReplaceUnchanged, nil
	}
	return rebuilt, ReplaceUpdated, nil
}

// findTable returns the byte offsets of the first run of `|`-prefixed lines
// within s. The returned range includes the trailing newline of the last
// table row.
func findTable(s string) (start, end int, ok bool) {
	i := 0
	for i < len(s) {
		lineStart := i
		nl := strings.IndexByte(s[i:], '\n')
		var lineEnd int
		if nl < 0 {
			lineEnd = len(s)
		} else {
			lineEnd = i + nl
		}
		line := s[lineStart:lineEnd]
		if strings.HasPrefix(line, "|") {
			// Found start of table. Scan forward for the end.
			start = lineStart
			end = lineEnd
			if nl >= 0 {
				end++ // include trailing newline
			}
			j := end
			for j < len(s) {
				nlj := strings.IndexByte(s[j:], '\n')
				var je int
				if nlj < 0 {
					je = len(s)
				} else {
					je = j + nlj
				}
				if !strings.HasPrefix(s[j:je], "|") {
					break
				}
				end = je
				if nlj >= 0 {
					end++
				}
				j = end
			}
			return start, end, true
		}
		if nl < 0 {
			break
		}
		i = lineEnd + 1
	}
	return 0, 0, false
}

// FormatParameterTable produces a 3-column Markdown table (Parameter |
// Description | Default) with aligned borders, matching the existing
// reliability-hub-db convention.
func FormatParameterTable(params []describe.Parameter) string {
	headers := []string{"Parameter", "Description", "Default"}
	rows := make([][]string, 0, len(params))
	for _, p := range params {
		rows = append(rows, []string{
			paramLabel(p),
			paramDescription(p),
			paramDefault(p),
		})
	}
	return renderAlignedTable(headers, rows)
}

func paramLabel(p describe.Parameter) string {
	if p.Label != "" {
		return escapePipe(p.Label)
	}
	return escapePipe(p.Name)
}

func paramDescription(p describe.Parameter) string {
	d := strings.TrimSpace(p.Description)
	if p.Required == nil || !*p.Required {
		if d == "" {
			d = "(optional)"
		} else {
			d = "(optional) " + d
		}
	}
	return escapePipe(d)
}

func paramDefault(p describe.Parameter) string {
	if p.DefaultValue == nil {
		return ""
	}
	return escapePipe(*p.DefaultValue)
}

func escapePipe(s string) string {
	return strings.ReplaceAll(s, "|", `\|`)
}

func renderAlignedTable(headers []string, rows [][]string) string {
	cols := len(headers)
	width := make([]int, cols)
	for i, h := range headers {
		if l := len(h); l > width[i] {
			width[i] = l
		}
	}
	for _, r := range rows {
		for i, c := range r {
			if l := len(c); l > width[i] {
				width[i] = l
			}
		}
	}
	var b strings.Builder
	writeData := func(cells []string) {
		b.WriteByte('|')
		for i, c := range cells {
			b.WriteByte(' ')
			b.WriteString(c)
			b.WriteString(strings.Repeat(" ", width[i]-len(c)))
			b.WriteByte(' ')
			b.WriteByte('|')
		}
		b.WriteByte('\n')
	}
	writeSep := func() {
		b.WriteByte('|')
		for i := 0; i < cols; i++ {
			b.WriteString(strings.Repeat("-", width[i]+2))
			b.WriteByte('|')
		}
		b.WriteByte('\n')
	}
	writeData(headers)
	writeSep()
	for _, r := range rows {
		writeData(r)
	}
	return b.String()
}
