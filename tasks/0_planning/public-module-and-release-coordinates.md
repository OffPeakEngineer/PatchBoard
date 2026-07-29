---
id: task-20260715-public-module-and-release-coordinates
owner: unassigned
tags:
  - patchboard
  - distribution
  - release
  - go
created: 2026-07-15
---

# Reconcile public module and release coordinates

## Problem

The Go module is currently named `ledoerr/patchboard`, which is not a remotely
resolvable module path. Users can run `go install ./cmd/patchboard` from a local
checkout, but cannot use the conventional remote form:

```text
go install <module>/cmd/patchboard@latest
```

The repository remote points to GitLab while parts of the GoReleaser and GitHub
Actions configuration assume GitHub-owned release destinations. Installation
documentation should not promise a canonical release URL until these
coordinates agree.

## Questions to work through

- Is GitLab the canonical public source and release host?
- Should the module path become
  `gitlab.com/off-peak.engineer/patchboard` or another stable public path?
- Should GitHub remain a release mirror, and if so, what is its canonical owner
  and repository name?
- Which GoReleaser publishers are genuinely supported and credentialed?
- Should remote `go install ...@latest` be the primary installation command or
  a secondary developer-oriented option?

## Done when

- One canonical public source URL and module path are chosen
- Go package imports and module metadata use the canonical path
- `go install <module>/cmd/patchboard@latest` works from outside a checkout
- GoReleaser publishes to documented release locations
- CI Go versions agree with `go.mod`
- Installation documentation links only to verified release surfaces
