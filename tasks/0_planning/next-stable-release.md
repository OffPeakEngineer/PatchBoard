---
id: task-20260729-next-stable-release
owner: andy
tags:
  - patchboard
  - release
  - distribution
created: 2026-07-29
---

# Prepare the next stable PatchBoard release

## Outcome

Produce a release candidate whose source, module, CI, artifacts, installation
instructions, and published release location all describe the same product.
Creating or publishing the tag is a separate, explicit final action.

## Dependency sequence

1. [x] Resolve the `main` integration and preserve the formalized CLI behind a
   shared root and `cmd/patchboard` entrypoint.
2. [x] Finish the ready browser-board usability work and return the ready lane
   to zero.
3. [ ] Decide the canonical public source, Go module path, and release host in
   `task-20260715-public-module-and-release-coordinates`.
4. [ ] Update imports, installation docs, CI, and GoReleaser to those coordinates
   and remove publishers that are not intentionally supported.
5. [ ] Run the release matrix with the `go.mod` Go version, race tests, vet,
   cross-builds, template/lint checks, and a GoReleaser snapshot.
6. [ ] Choose the stable version from the last published tag and prepare concise
   release notes from completed work.
7. [ ] Review the candidate, then explicitly tag and publish it.

## Scope decisions

- Done-task cleanup is not a blocker for this stable release. Its destructive
  behavior remains a separately designed follow-up.
- Browser support beyond the File System Access path remains progressive. The
  stable release must state tested capabilities and limitations without
  promising unverified write support.
- GitLab is the current repository remote and the leading canonical-host
  candidate, but the public module/release task remains the durable place to
  confirm that decision before changing public coordinates.

## Done when

- The canonical module and release coordinates are settled
- CI and release configuration contain no stale hosts or placeholder publishers
- Installation from the documented release surface is verified
- The full validation matrix passes from a clean candidate tree
- Release notes name browser capability limitations and migration details
- The tag/publish action happens only after explicit approval
