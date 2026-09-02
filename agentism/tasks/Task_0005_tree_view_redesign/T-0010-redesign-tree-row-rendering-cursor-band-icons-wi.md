---
id: T-0010
kind: ticket
title: "Redesign tree row rendering: cursor band, icons, width truncation"
task: Task_0005
status: IN_PROGRESS
priority: 10
tags:
  - ui
depends_on: []
contract_hash: 46009bb78eca4776
commits: []
updated_at: "2026-09-02T17:04:49Z"
changelog:
  - date: 2026-09-02
    kind: created
    summary: Created with task Task_0005
  - date: 2026-09-02
    kind: "status:IN_PROGRESS"
    summary: Selected by agentism next
---

# Redesign tree row rendering: cursor band, icons, width truncation

## Work

Replace the leading `> ` cursor marker with a full-width background band on
the selected row. Give each task row a `▸`/`▾` expand icon and each ticket
row a `└` branch icon. Truncate every row to the sidebar's inner width so a
long title can never wrap and break the tree's shape.

## Contract

### Inputs

- `Model.treeView(width int)`, called with the sidebar's inner width.
- `Model.rows()`, `Model.cursor`, `Model.expanded`.

### Outputs

- `treeView` returns the rows joined by `\n`, one line per row, each no
  wider than `width`.
- `renderRow(r row, selected bool, width int) string` renders one row: a
  plain styled line normally, or a full-width highlighted band when
  `selected` is true.

### Errors

- None. Pure rendering; no invalid input to reject.

## Tests First

- A task row shows its `▸`/`▾` icon matching `m.expanded[taskID]`, and a
  ticket row shows `└`, indented deeper than its task (existing coverage in
  `TestViewExpandShowsTicketsIndentedCollapseHidesThem` - extend or add
  alongside it).
- The cursor row's rendered line, stripped of ANSI, is no wider than the
  sidebar's inner width, at a normal width and at the narrowest width the
  layout permits (e.g. 10-20 columns).
- A ticket or task title longer than the sidebar width renders as one line,
  truncated, never as two lines (`strings.Count(line, "\n") == 0` after
  stripping ANSI, and no wrapped continuation of the title text appears).
- The cursor's highlighted line carries the background style across the
  full row, not just across the label text (check the raw ANSI output
  contains the background sequence up to the row's rendered width).

## Acceptance

- `go build ./...`, `go vet ./...`, and `go test ./...` all pass.
- In the running TUI, the selected row shows a solid full-width background
  band; task rows show an expand/collapse arrow; ticket rows show a branch
  glyph indented under their task.
- Long titles truncate cleanly at both normal and narrow sidebar widths,
  with no wrapped line and no broken border.

## Human Verification Steps

1. Run the TUI (`go run ./cmd/agentism-tui` or the project's normal launch
   command). Expect: the sidebar tree renders with a task row showing `▸`
   or `▾`, and its tickets (once expanded) showing `└`.
2. Press the down/up arrow keys to move the cursor across rows. Expect: the
   selected row shows a solid background band spanning the full sidebar
   width, not just behind the text.
3. Shrink the terminal window until the sidebar is narrow, or pick a task
   or ticket with a long title. Expect: the long title's row is cut off
   cleanly at the sidebar's edge, on one line, with the box border still
   straight and unbroken.
