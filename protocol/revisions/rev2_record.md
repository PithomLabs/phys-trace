# Representation Record v2 — rev2 (idea-centered, adjudicated, FROZEN)

Frozen M2 stress-test fixture. Do not mutate; errors become rev3 with
parent link. rev1 preserved unchanged.

## Identity

- Paper: Perez & Ribisi, arXiv:2101.11923 (energy-mass.pdf, 671,492 bytes)
- PDF SHA-256: `4cd872322ef0358a50c10a70281e76adcfaf66826b1707f757a8f883a563c598`
  (verified identical in TOE-v4/energy-mass.pdf and protocol/energy-mass.pdf)
- Revision: rev2 / Parent: `3ec9f52a…5a38` (rev1 exact bytes)
- Revision digest (exact bytes): `000f6eadb0b8d18ce05fb43436d82aa5dc1c0b9a2ba74ee7382c407569675280`
- Lineage: `[29095ee6…cfeb (rev0) → 3ec9f52a…5a38 (rev1) → 000f6ead…5280 (rev2)]`
  + review `0d8649c3…6167e` (declared-candidate mismatch disclosed in rev1
  adjudication; carried, not concealed)

## rev1 → rev2 delta (audit corrections)

1. PDF digest inserted (was null).
2. Scope reconciled: included `[[1,1],[4,14]]` + named regions; excluded
   `[[2,3]]` heuristic, `[[15,17]]` App. B/C routes + references. No
   overlaps. Ledger covers all 17 pages (p.17: App. C remainder + refs).
3. Cross-unit refs resolved via explicit shared-item semantics:
   `shared_items = [it_a1, it_e30, it_001]` with policy + per-unit
   `shared_dependencies` on D2. No implicit ownership.
4. Load-bearing links added: it_m7→it_d1 (defines), it_a2→it_e23
   (assumes), it_e23→it_e28 (constrains), it_s39→it_s40 (derives),
   explicit shared-pattern it_s40→it_e30 (non-required).

## Inventory

- X = 2 ideas (K1 Maxwell, K2 scalar), 2 units (D1 6 steps, D2 3 steps),
  17 items, 13 relations. Bijection holds; all refs resolve.
- Reason codes: `tensor_calculus` (new) + `integration`; 13 UNSUPPORTED,
  2 SUPPORTED-shape endpoints, assumptions NOT_APPLICABLE.
- Deferred: App. B/C routes, it_a2 uniformity criterion, it_m4 reading,
  review-session attestation (operator-attested only — M1 gate item
  honestly open; does not block M2 fixture use).

## Verdicts (derived)

- representation: PARTIAL (designed exclusions + open attestation)
- execution: PARTIAL / test: NOT_RUN
