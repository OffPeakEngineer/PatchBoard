# Contributing to Patchboard

Thank you for your interest in contributing! Please read this guide before submitting a pull request.

## Licensing & Authorship

By contributing to Patchboard, you agree to:

1. **License your changes under MIT OR Zlib** — the same dual-license terms as the project.
2. **Confirm you are the original author** — in the PR description, explicitly state that you authored the code and are contributing it under the MIT OR Zlib license.

## Pull Request Checklist

- [ ] I have authored this code myself
- [ ] I am contributing this work under the MIT OR Zlib license
- [ ] I have included the copyright header in new files
- [ ] I have tested my changes (`go test ./...`)

## Questions?

If you have questions, feel free to open an issue to discuss your contribution first.

## Commit messages

Use Conventional Commits: `fix: ...`, `feat: ...`, or `docs: ...`, with an
optional scope. Mark breaking changes with `!` or a `BREAKING CHANGE:` footer.
CI checks PR/MR commits; keep squash-merge titles conventional too. GitHub uses
these messages to choose versions and publish releases automatically after
default-branch checks pass. See [Releases](docs/releases.md).
