# s4 Divergence/Flux Capability Proposal (proposal only — no implementation)

Fixture: `protocol/revisions/rev2_candidate.json` (`000f6ead…5280`), D1 step
s4 (paper eqs. 24–29 + Appendix A boundary/flux). Branches:
`phys-math/phystrace/k1-tensor-foundations` (`cea7a7f`),
`phys-lib/phystrace/k1-maxwell-semantics` (`a0a9ef0`).
Status: K1 execution blocked at s4; s5/s6 replay as endpoint fragments;
representation/execution `PARTIAL`, test `NOT_RUN` (R10).

This document proposes the smallest executable capability that could
genuinely verify s4. It authorizes nothing: any new executable rule needs
separate explicit approval. No code, tests, commits, or releases here.

## 1. Exact scientific obligation (K1.s4)

Paper eqs. (24)–(29), pp. 7–8, with Appendix A wall results (A28/A29):

- (24): after stationarity (23), `E_box(τ) = γ(τ) ∫_Σ c²x̄⁻²T_ττ` — the
  energy integral whose τ-independence must be established.
- (25): `∇_a j^a_box = x̄⁻²T_τx̄` — divergence of the box current,
  using the Killing equation (14), field equations (A8), and wall
  orthogonality `σ·E_box = 0`.
- (26): `∫_R ∇_a j^a_box dV = ∫dτ(∫_Σ x̄⁻¹T_τx̄) = 0` — spacetime
  integral vanishes by stationarity.
- (27): `j^a_box N_a|_walls = 0` — no flux through walls (proved via
  A28/A29 in Appendix A.2).
- (28): Gauss theorem → flux only through Σ₁/Σ₂ spacelike caps, hence
  the conserved quantity is τ-independent.
- Required inputs: divergence relation (25), stationarity condition
  (23), wall-flux lemma (27/A29), domain R with boundary decomposition
  (Σ₁, Σ₂, walls), Killing/field-equation premises (14, A8).
- Conclusion: the s4 energy integral is τ-independent (feeds s5).

## 2. Current gap (three distinct absences)

- Representation: PRESENT. `Integral` (energy/spacetime integrals),
  `Domain` (Σ_τ, R — now with mandatory defining refs), `Contraction`
  (current pairing), wall/flux lemmas as gr assumptions + sketch refs.
- Structural validation: PRESENT. Shape, variance, grammar, SET
  uniqueness, symmetry checks all pass.
- Executable semantics: ABSENT. No Phase-1 rule transforms a divergence
  relation plus boundary data into a conservation conclusion.
  `Apply` (`phys-math/rules.go:42-58`) has exactly 7 cases; replay
  (`replay.go:18-69`) can only re-execute those. This — and only
  this — is what blocks s4.

## 3. Minimum capability: narrow flux-conservation vs general divergence

- **Option N (recommended): FLUX-CONSERVE-001**, a single closed rule:
  inputs = divergence relation artifact + boundary-flux lemma refs +
  stationarity condition ref + domain-with-decomposition refs;
  output = conservation equation artifact (integral over R equals cap
  flux difference, equal to zero). Preconditions verified structurally
  from refs (exact families, matching domain digests, cited assumptions
  present); the rule checks the implication shape, it does not discover
  physics. Reusable anywhere a box-region divergence + wall data +
  stationarity pattern occurs — not paper-specific, but deliberately
  narrower than Stokes in general.
- **Option G (rejected): general divergence-theorem execution** for
  arbitrary domains/forms. Unbounded scope (needs orientation theory,
  general Stokes, measure theory); invites CAS creep; cannot be
  validated within a bounded slice. Reject.
- Recommendation: Option N, with the rule's soundness resting on the
  explicitly cited boundary + stationarity premises (garbage in →
  fail-closed, never silent success).

## 4. Library placement

- PhysMath: `FLUX-CONSERVE-001` rule (pure implication over
  Integral/Domain/Relation artifacts) + parameter grammar + replay
  family contract + regression guard of the existing 7 rules. No
  physical vocabulary enters (purity test enforced).
