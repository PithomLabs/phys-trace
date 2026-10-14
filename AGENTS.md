# AGENTS.md — PhysTrace Engineering and Scientific Representation Principles

**Purpose:** This document teaches coding agents how to make sound decisions in the PhysTrace project. It captures the mission, architectural principles, mathematical standards, and implementation lessons established through project work. It is not merely a checklist: agents must understand *why* each rule exists and apply the underlying principle to new situations.

**Applies to:** `phys-math`, `phys-lib`, `phys-trace` / PhysTrace protocols, `phys-artifact`, their tests, manifests, derivation protocols, and future physics-domain packages.

**Authority:** Follow applicable repository `AGENTS.md` files and task-specific approved specifications alongside this document. This document states cross-project defaults. If instructions conflict, do not silently discard either: identify the concrete conflict, follow the higher-authority instruction, and report the issue. Do not use a local example to overrule a general project invariant.

---

## Physics-paper derivation specifications

Before implementing a new paper derivation or materially expanding an existing one, read `DERIVATION_SPEC_GUIDELINES.md` and follow its normative source-fidelity, quantity-mapping, assumption, independent-oracle, provenance, and readiness requirements. Use `GUIDELINES.md` as the additive record of lessons learned. If an adequate specification already exists, verify it against the current source and code rather than restarting planning. Do not let the specification gate become an excuse for endless planning: block only steps with material unresolved ambiguity and proceed with well-specified, useful increments under the authorized scope.

## 1. The mission: build scientific capability, not a collection of demos

PhysTrace is a research platform for representing physics papers faithfully and progressively executing and verifying their mathematical derivations. Its long-term reach includes quantum mechanics, cosmology, quantum gravity, and potentially a theory of everything.

The project must enable a scientist or researcher to trace a claim from a source paper through its assumptions, mathematical objects, operations, derivation steps, and resulting artifacts—and to see precisely which parts have and have not been established.

### Do

- Build the mathematical capabilities that real derivations require, even when they need a substantial, cross-package implementation.
- Prefer reusable foundations that unlock multiple papers and derivations over paper-specific workarounds.
- Move from research to implementation. Produce working capabilities and evidence instead of repeatedly generating plans for already understood tasks.
- Choose the next action by its scientific impact, reuse, testability, and dependency leverage.
- Keep the architecture extensible across theories and physics domains.
- Treat a tractable model, such as minisuperspace, as a useful starting point—not as a substitute for the more general theory it approximates or reduces.

### Do not

- Artificially restrict the project to existing operations when new mathematics is needed.
- Avoid tensor calculus, differential geometry, functional derivatives, curvature, constraint algebra, or other substantial capabilities merely because they are difficult or span packages.
- replace a missing capability with a paper-specific shortcut, opaque symbol, hard-coded result, or unexplained status flag.
- confuse a successful demo with broad scientific coverage.
- claim an entire paper or theory is verified because one component or reduced model works.

**Decision principle:** Pursue the capability the physics requires, at the scope the mission justifies, and make every scientific claim earn its status through evidence.

## 2. Be decisive, but not careless

Agents are expected to exercise architectural judgment. They should not ask the user to approve every routine design choice once the mission and implementation authority are clear.

For each task, satisfy two obligations:

1. **Deliver the immediate objective.** Fix the real blocker or implement the authorized capability, with tests and evidence.
2. **Identify the next high-impact action.** Recommend the most valuable reusable follow-on based on the actual dependency graph and remaining blockers. A recommendation does not automatically authorize unrelated work.

### Expected operating rhythm

1. Inspect live repository state, governing instructions, relevant code, and authoritative source material.
2. State a short execution plan when the task is substantial.
3. Make concrete architectural decisions using the principles here.
4. Implement a coherent vertical slice, including representations, computation, provenance, tests, and integration as appropriate.
5. Run the relevant tests and static checks; fix material defects discovered in the work.
6. Review the actual diff against the intended scope and invariants.
7. Commit locally when authorized and useful; report exact evidence and the next high-impact action.

Do not turn each step into another general planning cycle. Reopen a settled decision only when concrete evidence shows a defect, contradiction, missing prerequisite, or material regression. Avoid spending the project's time on style-only nitpicks when a meaningful scientific capability is ready to build.

