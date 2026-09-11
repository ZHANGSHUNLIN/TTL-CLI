---
name: personal-decision-records
description: Use when a ttl-cli change affects architecture, storage, configuration, file formats, synchronization, CLI compatibility, security, migration, privacy, or the local development workflow.
---

# Personal Decision Records

Use `docs/decisions/` to record why a meaningful engineering choice was made. A decision record is not a task status report, user manual, test log, or review transcript.

## Decide whether a record is needed

Create or update a record when the change:

- chooses or replaces a core implementation;
- changes a database, config, file format, sync protocol, or migration;
- changes user-visible defaults, commands, output, or compatibility;
- introduces security, privacy, migration, or data-loss risk;
- rejects a plausible alternative that may be proposed again;
- changes the local development or manual-review workflow.

Usually do not create one for a local bug fix, formatting-only change, comment cleanup, or behavior-preserving refactor whose rationale is obvious from the code and tests.

## Write the record

Copy `docs/decisions/TEMPLATE.md` and use a date-prefixed filename. Include:

- the problem and context;
- the adopted, proposed, or rejected decision;
- at least one real alternative;
- benefits, costs, risks, and maintenance impact;
- how the decision is verified.

Use present tense for an adopted decision. Name the affected command, package, file format, or workflow explicitly.

## Avoid duplication

Before adding a record, search `docs/decisions/` for an existing record about the same choice. Update or supersede the existing record when it already owns the fact. Do not create multiple records that explain the same decision.

During review, verify that the status matches the repository:

- `adopted`: the code or workflow already follows it;
- `proposed`: it is not yet authoritative;
- `rejected`: the alternative is explicitly not used and the reason remains useful.
