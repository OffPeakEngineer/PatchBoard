# Contributing to Patchboard

Thank you for your interest in contributing! Please read this guide before submitting a pull request.

## Licensing & Authorship

By contributing to Patchboard, you agree to:

1. **License your changes under MIT OR Zlib** — the same dual-license terms as the project.
2. **Confirm you are the original author** — in the PR description, explicitly state that you authored the code and are contributing it under the MIT OR Zlib license.

See [AGENTS.md](AGENTS.md) for the project's full philosophy on AI use and code ownership.

## Source File Headers

Please include a header comment in new Go files:

```go
// Copyright (c) 2026
// Distributed under the MIT OR Zlib license.
// See LICENSE, LICENSE.MIT, or LICENSE.zlib for full terms.
```

For other file types (YAML, shell scripts, etc.), use the appropriate comment syntax:

```yaml
# Copyright (c) 2026
# Distributed under the MIT OR Zlib license.
# See LICENSE, LICENSE.MIT, or LICENSE.zlib for full terms.
```

## Code Generation & AI

**AI-generated code is not permitted in contributions.** This includes:

- Code output from LLMs (ChatGPT, Claude, Copilot, etc.)
- Auto-generated code from AI coding assistants
- Any code whose origin or authorship is unclear

**Why?** AI-generated code may embed sources of unknown licensing origins, creating incompatibility with the MIT OR Zlib terms around misrepresentation and attribution.

**AI is welcome for:**

- Identifying potential issues or bugs
- Suggesting improvements to approach or design
- Code review feedback

However, **the actual code solution must be authored by you**. If you use AI feedback to identify an issue, please independently verify the problem exists and author the fix yourself.

## Pull Request Checklist

- [ ] I have authored this code myself
- [ ] I am contributing this work under the MIT OR Zlib license
- [ ] I have included the copyright header in new files
- [ ] I have tested my changes (`go test ./...`)
- [ ] I understand that AI-generated code cannot be included

## Questions?

See [AGENTS.md](AGENTS.md) for more on the project's philosophy, or open an issue to discuss your contribution first.
