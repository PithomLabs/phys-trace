# DERIVATION_SPEC_GUIDELINES.md — Normative Standard for Physics-Paper Derivation Specifications

**Status:** Normative project guidance for PhysTrace derivation ingestion, specification, implementation, and verification.

**Applies to:** PhysMath, PhysLib, PhysTrace, PhysArtifact, paper-specific protocols, manifests, fixtures, derivation tests, and future physics domains.

**Purpose:** Establish an implementation-ready, source-faithful specification for every physics-paper derivation that PhysTrace undertakes. The goal is not merely to make a coding agent follow a list of instructions. The goal is to ensure the agent understands the scientific claim, what the paper actually establishes, which assumptions make each step valid, how every quantity is represented, what evidence would verify each transformation, and what remains unknown or unsupported.

This document complements, rather than replaces:

- `AGENTS.md` — project mission, architecture, governing invariants, and agent operating rules.
- `GUIDELINES.md` — additive lessons learned from concrete failures and reviews.
- Package-level `AGENTS.md`, specifications, and existing protocol contracts — local implementation requirements.

If a package-specific contract supplies stricter requirements, follow it unless it conflicts with a higher-level project invariant. Do not import another repository's fixed operation counts, status enums, or package boundaries into PhysTrace without checking whether they apply.

The words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are normative.

---

## 1. The central rule: specify the science before specifying the code

A derivation specification MUST begin with the scientific meaning of the target, not with the function the agent plans to write or the expression tree that is easiest to construct.

For every target claim, the agent MUST be able to answer:

1. What exactly does the paper claim or derive?
2. Where does the source establish that claim?
3. What do the symbols and mathematical objects mean in that equation and framework?
4. Which premises, assumptions, conventions, regimes, and domain restrictions are required?
5. Which previous results are used, and are they derived in the paper, cited, stipulated, or merely asserted?
6. Which mathematical transformation connects the premises to the result?
7. What concrete representation will encode each quantity and operation?
8. What independent evidence could show that the implementation computes the intended mathematics?
9. What is still only represented, assumed, unverified, or empirically untested?

A specification that answers only “which equations should be stored” is incomplete. A specification that names the required Go functions without fixing their semantics is also incomplete.

### The source-to-result chain

Every significant result MUST be traceable through this chain:

`source claim → source-defined quantities → premises and assumptions → mathematical operations → intermediate results → target result → verification evidence → stated limitations`

The chain MUST be represented in the specification and, where relevant, reflected in durable artifacts such as the manifest, expression references, derivations, certificates, and tests.

## 2. Read the authoritative context before drafting a specification

Before specifying new work, the agent MUST inspect:

- the project's root `AGENTS.md` and `GUIDELINES.md`;
- relevant package-level agent instructions and architectural specifications;
- existing implementations and protocol examples that are closest to the target;
- current manifests, source references, capability maps, tests, and known limitations;
- the actual branch, commit, working-tree state, and module/workspace boundaries before implementation.

The agent MUST determine whether a target or part of its derivation already exists. Existing stored results can be dependencies or comparison targets, but a derivation MUST NOT use the stored target itself as an unexplained premise and then claim to have derived it.

The agent MUST inspect the relevant source paper directly. Extracted text is useful for navigation, but rendered source pages MUST be checked when extraction may have corrupted subscripts, superscripts, signs, fractions, indices, integrals, grouping, or operator order.

### Reuse the repository's contract without inheriting irrelevant restrictions

Existing package protocols provide useful patterns: human-oriented README, machine-readable manifest, strict manifest validation, executable derivation tests, intermediate checkpoints, independent expected values, provenance/assumption checks, and replay verification. Reuse the applicable contract.

Do NOT assume that another package's limited operation set, fixed rule count, or non-goals are permanent restrictions on PhysTrace. The mathematical foundation is allowed to grow when the derivation requires it.

## 3. Pin the source and record the evidence boundary

Every paper-specific specification MUST identify the authoritative source sufficiently to reproduce the interpretation:

- paper title and authors, when available;
- DOI, arXiv identifier/version, journal or publication version, or equivalent stable identifier;
- authoritative local source path or retrieval reference;
- source-file digest where available and relevant;
- relevant section, page, equation, figure, appendix, or theorem references;
- any difference between versions that affects the derivation;
- conventions and definitions established elsewhere in the paper that the target relies on.

For important source statements, distinguish these evidence classes:

