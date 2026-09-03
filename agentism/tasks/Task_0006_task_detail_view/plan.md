---
id: Task_0006
kind: task
title: Task detail view
status: COMPLETE
priority: 20
tags:
  - ui
locked: false
depends_on: []
version: v1
updated_at: "2026-09-03T06:25:43Z"
changelog:
  - date: 2026-09-03
    kind: created
    summary: Task created from the intent and the architecture
  - date: 2026-09-03
    kind: updated
    summary: "Added Task_0006: show a task's plan.md (Goal + Contract) in the main pane on Enter/Space, mirroring the existing ticket detail view."
  - date: 2026-09-03
    kind: complete
    summary: Every ticket is DONE
---

# Task detail view

## Goal

Pressing Enter or Space on a task row in the sidebar shows that task's
`plan.md` (its Goal and Contract) in the main pane, the same way Enter
already shows a ticket's detail there.

## Scope

- `Model.activateCursorRow` in `internal/ui/model.go`: a task row also sets
  `Model.selectedID`, alongside the existing expand/collapse toggle.
- `Model.mainView` and a new `Model.taskDetailView`: resolve `selectedID`
  against tasks as well as tickets, and render the task's Goal and
  Contract as markdown, reusing the existing `renderMarkdown` glamour path.

## Out Of Scope

- Editing a task or a ticket from the TUI. This view only reads, like the
  existing ticket detail view.
- Changing what a ticket's detail view shows.
- The tree view's rendering (rows, icons, cursor band). That is
  Task_0005, already done.

## Contracts

<!-- The source of truth. A ticket must not contradict this section.
     A change here marks finished tickets as STALE. -->

### Inputs

- `Model.activateCursorRow()` - called on Enter/Space. For a task row it
  keeps toggling `Model.expanded[taskID]` exactly as today, and also sets
  `Model.selectedID` to the task's ID.
- `Model.selectedID` - generalized to hold either a task ID (`Task_NNNN`)
  or a ticket ID (`T-NNNN`), or `""` for none. The two id formats never
  collide, so one field replaces the ticket-only field without ambiguity.
- `Model.project.Tasks[i].Goal`, `.Contract` - already parsed from the
  task's `plan.md` by `internal/store`. No new parsing.

### Outputs

- `Model.mainView`, when not showing the log, shows in this order: the
  selected ticket's detail if `selectedID` matches a ticket, else the
  selected task's detail if `selectedID` matches a task, else the
  existing placeholder text.
- `Model.taskDetailView` renders the task's title, status, priority, and
  its Goal and Contract sections as one markdown document, through the
  same `renderMarkdown` glamour renderer and `width` handling
  `detailView` already uses.
- Selecting a task replaces any previously selected ticket in the main
  pane, and selecting a ticket replaces any previously selected task:
  only one detail view shows at a time, as today.

### Errors

- None. `taskDetailView` is pure rendering, like `detailView`. A task
  with an empty Goal or Contract section (for example a freshly
  scaffolded task) renders its title and status with an empty section
  body; it does not error or panic.

### Invariants

- Enter/Space on a task row always toggles that task's expand/collapse
  state, exactly as it does today. This task adds the detail view; it
  does not remove or change the expand/collapse behavior.
- Enter/Space on a ticket row keeps behaving exactly as it does today.
- The task detail view wraps to `width` through the same renderer
  `detailView` uses, so it never breaks the main pane's border.

## Acceptance

- With the cursor on a task row, pressing Enter or Space shows that
  task's Goal and Contract text in the main pane, and expands or
  collapses its tickets in the sidebar, same as before.
- With the cursor on a ticket row, pressing Enter or Space still shows
  that ticket's detail, and it replaces any task detail shown before it.
- A task whose Goal or Contract section is empty still renders without
  error, showing its title and status.

## Design Notes

Reuse `selectedID` as one field for either kind of id, instead of adding
a second `selectedTaskID` field: `Task_NNNN` and `T-NNNN` never collide,
so a single lookup that tries tickets first, then tasks, is enough.
`taskDetailView` follows the same shape as `detailView` (build one
markdown string, hand it to `renderMarkdown`), so both share the same
wrapping and color behavior for free.

## Resolved Questions

- Q: Enter/Space already toggles expand/collapse on a task row. Should
  selecting a task for the detail view use that same key, a live preview
  on cursor move, or a new dedicated key?
  A: The same key. Enter/Space on a task both toggles expand/collapse and
  selects the task for the main pane, mirroring exactly how Enter already
  selects a ticket. No new keybinding, no change to how cursor movement
  or ticket selection works today.
