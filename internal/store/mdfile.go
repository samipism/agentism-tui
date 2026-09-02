// Package store loads an agentism project (its config, log, tasks, and
// tickets) from the markdown and JSON files on disk.
package store

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Doc is a markdown file split into its YAML frontmatter and its body.
type Doc struct {
	Meta map[string]any
	Body string
}

// ParseDoc splits text into frontmatter and body. A file that does not
// start with a "---" delimiter, or whose frontmatter has no closing
// delimiter, is not an error: the whole file becomes Body and Meta is
// empty, matching agentism's own mdfile.ts behaviour.
func ParseDoc(text string) (Doc, error) {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || lines[0] != "---" {
		return Doc{Meta: map[string]any{}, Body: text}, nil
	}

	closeIdx := -1
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			closeIdx = i
			break
		}
	}
	if closeIdx == -1 {
		return Doc{Meta: map[string]any{}, Body: text}, nil
	}

	meta := map[string]any{}
	fm := strings.Join(lines[1:closeIdx], "\n")
	if strings.TrimSpace(fm) != "" {
		if err := yaml.Unmarshal([]byte(fm), &meta); err != nil {
			return Doc{}, fmt.Errorf("parsing frontmatter: %w", err)
		}
	}

	body := strings.Join(lines[closeIdx+1:], "\n")
	return Doc{Meta: meta, Body: body}, nil
}

var (
	headingRe       = regexp.MustCompile(`^(#{1,6})\s+(.*)$`)
	headingPrefixRe = regexp.MustCompile(`^[\d\p{P}]+\s*`)
)

// normalizeHeading strips a leading number or punctuation prefix (as in
// "1. Contracts") and lower-cases the rest, so heading lookup is
// case-insensitive and ignores that prefix.
func normalizeHeading(s string) string {
	s = headingPrefixRe.ReplaceAllString(strings.TrimSpace(s), "")
	return strings.ToLower(strings.TrimSpace(s))
}

// Section returns the text under the named heading, up to the next
// heading at the same or a shallower level. It returns "", false when no
// heading matches name (case-insensitively, ignoring a leading number or
// punctuation prefix on either side).
func (d Doc) Section(name string) (string, bool) {
	target := normalizeHeading(name)
	lines := strings.Split(d.Body, "\n")

	for i, line := range lines {
		m := headingRe.FindStringSubmatch(line)
		if m == nil || normalizeHeading(m[2]) != target {
			continue
		}
		level := len(m[1])

		end := len(lines)
		for j := i + 1; j < len(lines); j++ {
			if mj := headingRe.FindStringSubmatch(lines[j]); mj != nil && len(mj[1]) <= level {
				end = j
				break
			}
		}
		return strings.TrimSpace(strings.Join(lines[i+1:end], "\n")), true
	}
	return "", false
}
