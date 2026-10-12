# GUIDELINES.md — PhysTrace Lessons Learned and Engineering Casebook

**Purpose:** Preserve the practical lessons learned from actual PhysTrace implementation and review work. This document complements `SKILLS.md`: `AGENTS.md` describes the governing mission and general engineering principles; `GUIDELINES.md` records concrete failure modes, why they occurred, and the prevention mechanisms future agents must apply.

**Scope:** Applies across PhysMath, PhysLib, PhysTrace protocols, PhysArtifact, manifests, fixtures, tests, and future physics domains.

**Maintenance rule — additive from now on:** Append new lessons as dated, numbered entries. Do not silently erase or rewrite earlier lessons. If a lesson must change, add a dated correction or superseding note that explains why. Link related lessons rather than duplicating them. A report or plan is not evidence that a fix was implemented; distinguish *identified*, *planned*, *implemented*, and *independently verified* states.

**Relationship to `AGENTS.md`:** Follow `AGENTS.md` for enduring principles and architecture. Use this file as the incident-derived companion: the examples show what can go wrong even when code compiles, types validate, or tests pass.

---

## Lesson 1 — Source equations define the meaning of quantities; variable names do not

**Incident:** In the Dürr–Struyve minisuperspace fixture, an expression derived from equation (76) was named `qexpr`/`qref` and treated as the quantum potential `Q`, although equation (76) defines the quantum energy density `rho_Q`. Pressure was then differentiated from that density, and the density was divided by `a^3` a second time.

For the test amplitude \(A(a)=a^2+1\), with the \(\phi\)-dependent contribution absent and \(W=(a\partial_a)^2A=4a^2\), the source-defined relationships are:

\[
\rho_Q=\frac{2\bar\kappa^2}{a^4(a^2+1)},\qquad
Q=a^3\rho_Q=\frac{2\bar\kappa^2}{a(a^2+1)},
\]

\[
p_Q=-\frac{1}{3a^2}\frac{\partial Q}{\partial a}
=\frac{2\bar\kappa^2(3a^2+1)}{3a^4(a^2+1)^2}.
\]

**Why it happened:** The implementation lost the identity of a source-defined quantity between paper equation, variable name, stored expression, dependency chain, test golden, and manifest. Several artifacts were internally consistent while encoding the wrong physical meaning.

**Required prevention:**

1. Before implementing a source formula, record the exact paper equation and the physical/mathematical quantity it defines.
2. Draw the dependency chain between named quantities before coding.
3. Give distinct source quantities distinct representations when they have distinct definitions. Here, `RhoQRef`, `QRef`, and `PQRef` must be distinct; `QRef` must represent the relation \(Q=a^3\rho_Q\), not merely be a renamed density expression.
4. Derive expected results independently from the paper, not by copying the implementation's intermediate expression.
5. Test the defining relationships, not just the final expression's shape or display name.
6. Audit downstream consumers, including \(T_Q^{00}\), whenever a source quantity is corrected.

**Regression rule:** Include an explicit negative test for the historical error path \(p_{wrong}=-(3a^2)^{-1}\partial_a\rho_Q\). A test suite must reject the specific old mistake, not merely a vaguely defined “wrong form.”

**Status note:** The defect was confirmed in the original implementation report and attached fixture. Do not claim the corrective implementation is complete unless the updated code, live tests, commits, and verification report establish that fact.

## Lesson 2 — A test oracle must be independent of the implementation's interpretation

**Incident:** The pressure test described the equation (76) density expression as `Q`, then used that same expression to build its expected pressure. The actual implementation and test golden passed together because both carried the same source-interpretation error.

**Why it happened:** “Hand-derived” was taken to mean “manually assembled,” rather than independently derived from the paper. Manual calculation can still reproduce an implementation's mistaken assumption.

**Required prevention:**

