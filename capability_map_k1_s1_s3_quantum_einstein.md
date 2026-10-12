# Capability Map: K1 s1–s3 Remaining Gaps + Quantum Einstein Forward-Compatibility
*Read-only analysis; no code, tests, commits, merges, tags, or releases. Report version for approval of the next bounded slice only.*

## 1. Executive recommendation

Implement one geometry-first slice: **(R1)** Rindler coordinate-map representation with gr-local opaque sinh/cosh nodes (fail-closed projection, no transform execution); **(R2)** diagonal-metric determinant computed with existing algebra (genuinely executable, no registry change); **(R3)** volume-density family with explicit signature/domain assumptions (`μ_g=√|det|`; Rindler `|x̄|→x̄` gated on `x̄>0`). Split contraction evaluation off as a separately approved executable. Use QE minisuperspace (FRW diagonal metric, §6) as the forward-compatibility test; do NOT squeeze QE's general ADM `h_ij` into diagonal form.

**Answer:** the smallest reusable capability is diagonal determinant + assumption-explicit volume density + schematic coordinate-map representation — it advances K1 s1 honestly, reuses in QE minisuperspace, and adds zero execution rules.

## 2. Verified repository and branch state

- `phys-math` on `phystrace/k1-tensor-foundations` (`cea7a7f` + uncommitted `M geometry.go` = DecodeDomain/DecodeIntegral support code from s4 work): **55 pass**, 0 fail.
- `phys-lib` on `phystrace/k1-maxwell-semantics` (`a0a9ef0` + uncommitted `M fixture/manifest/protocol-test`, new `fluxconserve.go` + test): **43 pass in 8 packages**, 0 fail.
- `phys`, `phys-gr`, `phys-artifact`: clean, master, untouched. Zero tags everywhere. `go.work` untracked at workspace root only; no committed `replace` (phys-gr historical exception pre-existing).
- K1 honest status: s1–s3 represented-not-executable; s4 schema-checked (FLUX-CONSERVE-001); s5–s6 replayed endpoint algebra; overall PARTIAL/NOT_RUN. K2 deferred.
- QE source: `TOE-v4/hackathon/dr_durr_2003.pdf` verified byte-identical to arXiv:2003.03839v2 (CQG 37, 135002 (2020), 17pp, read in full for this report).

## 3. K1 s1–s3 dependency and capability map (dependency order)

Paper order governs: eqs. (4)–(17) establish geometry before §IV currents (18+). K1 stays one idea/unit.

### s1 — Rindler geometry (paper (4)–(17))
| Need | Existing | Missing | First verification obligation |
|---|---|---|---|
| Minkowski metric (4) | SymmetricBilinear + symmetry check (diagonal stored) | nothing structural | already storable |
| Rindler transform (6), sinh/cosh | nothing (GRExpr has Sin/Cos local-only precedent, no sinh/cosh) | gr-local opaque sinh/cosh nodes + TransformMap certificate (source/target charts, component refs) | projection-failure test; transform-law execution explicitly unclaimed |
| Rindler metric (7), differentials | SymmetricBilinear storage | metric-transformation validation (NOT execution) | represent both metrics + transform refs; assert transformation unverified |
| Volume form (11) `x̄dxdydzdτ`, dΣ (12) | Domain/Integral families | volume-density semantics (√|det| + x̄>0 assumption) | det + density validated; measure evaluation unclaimed |
| Frame relations (8)–(10),(13)–(17) | GRExpr projection for bare algebra; frame definitions | indexed vector relations | represent, not execute |

Earliest s1 blocker: coordinate-map + hyperbolic representation (precedes all use of Rindler objects downstream).

### s2 — Current + integral (paper (18)–(24))
| Need | Existing | Missing | Obligation |
|---|---|---|---|
| `j_a=-T_ab ξ^b` contraction | Contraction certificate (validates pairing) | contraction EVALUATION (separate approval track) | certificate ≠ evaluation, tested |
| Hypersurface integral, measure | Integral/Domain families | measure evaluation; determinant/volume semantics (this slice) | represent + validate refs |
| Energy definition (19)–(21) | Equation/Relation wrappers | term-denotation linkage | caller-asserted, recorded |

