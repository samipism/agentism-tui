---
kind: architecture
version: v0
status: FINAL
updated_at: "2026-09-02T10:23:11Z"
tags: []
changelog:
  - date: 2026-09-02
    kind: created
    summary: First version from the architecture questions
  - date: 2026-09-02
    kind: finalised
    summary: Human approved the document
  - date: 2026-09-02
    kind: updated
    summary: Fixed the Tags section; a wrapped continuation line had been misread as a bogus tag.
  - date: 2026-09-02
    kind: updated
    summary: "Fixed the Tags section: a wrapped continuation line under the parser tag had been misread as a bogus 'section' tag by the tag extractor. Collapsed to one line per tag; tag list is now the intended five: parser, ui, cli, release, docs."
  - date: 2026-09-02
    kind: updated
    summary: Added glamour (github.com/charmbracelet/glamour) to render a ticket or task section's markdown into styled terminal text in the detail pane, replacing the earlier plain-text-only call. Reverses the no-rendering decision from the first pass.
---

# Architecture Decisions

## Domain Education

**Frontmatter.** A block of `key: value` lines between two `---` lines at
the top of a markdown file. agentism stores task and ticket metadata (id,
status, tags) there. It is a small YAML subset: scalars, flat lists, and
lists of flat mappings.

**Bubble Tea / the Elm architecture.** A Bubble Tea program holds one
`Model` struct. `Update` takes the current model and an event (a key press,
a tick) and returns a new model. `View` turns the model into the text the
terminal shows. The loop repeats. This keeps all state in one place and
makes the screen a pure function of that state.

**GoReleaser.** A tool that reads one YAML file and, for each target
platform, compiles the binary, packs it into an archive, and uploads it to
a GitHub Release. It replaces a hand-written build matrix.

**Glamour.** A library that turns markdown text into styled terminal
output: a heading becomes bold and colored, a list gets its bullets, a
code block gets a background and syntax color. It picks a light or a
dark style to match the terminal automatically.

## Language And Runtime

Go, latest stable release line (1.23 or newer). A single static binary,
`CGO_ENABLED=0`, no runtime dependency beyond the OS.

## Frameworks And Libraries

- `github.com/charmbracelet/bubbletea` - the Model-Update-View loop and
  the terminal event loop.
- `github.com/charmbracelet/bubbles` - the list component for the task and
  ticket tree, the viewport component for the scrollable detail pane.
- `github.com/charmbracelet/lipgloss` - borders, color, and layout for the
  status bar, tree, and detail pane.
- `gopkg.in/yaml.v3` - parses task and ticket frontmatter. Agentism's
  frontmatter is a subset of real YAML, so a full YAML parser reads it
  without a hand-written one.
- `github.com/charmbracelet/glamour` - renders a section's markdown
  (headings, lists, code blocks, bold/italic) into styled terminal
  text for the detail pane, auto-matching the terminal's light or dark
  style.
- `encoding/json` (standard library) - parses `.agentism/config.json` and
  `.agentism/log.jsonl`.

No other third-party library. No HTTP client, no database driver: this
tool only reads local files.

## Data Storage

None. agentism-tui holds no state of its own. On launch, and again when
the user presses 'r', it reads:

- `.agentism/config.json` - the project phase, version, and tag list.
- `.agentism/log.jsonl` - the activity log, one JSON object per line.
- `agentism/tasks/Task_NNNN_<slug>/plan.md` - one task's frontmatter
  (id, title, status, priority, tags, depends_on) and its `Contracts`
  section.
- `agentism/tasks/Task_NNNN_<slug>/T-NNNN-<slug>.md` - one ticket's
  frontmatter (id, title, task, status, priority, tags, depends_on,
  contract_hash) and its `Contract`, `Work`, and `Acceptance` sections.
- `agentism/versions/<version>/Task_NNNN_<slug>/` - archived tasks, same
  file shape as above.

agentism-tui never reads `.agentism/state.json`. That file is the CLI's
own cache and can be stale; agentism-tui rebuilds the same tree from the
markdown files every time it reads, which is what the CLI's own
`reconcile()` does. This keeps one code path and matches FR-7's rule
against trusting a cache.

