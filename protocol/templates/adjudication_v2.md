# Adjudication Record v2 — rev0 → rev1 (idea-centered)

Human decides; assistance may recommend. Findings never auto-modify;
`accept` rows with stated edits create the next revision.

- Candidate digest / findings digest / adjudicator / date

## Idea findings (split/merge/scope first — X is decided here)

| finding_id | kind | decision | rationale | edit (merge/split/scope) |
|------------|------|----------|-----------|--------------------------|
| … | bad_split | accept | transformation+ID one argument | merge D1/D2→K1/D1 |

## Derivation findings

| finding_id | severity | location | decision | rationale | edit |
|------------|----------|----------|----------|-----------|------|

## Frozen inventory

- Accepted X = ___ / units = ___ (bijection checked)
- Approved digest / parent digest / change summary
- Unsupported items preserved (≥1 genuine): `<list>`
- Deferred questions stay explicit, never recorded as results.
- Adjudication = declared provenance, not cryptographic proof.

## Lineage for package metadata

`[candidate → … → approved]` digests + this record. Compile recomputes
digests and confirms the approved tip; corrections → new revision →
revalidate + recompile.
