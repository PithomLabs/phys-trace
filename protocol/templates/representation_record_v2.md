# Representation Record v2 — revN (idea-centered)

One copy per revision. rev0 = candidate; accepted corrections create
revN+1 with parent link. v1 equation-level records preserved as history.

## Identity

- Paper: Perez & Ribisi, arXiv:2101.11923
- PDF SHA-256: `<hex of exact bytes>`
- Revision: `rev0` / Parent digest: `null` or `<hex>`
- Canonical digest rule: SHA-256 over canonical JSON (keys sorted
  recursively, UTF-8, no insignificant whitespace, shortest round-trip
  numbers). Revision digest: `<hex>`
- Provenance: prompt `extract_v2` / model `<name/version>` / `<date>`

## Declared scope and ledger

- Included: `<pages/regions>` / Excluded: `<region — reason>`
- Page ledger table: page | status | region note

## Core-idea inventory (X = ___)

| idea_id | claim | distinctness (one line) | regions | unit |
|---------|-------|-------------------------|---------|------|
| K1 | Maxwell mass-energy | transformation+ID one argument | Sec. III–IV, App. A* | D1 |
| K2? | scalar extension | candidate, pending scope | Sec. V, App. C* | D2? |

## Derivation units (count MUST equal X)

| unit_id | idea | conclusion | steps (required marked *) | unsupported portions | obligations |
|---------|------|------------|---------------------------|----------------------|-------------|
| D1 | K1 | E(0)=mc² | s1 setup … s4* integral … s6* limit | s4–s5 integration/tensors | s6 (algebra) |

## Items (components of units — never units themselves)

| item_id | unit | kind | label | pages | capture | compiler (reason) | test |
|---------|------|------|-------|-------|---------|-------------------|------|

## Assumptions register (author-stated vs inferred flagged)

## Ambiguities / deferred

## Verdicts (derived at compile, never hand-set)

- representation: PARTIAL / execution: PARTIAL / test: NOT_RUN
## Attachments

- candidate bytes, findings bytes, adjudication record ref