## 3. Architectural responsibilities

Keep responsibilities distinct so that shared mathematics can grow without absorbing the ontology of a particular theory.

### PhysArtifact

- Own and preserve the artifact-level contracts, identities, canonical encodings, and reference semantics assigned to this layer.
- Preserve traceability between artifacts and their dependencies.
- Do not let downstream packages quietly change the meaning of an existing artifact family.
- Prefer explicit versioned representation changes over silently breaking persisted artifacts.

### PhysMath (`phys-math`)

- Own theory-neutral mathematical objects and executable mathematical operations: expressions, functions, indices, tensors, derivatives, algebra, and other general foundations as justified by real derivations.
- Define operations precisely enough that independent expected results can be computed and tested.
- Keep physics ontology, theory interpretation, and model-specific assumptions out of the shared mathematical layer.
- Expand the mathematical layer when a second or broader use justifies a reusable abstraction; do not artificially preserve an obsolete minimal rule set.

### PhysLib (`phys-lib`)

- Own physical frameworks, models, ontology, conventions, assumptions, and theory-specific composition of shared mathematics.
- Put Bohmian-framework concepts in an appropriate `phys-lib/bohm` package; keep GR-specific semantics in GR packages.
- Reuse `phys-math` rather than duplicating mathematical machinery inside a theory package.
- Promote a capability into a shared core only when multiple use cases establish that it is genuinely common. Do not promote code based on speculation alone.

### PhysTrace and paper protocols

- Map source claims to core ideas, derivation units, assumptions, dependencies, mathematical artifacts, and evidence.
- Preserve source identity and explain what is covered, excluded, represented, executed, or unsupported.
- Keep paper-ingestion and protocol logic aligned with the actual mathematical capabilities; never imply a derivation is executable merely because a schema can store it.

### Dependency discipline

Preserve the established one-way architecture and verify real module edges before changing them. The current intended lower-level direction is `phys-lib → phys-math → phys-artifact`; PhysTrace coordinates paper/protocol workflows around these capabilities. Do not introduce reverse dependencies, cycles, or framework-specific imports into theory-neutral code to save a few lines.

## 4. Represent mathematical meaning, not paper notation

This is a project-wide invariant across PhysMath, PhysLib, PhysTrace, and PhysArtifact:

> **Represent what the mathematics means, not merely what the paper prints.**

Do not store a derivative, function application, tensor operation, or other mathematical operation as a standalone symbol or source-text string when it can and should be represented structurally.

### Preferred representation patterns

| Mathematical meaning | Structural representation principle |
|---|---|
| `a(t)` | Function application with an explicit function and argument. |
| `da/dt` or `ȧ` | A total-time-derivative operation applied to `a(t)`. |
| `d²a/dt²` or `ä` | A nested total-time-derivative operation. |
| `∂U/∂φ` | A partial-derivative operation applied to the represented function/expression `U(φ)`. |
| `d/dt (a³ φ̇/N)` | A total derivative applied to a quotient expression, not a special symbol for the entire term. |
| `a³` | A power node with base `a` and exponent `3`. |
| `T_ab ξ^b` | Typed tensor/vector operands plus an explicit contracted slot and a computed result. |
| A condition such as `a > 0` | A represented condition connected to the operation or reduction it licenses. |

A genuine independent variable may be a symbol. The prohibition is against using a symbol to conceal a computable operation or a structured relationship.

### Do

- Preserve operand identities, argument lists, index positions, dimensions, variance, space identities, and dependencies.
- Keep source notation and equation numbers as provenance metadata, not as replacements for the math.
- Represent opaque functions explicitly when their internal definition is unavailable; state the limit rather than pretending the function was evaluated.
- Keep **partial derivatives** and **total derivatives** semantically distinct.
- Represent higher derivatives compositionally, and implement reusable chain, product, quotient, and power rules when required.

### Do not

- encode `φ̇`, `ä`, `∂Q/∂a`, or a tensor contraction as an unexplained symbol just to make a paper protocol pass;
- use equation strings as the canonical mathematical object;
- hard-code a paper's final expression or expected test output into the implementation path;
- treat a derivative-looking symbol as evidence that differentiation occurred.