- **EXPLICIT_SOURCE:** directly stated or derived by the paper at the cited location;
- **SOURCE_CITES_EXTERNAL_RESULT:** relies on a referenced source not independently derived in the target paper;
- **SOURCE_ASSERTS_OR_SKETCHES:** claimed or described without enough displayed derivation to reproduce it;
- **SPEC_INFERENCE:** reconstructed or inferred by the specification author;
- **UNRESOLVED:** ambiguous, inconsistent, or unsupported by the inspected source.

These labels MAY be mapped to the project's existing schema, but their semantic distinctions MUST be preserved. Do not silently turn a citation, author assertion, or reviewer inference into a derivation established by PhysTrace.

If two available paper copies differ, record the discrepancy and specify which version governs. Do not quietly blend equations or conventions from different versions.

## 4. Inventory scientific ideas before splitting derivations

A paper is not a bag of equations. Before implementation, create a **proposed core-idea inventory** and a **derivation-unit map**.

### 4.1 Core-idea inventory

For each proposed core idea, record:

- a stable provisional ID;
- the physical or mathematical claim;
- why it merits representation as a distinct idea;
- source locations;
- its framework, regime, and interpretation;
- assumptions and prerequisites;
- target result(s);
- known limits, anomalies, caveats, or falsification conditions where relevant.

The proposed number of core ideas MUST remain provisional until adjudicated where human adjudication is part of the research workflow. Do not fabricate certainty merely to fill a schema.

### 4.2 Derivation units

Map each accepted core idea to one or more coherent derivation units. A derivation unit is a scientific reasoning unit, not automatically one source equation or one code function.

Do not create one unit per equation by default. Algebraic manipulations, coordinate substitutions, intermediate lemmas, and component expansions normally belong inside a larger derivation unless they establish an independently meaningful result.

Conversely, keep genuinely distinct derivation routes separate when they differ in premises, reasoning, assumptions, or verification obligations—even when they reach the same target equation. A shared conclusion does not prove each route independently.

For every unit, specify:

- stable ID and parent core idea;
- target statement and source location;
- ordered premises and dependencies;
- assumptions and applicable regime;
- substeps and intermediate results;
- expected output and its mathematical identity;
- evidence class of each substep;
- representation/computation/verification status;
- acceptance tests and explicit exclusions.

## 5. Preserve quantity identity and source-defined relationships

**Quantity identity is part of the specification.** Variable names alone do not establish mathematical meaning.

For every important quantity, record its source definition, not merely its displayed symbol. Related quantities MUST remain distinct when the paper defines distinct mathematical roles for them.

For example, if a paper defines

\[
\rho_Q = Q/a^3, \qquad p_Q = -(3a^2)^{-1}\partial_a Q,
\]

the specification MUST distinguish `rho_Q`, `Q`, and `p_Q`, record the relationship `Q = a^3 rho_Q`, and state explicitly that the pressure derivative acts on `Q`, not on `rho_Q`.

### Required quantity map

For each significant source-derived quantity, record the following fields directly or by an equivalent unambiguous schema:

| Field | Required meaning |
|---|---|
| `quantity_id` | Stable identifier distinct from its display label |
| `source_equation` | Equation/theorem/definition that establishes the quantity |
| `source_quantity` | Mathematical name and meaning in context |
| `defining_relation` | Actual source-defined equation or construction |
| `representation` | Expected artifact type and structural expression/function/operator form |
| `dependencies` | Exact input quantities and premise identities |
| `assumptions` | Conditions that license its definition or transformation |
| `expected_result` | Independent symbolic or numerical target, where applicable |
| `evidence_status` | What is represented, computed, replayed, or still unsupported |
| `limitations` | Gaps between source claim and implementation evidence |

The field names MAY follow existing manifest conventions, but the information MUST be recoverable without guessing.

### Avoid semantic drift

The agent MUST NOT:

- name an expression after a physical quantity it does not represent;
- conflate related quantities just because they share much of an expression;
- copy implementation naming into an expected-value test oracle;
- state a relation in prose while leaving its operand identities unchecked;
- derive a downstream result from the wrong intermediate quantity;
- treat matching display labels as matching mathematical variable identities.

Audit neighboring quantities for the same class of error when a semantic misidentification is found.

## 6. Represent mathematical meaning, not paper notation

Paper typography is source metadata, not the canonical representation.