- Write the expected formula from the paper's definitions before copying or inspecting the implementation's expression construction in detail.
- For a quantity chain, test intermediate identities as well as the final result. For the example above, test `W`, `rho_Q`, `Q = a^3*rho_Q`, `dQ/da`, `p_Q`, and the operands consumed by `T_Q^00`.
- Use independent hand-calculated golden values for small exact examples. If a small rational-point oracle is practical, use it as supplementary evidence rather than as a replacement for symbolic and provenance tests.
- Construct anti-goldens for known historical errors.
- Do not create the expected expression by calling the same helper that creates the actual expression.
- A passing test establishes agreement with its oracle; it does not establish that the oracle expresses the intended physics.

**Review question:** Ask both “Did the engine perform the declared operation correctly?” and “Was that the right operation for the quantity defined by the source?”

## Lesson 3 — Exact AST equality, mathematical equality, and source fidelity are different claims

PhysMath deliberately supports a bounded set of algebraic transformations. It does not automatically provide unrestricted distributivity, polynomial collection, or rational-function normalization.

**Failure modes:**

- Requiring exact AST equality for an equivalence the engine is not designed to normalize can make a correct implementation appear wrong or create pressure to add unrelated rewrite rules.
- Conversely, weakening the assertion until two expressions “look plausible” can conceal a genuine mathematical mistake.
- A comment stating that two expressions are algebraically equivalent is not execution evidence.

**Required prevention:**

1. Use exact structural equality when the canonical form is part of the operation's contract.
2. Use independent algebraic derivation to justify equivalence when the current engine cannot prove it.
3. Label the result accurately: executed/canonicalized, structurally validated, or hand-verified equivalent.
4. For important formulas, combine local exact tests, full-chain provenance checks, anti-goldens, and independent numeric checks at valid points when feasible.
5. Do not implement broad distributivity or term collection merely to satisfy one test golden. Track it as a reusable capability and prioritize it based on multiple real derivation needs.

**Example:** The composed operator \((a\partial_a)^2\) must preserve its product-rule term. A factored form such as \(a(2a+2a)\) may be honest evidence of the executed composition when the engine cannot collapse it to \(4a^2\), but only if the operation order and equivalence are independently established. It is not permission to claim the collapse was executed.

## Lesson 4 — Test objects must preserve variable identity, not just display labels

**Incident:** The FRW determinant anti-golden constructed \(a^6\) using an ephemeral symbol named `a`, while the fixture metric used a bound symbol tied to a `ParameterID`. Equal labels do not imply equal mathematical variables, so the comparison could never detect the intended collapse reliably.

**Required prevention:**

- Construct test alternatives using the exact same bound-variable identity and relevant parameter scope as the fixture.
- Add controls demonstrating the test distinguishes equal labels with different identities from the same variable with the wrong mathematical result.
- Compare decoded structures or identity-aware canonical expressions rather than relying only on artifact-reference inequality or matching printed names.
- Audit anti-goldens for whether they could actually fail under the targeted bug. A test that can never observe the intended failure is not a regression test.

## Lesson 5 — Derivative-generated intermediate forms must not invent domain restrictions

**Incident:** Computing \(\partial_\phi^2(a^2+\phi^2+1)=2\) required a `phiNonzero` condition because the first-power derivative route manufactured a `Pow(phi,0)` node, after which generic zero-power simplification required a nonzero-base condition.

**Why it matters:** The derivative is valid at \(\phi=0\). A representation artifact must not impose a false physical or mathematical domain restriction.

**Required prevention:**

- For an exact first power, implement the direct identity \(d(u^1)/dx=du/dx\) without constructing \(u^0\).
- Return \(1\) only where the derivative of the base is actually \(1\), such as a matching bound variable.
- Do not globally weaken generic zero-power or negative-power safety checks.
- Preserve legitimate conditions for reciprocal and negative-power expressions.
- Test the polynomial derivative over its true domain, including \(\phi=0\), while retaining negative tests for genuine denominator restrictions.
- Do not convert an engine limitation into an assumption in the paper's manifest.

## Lesson 6 — Every verification entry point must cover every evidence item it claims to verify

