# Adjudication Record v2 — rev0 → rev1 (idea-centered)

- Candidate file: `protocol/prompts/extract_v2.json`
  exact-byte SHA-256: `29095ee6d3d943bfc4b8642b47541953e0a0a4dde3736e9d5b54b8378838cfeb`
- Review file: `protocol/prompts/review_v2.json`
  exact-byte SHA-256: `0d8649c3732dd8ab7f82e7254cb4e46b82a89f4e5ab2c6d109b3148f3676167e`
- Review DECLARED candidate digest: `7f8b2c1a0e9d8c7b6a5f4e3d2c1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b`
  → MISMATCH with actual candidate bytes. Findings usable, not
  cryptographically bound. Recorded per R11; do not copy declared digest.
- Adjudicator: human researcher. Assistance: ChatGPT analysis; human decides.

## X decision

X=2. K1 single Maxwell derivation (transformation + mass ID as stages);
K2 scalar extension included as independently meaningful universality
check. D1/D2 split retired. Prompts must not hard-code X.

## Idea findings

| finding_id | kind | decision | rationale |
|------------|------|----------|-----------|
| f_001 missing ideas | partial accept | Accept K2 scalar extension; reject heuristic (Sec. II, motivation) and Rindler setup (Sec. III, internal to K1) as separate ideas |
| f_002 bad merge | complaint accepted, remedy rejected | D1 was skeletal (one step); fixed by filling steps (f_010), not by splitting into topic units — splitting is the forbidden artificial split |
| f_003 scope | partial accept | Needed App. A boundary/flux material inside D1 (it_a28); App. B/C alternative routes explicitly excluded from rev1 |

## Derivation findings

| finding_id | severity | decision | rationale |
|------------|----------|----------|-----------|
| f_010 omission | accept | D1 filled: s1 Rindler setup, s2 current+integral, s3 stationarity, s4 conserved integral+flux, s5 transformation, s6 limit+ID; D2 filled: t1-t3 |
| f_011 support | accept with correction | Shape-vs-chain rule: it_001/it_e30 SUPPORTED-as-shape; it_d1 + all tensor/integral items UNSUPPORTED (`tensor_calculus` new code alongside `integration`); obligations only s5/s6/t3; execution PARTIAL |

## Integrity notes (also recorded)

- rev0 candidate file is NOT strict-JSON-parseable (unescaped LaTeX
  backslashes, e.g. `\e`); lineage therefore uses exact-byte digests.
  rev1 is strict-valid JSON; future digests use canonical form.
- K1 source_regions corrected (Sec. V claim removed); K2 owns Sec. V.
- Review's own omission list partial (missed main-text eqs 8-14, 16-23,
  25-29, 31, 33-40 beyond flagged items); rev1 covers the load-bearing
  chain; residual transcription risk stays explicit, not hidden.
- Adjudication = declared provenance, not cryptographic proof of review.

## Frozen rev1

- Approved digest (exact bytes): `3ec9f52ac7d279f07a9a06ba5b9284761f3f116150247f4603301ad52a4a5a38`
- Parent digest: `29095ee6…cfeb` (rev0 exact bytes)
- Lineage: `[29095ee6…cfeb → 3ec9f52a…5a38]` + review `0d8649c3…6167e`
- Unsupported preserved (genuine): tensor_calculus + integration across
  s1-s4/t1-t2. Deferred: App. B/C routes, it_a2 uniformity criterion,
  it_m4 symbolic-reading question.