Operations with mathematical meaning MUST be represented as structured expression, function, operator, tensor, relation, condition, or other correctly typed objects. Do not encode a computable operation as an unexplained symbol or paper-text string merely because that is convenient.

Examples:

- `a(t)` is a function application with an explicit function and argument.
- `da/dt` is a total-time-derivative operation on `a(t)`.
- `d²a/dt²` is a nested total derivative.
- `∂U/∂φ` is a partial-derivative operation with an explicit differentiation variable.
- a tensor contraction binds typed operands, slot positions, index metadata, and a computed output;
- a domain condition such as `a > 0` is a structured premise attached to the reduction it licenses.

A genuine independent variable MAY be represented as a symbol. An opaque function MAY remain uninterpreted if its definition is absent, but its application and dependencies must still be represented structurally wherever the model supports them.

### Mathematical structure that must not be flattened away

The specification MUST preserve, where applicable:

- operator composition and ordering;
- function arguments and bound-variable identity;
- partial versus total versus functional derivatives;
- tensor rank, ordered slots, variance, index labels, index spaces, and component order;
- dimensions, physical kinds, units, signatures, coordinate charts, and orientation;
- integration domains, measures, boundaries, and normalization factors;
- branches, nonzero denominators, positivity, differentiability, and other domain assumptions;
- symmetry-reduction ansätze and the maps/Jacobians they induce;
- formal versus physical and empirical meaning.

For example,

\[
(a\partial_a)^2 f = a\partial_a(a\partial_a f)
\]

is an operator composition and is generally not the same as \(a^2\partial_a^2 f\). The specification MUST state the correct composition and require tests that distinguish it from the plausible incorrect form.

A general tensor identity and a coordinate selector are not interchangeable. A rank-two Kronecker identity \(\delta^\mu_\nu\) is not a rank-one selector \(\delta^\mu_0\). The representation must match the object's rank and semantics.

## 7. Specify assumptions, conventions, and domains explicitly

Every derivation unit MUST list the assumptions and conventions that affect its validity. Relevant items include, where applicable:

- framework/interpretation and scope;
- physical regime and approximations;
- metric signature and index conventions;
- coordinate chart, foliation, gauge, boundary conditions, topology, or orientation;
- operator ordering and quantization choices;
- positivity, nonzero denominators, square-root branches, differentiability, integrability, and regularity;
- normalization, fiducial volume, measure, units, and constants;
- symmetry-reduction assumptions and whether a reduced theory is derived from or separately postulated relative to the full theory.

Each assumption MUST be marked as source-stated, sourced externally, inferred, or unresolved. Structured conditions MUST be represented as structured conditions if the system supports them; prose can explain an assumption but must not replace a checkable condition where one is required.

Assumptions must be propagated through derivations according to the project's semantic rules. They must not be silently added, dropped, strengthened, or weakened.

A condition introduced accidentally by an implementation artifact is not automatically a paper assumption. For example, if differentiating a polynomial generates a `Pow(x, 0)` intermediate and causes an unnecessary `x != 0` requirement, the specification should identify that as a derivative-semantics defect, not retroactively restrict the polynomial's domain.

Likewise, a reduction such as \(\sqrt{a^2}\to a\) requires an explicit sign assumption. A spatial integration or fiducial-volume factor must not be treated as an evaluated constant unless the source and representation justify that treatment.

## 8. Specify derivation routes and substeps as computational contracts

For each derivation unit, write an ordered dependency graph or equivalent structured trace.

Each step MUST identify:

1. step ID and position;
2. mathematical operation or source inference;
3. exact premise quantities/artifacts;
4. required input shapes and types;
5. assumptions and conditions;
6. expected intermediate/result quantity;
7. expected output type and relevant identity invariants;
8. verification method and independent expected result;
9. known limitations or reason for non-execution.

Classify each graph edge. Recommended meanings:

- **PREMISE:** a source-declared or externally established input;
- **ALGEBRAIC/COMPUTATIONAL:** an operation that should calculate a result;
- **ASSUMPTION:** a condition accepted as an input, not derived in this step;
- **CITED DEPENDENCY:** a result whose justification is in another source;
- **ASSERTED/UNSHOWN:** claimed or sketched by the paper but not reproduced in the available derivation;
- **INFERRED BY SPECIFICATION:** reconstructed by the specification author and requiring explicit review;
- **UNVERIFIED TRANSITION:** required mathematics that current evidence does not establish.

The project's schema MAY encode these using different names, but the distinction MUST remain clear.