**Incident:** The fixture minted the FRW determinant derivation into `FRWDerivRef`, but `Fixture.Verify()` replayed only derivations in `f.Derivs`. The manifest's statement that all stored derivations were replayed was stronger than the actual implementation.

**Required prevention:**

- Maintain an explicit, auditable collection of all derivations claimed as executable evidence, including artifacts held in dedicated fields.
- Ensure `Verify()` resolves, decodes, and replays each one.
- Ensure store-aware checks for derived records still run; structural validity does not replace derivation replay.
- Test failure after deleting, corrupting, or substituting each important evidence reference, including determinant, Q, and pressure evidence.
- Keep the manifest's verification claim exactly aligned with what the authoritative verification function actually checks.
- Audit newly added derivation references when adding a field: minting an artifact without wiring it into verification is an easy omission.

## Lesson 7 — A stored relation is not a verified relation unless its operands are bound and checked

Adding a first-class `QRef` improves traceability only if the verifier establishes that the stored Q expression is related to the exact stored `RhoQRef` through \(Q=a^3\rho_Q\).

**Required prevention:**

- Bind the relation to the exact artifact identities or decoded expressions, the same bound variable identity, and the specified operation/construction.
- Verify the relation in the fixture's authoritative verification path.
- Reject substituted density, potential, pressure, or result references.
- Do not treat a relation written in a comment, manifest, or premise string as verification evidence.
- Where the existing engine cannot replay a relation with the available operations, label the check as a limited store-aware construction check; do not falsely claim a replayed mathematical proof.

This applies to all derived quantities, not only \(Q\): transformation outputs, contractions, determinants, volume densities, reduced variations, curvature, and future functional derivatives need the same discipline.

## Lesson 8 — Source-to-representation mapping belongs in durable project artifacts

Developer discipline alone did not prevent the Q/density confusion. The artifact set should make source traceability inspectable.

For source-derived quantities with meaningful identity, the manifest or equivalent protocol metadata should record, where appropriate:

- source equation and page/section;
- source quantity's mathematical name and meaning;
- its defining relation;
- the fixture field or artifact identity;
- direct dependencies and assumptions;
- whether the result is represented, executed, replayed, or only hand-verified.

At minimum, the minisuperspace mapping should distinguish \(\rho_Q\) (76), \(Q=a^3\rho_Q\), \(p_Q\) (77), and \(T_Q^{00}\) (78).

Preserve exact source identity, version, PDF digest, and relevant equation numbers. If source notation is ambiguous in extracted text, inspect the rendered source rather than guessing.

## Lesson 9 — Preserve rank, index, and space semantics in mathematical fixtures

Similar notation can represent different mathematical objects. For example, the Kronecker identity \(\delta^\mu_\nu\) is rank two, while the coordinate selector \(\delta^\mu_0\) is rank one.

**Required prevention:**

- Validate rank, index slots, variance, dimensions, and resolved space identities.
- Do not use an object of the wrong rank just because its printed name appears relevant.
- Do not infer equal spaces from equal dimensions or labels.
- Require result free-index identity and component order to match the operation contract.
- Ensure test fixtures instantiate actual typed structures with real identity references, not label-shaped substitutes.

## Lesson 10 — Provenance and replay must establish input-to-output linkage

The MINOR-04 and CONTRACT-EVAL work established a broader rule: a correct-looking result or a derivation object is insufficient unless its precise premises, operation, steps, and output are bound together.

For derived artifacts, validate where relevant:

- exact input artifact references and represented mathematical content;
- operation identity, parameters, and ordered step sequence;
- assumptions/conditions and their applicability;
- derivation linkage and replay;
- exact output reference and component order;
- any post-processing transformation as separate evidence.

A tensor contraction proves only the computed contraction, not a subsequent negation or integral. An FRW volume density does not prove its determinant derivation unless that evidence is checked. A reduced variation proves only the explicitly implemented projection, not general functional calculus.

