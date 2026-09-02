---
id: T-0006
kind: ticket
title: Status bar and task/ticket tree
task: Task_0003
status: IN_PROGRESS
priority: 10
tags:
  - ui
depends_on:
  - T-0002
contract_hash: 98405ba0ae5f4932
commits: []
updated_at: "2026-09-02T12:49:18Z"
changelog:
  - date: 2026-09-02
    kind: created
    summary: Created with task Task_0003
  - date: 2026-09-02
    kind: contract-change
    summary: Detail pane now renders section markdown through glamour instead of printing it raw, per the architecture update adding that library. Updated Scope, Out Of Scope, Contracts, and Design Notes; T-0007's contract restamped to match.
  - date: 2026-09-02
    kind: "status:IN_PROGRESS"
    summary: Selected by agentism next
---

# Status bar and task/ticket tree

## Work

The `ui.Model`'s `Init`, `Update`, and `View` for the top half of the
screen: the status bar (phase, version, ticket-status counts) and the
task/ticket tree with cursor movement and expand/collapse. No detail
pane and no log view yet; those are T-0007.

## Contract

### Inputs

- A `*store.Project`, passed into `ui.New(project *store.Project)
  tea.Model`.
- `tea.KeyMsg` for up/down/j/k (move the cursor) and enter/space
  (expand or collapse the task under the cursor).
- `tea.WindowSizeMsg` for the terminal size.

### Outputs

- `View()` renders the status bar as one line: phase, version, and the
  `DONE`/total ticket count plus any non-zero non-DONE status counts.
- `View()` renders the tree below it: one line per task, its tickets
  indented under it when the user expands the task, each line showing
  an id, a title, and a status.
- The cursor highlights the current line. Moving past the last line or
  before the first line does nothing (no wraparound).

### Errors

- None. This ticket does not read files; it only renders the
  `*store.Project` the caller passed in.

## Tests First

- `View()` on a project with two tasks (one expanded, one collapsed)
  shows the expanded task's tickets indented and the collapsed task's
  tickets absent.
- Pressing down `n` times moves the cursor `n` lines, stopping at the
  last visible line.
- Pressing enter on a task line toggles that task between expanded and
  collapsed.
- The status bar text contains the phase and version from the
  `*store.Project` and a count that matches `Project.Counts()`.

## Acceptance

Running a small test harness (or the ticket's own tests) against a
sample `*store.Project` shows the phase, the version, the ticket
counts, and an expandable tree matching that project's tasks and
tickets.

## Human Verification Steps

1. Build the binary and run it in this repository's root (it has no
   tickets yet at the time this ticket ships, so run it against a
   fixture project with at least one task and one ticket instead, or
   revisit after later tickets exist here).
2. Confirm the status bar shows `TASKS` (or whatever phase the fixture
   project is in) and the right ticket counts.
3. Press down/up (or j/k) and confirm the highlighted line moves.
4. Press enter on a task line and confirm its tickets appear; press
   enter again and confirm they disappear.