For each step, state which operations actually run. A stored equation, a constructed expression tree, a successful type check, a successful replay, and a proven algebraic equality are separate claims.

If a paper presents two derivation routes, specify both separately, including their shared premises and target. Do not merge them simply because their conclusions coincide.

## 9. Audit the mathematical capability baseline

Before proposing new implementation, map each required substep against the actual code—not against what would be convenient if the system had a feature.

For each operation, record:

- existing representation and artifact family;
- current executable rule/algorithm and its exact scope;
- encoding/decoding and canonicalization support;
- assumption and domain-condition behavior;
- provenance and replay support;
- independently tested edge cases;
- missing semantics or infrastructure;
- the smallest honest unblocker and broader reusable options.

Use the project's current status vocabulary, with its operational definitions. A practical default mapping is:

- **SUPPORTED:** the claim within the declared scope has the required executed/replayed evidence or the exact validation guarantee the specification requires;
- **PARTIALLY SUPPORTED:** important substeps are supported but a material transition remains unexecuted or unverified;
- **REPRESENTATIONAL ONLY:** the mathematical object is faithfully encoded/validated, but no execution path establishes the requested transformation;
- **UNSUPPORTED:** required semantics or representation are missing, so even a faithful representation of the target step is not currently established.

Follow the authoritative definition in the relevant repository/specification if it differs. Do not classify an entire derivation from its easiest substep.

A missing capability is a valid implementation target. If the physics genuinely requires total derivatives, functional derivatives, tensor calculus, covariant derivatives, curvature, integration, constraint algebra, or another significant foundation, specify and implement that capability when authorized. Do not shrink the scientific claim merely to fit today's operation registry.

Conversely, do not implement an enormous generic subsystem when a well-defined reusable capability suffices. Decide from the dependency graph and expected cross-domain reuse, not habit or a fixed operation count.

## 10. Independent expected results are a required deliverable

Before implementation begins, derive test expectations independently from the source equations and mathematical definitions.

For each executable substep, include as appropriate:

- hand-calculated exact rational/scalar examples;
- symbolic derivatives and operator-composition examples;
- dimensional and index-structure checks;
- boundary/domain cases;
- special cases that distinguish plausible but wrong algorithms;
- cross-representation consistency checks;
- independent numerical spot checks for supported expressions, when useful;
- replay/provenance-mutation tests.

The test oracle MUST NOT be generated solely by the implementation under test. A test comment saying “hand-derived” does not establish independence if it simply copies the implementation's assumptions and expression-building sequence.

### Use layered evidence

Prefer more than one kind of check for important mathematics:

1. **Source-level oracle:** independently establish the formula and operation ordering from the paper.
2. **Structural contract:** check operand identities, typed shapes, ordering, and assumptions.
3. **Local golden:** compare a supported operation against a hand-derived exact result.
4. **End-to-end derivation:** establish that the intermediates compose into the source-defined target.
5. **Anti-golden/mutation test:** deliberately introduce a plausible wrong formula, swapped input, wrong variable identity, dropped condition, or substituted result and show that the relevant gate rejects it.
6. **Replay:** verify that the recorded trace reproduces outputs from the declared premises.
7. **Manifest/protocol check:** confirm that source mapping and claimed verification match the actual implementation.

Not every step needs every test technique, but every important result MUST have an independent route to checking correctness.

### Exact structure versus mathematical equivalence

Exact AST or byte equality is appropriate when canonical representation is itself part of the contract. It is not a universal test for algebraic equality.

If the engine does not implement distributivity, term collection, rational-function normalization, or another equivalence procedure:

- test the structural form the engine actually promises;
- independently derive the expected mathematics;
- use supported rewrite/equivalence mechanisms if available;
- optionally use independent exact numerical checks at valid sample points;
- document the remaining equivalence gap accurately.

Do not claim the engine proved an equivalence that only a human reviewer checked. Also do not expand general algebra opportunistically just to force one test into a preferred normal form; track that work as a reusable capability and prioritize it deliberately.

### Negative tests must challenge the actual risk

For each critical invariant, include an anti-golden or mutation test that exercises the failure the invariant is intended to prevent. Examples include:

- deriving pressure from density instead of potential;
- confusing two variables with identical display names but different identities;
- applying \((a\partial_a)^2\) as \(a^2\partial_a^2\);
- losing a required nonzero or positivity condition;
- accepting a tensor contraction across distinct index spaces;
- pairing a result with evidence from different inputs;
- omitting a minted derivation from the verifier's evidence set.