### Differential operators must preserve their actual meaning

For example:

\[
(a\,\partial_a)^2 f = a\,\partial_a(a\,\partial_a f),
\]

which is generally **not** the same as \(a^2\partial_a^2f\). Tests must catch this distinction and verify the product-rule contribution.

A total derivative must operate on time-dependent functions such as `a(t)`, `φ(t)`, and `N(t)`. It must not be represented as a partial derivative merely because both eventually use derivative machinery. Conversely, a partial derivative must preserve which variable is held independent.

## 5. Mathematics is allowed to grow

Do not assume the existing mathematical operation registry or type system is complete. It is an evolving implementation of reusable mathematics.

When a derivation needs a missing capability, determine its mathematical semantics and implement it at the appropriate layer. Candidates may include:

- negative-integer powers with valid nonzero-base conditions;
- total-time derivatives and higher derivatives;
- index-aware tensor contraction and general tensor ranks;
- Kronecker selectors and identity tensors with the correct ranks;
- symmetry-reduced variation and, when justified, general functional derivatives;
- covariant derivatives, curvature, differential geometry, and constraint-related operations.

These are examples, not a mandate to build every subsystem at once. Choose based on concrete derivation dependencies, but do not reject a capability solely because its implementation is substantial.

### Mathematical rules for new capabilities

Every executable operation must specify:

1. Accepted operand families and shapes.
2. Exact semantics and output shape.
3. Index, dimension, variance, and domain requirements where applicable.
4. Preconditions and which assumptions discharge them.
5. Deterministic output and canonical representation.
6. Failure behavior for unsupported or invalid inputs.
7. Provenance and replay guarantees.
8. Independent positive and negative tests.

A new `RuleID` is appropriate when a genuinely new executable semantic operation is needed. Do not force a new operation into a misleading existing rule merely to avoid updating a registry. Conversely, do not add a new rule when the operation can be correctly expressed using existing compositional operations and driver logic. Decide based on mathematical semantics and replay evidence, not registry-count preference.

### Domain and assumptions are part of the mathematics

- Negative powers require a nonzero-base condition where mathematically necessary.
- A quotient needs a nonzero denominator condition.
- A reduction from `|a|` or `sqrt(a²)` to `a` requires the appropriate sign assumption.
- An index contraction requires compatible spaces, dimensions, labels, and variance under the declared semantics.
- A symmetry-reduced variation must derive its factors from the represented ansatz, spatial measure, and normalization conventions.

Never silently assume a condition merely because a test fixture uses values for which it happens to hold.

## 6. Tensor and index semantics must be real

Tensor computations must be driven by typed structure and explicit index semantics, not by visually plausible notation.

- Preserve tensor rank, ordered index slots, variance, component ordering, and dimension.
- Require compatible index spaces where the operation requires them; equal dimension alone does not imply that two spaces are identical.
- Fail closed on unspecified or incompatible variance when the operation needs variance to determine meaning.
- Preserve free-index identity in outputs.
- Check row-major layout formulas explicitly and test both contraction-slot orientations.
- Reject invalid rank, slot, shape, space, or component references at the intended boundary.

Do not confuse similarly named mathematical objects. For example, `δᵘᵥ` is a rank-two identity tensor; `δᵘ₀` is a rank-one coordinate selector/vector for the chosen basis. A rank-two identity must not be used as a rank-one selector merely because both use the word “Kronecker.”

A generic tensor evaluator is a legitimate future foundation when it supports real use cases. It must still have precise index semantics, bounded resource behavior, independent goldens, and traceable output identities.

## 7. Paper ingestion: source fidelity and idea-centric representation

A physics paper is not a bag of equations. Represent its scientific claims and the derivations that support them.

### For each paper

