---
id: T-0011
kind: ticket
title: Select and render a task's plan.md in the main pane
task: Task_0006
status: IN_PROGRESS
priority: 10
tags:
  - ui
depends_on: []
contract_hash: 05e510848cac71a6
commits: []
updated_at: "2026-09-03T06:20:08Z"
changelog:
  - date: 2026-09-03
    kind: created
    summary: Created with task Task_0006
  - date: 2026-09-03
    kind: "status:IN_PROGRESS"
    summary: Selected by agentism next
---

# Select and render a task's plan.md in the main pane

## Work

Make Enter/Space on a task row also select that task for the main pane,
and add `Model.taskDetailView` to render its Goal and Contract there.

## Contract

### Inputs

- `Model.activateCursorRow()`: for a task row (`r.isTask == true`), keep
  the existing `m.expanded[r.taskID] = !m.expanded[r.taskID]` toggle, and
  add `m.selectedID = r.id`.
- `Model.selectedID`: generalized to hold a task ID or a ticket ID.

### Outputs

- A `Model.selectedTask()` method, mirroring `selectedTicket()`: looks up
  `selectedID` against `m.project.Tasks`, returns `*store.Task` or nil.
- `Model.mainView`: when not showing the log, try `selectedTicket()`
  first; if nil, try `selectedTask()`; if that is also nil, keep the
  current placeholder text.
- `Model.taskDetailView(task *store.Task, width int) string`: builds one
  markdown string with the task's title, status, priority, and its Goal
  and Contract sections, and renders it through the existing
  `renderMarkdown`, the same way `detailView` does for a ticket.

### Errors

- None. Same as `detailView`: pure rendering, no error return. An empty
  Goal or Contract section renders an empty section body, not an error.

## Tests First

- `TestActivateCursorRow_TaskRowSelectsAndToggles`: cursor on a task row,
  call `activateCursorRow`, assert both `m.expanded[taskID]` flipped and
  `m.selectedID` equals the task's ID.
- `TestActivateCursorRow_TicketRowUnchanged`: cursor on a ticket row,
  call `activateCursorRow`, assert `m.selectedID` equals the ticket's ID
  and no task's `expanded` entry changed (existing behavior, regression
  guard).
- `TestMainView_ShowsTaskDetailWhenTaskSelected`: `selectedID` set to a
  task ID, assert `mainView` output contains that task's Goal text.
- `TestMainView_TicketSelectionWinsOverStaleTaskID`: `selectedID` set to
  a ticket ID, assert `mainView` shows the ticket detail (not a task
  detail, not the placeholder).
- `TestTaskDetailView_EmptyGoalAndContractRendersNoError`: a task with
  empty `Goal` and `Contract` fields renders without panicking and
  includes the task's title.

## Acceptance

- Cursor on a task row, press Enter: sidebar expands/collapses that
  task's tickets, and the main pane shows the task's Goal and Contract.
- Cursor on a ticket row, press Enter: main pane shows the ticket detail,
  unchanged from before this ticket.
- `go test ./...` passes, including the new cases above.

## Human Verification Steps

1. Run the TUI (`go run ./cmd/agentism-tui`) against this project.
2. Move the cursor to a task row and press Enter.
   Expect: the task's tickets expand or collapse in the sidebar, and the
   main pane shows that task's Goal and Contract text.
3. Press Enter on the task row again.
   Expect: the tickets collapse or expand back, and the main pane still
   shows that task's detail (selection does not clear on re-toggle).
4. Move the cursor to a ticket row and press Enter.
   Expect: the main pane switches to that ticket's detail (Contract,
   Work, Acceptance), replacing the task detail.
5. Resize the terminal narrower.
   Expect: the task detail text wraps inside the main pane and the
   border does not break.
