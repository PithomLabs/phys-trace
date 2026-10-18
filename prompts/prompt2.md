# Task: Solve a Physics Problem Using Phys Infrastructure

Given the supplied textbook physics problem, solve it and produce a reproducible **Type-B benchmark-solution package** under `phys-demo`, following `phys-demo/PROTOCOL.md`.

## Required instructions

1. Read `phys-system/AGENTS.md`, `phys-demo/PROTOCOL.md`, `TOE-v4/plan2/adr_package_types.md`, the applicable repository `AGENTS.md` files, and the relevant Phys infrastructure API documentation and tests.

2. Inspect the problem statement and identify its source, requested result, parameters, units or unit conventions, assumptions, physical regime, governing equations and required derivation steps.

3. Before implementing, audit the capabilities that actually exist in PhysArtifact, PhysMath and the relevant `phys-*` libraries. Reuse existing APIs to represent supported mathematical expressions, equations, exact quantities, derivation steps and references where appropriate.

4. Follow the Type-B manifest schema and package requirements in `phys-demo/PROTOCOL.md`. Record meaningful infrastructure reuse in `infrastructure_use`, including the exact API, executing code path and test. If a suitable API does not exist or cannot faithfully represent the problem, provide the required justified `deviations` record. **Imports alone do not constitute infrastructure adoption.**

5. Keep numerical computation outside PhysMath and PhysLib. If the problem requires numerical calculations, perform them in solver-local code and document their inputs, methods, conventions, outputs and independent checks. Never claim that Phys infrastructure performed a calculation it did not execute.

6. Do not force the problem into an unrelated physical framework. Do not invent APIs, artifact families or executable rules. Distinguish supported operations from partially supported, represented-only, deferred and unavailable capabilities.

7. Add meaningful tests for the actual entry points, supported library operations, results, units where applicable, manifest conformance, and justified deviations. Compute the package's conformance status from the protocol requirements; unresolved mandatory gaps must prevent a conformant status.

8. If any contained artifact claims formalization under Representation Specification v1.2 §1.4, follow the Type-A requirements and the mixed-purpose rules in `phys-demo/PROTOCOL.md` §3. Do not use a Type-B label to evade the formalization contract.

## Definition of done

Deliver the solution, its mathematical representation through applicable Phys libraries, independently checked results, required manifest and tests, actual infrastructure APIs reused, solver-local computations, and explicit limitations. Make it possible for another user or agent to reproduce the work and determine precisely what Phys represented or verified versus what was calculated outside it.