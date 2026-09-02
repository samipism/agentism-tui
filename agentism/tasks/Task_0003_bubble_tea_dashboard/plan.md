---
id: Task_0003
kind: task
title: Bubble Tea dashboard
status: COMPLETE
priority: 30
tags:
  - ui
locked: true
depends_on: []
version: v0
updated_at: "2026-09-02T16:06:00Z"
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
  - date: 2026-09-02
    kind: complete
    summary: Every ticket is DONE
  - date: 2026-09-02
    kind: updated
    summary: "Dashboard read as plain text, not a dashboard. Replaced the status-bar-plus-two-columns layout with four bordered regions (header, sidebar, main, footer), added a fixed status-to-color lookup, and added a completion bar in the header and per task row. No new library: lipgloss already covered borders and color."
  - date: 2026-09-02
    kind: complete
    summary: Every ticket is DONE
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
- Header region: current phase, version, the ticket-status counts
  from `store.Project.Counts()`, and an overall completion bar
  (DONE tickets over total tickets).
- Sidebar region: the task/ticket tree, one line per task, its
  tickets indented under it when expanded. Each line shows an id, a
  title, and a status. A task line also shows a small completion bar
  for that task's own tickets. A ticket's or a task's status text
  renders in a fixed color for its status value.
- Main region: the detail pane or the log view, whichever is active.
- Detail pane: fills when the user selects a ticket, showing the
  ticket's title, status, task, and its `Contract`, `Work`, and
  `Acceptance` text, rendered through `glamour` so headings, lists,
  and code blocks show as formatted text, not raw markdown source.
- Activity log view: a scrollable list built from `store.Project.Log`.
- Footer region: a one-line keybinding hint (the active keys and
  what each one does).
- Every region (header, sidebar, main, footer) draws inside its own
  bordered box.
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
- Resizing, moving, or hiding a region. The four regions (header,
  sidebar, main, footer) always show, in a fixed layout.
- A configurable color palette or theme file. Colors are fixed in
  code, not read from a config or an environment variable.

## Contracts

<!-- The source of truth. A ticket must not contradict this section.
     A change here marks finished tickets as STALE. -->

### Inputs

- A `*store.Project` (from Task_0001), passed in when the model is
  built.
- Terminal key events, delivered by Bubble Tea as `tea.KeyMsg`.
- The terminal window size, delivered as `tea.WindowSizeMsg`.

### Outputs

- A rendered frame (a string) on every `View()` call: a bordered
  header region on top, a bordered footer region on the bottom, and
  a bordered sidebar (the tree) and a bordered main region (the
  detail pane or the log view) side by side between them.
- Each status value (a ticket's or a task's) renders in one fixed
  color for that value: green for DONE, yellow for IN_PROGRESS and
  IN_REVIEW, red for BLOCKED and NOT_ACCEPTED, a neutral color for
  every other value (TODO, PLANNED, STALE, and so on).
- The header shows an overall completion bar: a string of filled and
  empty blocks plus a percentage, built from `Project.Counts()`. Each
  task row in the tree shows the same kind of bar, built from that
  task's own ticket counts.
- The detail pane's text runs through `glamour.Render` before display,
  using a style that matches a dark or a light terminal background.
- On 'r', a re-read `store.Project` replaces the model's data and the
  next frame reflects it. A load error on refresh shows as a one-line
  message in the header region; it does not crash the program and it
  keeps showing the last good data.

### Errors

- A `store.Load` error during a refresh never exits the program. It
  shows in the header region and the previous data stays on screen.

### Invariants

- No keybinding calls a function that writes a file or spawns a
  process.
- The tree always shows every task from `Project.Tasks`, in the order
  `store.Load` returns (priority then id).
- Every color and every border comes from a fixed `lipgloss` style
  defined in code. There is no runtime color or theme configuration.

## Acceptance

Running the binary in a project with tasks and tickets shows four
bordered regions: a header with the phase, version, counts, and an
overall completion bar; a sidebar tree with colored statuses and a
per-task completion bar; a main pane that fills with a rendered
ticket or the log view; and a footer that names the active keys.
Pressing 'r' after an external change to the project (for example a
ticket's status file edited by hand) shows the new status without
restarting the program. Pressing 'q' exits cleanly.

## Design Notes

Uses `bubbles/list` for the tree (a flat list where a task's ticket
rows insert or remove on expand/collapse) and `bubbles/viewport` for
the scrollable detail pane and log view. `lipgloss` draws a bordered
box per region and joins the four with `JoinHorizontal` (sidebar next
to main) and `JoinVertical` (header, that row, footer). A small
`status -> lipgloss.Color` lookup gives every status its color. A
small `progressBar(done, total int) string` helper builds the block
string used by both the header's overall bar and each task row's bar,
so the two never drift apart. `glamour` renders the detail pane's
markdown text before it goes into the viewport; the viewport itself
stays a plain scroll container over already-rendered text.

## Resolved Questions

- Q: Does a ticket's status get a color, and if so from what palette?
  A: Yes, a small fixed color per status (for example green for DONE,
  red for BLOCKED and NOT_ACCEPTED, yellow for IN_REVIEW). Exact hex
  values are a detail for the ticket that implements styling, not an
  architecture decision.
- Q: Should the detail pane show raw markdown or render it?
  A: Render it, with `glamour` (added to ArchDecision.md in this
  update). The human asked for formatted markdown in the terminal.
- Q (2026-09-02 update): The dashboard reads as plain text, not a
  dashboard. What should change, and how far?
  A: A full multi-region layout: a bordered header, sidebar, main
  pane, and footer, replacing the single status-bar-plus-two-columns
  layout.
- Q (2026-09-02 update): Should status colors and completion bars
  land in this same update, or wait for a later one?
  A: Land them together with the layout change. One task update, one
  round of re-verification, instead of two passes over the same
  files.