Keep structural `Validate()` separate from store-aware verification. Pure decoders must not execute or certify mathematics. A verifier must fail closed on missing, wrong-family, mismatched, or corrupted evidence.

## Lesson 11 — Reduced models and partial computations must retain their limits

Minisuperspace is a symmetry-reduced toy model in the source paper, not a proof of the full Wheeler–DeWitt theory or general ADM formulation.

Similarly:

- applying the Wheeler–DeWitt operator to a test function is not solving the Wheeler–DeWitt equation;
- computing the guidance equation's partial-derivative right-hand side does not derive its time-evolution left-hand side;
- assembling \(T_Q^{00}\) does not establish every downstream physical identity;
- representing \(\sqrt{|\det h|}\) does not mean square-root or absolute-value evaluation occurred;
- an ansatz-derived variation projection is not a general functional derivative;
- a replayed derivation does not establish empirical confirmation.

Manifests and reports must state these limitations specifically, not use broad phrases such as “verified” when only fragments were executed.

## Lesson 12 — Conditions are mathematical premises, not test conveniences

Conditions such as \(a\ne0\), \(a>0\), \(N\ne0\), or an amplitude being nonzero must be introduced for a mathematical reason and bound to the operation they license.

**Do:** test both required-condition success and condition-removal failure. Check condition propagation through differentiation, simplification, substitution, and replay.

**Do not:** add a condition just to satisfy an intermediate representation artifact; reuse unrelated fixture conditions merely to make a rule pass; or claim a condition was discharged when the code only threaded its reference through the derivation.

For reductions such as \(|a|^3\to a^3\), record the sign condition and verify that the reduction is licensed by it. A condition being present is not the same as the condition being mathematically sufficient.

## Lesson 13 — Test verification coverage, not merely that `Verify()` returns nil

A happy-path `Verify()` test proves little if the evidence inventory itself omits a stored derivation.

Tests should mutate or remove important evidence and confirm the authoritative entry point fails. At minimum, include missing/corrupted FRW determinant evidence, missing/corrupted Q-to-density or pressure derivations, and an inconsistent Q relation.

Use an explicit evidence inventory or otherwise prove that all declared executable derivations are included. Keep the manifest's advertised coverage synchronized with this inventory.

## Lesson 14 — Changes to shared mathematical rules need collateral regression coverage

A narrow calculus fix can affect every later physics paper. The `d(u^1)/dx` fix must not weaken generic `Pow(u,0)` checks; negative-power differentiation must retain nonzero-base conditions; and a new rewrite must not silently change old derivation results.

Whenever changing shared PhysMath semantics:

1. Test the direct identity being added.
2. Test composite and repeated operations.
3. Test invalid/missing conditions.
4. Run existing expression, replay, determinism, and provenance tests.
5. Run cross-package physics fixtures that depend on the affected semantics.
6. Inspect result changes in the complete diff; do not rely on green tests alone.

A new rule is justified by new, precisely defined executable semantics, not by wanting a test to reach a particular output form. Equally, do not forbid needed rules merely to keep a historic registry count unchanged.

## Lesson 15 — Keep code, tests, manifests, and review reports in sync

A protocol can become misleading if implementation changes but manifest limitations do not, or if a report claims more derivations are replayed than `Verify()` covers.

Whenever a mathematical capability or source-derived quantity changes, audit all four:

1. implementation and typed artifact fields;
2. independent tests and negative/anti-golden cases;
3. manifest/source mapping and limitation records;
4. verification entry point and final report.

The report must distinguish “planned,” “implemented,” “tested,” “independently reviewed,” and “published.” Do not infer one status from another.

## Lesson 16 — A disciplined corrective workflow

For a source-fidelity defect, use this sequence:

