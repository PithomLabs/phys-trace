# M1 Exit Gate — Strict Checklist (from plan1.5.md §12)

M1 passes only when EVERY box is checked with evidence. A documented
trial run with open boxes is useful progress, not a pass. Failed boxes
are recorded with criterion + evidence + corrective action; the gate is
never lowered.

- [ ] 1. Adjudicated revision exists: frozen revN JSON + canonical digest
  + parent chain back to rev0 candidate. Evidence: `________________`
- [ ] 2. Full page ledger: every in-scope page has an explicit status; no
  `UNASSESSED`/`UNREADABLE` in scope unless gate stays PARTIAL with
  reason. Evidence: `________________`
- [ ] 3. Triple classification: every in-scope item has capture +
  compiler (+reason) + test status; no silent omissions.
  Evidence: `________________`
- [ ] 4. Genuine unsupported item: ≥1 item preserved as `UNSUPPORTED`
  with correct reason (expected: integration / vector identities /
  tensor operations in the EM-box argument). Evidence: `________________`
- [ ] 5. All findings dispositioned: every `f_NNN` accepted / rejected /
  deferred with rationale; deferred questions stay explicit.
  Evidence: `________________`
- [ ] 6. Cross-session proof: extraction output and review input conform
  to the same contract version; fresh session reviewed paper +
  candidate without extraction context. Evidence: `________________`
- [ ] 7. Evaluation recorded: omissions, mistranscriptions, ambiguities,
  unsupported ops, disagreements, human corrections — all logged.
  Evidence: `________________`

## Verdict semantics reminder

- `COMPLETE` representation = accounted-for within scope, not correctness.
- Execution `COMPLETE` only with zero relevant omitted/ambiguous/
  unassessed/unsupported items; else `PARTIAL`.
- `test_execution` = `NOT_RUN` until external `GOWORK=off` replay runs.
- Passing checks never promote partial coverage; findings never auto-edit.

## On pass

Approved revision + lineage becomes the Stage B compiler input fixture.
M2 begins only on M1 pass or your explicit gate change.
