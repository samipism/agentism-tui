---
id: Task_0005
kind: task
title: Tree view redesign
status: COMPLETE
priority: 10
tags:
  - ui
locked: true
depends_on: []
version: v1
updated_at: "2026-09-02T17:11:50Z"
changelog:
  - date: 2026-09-02
    kind: created
    summary: Task created from the intent and the architecture
  - date: 2026-09-02
    kind: locked
    summary: Planning finished
  - date: 2026-09-02
    kind: complete
    summary: Every ticket is DONE
---

# Tree view redesign

## Goal

The sidebar tree view shows the cursor row as a full-width highlight band.
It shows a task row with an expand or collapse arrow. It shows a ticket row
with a branch glyph, indented under its task. Every row fits inside the
sidebar width. A row never wraps to a second line.

## Scope

- `Model.treeView` and a new `Model.renderRow` in `internal/ui/model.go`.
- The cursor style: a full-width background band instead of a leading `> `
  marker.
- The row icons: `▸` or `▾` for a task, `└` for a ticket.
- Row truncation to the sidebar's inner width.

## Out Of Scope

- The main pane and its scrolling.
- The footer.
- Status colors and the task completion bar. These stay as they are.

## Contracts

<!-- The source of truth. A ticket must not contradict this section.
     A change here marks finished tickets as STALE. -->

### Inputs

- `Model.treeView(width int)` - the sidebar's inner content width, in
  columns.
- `Model.rows()` - the visible rows, each a task or a ticket.
- `Model.cursor` - the index of the selected row.
- `Model.expanded` - which tasks show their tickets.

### Outputs

- One rendered line per row, joined by `\n`.
- The cursor's row: one line, bold, with a full-width background band.
  `renderRow` truncates the text to `width` first, then draws the band at
  `width`.
- A task row: a `▸` (collapsed) or `▾` (expanded) icon, bold text, and its
  completion bar appended when the task has one.
- A ticket row: a `└` icon, indented one level under its task.

### Errors

- None. `treeView` and `renderRow` are pure rendering. They render
  whatever `width` and rows the caller passes.

### Invariants

- A rendered row never exceeds `width` visible columns. `renderRow`
  truncates it; it never wraps the row, so the sidebar's border never
  breaks.
- The cursor's highlight band spans the full `width`, regardless of how
  short or long the row's label is.
- A ticket row's text starts to the right of its task's text, at every
  width the layout can size the sidebar to.

## Acceptance

- In the running TUI, moving the cursor up and down shows a solid
  background band across the full sidebar width on the selected row only.
- A ticket title longer than the sidebar width shows as one truncated line.
  It does not wrap and does not push the border out of shape.
- At the narrowest sidebar width the layout allows, long titles still
  truncate cleanly and the tree keeps its shape.

## Design Notes

`lipgloss.Style.Width` pads short text but word-wraps text that overflows;
it does not cut it. `renderRow` calls `MaxWidth` first to truncate, then
applies `Width` (for the cursor band) or leaves the truncated text as is.

## Resolved Questions

- Q: Is there more UI work planned for v1 beyond this rendering change?
  A: No. This task's scope is exactly the tree view rendering change.
- Q: Any edge case to call out beyond the existing `model_test.go` suite?
  A: Yes - explicitly check narrow sidebar widths and long titles for clean
     truncation.