A negative test that only asserts “an error occurred” is weak if the error could be unrelated. Verify the intended rejection boundary.

## 11. Provenance and authoritative verification are part of the specification

For every derived artifact, specify what connects the result to its inputs and operations.

Where relevant, require binding of:

- source quantity and equation mapping;
- exact input artifact references and mathematical identities;
- operation identity, version, parameters, and conditions;
- intermediate steps and their dependency links;
- ordered tensor components, free indices, and output identity;
- declared assumptions and their applicability;
- final result reference;
- evidence status actually established.

Keep validation layers separate:

- **Structural validation** checks local grammar and invariants without hidden store I/O.
- **Pure decoding** validates canonical representations; it does not prove computation.
- **Store-aware checking** resolves the actual artifacts, checks families and linkage, verifies conditions, replays the computation, and recomputes required premises/results where needed.
- **Consumer gating** ensures downstream work cannot treat merely structural data as verified provenance.

The authoritative `Verify()`/protocol entry point MUST cover every derivation and evidence item it claims to verify, including items stored outside the main derivation slice. Enumerate or otherwise make the evidence set explicit enough that omissions are testable.

A relation record that merely states `Q = a^3 rho_Q` does not prove that its Q and density references satisfy the relation. The verifier must check the actual stored operands and relation.

Never allow an expected-premise recomputation to silently repair missing evidence before replay or availability checks. If the content-addressed store cannot roll back intermediate artifacts, document that property and ensure failure cannot produce an accepted certificate or valid-looking partial result.

Describe provenance guarantees accurately. Binding the component content that determines a result is not necessarily the same as binding a unique higher-level model identity. State precisely what is and is not guaranteed.

## 12. The package-level derivation contract

Unless an existing repository standard specifies an equivalent structure, a paper-specific derivation package SHOULD contain:

```text
<derivation-package>/
  README.md
  manifest.json
  manifest_test.go
  protocol/
    derivation_test.go
  SPECIFICATION.md       # or a canonical project-level spec referenced here
```

Use the actual repository's established path convention; do not create duplicate specifications without a reason.

### README.md — human orientation

Explain:

- scientific question and target claim;
- framework, interpretation, regime, and premises;
- assumptions and conventions;
- overview of the derivation routes;
- mathematical and implementation boundaries;
- known limitations, anomalies, and unresolved questions.

README prose is explanatory, not the machine-checkable mathematical result.

### manifest.json — machine-readable derivation metadata

Use strict parsing and the schema already established by the project. At minimum, capture equivalent fields for:

- stable derivation ID and framework/package ID;
- source identifier, version, digest, and scope note;
- entry points and verification path;
- proposed core idea and derivation-unit IDs;
- source-to-quantity mappings;
- premises, assumptions, conventions, regime, and dependencies;
- expected targets and evidence status;
- limitations, anomalies, unresolved claims, and empirical status when applicable.

Reject unknown fields if that is the prevailing schema contract. Add a new field deliberately rather than relying on silently ignored metadata. Tests MUST ensure manifest metadata agrees with implementation constants and actual verification behavior.

For source mapping, use an existing compatible schema or the equivalent of:

```json
{
  "source_mappings": [
    {
      "source_equation": "(76)",
      "source_quantity": "rho_Q",
      "defining_relation": "rho_Q = Q / a^3",
      "fixture_ref": "RhoQRef",
      "dependencies": ["QRef", "a"]
    }
  ]
}
```

This example illustrates required information, not a command to add this exact JSON shape regardless of the package schema. Ensure the declared dependency graph agrees with the mathematical direction and actual code.

### manifest_test.go — package contract validation

Where the package convention provides this test, validate:

- manifest exists, parses strictly, and contains required fields;
- source identity and digest are present and consistent;
- entry points and required protocol files exist;
- unknown fields are rejected as required by the schema;
- source mappings, fixture identifiers, and implementation metadata agree;
- verification and limitation statements do not overclaim.

### protocol/derivation_test.go — executable evidence

The protocol test MUST exercise actual constructors and mathematical operations. It SHOULD include intermediate golden states, independent mathematical targets, critical anti-goldens, assumption/domain checks, provenance checks, replay, and `Verify()` coverage. It MUST NOT simply assert that a fixture returned without error.

## 13. Separate mathematical status from physical and empirical status

