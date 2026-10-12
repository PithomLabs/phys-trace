# M1 Exit Gate v2 — Strict Checklist (spec v1.2.3 §12)

Every box checked with evidence, or M1 has not passed. Trial ≠ pass.
Failures recorded (criterion + evidence + correction); gate never lowered.

- [ ] 1. Versioned approved revision: frozen JSON + canonical digest +
  parent chain to rev0. Evidence: `____`
- [ ] 2. Bijection: X accepted ideas, X units, each idea ↔ one unit;
  no equation-born units, no idea without a unit. Evidence: `____`
- [ ] 3. Full ledger: every in-scope page statused + region notes;
  exclusions explicit. Evidence: `____`
- [ ] 4. Triple classification incl. intra-unit portions; no silent
  omissions. Evidence: `____`
- [ ] 5. ≥1 genuine unsupported item preserved with reason.
  Evidence: `____`
- [ ] 6. Every finding (incl. split/merge/scope proposals)
  accepted/rejected/deferred with rationale. Evidence: `____`
- [ ] 7. Cross-session proof: same contract version both sessions;
  fresh review without extraction context. Evidence: `____`
- [ ] 8. Evaluation logged: omissions, errors, ambiguities, unsupported
  ops, disagreements, corrections. Evidence: `____`

Reminders: COMPLETE = accounted-for in scope. Execution COMPLETE only
with zero relevant omitted/ambiguous/unassessed/unsupported required
content. test = NOT_RUN until external GOWORK=off replay. Passing
checks never promote partial coverage; findings never auto-edit.
On pass, approved revision + lineage is the Stage B fixture. M2 needs
M1 pass or explicit gate change.