1. **Locate the source equation.** Verify the equation number, notation, signs, operator grouping, and definitions from the authoritative source.
2. **Name every quantity.** Write the source-defined dependency graph before inspecting how the implementation currently names it.
3. **Derive independently.** Work a small symbolic or rational example by hand.
4. **Inspect the live baseline.** Record branches, commits, working trees, and current test outputs.
5. **Fix the root semantic error.** Correct the represented quantities and dependencies, not just the final output or golden.
6. **Build layered tests.** Use independent positive goldens, historical-error anti-goldens, domain-condition negatives, provenance mutations, and a small exact numerical oracle where feasible.
7. **Complete verification coverage.** Ensure every claimed derivation is included in the authoritative verification path.
8. **Reconcile protocol metadata.** Update source mappings and limitations to match the actual implementation.
9. **Run regression and adversarial checks.** Cover shared mathematics and cross-package consumers.
10. **Commit only after evidence is green.** Report actual output, exact commits, and remaining unsupported steps.

## Lesson 17 — Structural validity does not establish which variable was differentiated (2026-10-12)

**Observed failure:** `DIFFERENTIATE-001` and `TOTAL-DERIVATIVE-001` matched symbols by bare `ParameterID`. Parameter IDs are ordinal-per-function and globally reused, so a target-PID occurrence carrying a foreign display (an `a` sitting in the time slot, say) differentiated to 1 exactly as if it were the intended `t` — a wrong-variable derivation that replayed and verified structurally.

**Underlying cause:** the rules validated shape (declared slot, order, conditions) but never that a symbol carrying the target PID *denotes the declared variable*. Display names lived only inside expression payloads, with no authoritative per-function vocabulary to check them against.

**Why existing tests missed it:** each fixture used a single schematic vocabulary per function and never reused an ordinal PID across functions with different vocabularies; a wrong-variable request is byte-for-byte indistinguishable from a correct one when the fixture vocabulary happens to agree.

Whenever implementing or reviewing a differentiation boundary:

1. Fail closed without metadata: trusted differentiation requires `Function.ParamNames` and rejects its absence with `ErrMissingIdentity`. Never infer names from expression bodies.
2. Enforce the check once, in the shared `checkDiffTarget`, at both rule entry points; replay then inherits the contract automatically because it re-executes `Apply` per step.
3. Validate the *expanded* target tree under `TOTAL-DERIVATIVE-001`: substituted call arguments are exactly where a foreign display enters.
4. Keep vacuous differentiation valid: a validly declared variable absent from the target differentiates to 0. Body membership is never consulted, and non-target PIDs are independent quantities the check does not examine.
5. When migrating fixtures, supply explicit names from the mathematical variable vocabulary the targets and charts actually use; correct schematic bodies that disagree with it rather than loosening the check.

**Required regression:** `phys-math/diff_identity_test.go` pins the rejection reasons (wrong display at root and nested, missing names, PID reuse under a foreign vocabulary, substituted collisions), the acceptance set (vacuous zero, off-coordinate polar zero, independent quantities), the replay negatives, the codec round-trip with legacy tolerance, and the residual limitation below.

**Scope/limitations:** this lesson does not solve — it documents — the residual boundary that identical `(display, PID)` bytes under different declaring functions are indistinguishable to the check; that case is pinned as accepted behavior in `TestDifferentiateIdentityResidualLimitation`. It authorizes no PID redesign, no symbol-scope enrichment, and no body-membership requirement.

## How to add new lessons

Append a numbered lesson rather than overwriting historical entries. For each new incident, record:

1. **Observed failure:** what code/test/protocol did.
2. **Underlying cause:** the violated mathematical or architectural assumption.
3. **Why existing tests/reviews missed it.**
4. **Prevention rule:** what agents must do differently.
5. **Required regression:** the test or validation that would expose recurrence.
6. **Scope/limitations:** what the lesson does not prove or authorize.

Prefer evidence-backed, reusable lessons. Do not turn one incident into a universal rule without explaining the conditions under which it applies. If a lesson is superseded, append a correction with rationale and references; never silently erase the historical record.

---

## Closing principle

**A scientifically correct system must preserve the source's meaning across representation, implementation, independent tests, provenance, replay, manifest, and final claims. Agreement among components is not sufficient if they all inherit the same mistaken interpretation.**