A formally valid expression or derivation does not, by itself, establish physical truth.

The specification MUST separately state:

- the status of the mathematical representation;
- the status of each executed/replayed derivation step;
- the physical framework and interpretation in which the result is used;
- what follows only conditionally from declared assumptions;
- whether any empirical evidence is actually present.

Do not inherit corpus status, provenance status, or physical interpretation implicitly from a related package. Do not promote a hypothesis or a paper's claim to an established fact merely because a computation succeeds. Respect the current repository's authoritative status model where one exists; do not blindly copy another project's enums.

For scientific consequences such as conservation, equivalence, or recovery of a known limit, specify the exact premises and mathematical operations that establish the consequence. If the argument relies on a theorem not represented or verified, record it as a dependency or unverified transition.

## 14. Specification-readiness gate

A derivation specification is ready for implementation only when the following are sufficiently concrete for an independent agent to implement without guessing about the science.

### Required before implementation

- [ ] Authoritative source/version/scope identified.
- [ ] Core-idea inventory and derivation-unit mapping proposed, with unresolved adjudication points labelled.
- [ ] Each important source quantity mapped to its definition, representation, and dependencies.
- [ ] Premises, conventions, regimes, and mathematical domain conditions listed.
- [ ] Derivation routes and intermediate steps ordered and typed.
- [ ] Operation semantics specified sufficiently to distinguish correct computation from plausible wrong forms.
- [ ] Existing capability support checked against actual code and tests.
- [ ] Missing mathematical foundations identified, with implementation-sized work items and dependencies.
- [ ] Independent positive expectations and critical anti-goldens derived before coding.
- [ ] Provenance, replay, consumer checks, and verifier coverage specified.
- [ ] Manifest and protocol updates defined.
- [ ] Unsupported claims, risks, and acceptance criteria explicit.
- [ ] Repository boundaries and authorized operations stated.

The gate MUST block an affected derivation step when a material mathematical ambiguity prevents a sound implementation. It MUST NOT become an excuse for endless planning when the specification is adequate for a useful implementation increment. Unresolved items that do not block a coherent substep should remain explicit limitations while implementation proceeds on the verified dependency chain.

### Required implementation handoff

The handoff MUST state:

1. What to implement and why the mathematics requires it.
2. Which package owns each semantic operation and why.
3. Ordered workstreams and their dependencies.
4. Exact acceptance criteria and independent goldens.
5. Required negative/adversarial tests.
6. Evidence the verifier must cover.
7. Completion report format and repository-operation restrictions.

A good specification allows broad implementation where justified, but leaves little ambiguity about what scientific claim the code is supposed to establish.

## 15. Recommended implementation flow after specification readiness

Once the specification is implementation-ready and authorized:

1. Record live branches, commits, working trees, and baseline test results.
2. Implement a coherent, reusable mathematical capability—not merely a paper-specific output.
3. Add structural representation, execution, provenance, replay, and package integration appropriate to the claim.
4. Add independent mathematical tests and adversarial tests as the capability is built.
5. Update source mappings, manifest limitations, and capability classifications alongside code.
6. Run focused tests, full suites, static checks, formatting, and applicable purity checks.
7. Audit the diff and perform an independent adversarial review at meaningful milestones.
8. Correct material findings before declaring completion.
9. Commit locally where authorized and report exact commits and test results.
10. Recommend exactly one next high-impact reusable capability based on verified downstream dependencies.

A missing capability is an implementation opportunity, not automatically a reason to weaken a derivation. A useful bounded increment is acceptable when it genuinely advances a connected derivation and its limits are precise. Do not claim that partial work verifies the full paper.

## 16. Adversarial review requirements for a completed derivation slice

The review MUST ask two independent questions:

1. Did the engine perform the declared mathematical operation correctly?
2. Was that the correct mathematical operation and quantity according to the paper?

Review at least:

- source equation and quantity identity;
- coefficients, signs, exponents, operator ordering, indices, and measures;
- variable identity and expression dependencies;
- assumptions, domain conditions, and their propagation;
- provenance completeness and replay validity;
- oracle independence and anti-golden effectiveness;
- complete `Verify()` evidence coverage;
- manifest/implementation/test consistency;
- resource bounds and failure behavior;
- unsupported claims and status inflation;
- architecture, dependency, and repository boundaries.

A green full test suite does not override a demonstrated mathematical source-fidelity defect. A clean manifest does not prove code behavior. A passing replay proves only what the recorded operations and premises actually establish.