- phys-lib/gr: wall-flux lemma artifacts (A28/A29 content as
  definitions), stationarity assumption threading (existing
  `stationarity_zero_momentum`), domain-decomposition definitions
  (R, Σ₁, Σ₂, walls as defined domains), and the K1 sketch status
  flip s4 → SUPPORTED only via cited replay (anti-vacuity preserved).
- No `core/` promotion (single-framework use); no frozen contact.

## 5. Preconditions and failure behavior

Execution succeeds ONLY when ALL hold (else fail closed, no output):
P1 divergence-relation input resolves to `relation/1` over the SAME
domain digest cited by the integral; P2 every boundary piece named by
the domain decomposition has a cited flux lemma (wall pieces cite
vanishing-flux lemmas; caps cite evaluated flux terms); P3 a
stationarity condition (`condition/1`) is threaded and matches the bulk
term; P4 framework assumptions (stationarity, walls, chart) resolve and
are in-closure; P5 output equation bytes equal the recomputed
conservation statement exactly (replay-equality pattern).
Missing/dangling/mismatched refs, wrong families, absent stationarity,
or unsigned boundary pieces → error, no artifact minted.

## 6. Validation battery (to be written only if approved)

Positive: K1 s4 exact-shape success; minimal box happy path; cap-flux
nonzero-but-equal case (flux difference, not just zero).
Malformed: wrong-family inputs; digest-mismatched domain; duplicate
boundary pieces; non-closed assumption set.
Missing-data: no stationarity condition; wall piece without flux lemma;
absent domain decomposition.
Mutation: remove stationarity ref → fail; swap wall lemma → fail;
corrupt output bytes → replay mismatch fail; drop s4 step → downstream
s5/s6 still replay (independence preserved) but chain-status test
reports s4 unverified.
Anti-vacuity: rule invoked on non-conservation inputs → fail; sketch
status flip requires cited replay digest; s5/s6-only suite still
reports chain PARTIAL until s4 replays.

## 7. K1 acceptance criteria

s4 counts as verified IFF: replayed `FLUX-CONSERVE-001` step with all
P1–P5 satisfied; conclusion bytes equal the paper's (28)-content
relation; sketch s4 → SUPPORTED citing the replay; s1–s3 remain
represented-but-unexecuted (execution still PARTIAL overall —
s4 verification does NOT promote them); s5/s6 unchanged. K1 is
end-to-end verified only when s1–s3 also gain verified execution or
are explicitly accepted as representational with documented rationale —
a separate decision, not implied by s4.

## 8. Explicit exclusions

No general integral evaluation; no arbitrary Stokes/divergence support;
no trigonometry; no CAS; no new rules beyond the single proposed one;
no K2 implementation; no `core/` promotion; no tags/releases/merges;
no frozen contact; no committed `go.work`/`replace`.

## 9. Decision and risks

- Choices needing approval: (a) Option N rule semantics (is the
  implication shape sound, or does it beg the question by assuming the
  flux lemmas that do the real work?); (b) exact parameter grammar;
  (c) whether s4-success plus s5/s6 suffices to call K1 "verified
  modulo s1–s3 representation" or K1 stays PARTIAL regardless.
- Risk: the rule could degenerate into a paper-specific axiom
  (mitigation: require the nonzero-flux-difference positive test and a
  second, non-Maxwell use case before acceptance).
- Risk: scope creep into execution of s1–s3 (mitigation: this proposal
  covers s4 only; further execution needs its own proposal).
- Downstream: on approval, implement on the existing feature branches,
  extend the sketch + fixture, rerun full batteries, report.

**Bounded approval request:** approve (a) Option N semantics, (b) the
§6 test battery as the acceptance bar, and (c) implementation confined
to one rule + gr-side lemma/definition artifacts + tests. Everything
else stays as it is. No implementation begins without this approval.