Earliest s2 blocker: contraction evaluation (after s1 geometry; an early primitive must not imply s1 satisfied).

### s3 — Stationarity (paper (23))
Needs s2 machinery + condition-discharge semantics (proving a condition holds vs threading it — new executable semantics, out of this slice). Earliest blocker: s2 evaluation. Assumption recording already exists (regularity category).

## 4. Source-grounded QE derivation inventory (v2, §§1–8)

- **D-QE1 action route (§§3.2–4, eqs. 35–55):** canonical momenta (35), quantum Hamilton-Jacobi + V+Q (37), Q=QG+QM with functional derivatives δ/δh, δ/δφ and amplitude |Ψ| (38–40), action SQ (46–47), TQ via δSQ/δg (49) with components (50–53), scalar eq with δQ/δφ (54–55).
- **D-QE2 direct route (§5, eqs. 56–66):** phase-space split HG/HM/HQ + constraints (56–59), F-tensors (60–63), Einstein components (64), TM/TQ read-off (65–66).
- **D-QE3 minisuperspace (§6, eqs. 67–78):** FRW diagonal metric (67), perfect fluid (68–69), Friedmann (70–71), WDW (74), guidance (75), ρQ/pQ (76–77), TQ fluid (78). **This is the forward-compat test target** (diagonal metric + volume a³ + scalar algebra).
- Supporting: classical ADM setup (11–27), WDW operators with LB ordering (28–31), guidance (32–33), conservation consequence §7 (follows from Bianchi, not an independent proof obligation).

## 5. Shared-versus-specific matrix (status: 1=implemented+tested, 2=representable, 3=unsupported, 4=spec-blocked)

| Capability | K1 s1 | K1 s2 | K1 s3 | QE derivations | PhysMath now | PhysLib now | Missing repr. | Missing exec | Placement | Tests needed | Priority |
|---|---|---|---|---|---|---|---|---|---|---|---|
| Diagonal metric storage | (4),(7) | — | — | D-QE3 (67) | 1 Tensor/SymBilinear | 1 frame defs | — | — | — | existing | done |
| Diagonal determinant | vol (11) | measure | — | D-QE3 a³; LB √g (2) | 2 (Mul exists, no det convention) | — | det certificate | product via existing Mul | PhysMath | correctness + non-diagonal reject | **P0 (this slice)** |
| Volume density √|det| + assumptions | (11),(12) | measure | — | D-QE1/2 √h,√−g; D-QE3 | 2 Domain | VolumeDensity family + x̄>0 assumption | — | PhysMath family; gr assumption | assumption + misuse tests | **P0 (this slice)** |
| Hyperbolic fns | (6) | — | — | — (K1-specific) | 3 | 3 (Sin/Cos precedent only) | opaque sinh/cosh nodes | evaluation (excluded) | gr-local | projection-failure; shape | **P0 (this slice)** |
| Coordinate-map validation | (6)→(7) | frame use | — | ADM foliation (14–15) idea-level | 3 | 3 | TransformMap cert | transform law (excluded) | PhysMath cert; gr chart | grammar/dim-match; execution-unclaimed | **P0 (this slice)** |
| Contraction evaluation | — | Tξ (18–21) | via s2 | everywhere (QG/QM, constraints) | 2 (certificate only) | — | — | CONTRACT-EVAL rule | PhysMath (separate approval) | eval + negative + mutation | **P1 separate track** |
| General determinant | — | — | — | general h_ij (21) | 3 | — | general det | cofactor exec (excluded) | — | — | excluded |
| Covariant derivative exec | — | — | — | LB (2), Di, Kij (18) | 2 (sentinel only) | — | — | excluded | — | sentinel retained | excluded |
| Curvature (Riemann/Ricci/R³, Kij) | — | — | — | (12),(18),(22) | 3 | 3 | curvature families | curvature computation | future proposal | — | deferred |
| Functional derivatives | — | — | — | (39),(40),(49),(54) | 3 | 3 | functional-derivative repr | functional calculus | future proposal | — | deferred |
| Amplitude/phase \|Ψ\|,S; Q terms | — | — | — | (35–40) | 3 | 3 | wavefunctional repr | Q evaluation | future proposal | — | deferred |
| Constraints/Hamiltonian struct | — | — | — | (21),(27),(42),(56–59) | 2 (conditions thread) | 2 (assumptions) | — | discharge semantics | future proposal | — | deferred |

