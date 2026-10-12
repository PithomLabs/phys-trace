# Representation Record v2 — rev1 (idea-centered, adjudicated)

## Identity

- Paper: Perez & Ribisi, arXiv:2101.11923 (energy-mass.pdf verified)
- PDF SHA-256: `<to fill from exact uploaded bytes>`
- Revision: rev1 / Parent: `29095ee6…cfeb` (rev0 exact bytes)
- Revision digest (exact bytes): `3ec9f52ac7d279f07a9a06ba5b9284761f3f116150247f4603301ad52a4a5a38`
- Provenance: extract_v2 → review_v2 → human adjudication (this record)

## Scope and ledger

- Included pp. 1-16; excluded: p.17 refs; pp.2-3 heuristic (motivation);
  App. B/C routes; App. A remainder. Full per-page ledger in rev1 JSON.
- X = 2 ideas, 2 units (bijection holds).

## Inventory summary

- K1/D1: 6 ordered steps (s1-s4 UNSUPPORTED, s5-s6 SUPPORTED-shape);
  13 items incl. App. A lemma it_a28 inside s4.
- K2/D2: 3 ordered steps (t1-t2 UNSUPPORTED, t3 SUPPORTED-shape);
  4 new items (scalar eqs 33,34,39,40).
- Reason codes: `tensor_calculus` (new, adjudicated) + `integration`.
- Ambiguities: it_a2 uniformity criterion; it_m4 symbolic reading.

## Verdicts (derived)

- representation: PARTIAL (appendices B/C + heuristic out of scope by
  design; rev1 accounts for declared scope)
- execution: PARTIAL (tensor/integral steps unsupported)
- test: NOT_RUN (no backend execution per R10; external replay only)

## Attachments / pointers

- `revisions/rev1_candidate.json` (digest above)
- `revisions/rev1_adjudication.md` (this file's companion)
- rev0: `prompts/extract_v2.json` (29095ee6…); review: `prompts/review_v2.json` (0d8649c3…)
- Declared-but-unmatched review digest 7f8b2c1a… recorded as mismatch.