## Design Patterns

One loader function, `store.Load(root string) (*Project, error)`, walks
the two folders above and returns a plain struct tree (`Project` ->
`Task` -> `Ticket`). No interface, no repository layer: there is exactly
one data source and one implementation, so an interface would only add a
name to jump through.

The UI follows Bubble Tea's Model-Update-View shape, because that is how
the chosen framework works, not an added pattern.

## Test Strategy

- **Unit tests** (`go test ./...`) cover the loader: frontmatter parsing,
  section extraction, task/ticket assembly, and the ticket-status counts
  used by the status bar. Table-driven tests with fixture files under
  `internal/store/testdata/`.
- **No integration or contract tests.** There is one data source (the
  local filesystem) and one consumer (the loader); a unit test against
  fixture files covers that boundary.
- **No end to end test.** A human runs the built binary against a real
  agentism project and checks the screen during verification. This is
  not agent work.
- **Gate.** Before a ticket goes to review: `gofmt -l` reports no files,
  `go vet ./...` is clean, and `go test ./...` passes.

## Tooling

- `gofmt` - formatting.
- `go vet` - static checks. Both ship with the Go toolchain; this project
  skips a separate linter, because these two catch what its size needs.
- `go build ./...` - compile check.
- `goreleaser build --snapshot --clean` - local dry run of the release
  build, without publishing.

## Tags

- parser - reading agentism files, frontmatter, and task/ticket assembly
- ui - the Bubble Tea screens, the tree, the detail pane, styling
- cli - argument parsing, the entry point, the no-project error path
- release - the GoReleaser config and the release workflow
- docs - the README and usage notes

## Repository Layout

```
agentism-tui/
  go.mod
  cmd/agentism-tui/
    main.go          -- entry point: find project, load, run the TUI
  internal/store/
    store.go          -- Load(root) walks .agentism/ and agentism/
    mdfile.go          -- frontmatter + section parsing
    types.go            -- Project, Task, Ticket, LogEntry structs
    store_test.go
    testdata/
  internal/ui/
    model.go            -- Bubble Tea Model, Update, View
    style.go             -- lipgloss styles
  .goreleaser.yaml
  .github/workflows/release.yml
  README.md
```

## Decision Record

| Decision | Rejected | Reason |
|---|---|---|
| Bubble Tea + Bubbles + Lip Gloss | tview | Idiomatic Go, active project, matches the loop-based UI this tool needs |
| gopkg.in/yaml.v3 for frontmatter | Hand-written YAML parser | Agentism's frontmatter is a YAML subset; a real parser reads it, no need to port the TypeScript one |
| Read markdown files directly, ignore state.json | Reading `.agentism/state.json` | That file is a cache the CLI itself may not have refreshed; reading the source markdown is the only way to guarantee current data |
| glamour for the detail pane | Plain wrapped text (the original call) | The human asked for markdown to read as formatted text, not raw source; glamour is Charm's own renderer, so it matches the rest of the stack |
| GoReleaser + GitHub Actions | Hand-written build matrix | One YAML file instead of hand-rolled matrix scripting, standard tool for this exact job |
| No repository/interface layer | A `Store` interface with one implementation | One data source; an interface here has no second implementation to justify it |

## Task Sketch

Build order for the TASKS phase, not yet created:

1. **parser** (tag: parser) - Load `.agentism/config.json`,
   `.agentism/log.jsonl`, and the `agentism/tasks` tree into `Project`,
   `Task`, and `Ticket` structs. Unit tests against fixture files.
2. **cli** (tag: cli) - `main.go`: find the project in the current
   directory, print the no-project error and exit 1 when it is missing,
   otherwise call the loader and start the UI.
3. **ui** (tag: ui) - The Bubble Tea dashboard: status bar, task/ticket
   tree, detail pane, and the keybindings (arrows/j-k, enter/space to
   expand, 'r' to refresh, 'q' to quit).
4. **release** (tag: release) - `.goreleaser.yaml`, the GitHub Actions
   release workflow, and the README.

## Open Questions
