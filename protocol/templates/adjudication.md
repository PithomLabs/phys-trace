# Adjudication Record — candidate rev0 → approved rev1

One row per review finding. The human decides; ChatGPT may recommend.
Findings never auto-modify the representation — only rows marked
`accept` with a stated edit create the next revision.

- Candidate digest: `<hex of exact candidate bytes reviewed>`
- Review output digest: `<hex of exact findings bytes>`
- Adjudicator: `<name>` / Date: `<YYYY-MM-DD>`

## Findings

| finding_id | severity | location | decision | rationale | revision edit |
|------------|----------|----------|----------|-----------|---------------|
| f_001 | omission | p.5 | accept | genuine gap; eq. (4.2) missing | add it_014, ledger p.5 stays WITH_CONTENT |
| f_002 | support | it_003 | reject | surface integral is explicit in text; stays UNSUPPORTED | none |
| f_003 | assumption | it_009 | defer | boundary regularity unclear; needs source recheck | mark AMBIGUOUS, keep in rev1 |
| … | … | … | … | … | … |

Decisions: `accept | reject | defer` + rationale (required).
Deferred/unresolved questions stay explicit — never recorded as results.

## Approved revision

- Approved digest: `<hex over canonical rev1 bytes>`
- Parent digest: `<candidate hex>`
- Change summary: `<added/removed/corrected items; ledger updates>`
- Unsupported items preserved: `<list — must be ≥1 genuine item>`
- Adjudication is declared provenance, not cryptographic proof of review.

## Lineage appended to package metadata

`[candidate_digest → … → approved_digest]` + this record (or its digest
with external storage pointer). Compile MUST recompute digests and
confirm the approved revision is the chain tip.
