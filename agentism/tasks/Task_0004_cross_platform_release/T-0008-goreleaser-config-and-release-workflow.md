---
id: T-0008
kind: ticket
title: GoReleaser config and release workflow
task: Task_0004
status: TODO
priority: 10
tags: [release, docs]
depends_on: [T-0005, T-0007]
contract_hash: e1394d3a13fae418
commits: []
updated_at: 2026-09-02T10:19:42Z
changelog:
  - date: 2026-09-02
    kind: created
    summary: Created with task Task_0004
---

# GoReleaser config and release workflow

## Work

`.goreleaser.yaml` plus `.github/workflows/release.yml`, so a pushed
`v*` tag builds and publishes the four target binaries with no manual
step.

## Contract

### Inputs

- `.goreleaser.yaml`: builds for darwin/amd64, darwin/arm64,
  linux/amd64, linux/arm64, `CGO_ENABLED=0`, `.tar.gz` archives, a
  checksums file, a GitHub Release.
- `.github/workflows/release.yml`: triggers on push of a tag matching
  `v*`; steps: checkout, set up Go, run
  `goreleaser release --clean` with `GITHUB_TOKEN` from
  `secrets.GITHUB_TOKEN`.

### Outputs

- `goreleaser build --snapshot --clean` run locally produces four
  binaries under `dist/`, one per target platform.
- A pushed `v*` tag, once run through the workflow, produces a GitHub
  Release named after the tag with four `.tar.gz` archives and one
  checksums file attached.

### Errors

- A compile failure for any one target fails the whole `goreleaser`
  run; no partial release publishes.

## Tests First

No unit test: this ticket is configuration, not Go code. The check is
`goreleaser build --snapshot --clean` (or `goreleaser check` for the
config's own syntax) run locally and inspected by hand.

## Acceptance

`goreleaser check` reports the config valid. `goreleaser build
--snapshot --clean` produces four binaries under `dist/`, one for each
of darwin/amd64, darwin/arm64, linux/amd64, linux/arm64.

## Human Verification Steps

1. Install `goreleaser` locally (or use its Docker image). Run
   `goreleaser check` in the repository root. Expect no error.
2. Run `goreleaser build --snapshot --clean`. Expect `dist/` to
   contain four binaries, one per target platform, and expect each one
   to run and print something sane when copied to a matching machine
   (or inspected with `file` for its target architecture).
3. Read `.github/workflows/release.yml` and confirm it triggers only
   on a `v*` tag push.
