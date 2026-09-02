---
id: Task_0004
kind: task
title: Cross-platform release
status: ARCHIVED
priority: 40
tags:
  - release
  - docs
locked: true
depends_on: []
version: v0
updated_at: "2026-09-02T16:15:04Z"
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
  - date: 2026-09-02
    kind: archived
    summary: Shipped in v0
---

# Cross-platform release

## Goal

A version-tag push builds `agentism-tui` for every target platform and
publishes the binaries to a GitHub Release, with no manual build step.

## Scope

- `.goreleaser.yaml`: build for darwin/amd64, darwin/arm64,
  linux/amd64, linux/arm64, with `CGO_ENABLED=0`. Archive each binary
  as a `.tar.gz` and include a checksums file.
- `.github/workflows/release.yml`: on a pushed tag matching `v*`,
  check out the repo, set up Go, and run
  `goreleaser release --clean`.
- `README.md`: what the tool does, the current-directory requirement,
  the keybindings, and the install instructions (download a release
  binary, or `go install`).

## Out Of Scope

- Windows binaries.
- Publishing to a package manager (Homebrew, apt, and so on).
- A release triggered by anything other than a `v*` tag push.

## Contracts

<!-- The source of truth. A ticket must not contradict this section.
     A change here marks finished tickets as STALE. -->

### Inputs

- A git tag matching `v*` (for example `v0.1.0`), pushed to the
  repository's default remote.
- `GITHUB_TOKEN`, provided by GitHub Actions to the workflow.

### Outputs

- A GitHub Release named after the tag, holding four `.tar.gz`
  archives (one per platform) and a checksums file.
- Each archive holds one binary named `agentism-tui` (or
  `agentism-tui.exe`, not applicable here since Windows is out of
  scope).

### Errors

- A `goreleaser` build failure fails the workflow run and publishes no
  release. GoReleaser's own default behaviour: partial output is not
  published.

### Invariants

- `go build ./...` and `goreleaser build --snapshot --clean` both
  succeed locally before the human marks a ticket in this task done,
  so a tag push never tries the build for the first time.

## Acceptance

Pushing a `v*` tag on the repository produces a GitHub Release with
four downloadable archives, one per target platform, each containing a
binary that runs on that platform. `goreleaser build --snapshot
--clean` succeeds locally without a tag or a token.

## Design Notes

The workflow needs no test step of its own: Task_0001-0003's own
tickets already gate `go test`, `go vet`, and `gofmt` before merge.
This task's job is packaging and publishing what already works.

## Resolved Questions

- Q: Which GitHub repository does the release publish to?
  A: Whichever remote the maintainer pushes the tag to. The workflow
  uses the repository GoReleaser runs in; `.goreleaser.yaml` never
  names an owner or a repo.
