# Representation Record — revN

Fill one copy per revision. Revision 0 = extraction candidate; each
adjudicated correction creates revN+1 with a parent link.

## Identity

- Paper: Perez & Ribisi, arXiv:2101.11923
- PDF SHA-256: `<hex of exact bytes processed>`
- Revision: `rev0`
- Parent revision digest: `null` (rev0) / `<hex>` (later revisions)
- Canonical digest rule: SHA-256 over the canonical JSON form of the
  revision payload — keys sorted recursively, UTF-8, no insignificant
  whitespace, numbers in shortest round-trip form. Digest recorded below
  after canonicalization.
- Revision digest: `<hex>`
- Provenance: extraction prompt `extract_v1` / model `<name/version>` /
  date `<YYYY-MM-DD>` / operator `<name>`

## Declared scope

- Included pages: `<e.g. 1–12>`
- Excluded: `<e.g. p.13 references — reason>`

## Coverage ledger

| Page | Status | Note |
|------|--------|------|
| 1 | INSPECTED_WITH_CONTENT | |
| … | … | |

Statuses: `INSPECTED_WITH_CONTENT | INSPECTED_NO_FORMAL_CONTENT |
UNREADABLE | UNASSESSED`.

## Item inventory

| item_id | kind | label | pages | capture | compiler (reason) | test | required-by |
|---------|------|-------|-------|---------|-------------------|------|-------------|
| it_001 | equation | (2.1) | 3 | CAPTURED | UNSUPPORTED (integration) | NOT_RUN | it_004 |
| … | … | … | … | … | … | … | … |

## Assumptions register

| item_id | text/passage | author-stated? | disposition |
|---------|--------------|----------------|-------------|
| … | … | yes/no | accepted/… |

## Relations (required links)

- `it_001 --derives--> it_002 (required)`

## Ambiguities / deferred questions

- `<item_id>: <note>`

## Verdicts (derived, recomputed at compile — never hand-set)

- `representation_verdict`: PARTIAL (until gate in M1_checklist passes)
- `execution_coverage`: PARTIAL
- `test_execution`: NOT_RUN

## Attachments

- `candidate_rev0.json` (exact bytes reviewed)
- `review_findings.json` (exact bytes from review session)
