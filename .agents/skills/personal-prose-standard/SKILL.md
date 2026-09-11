---
name: personal-prose-standard
description: Use when writing or reviewing ttl-cli README text, command help, comments, Go documentation, error messages, or localized user-facing strings.
---

# Personal Prose Standard

Write the smallest complete explanation of the current behavior. Prose must help a user or maintainer act; it must not preserve an authoring-session transcript.

## Keep

- Preconditions, postconditions, failure behavior, compatibility promises, data ownership, and security or privacy constraints.
- Concrete examples that match the actual CLI and configuration.
- Comments for non-obvious invariants or decisions that code cannot express.

## Remove or rewrite

- Control-flow narration that only repeats the code.
- Review choreography such as “the reviewer requested”.
- Temporary plans, vague hedges, and references to uncommitted drafts.
- Historical wording such as “used to” when present behavior is enough.
- Comments that justify obvious syntax instead of naming a real invariant.

## Project-specific checks

- User-visible CLI strings belong in `i18n/locales/` when the surrounding command supports localization.
- README and command help must describe the actual flags, defaults, errors, and examples.
- Go exported identifiers need concise documentation when they are part of a package-facing API.
- Keep Chinese and other locale changes aligned with the behavior change; do not invent a translation for an unsupported feature.
- Link to `docs/decisions/` for rationale instead of duplicating the entire decision in several files.

Review prose against the code and tests. A successful documentation check does not prove that the text is accurate.