1. Pin the authoritative source: identifier/version, exact source file, digest when available, and relevant source locations.
2. Inspect mathematical notation carefully. Check rendered pages when extracted text may have lost subscripts, superscripts, signs, integral measures, or operator grouping.
3. Propose a core-idea inventory. Explain why each proposed idea deserves separate representation.
4. Map each accepted core idea to its derivation unit(s). Do not create one unit per equation by default.
5. Record assumptions, regimes, interpretations, predictions where relevant, dependencies, limitations, and unsupported transitions.
6. Distinguish what the authors explicitly derive from what they assert, cite, omit, or call straightforward.
7. Preserve different derivation routes when the paper presents them separately. A shared conclusion does not prove both routes independently.
8. Treat model reductions, such as minisuperspace, as reductions with declared assumptions—not as proof of the unreduced theory.

A candidate core-idea count is provisional until adjudicated. Do not manufacture certainty about the inventory merely to complete a manifest.

### Source fidelity rules

- Never silently “fix” a paper's equations or conventions. Record real ambiguities and resolve them only with evidence.
- Preserve the paper's sign conventions, operator ordering, coordinate/chart assumptions, and normalization choices.
- Distinguish source-declared dependencies from reviewer-inferred dependencies.
- Identify omitted derivation steps and important cited results as dependencies instead of falsely claiming that they were independently derived.
- Do not claim experimental agreement where no empirical evidence has been represented.

## 8. Separate representation, computation, verification, interpretation, and reality

These are different evidence levels, not interchangeable labels.

1. **Represented:** a mathematical object or claim has a faithful structural representation.
2. **Structurally validated:** schema, types, references, shapes, or index invariants have passed.
3. **Executed:** a supported operation actually computed a result.
4. **Replayed/verified:** recorded steps reproduce the result from the bound premises under the specified operation semantics.
5. **Physically interpreted:** the mathematical result is connected to a physical model under explicit assumptions and interpretation.
6. **Empirically supported:** relevant observations or experiments support the claim.

Do not advance a status solely because an artifact exists or a schema test passes. A structural certificate is not a mathematical proof; replay proves only what the operations and premises actually justify; a derivation is not empirical confirmation.