Generality test (≥2 distinct uses): diagonal det — K1 Minkowski/Rindler + QE FRW (distinct papers, same math). Volume density — K1 dΣ/R + QE √h/√−g measures. Hyperbolic nodes — K1-only (honestly marked specific, not forced shared). Nothing promoted to `core/` (single-framework evidence only).

## 6. Placement decisions

PhysMath: diagonal-det certificate convention, VolumeDensity family, TransformMap certificate, hyperbolic-free (gr owns sinh/cosh). gr: opaque sinh/cosh nodes + fail-closed projection, Rindler/Minkowski/FRW frame definitions, x̄>0 + signature/convention assumptions, sketch updates. No registry change, no Apply/replay change, no core promotion, no frozen contact.

## 7. Proposed slice and exclusions

Implement R1–R3 + tests; contraction evaluation, general det, curvature, functional derivatives, amplitude/phase, discharge semantics explicitly excluded. s1 becomes represented-with-validated-measure (transformation still schematic); s2/s3 unchanged in status. QE minisuperspace gains storable det + volume; general ADM untouched.

## 8. Tests and acceptance criteria

Det correctness (diag(-1,1,1,1)→-1; diag(-x̄²,1,1,1)→-x̄² symbolic via existing Mul/Canon); non-diagonal + non-square rejection; volume assumption tests (missing x̄>0 → fail; wrong signature use flagged); hyperbolic projection-failure (all inputs); transform cert grammar/dimension + execution-unclaimed assertion; full 7-rule regression; required-dependency mutation (drop determinant ref → fail); anti-vacuity (det ≠ transform proof; volume ≠ measure evaluation). Acceptance: all green + K1 sketch s1 sub-status upgraded to measure-validated (s1 overall still UNSUPPORTED) + QE minisuperspace det/volume storable in a scratch test (not a fixture).

## 9. Remaining blockers and unresolved assumptions

Contraction evaluation (separate approval); transform-law execution; s3 discharge semantics; general det/curvature/functional-derivative/amplitude families; x̄>0 and signature assumptions need operator attestation at M1-style review; QE general-ADM support untouched by design.

## 10. Bounded approval request

Approve R1–R3 representation/validation slice (no new executable rule; contraction evaluation tracked separately) with the §8 battery as acceptance. On approval: implement on existing feature branches, rerun 55+43 batteries, update sketch sub-statuses, report. No merge/tag/release; K2 and QE derivations stay deferred.

## As-built record (R1–R3 implemented, uncommitted on feature branches)

- R2 detours through CANON+SIMPLIFY (Canon alone does not fold
  constants); Minkowski det exactly -1; Rindler symbolic det
  deterministic; FRW product form honestly distinct from a^6.
- InverseRef semantics corrected mid-slice: nil makes no inverse claim
  (was: self-inverse assertion); Minkowski/Rindler link no inverse.
- VolumeDensity with typed signature classes; reduction gated on cited
  conditions; sqrt/abs evaluation nowhere implemented.
- R1 GSinh/GCosh opaque nodes + TransformMap (validated structure, no
  execution entry point by construction).
- phys-math 61 pass; phys-lib 47 pass in 8 packages; 7 rules untouched;
  frozen repos untouched; no merges/tags/releases.
- K1 s1 now measure-validated (metrics/dets/volumes) but sketch stays
  UNSUPPORTED (transform unverified); execution still PARTIAL.
