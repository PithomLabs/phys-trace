# PhysTrace Protocol — Human-Run Manual Validation (M1)

Stage A manual-validation assets for PhysTrace Phase I,
per `plan1.7.md` (spec v1.2.3, in `TOE-v4/plan1.7.md`).

## Version status (v1.2.3)

- **Canonical for M1 (idea-first, v2):** `prompts/extract_v2.md`,
  `prompts/review_v2.md`, `templates/representation_record_v2.md`,
  `templates/adjudication_v2.md`, `M1_checklist_v2.md`.
- **Historical (equation-level, v1):** `prompts/extract_v1.md`,
  `prompts/review_v1.md`, `templates/representation_record.md`,
  `templates/adjudication.md`, `M1_checklist.md`, plus
  `prompts/extract_v1.json` and `prompts/adv_v1.json` (rev0 review
  artifacts for arXiv:2101.11923). Preserved unchanged; do not use
  for new M1 work.
- Governing invariant: X accepted core ideas → exactly X derivation
  units. Starting proposal: K1 single Maxwell derivation (transformation
  + mass identification as stages); K2 scalar extension as candidate.
  Final X adjudicated in M1.

## M1 paper

Perez & Ribisi, arXiv:2101.11923 (Am. J. Phys. 2022) — energy-mass equivalence
from Maxwell equations (EM radiation confined in a box). This paper is the M1
**representation benchmark**, not the compiler's golden fixture.

## Workflow (human-operated)

1. **Propose ideas + derivations.** Open Gemini Web, upload the paper PDF,
   run `prompts/extract_v2.md`. Save the exact candidate JSON output.
2. **Review.** Open a **fresh** Gemini Web conversation (no shared context).
   Supply the same paper PDF + the candidate JSON. Run
   `prompts/review_v2.md`. Save the exact findings output. The review
   challenges split/merge/scope as well as fidelity.
3. **Adjudicate.** For each finding, record accept / reject / defer +
   rationale in `templates/adjudication_v2.md`. Accepted corrections
   produce a new revision (new digest, parent link). Freeze X and the
   approved revision.
4. **Evaluate.** Check off `M1_checklist_v2.md`. M1 passes only on the
   strict gate; otherwise record the failed criterion and required
   correction.

## Transferring artifacts between sessions

Gemini Web sessions are not programmatically linked. The human carries:

- Paper PDF (identical bytes both sessions; record its SHA-256 once).
- Candidate JSON (paste or upload into the fresh review session).
- Review findings (bring back for adjudication assistance).

The approved revision + lineage becomes the Stage B compiler input fixture.

## Key rules

- Representation is independent of the session that produced it.
- Findings never auto-modify the representation; only recorded human
  adjudication creates a new revision.
- Unsupported math is preserved and classified, never coerced or dropped.
- `COMPLETE` = accounted-for within declared scope, never correctness/truth.
- No unit exists because an equation was extracted; no accepted idea is
  omitted for exceeding compiler power.
- The GR equivariance fixture (`phys-lib/gr/equivariance/`) is a synthetic
  compiler-plumbing test (APPLY+CANON only, no paper source) — never the
  paper demo.