Use project capability states consistently (`SUPPORTED`, `PARTIALLY SUPPORTED`, `REPRESENTATIONAL ONLY`, `UNSUPPORTED`, or the specific schema's equivalent). If classifications are operationally defined in the relevant spec, follow that definition. A step with executable fragments but an unimplemented transition should remain partial or unsupported as appropriate.

Never use a `verified: true` flag or caller assertion in place of verifiable derivation evidence.

## 9. Derived-artifact provenance is mandatory

A result produced from inputs must be traceable to the exact inputs, operation, steps, and output. This principle applies to determinants, contractions, derivatives, variations, curvature, and future derived artifacts.

Where applicable, the provenance contract must bind:

- exact operand artifact references and relevant revisions;
- operation identity, version, and parameters;
- required assumptions and conditions;
- intermediate derivation steps and their dependency links;
- exact result reference and component ordering;
- the status actually established by execution and replay.

### Validation layers

Keep these responsibilities separate:

- `Validate()` (or its equivalent) checks structural grammar and local invariants. It performs no hidden store I/O.
- Pure decoders parse canonical representations and reject malformed, incomplete, or noncanonical payloads. Decoding alone does not establish provenance.
- Store-aware checkers resolve references, verify families, check links, replay derivations, recompute required premises/results, and enforce provenance.
- Consumer gates call the required store-aware validation before accepting evidence as valid.

Do not accept a result merely because a caller supplied a plausible output reference. Compare the exact bound output with the result actually generated or reconstructed from the declared inputs.

### Exact guarantees, not exaggerated ones

State precisely what provenance establishes. A determinant check that recomputes diagonal-expression content binds that content; it does not necessarily establish unique identity of a higher-level metric artifact when distinct artifact references contain identical relevant components. Do not promise stronger identity semantics than the evidence encodes.

### Failure and storage behavior

- Fail closed on missing, unresolved, wrong-family, malformed, or mismatched dependencies.
- A failed evaluation must never produce an accepted certificate or valid-looking result status.
- If the store is content-addressed but nontransactional, do not pretend it can roll back minted intermediate artifacts. Control the acceptance boundary, enforce resource limits before minting where specified, and disclose orphaned intermediate artifacts when relevant.
- Never let recomputing expected evidence silently repair missing required evidence before replay or availability checks have had a chance to reject it.

## 10. Derivation and replay discipline

A derivation is an explicit, ordered computational trace—not prose and not a list of expected outputs.

- Verify the actual semantics of `Apply` and `Replay` before designing a derivation format.
- Do not assume batched execution if each rule accepts exactly one expression.
- Link every step to its actual premise or prior output.
- In a flat multi-output derivation, explicitly verify which indexed step pair produces each output component and ensure result order matches tensor component order.
- Reject missing, extra, reordered, duplicated, swapped, or substituted steps when they violate the declared operation contract.
- Verify full reference equality where required, not only a digest or a caller-supplied name.
- Keep calculation of a result distinct from later transformations of that result. For example, contraction evidence for `wₐ = Tₐᵦ ξᵇ` does not by itself prove `jₐ = −wₐ`; the negation needs its own linked evidence.

## 11. Tests must establish mathematical correctness

Passing tests are necessary, not sufficient. Tests must be able to falsify an incorrect implementation.

### Positive tests

- Use independently hand-derived expected values or expressions.
- Cover different shapes, both relevant operation orientations, symbolic inputs, boundary values, and repeated execution where applicable.
- Verify exact result references and ordered components, not just non-nil outputs.
- Verify replay from the actual dependencies.
- Test operator composition explicitly (e.g. `(a·∂ₐ)^2`, not just a first derivative).

### Negative/adversarial tests

Test every semantic gate and fail at the intended boundary:

- wrong ranks, slots, variance, labels, spaces, dimensions;
- missing or wrong-family references;
- missing assumptions or invalid domains;
- substituted operands or output references;
- stale evidence after a relevant input changes;
- swapped component/step mappings;
- malformed, truncated, trailing, or noncanonical encodings;
- reordered, extra, missing, or corrupted derivation steps;
- exceeded computational limits.

A negative test that merely checks “an error occurred” is weak if it could fail for an unrelated reason. Assert the specific failure boundary or observable invariant.

### Independent oracles

- Do not generate expected values solely by invoking the implementation under test.
- Use small hand-calculated examples, identities, limiting cases, and mathematically justified invariants.
- For symbolic operators, test cases designed to distinguish plausible but incorrect implementations.
- Test the same shared capability in more than one consumer when feasible, but keep theory-specific assumptions out of the shared math tests.

### Verification report

Run the actual commands and report results from the checked-out commits. Include full-suite tests, focused tests, static checks such as `go vet`, formatting checks, and any project-established purity checks where relevant. Reconcile inconsistent counts from earlier reports rather than guessing. Do not report unrun tests as passing.

## 12. Resource bounds and computational honesty

Reusable math can be broad without being unbounded.

- Estimate time, memory, artifact count, and derivation growth for new operations.
- Enforce documented bounds before expensive intermediate construction or minting.
- Use deterministic failure behavior and meaningful errors.
- Test both sides of each resource boundary.
- Avoid claims that a dimension limit makes an operation cheap unless the worst-case work actually supports that claim.
- Do not silently truncate, approximate, or simplify beyond the operation's declared mathematics.

If a general capability needs an explicit limit in its initial version, implement the useful bounded capability while documenting the boundary. Do not use the existence of limits as a reason to avoid the foundation altogether.

## 13. Implementation and review workflow

### Before code changes

- Read package and workspace instructions (`AGENTS.md`) and relevant governing specs.
- Inspect the live branches, commit history, working-tree state, and exact relevant code paths.
- Establish baseline tests and record actual results.
- Identify the scientific claim being enabled and its dependency chain.
- Decide which layer owns the new semantics and why.

### During implementation

- Prefer a coherent end-to-end slice: representation, computation, provenance, protocol integration, tests, and classification update where appropriate.
- Reuse existing mathematical primitives when they truly implement the required semantics.
- Add abstractions where they create real reuse; avoid both duplicated one-off code and speculative frameworks with no consumer.
- Do not casually modify frozen repositories or unrelated modules.
- Preserve compatibility deliberately. If artifact encodings or registries must change, document and test the consequences.
- Implement the broad authorized task rather than stopping at the first missing primitive or returning a plan-only response.

### After implementation

- Format and run focused tests, full test suites, and applicable static checks.
- Inspect the actual diff for accidental scope changes, weakened assertions, dependency reversals, and claims stronger than the code supports.
- Perform an adversarial review at meaningful milestones; fix material findings rather than leaving known correctness blockers.
- Commit locally on the authorized feature branches when useful and authorized.
- Provide exact commit hashes, test outputs, changed paths, remaining limitations, and the next highest-impact capability.

## 14. Repository operations and authorization boundaries

Default to protecting repository state and avoiding irreversible external actions.

- Local implementation and local commits are permitted when authorized by the task and repository instructions.
- Do **not** push, merge, tag, publish, release, or modify protected remote branches without explicit separate authorization.
- Do not commit `go.work`, local-only `replace` directives, or developer-specific workspace configuration unless a task explicitly changes that policy and the change is justified.
- Verify whether a branch is ahead of or behind its remote; “tree clean” and “in sync” are different claims.
- Keep frozen legacy repositories read-only unless explicitly reopened.
- Never claim a remote operation succeeded without observing its actual result.

## 15. Architecture decisions and scope control

When a new capability could be implemented in several ways:

1. Identify the exact mathematical obligation and the downstream derivations it enables.
2. Compare approaches by correctness, reuse, testability, provenance, integration cost, and future generalization.
3. Prefer the smallest **complete and mathematically honest** capability, not the smallest diff.
4. Generalize when several real derivations share the same semantic need, or when a general primitive is clearly required by the mission.
5. Avoid broad frameworks whose meaning is underspecified or whose correctness cannot yet be tested.
6. Record legitimate trade-offs and continue implementation within the authorized scope.

Do not reopen a settled architecture because another style is imaginable. Reopen it when evidence shows a false guarantee, missing semantic prerequisite, incompatible contract, or material lack of reuse.

## 16. Example lessons from project work

These are enduring lessons, not special-case rules tied to one paper.

- **Determinant provenance:** a structurally valid volume-density record is not enough. Bind it to replayable determinant evidence and verify the relevant input content. Be honest about whether the guarantee is component-content binding or exact parent-artifact identity.
- **Contraction provenance:** an operation record must bind the actual tensor, vector, contracted slot, output tensor, and derivation. An older certificate that merely describes two slots in one tensor cannot stand in for the complete tensor-vector operation.
- **Multi-output derivations:** a flat sequence of steps is safe only if the checker verifies every component's exact step linkage and the final ordered output mapping.
- **K1 dual evidence:** computed contraction `wₐ` and negated current `jₐ = −wₐ` are separate results with separate evidence; neither proves the unsupported energy integral.
- **Minisuperspace:** a reduced model is a useful executable target but must not be confused with the full Wheeler–DeWitt or general ADM theory.
- **Kronecker objects:** distinguish a rank-two identity tensor from a rank-one coordinate selector. Mathematical names are not substitutes for rank and index semantics.
- **Symmetry-reduced variation:** derive factors from the represented ansatz, spatial measure, and normalization. Do not hard-code the answer's volume factors as magic constants.
- **Bohmian package placement:** framework-level Bohmian assumptions belong in a reusable `phys-lib/bohm` layer; shared derivative and tensor semantics belong in PhysMath.
- **Registry growth:** a fixed historic rule count is not a goal. Add rules when new executable semantics require them and verify that they are reusable, precise, and replayable.

## 17. Minimum quality bar for “done”

A task is not complete merely because files were edited or tests compile. The completion report must answer:

- What scientific capability is now represented or executable?
- What exact semantics and assumptions govern it?
- Which inputs and results are bound by provenance?
- Which independent mathematical examples establish correctness?
- Which invalid cases are rejected, and at which boundary?
- Which tests and checks were actually run, and what were the observed results?
- What remains unsupported, and why?
- What is the next highest-impact reusable action?
- Were repository and publishing boundaries respected?

The agent must not mark an entire theorem, paper, physical model, or derivation verified when only a subset is covered.

---

## Final operating principle

Build what the physics requires, at the scope the mission justifies. Represent mathematical meaning structurally. Implement reusable mathematics instead of paper-specific shortcuts. Bind every result to its inputs and derivation. Test with independent evidence. State limitations honestly. Keep moving toward connected, executable physics derivations without unnecessary planning cycles.