The review should recommend one verdict from the project's established review vocabulary and supply concrete evidence. It should not reopen settled design choices for stylistic reasons; it must reopen them when a real contradiction, invalid guarantee, or scientific defect is shown.

## 17. How to maintain this standard

This document is normative and changes deliberately. Unlike `GUIDELINES.md`, which is an additive casebook of lessons learned, this file defines the baseline method that new derivation specifications MUST follow.

When a real failure reveals a missing requirement:

1. Record the concrete failure and root cause in `GUIDELINES.md`.
2. Determine whether it reveals a general specification obligation rather than a one-off defect.
3. Amend this document with the durable normative requirement if warranted.
4. Add or update a template, manifest validation, protocol test, or check where practical so the requirement is enforceable.
5. Explain the change and retain historical context where it matters.

Do not duplicate every lesson verbatim. `GUIDELINES.md` preserves the incident and its lesson; this document incorporates the general rule that prevents recurrence.

## 18. Minimum derivation specification template

Use this outline for a new derivation specification, adapting headings only when the information remains explicit elsewhere.

```markdown
# Derivation Specification — <paper / target result>

## 1. Status and decision boundary
- Specification version/status:
- Requested scope:
- Implementation authorization/boundary:
- Open decisions and blockers:

## 2. Source identity and fidelity
- Title/authors/version/identifier:
- Canonical source path and digest:
- Exact sections/pages/equations in scope:
- Version differences or source ambiguities:
- Conventions and notation confirmed from rendered source:

## 3. Scientific question and regime
- Question/target claim:
- Framework and interpretation:
- Regime/domain/symmetry reduction:
- What the paper establishes vs cites/asserts:

## 4. Proposed core-idea inventory
- ID, claim, rationale, source, dependencies, status:
- Provisional count and unresolved adjudication:

## 5. Quantity/source map
For every significant quantity:
- stable ID and source equation;
- mathematical meaning and defining relation;
- typed representation / target artifact;
- exact source dependencies and assumptions;
- expected result and independent derivation;
- verification status and limitation.

## 6. Derivation units and dependency graph
- Unit IDs and core-idea mapping:
- ordered premises and intermediate steps:
- separate alternative derivation routes:
- typed edges: premise, computation, assumption, cited,
  asserted/unshown, inferred, or unverified:

## 7. Assumptions, conventions, and domain
- framework, regime, signature, coordinates, gauge/foliation:
- boundary conditions, normalization, units, operator ordering:
- nonzero/positive/smooth/integrable conditions:
- source-stated vs inferred assumptions:

## 8. Mathematical semantic contracts
For each operation:
- accepted input families and shapes:
- exact semantics and output invariants:
- index/dimension/domain rules:
- conditions and failure behavior:
- canonical representation and composition:

## 9. Current capability map
- Existing implementation and tests:
- Supported/partial/representational/unsupported substeps:
- Missing reusable foundations:
- Proposed package placement and rationale:

## 10. Independent expected results
- Hand-derived positive cases:
- Plausible wrong forms and anti-goldens:
- Boundary/domain cases:
- Exact structural checks vs equivalence checks:
- Independent numerical checks, if suitable:

## 11. Provenance and verification contract
- Inputs, operation parameters, conditions, intermediate refs:
- Result identity and ordered-component invariants:
- Replay and store-aware validation:
- Complete evidence set covered by Verify():

## 12. Package and manifest contract
- README/manifest/validation/protocol paths:
- Source mappings and metadata consistency:
- Limitations/anomalies/falsification/empirical status:

## 13. Acceptance criteria
- Required tests and failure boundaries:
- Required verification claims:
- Explicit exclusions and known gaps:

## 14. Implementation sequence
- Ordered, dependency-aware increments:
- Baseline/test/review requirements:
- Repository boundaries and commit/publish policy:

## 15. Remaining questions and next capability
- Material blockers:
- Non-blocking limitations:
- One next high-impact reusable capability:
```

## Final principle

**A derivation specification is the contract between the scientific source and the implementation.** It must preserve what each quantity means, why each transformation is valid, what assumptions it relies on, what results are independently expected, and exactly what the resulting evidence does and does not prove.

Build what the physics requires. Specify it so another agent can implement it without guessing about the mathematics. Verify it so the implementation cannot merely agree with its own mistaken interpretation. Keep every scientific claim bounded by the evidence actually established.
