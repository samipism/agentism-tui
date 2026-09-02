---
id: Task_0003
kind: task
title: Bubble Tea dashboard
status: PLANNED
priority: 30
tags:
  - ui
locked: true
depends_on: []
version: v0
updated_at: "2026-09-02T10:24:36Z"
changelog:
  - date: 2026-09-02
    kind: created
    summary: Task created from the intent and the architecture
  - date: 2026-09-02
    kind: locked
    summary: Planning finished
  - date: 2026-09-02
    kind: updated
    summary: Detail pane now renders section markdown through glamour instead of printing it raw, per the architecture update adding that library. Updated Scope, Out Of Scope, Contracts, and Design Notes; T-0007's contract restamped to match.
---

# Bubble Tea dashboard

## Goal

The single-screen dashboard: a status bar, a task/ticket tree, and a
detail pane, driven by a `store.Project` and Bubble Tea's Model-Update-
View loop. View-only: no key ever changes project data.

## Scope

- A `ui.Model` holding the loaded `store.Project`, the tree's
  cursor and expand/collapse state, and which ticket (if any) is
  selected for the detail pane.
- Status bar: current phase, version, and the ticket-status counts
  from `store.Project.Counts()`.
- Task/ticket tree: one line per task, its tickets indented under it
  when expanded. Each line shows an id, a title, and a status.
- Detail pane: fills when the user selects a ticket, showing the
  ticket's title, status, task, and its `Contract`, `Work`, and
  `Acceptance` text, rendered through `glamour` so headings, lists,
  and code blocks show as formatted text, not raw markdown source.
- Activity log view: a scrollable list built from `store.Project.Log`.
- Keybindings: up/down or j/k move the cursor; enter or space
  expands/collapses a task, or selects a ticket; 'r' re-runs
  `store.Load` and replaces the model's project data; 'l' toggles the
  log view; 'q' or Ctrl+C quits.

## Out Of Scope

- Any key or action that writes a file or calls the `agentism` binary.
- Any project path other than the one `main.go` already loaded.
- Mouse support.
- A markdown feature `glamour` does not already render on its own
  (for example, a live table of contents, or clickable links).

## Contracts

<!-- The source of truth. A ticket must not contradict this section.
     A change here marks finished tickets as STALE. -->

### Inputs

- A `*store.Project` (from Task_0001), passed in when the model is
  built.
- Terminal key events, delivered by Bubble Tea as `tea.KeyMsg`.
- The terminal window size, delivered as `tea.WindowSizeMsg`.

### Outputs

- A rendered frame (a string) on every `View()` call: status bar on
  top, the tree and the detail pane side by side below it.
- The detail pane's text runs through `glamour.Render` before display,
  using a style that matches a dark or a light terminal background.
- On 'r', a re-read `store.Project` replaces the model's data and the
  next frame reflects it. A load error on refresh shows as a one-line
  message in the status bar; it does not crash the program and it
  keeps showing the last good data.

### Errors

- A `store.Load` error during a refresh never exits the program. It
  shows in the status bar and the previous data stays on screen.

### Invariants

- No keybinding calls a function that writes a file or spawns a
  process.
- The tree always shows every task from `Project.Tasks`, in the order
  `store.Load` returns (priority then id).

## Acceptance

Running the binary in a project with tasks and tickets shows the
status bar, an expandable tree, and a detail pane that fills in on
selecting a ticket. Pressing 'r' after an external change to the
project (for example a ticket's status file edited by hand) shows the
new status without restarting the program. Pressing 'q' exits cleanly.

## Design Notes

Uses `bubbles/list` for the tree (a flat list where a task's ticket
rows insert or remove on expand/collapse) and `bubbles/viewport` for
the scrollable detail pane and log view. `lipgloss` styles the borders
and the status colors (for example, `DONE` in one color, `BLOCKED` in
another). `glamour` renders the detail pane's markdown text before it
goes into the viewport; the viewport itself stays a plain scroll
container over already-rendered text.

## Resolved Questions

- Q: Does a ticket's status get a color, and if so from what palette?
  A: Yes, a small fixed color per status (for example green for DONE,
  red for BLOCKED and NOT_ACCEPTED, yellow for IN_REVIEW). Exact hex
  values are a detail for the ticket that implements styling, not an
  architecture decision.
- Q: Should the detail pane show raw markdown or render it?
  A: Render it, with `glamour` (added to ArchDecision.md in this
  update). The human asked for formatted markdown in the terminal.
