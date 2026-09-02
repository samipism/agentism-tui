---
id: T-0009
kind: ticket
title: README
task: Task_0004
status: DONE
priority: 20
tags:
  - release
  - docs
depends_on:
  - T-0007
  - T-0008
contract_hash: e1394d3a13fae418
commits:
  - dcd7ce3ad6d8b792c99abb8cd626c956fc493cec
updated_at: "2026-09-02T16:13:54Z"
changelog:
  - date: 2026-09-02
    kind: created
    summary: Created with task Task_0004
  - date: 2026-09-02
    kind: "status:IN_PROGRESS"
    summary: Selected by agentism next
  - date: 2026-09-02
    kind: "status:IN_REVIEW"
    summary: Moved from IN_PROGRESS
  - date: 2026-09-02
    kind: "status:DONE"
    summary: Human accepted the result
---

# README

## Work

Write `README.md`: what agentism-tui does, how to install it, and how
to use it.

## Contract

### Inputs

- The finished behaviour from T-0005 (entry point), T-0006/T-0007
  (dashboard and keys), and T-0008 (release artifacts).

### Outputs

`README.md` with these sections: what the tool is (a read-only viewer
for an agentism project's plan), install (download a release binary
for your platform, or `go install
github.com/<owner>/agentism-tui/cmd/agentism-tui@latest`), usage (run
`agentism-tui` inside a project folder that has `.agentism/`), and the
full keybinding list (arrows/j-k, enter/space, 'r', 'l', 'q').

### Errors

None. This is documentation.

## Tests First

No automated test. The check is a human read-through against the
actual keybindings and install path once T-0005 through T-0008 exist.

## Acceptance

A person who has never seen the tool can install it and start it from
the README alone, without asking a question.

## Human Verification Steps

1. Read `README.md` top to bottom.
2. Follow its install step on a machine that has never had the binary.
3. Follow its usage step in a folder with an agentism project. Confirm
   every keybinding it lists matches what the running program accepts.
