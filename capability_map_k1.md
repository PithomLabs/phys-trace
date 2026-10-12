# K1 Capability Map — Maxwell Mass-Energy Derivation (M2 Gate Deliverable)

Fixture: `protocol/revisions/rev2_candidate.json` (digest
`000f6ead…5280`), paper arXiv:2101.11923. Branches:
`phys-math/phystrace/k1-tensor-foundations`,
`phys-lib/phystrace/k1-maxwell-semantics` (no commits yet; baselines with
untracked workspace `/home/chaschel/Documents/go/go.work`: phys-math 43
pass, phys-artifact 9 pass, phys-lib 30 pass in 6 packages).
Without workspace, all three repos fail to build (missing requires /
`v0.0.0` placeholder) — R1/R9 unchanged.

Rule used throughout: representable ≠ executable ≠ verified.

## Per-step map

| Step | Needs | Existing support | Missing representation | Missing operation | Placement |
|---|---|---|---|---|---|
| D1.s1 Rindler setup (4,6,7) | metric tensors, coords, cosh/sinh transforms | `Tensor` (dense, `families.go:390`), `Index` + COVARIANT/CONTRAVARIANT variance (296-356), scalar algebra | `Metric` family (symmetric rank-2 + inverse relation); coordinate chart; hyperbolic functions (deliberately: GUnknown-style local only, never projected) | none proposed (represent-only) | PhysMath: Metric + chart families |
| D1.s2 current+integral (18-21) | index contraction T_ab ξ^b; hypersurface integral | GRExpr→Expr projection for bare algebra only (`gr/algebra.go:120-175`); `Relation`/`Equation` wrappers | `Contraction` certificate (validated pairing, not evaluation); `Integral` family (domain + integrand refs) | contraction evaluation; integration execution — NOT proposed | PhysMath: representation families + structural validation |
| D1.s3 stationarity (23) | condition over tensor integral | `MathCondition` wraps relation (`families.go:680`); `bridge.CaptureCondition` → assumption | integral-valued predicate support (relation over Integral family) | none (condition threading already enforced by replay) | PhysMath: allow condition predicates over new families; gr: assumption |
| D1.s4 conserved integral (24-29) + App.A flux | divergence theorem, Stokes/Gauss, wall conditions | `Derivation` ordered-trace machinery; `CheckTheoremConclusion`; gr framework/assumption/interpretation types | `CovariantDerivative` as reserved non-executable operator ID; divergence-theorem as named framework-local lemma (NOT a PhysMath rule) | full chain execution — NOT proposed | PhysMath: operator ID reservation; gr: local derivation with explicit projection boundary (§9 pattern) |
| D1.s5/s6 endpoints (30,32) | gamma algebra, limit correspondence | FULLY SUPPORTED today: Mul/Pow/Rational + SIMPLIFY/CANON + Replay | none | none | — (genuine executable obligations) |
| D2.t1-t3 scalar | wave operator, scalar T_ab, same integral pattern | same as above | same families (proves generality if reused unchanged) | same exclusions | same |

## Proposed bounded slice (approval requested — DO NOT implement yet)

**A. PhysMath, representation-only + structural validation (no new execution):**
A1. `Metric` family: symmetric rank-2 tensor + inverse-pair validation + canonical encoding (pure bilinear form; NO signature/convention — those stay in gr).
A2. `Contraction` certificate: validated index-pairing artifact (variance opposition, label match, single-use); validation only, no evaluation rule.
A3. `Integral` family: domain ref + integrand refs + measure note; structural validation + canonical bytes; explicitly non-executable.
A4. Reserve `covariant_derivative` operator ID as representable-but-non-executable (`Apply` rejects with explicit sentinel; `Validate` accepts shape).
A5. Thread new families through `Derivation` premise/result sets, `TargetChecks`, replay family contracts, digest/encode rules. Regression-guard all 7 existing rules.

**B. phys-lib/gr, physical semantics only:**
B1. Rindler/Minkowski frame artifacts + coordinate relations (no math ops duplicated).
B2. Maxwell F_ab / T_ab structures as gr-local representations with one-way fallible projection (extends `algebra.go` pattern; trig/hyperbolic/unknown stay local).
B3. Stationarity assumption (regularity category) via `CaptureCondition`; conducting-wall assumptions + wall-flux conditions; energy-current conservation as framework-local derivation referencing PhysMath fragments.
B4. Endpoint equations (30)/(32) as `PhysicalEquation`s; scalar-field mirrors for K2 later.

**C. Validation (with slice):** unit + negative (bad contraction, variance mismatch, rank/dim violation, malformed integral, missing assumption) + determinism/canonicalization + full 7-rule regression + cross-layer projection tests + required-dependency mutation (break s5/s6 link → fail) + K1 chain-status test asserting s1–s4 UNSUPPORTED / s5–s6 executable (anti-vacuity: test fails if chain ever claimed fully verified).

## Explicit exclusions (not in slice)

General CAS; symbolic integration/divergence execution; Stokes/Gauss as PhysMath rules; trig/hyperbolic projection; new executable rewrite rules; sparse tensors; `core/` promotion; endpoint-only pass claims; tags/releases/merges.

## Remaining gaps after slice (honest)

Full K1 replay still impossible (s1–s4 non-executable by design); K1 verdict stays execution-PARTIAL with executable s5/s6 fragments. K2 reuses the same families (generality test, not new math). Portable replay still R1/R9-blocked. First-blocker if slice approved and built: s4 divergence-theorem execution (requires specification decision, not just code).

## As-built record (slice implemented, uncommitted on feature branches)

- A1 landed as `FamilySymmetricBilinear` (`symmetric_bilinear/1`), not
  `metric/1`: the purity test bans physical vocabulary in exported
  identifiers, and the rename keeps the type honestly mathematical
  (pure bilinear form + optional self-inverse assertion; symmetry via
  `CheckSymmetricBilinear`, signature/conventions stay framework-owned).
  Inverse is `*ArtifactRef`, nil = explicit self-inverse claim.
- A4 sentinel: `ValidOperatorID` admits {differentiate,
  covariant_derivative}; `Executable()` + `ExecuteOperator` fail closed
  with `ErrOperatorNotExecutable` even on perfect params; existing
  `families_test.go` expectation updated with justification comment.
- Registry 16 -> 20; premise/result sets extended with 6 context
  families (carried-only; step outputs remain expression/1).
- Tests: phys-math 54 pass (43 baseline incl. updated operator/registry
  expectations + 11 new); phys-lib 39 pass in 8 packages (30 baseline +
  9 new: projection report, frames, assumptions, sketch validation,
  5 maxwellbox protocol tests).
- K1 attempt (`gr/maxwellbox`, rev2 fixture): s1 metric symmetry
  checked + integral/contraction represented; s5/s6 endpoint fragments
  replayed with theorem-conclusion gates; s1-s4 recorded UNSUPPORTED in
  the sketch. FIRST BLOCKER (as pre-declared): s4 divergence/flux
  execution — needs a specification decision, not just code. K1 is
  faithfully represented with executable endpoint fragments; it is NOT
  verified end to end. K2 deferred.

## Stop

No new capability family to be implemented until slice A–C is approved. Routine detail within an approved slice needs no further gate; any material family beyond A1–A4/B1–B4 returns for approval.
