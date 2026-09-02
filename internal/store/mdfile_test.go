package store

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseDoc(t *testing.T) {
	t.Run("scalar, flat list, and list of flat mappings", func(t *testing.T) {
		text := "---\n" +
			"id: T-0001\n" +
			"tags: [parser, ui]\n" +
			"changelog:\n" +
			"  - date: 2026-09-02\n" +
			"    kind: created\n" +
			"  - date: 2026-09-03\n" +
			"    kind: updated\n" +
			"---\n" +
			"body text\n"

		doc, err := ParseDoc(text)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if doc.Meta["id"] != "T-0001" {
			t.Errorf("id = %v, want T-0001", doc.Meta["id"])
		}
		wantTags := []interface{}{"parser", "ui"}
		if !reflect.DeepEqual(doc.Meta["tags"], wantTags) {
			t.Errorf("tags = %#v, want %#v", doc.Meta["tags"], wantTags)
		}
		changelog, ok := doc.Meta["changelog"].([]interface{})
		if !ok || len(changelog) != 2 {
			t.Fatalf("changelog = %#v, want a 2-item list", doc.Meta["changelog"])
		}
		first, ok := changelog[0].(map[string]interface{})
		if !ok || first["kind"] != "created" {
			t.Errorf("changelog[0] = %#v, want kind=created", changelog[0])
		}
		if doc.Body != "body text\n" {
			t.Errorf("Body = %q, want %q", doc.Body, "body text\n")
		}
	})

	t.Run("no frontmatter", func(t *testing.T) {
		text := "# Title\n\nJust body text.\n"
		doc, err := ParseDoc(text)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(doc.Meta) != 0 {
			t.Errorf("Meta = %#v, want empty", doc.Meta)
		}
		if doc.Body != text {
			t.Errorf("Body = %q, want %q", doc.Body, text)
		}
	})

	t.Run("unterminated frontmatter delimiter", func(t *testing.T) {
		text := "---\nid: T-0001\nno closing delimiter here\n"
		doc, err := ParseDoc(text)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(doc.Meta) != 0 {
			t.Errorf("Meta = %#v, want empty", doc.Meta)
		}
		if doc.Body != text {
			t.Errorf("Body = %q, want %q", doc.Body, text)
		}
	})

	t.Run("frontmatter breaks YAML syntax", func(t *testing.T) {
		text := "---\nid: [unclosed\n---\nbody\n"
		doc, err := ParseDoc(text)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
		if !reflect.DeepEqual(doc, Doc{}) {
			t.Errorf("doc = %#v, want zero value", doc)
		}
	})
}

func TestSection(t *testing.T) {
	body := "# Ticket\n\n" +
		"## Contract\n\n" +
		"contract text\n\n" +
		"## Contracts\n\n" +
		"top level contracts text\n\n" +
		"### Inputs\n\n" +
		"inputs text\n\n" +
		"### Outputs\n\n" +
		"outputs text\n\n" +
		"## Acceptance\n\n" +
		"acceptance text\n"

	doc, err := ParseDoc(body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("nested headings come back together", func(t *testing.T) {
		got, ok := doc.Section("Contracts")
		if !ok {
			t.Fatal("Section(\"Contracts\") = false, want true")
		}
		for _, want := range []string{"top level contracts text", "Inputs", "inputs text", "Outputs", "outputs text"} {
			if !strings.Contains(got, want) {
				t.Errorf("Section(%q) = %q, missing %q", "Contracts", got, want)
			}
		}
		if strings.Contains(got, "acceptance text") {
			t.Errorf("Section(%q) = %q, should stop before the next ## heading", "Contracts", got)
		}
	})

	t.Run("case and numeric prefix are ignored", func(t *testing.T) {
		want, _ := doc.Section("Contracts")
		lower, ok := doc.Section("contracts")
		if !ok || lower != want {
			t.Errorf("Section(\"contracts\") = %q, %v; want %q, true", lower, ok, want)
		}
		numbered, ok := doc.Section("1. Contracts")
		if !ok || numbered != want {
			t.Errorf("Section(\"1. Contracts\") = %q, %v; want %q, true", numbered, ok, want)
		}
	})

	t.Run("missing heading", func(t *testing.T) {
		got, ok := doc.Section("Missing")
		if ok || got != "" {
			t.Errorf("Section(\"Missing\") = %q, %v; want \"\", false", got, ok)
		}
	})
}
