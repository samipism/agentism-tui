---
id: T-0001
kind: ticket
title: Frontmatter and section parser
task: Task_0001
status: IN_PROGRESS
priority: 10
tags:
  - parser
depends_on: []
contract_hash: 58780ac705de007e
commits: []
updated_at: "2026-09-02T12:02:07Z"
changelog:
  - date: 2026-09-02
    kind: created
    summary: Created with task Task_0001
  - date: 2026-09-02
    kind: "status:IN_PROGRESS"
    summary: Selected by agentism next
---

# Frontmatter and section parser

## Work

Two small functions that every later ticket in this task builds on:
one that splits a markdown file into its YAML frontmatter and its
body, and one that pulls one named section's text out of that body.

## Contract

### Inputs

- `ParseDoc(text string) (Doc, error)` takes a whole file's text.
- `Doc.Section(name string) (string, bool)` takes a heading name, for
  example `"Contracts"`.

### Outputs

- `Doc{ Meta map[string]any; Body string }`. `Meta` holds the parsed
  frontmatter (empty map when the file has none).
- `Section` returns the text between the named heading and the next
  heading at the same or a shallower level, trimmed, plus `true`. It
  returns `"", false` when no heading matches.
- A heading matches case-insensitively and ignores a leading number or
  punctuation prefix, so `"## Contracts"` and `"## 1. contracts"` both
  match `Section("Contracts")`.

### Errors

- Frontmatter between two `---` lines that is not valid YAML: `Doc{}`
  and a non-nil error naming the problem.
- No closing `---` line: not an error. The whole file becomes `Body`,
  and `Meta` is empty, matching agentism's own `mdfile.ts` behaviour.
- A file that does not start with `---`: not an error, same as above.

## Tests First

- Frontmatter with a scalar, a flat list, and a list of flat mappings
  parses into the matching Go types.
- A file with no frontmatter returns an empty `Meta` and the full text
  as `Body`.
- A file with an unterminated frontmatter delimiter returns an empty
  `Meta` and the full text as `Body`.
- Frontmatter that breaks YAML syntax returns a non-nil error.
- `Section("Contracts")` on a body with `## Contracts` followed by
  `### Inputs` and `### Outputs` returns all three headings' text
  together, matching `mdfile.ts`'s nesting rule.
- `Section("contracts")` (lower case) and `Section("1. Contracts")`
  match the same heading as `Section("Contracts")`.
- `Section("Missing")` returns `"", false`.

## Acceptance

`go test ./internal/store/... -run TestParseDoc` and `-run TestSection`
both pass, covering every case above.

## Human Verification Steps

1. Run `go test ./internal/store/... -run TestParseDoc -v` in the
   repository root. Expect every subtest to print `PASS`.
2. Run `go test ./internal/store/... -run TestSection -v`. Expect
   every subtest to print `PASS`.
