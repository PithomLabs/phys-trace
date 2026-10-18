# Task: Represent a Physics-Paper Derivation Using Phys Infrastructure

Given the supplied physics paper and target derivation, produce a reproducible, conformant **Type-A formalization package** using the existing Phys infrastructure.

## Required instructions

1. Read `phys-system/AGENTS.md`, the applicable repository and package `AGENTS.md` files, `phys-trace/DERIVATION_SPEC_GUIDELINES.md`, Architecture v1.5, Representation Specification v1.2, and the relevant sections of `phys-trace/GUIDELINES.md`.

2. Inspect the actual public APIs, implementation and tests of PhysArtifact, PhysMath and the relevant `phys-*` library before writing code. Use existing capabilities wherever they faithfully represent the mathematics and physical assumptions.

3. Follow Representation Specification v1.2 §1.4 exactly: use the canonical `AGENTS.md`, `manifest.json`, and `protocol/derivation_test.go` package structure and the exact GR-example manifest schema and loading convention.

4. Pin the source, identify the relevant equations, define the assumptions and conditions, and represent the derivation using existing PhysArtifact, PhysMath and appropriate `phys-*` artifacts. Record supported derivation steps and execute and replay them through the existing protocol.

5. Consult the live PhysMath rule registry to determine what actually executes. **Never invent APIs, rule IDs, artifact families or mathematical capabilities.** Unsupported steps must remain explicitly represented-only or be reported as a `PROTOCOL GAP` under the governing specification. Do not silently extend the manifest schema.

6. Numerical computation is outside PhysMath and PhysLib scope. Do not introduce numerical solvers, dimensional-semantics systems or new executable rules. Any excluded work must remain outside the representation layer.

7. Run meaningful tests that exercise the actual formalization and its entry points. Record exactly what was represented, structurally validated, executed and replayed, and what remains unverified or open. Do not equate a passing test or replay with scientific truth.

## Definition of done

Deliver the conformant package, reproducible test results, exact APIs and operations reused, unsupported capabilities and limitations, and a clear distinction between the formalized mathematical claim and what has not been established.