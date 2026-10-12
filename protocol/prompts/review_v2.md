# PhysTrace Idea-First Adversarial-Review Prompt — v2 (canonical for M1)

Run in a FRESH Gemini Web conversation with the SAME paper PDF plus the
candidate v2 JSON. No inherited extraction context. You review ideas AND
derivations — not equation counts.

---

The candidate claims X core ideas with exactly X derivation units. Break
that claim: missing ideas, false splits/merges, missing dependencies,
wrong exclusions, transcription errors, false SUPPORTED labels.

## Challenges (all mandatory)

1. **Missing idea.** Is any independently meaningful contribution absent
   (e.g. scalar-field extension, heuristic with standalone content)?
2. **Artificial split.** Is one coherent argument spread over ≥2 units
   (e.g. transformation result vs mass identification of the same proof)?
   Propose the merge.
3. **Unjustified merge.** Are genuinely distinct contributions folded
   into one unit? Propose the split.
4. **Intra-derivation gaps.** For each unit: missing equations,
   assumptions, lemmas, boundary arguments, or load-bearing links needed
   to reconstruct the reasoning.
5. **Wrong exclusions.** Is needed appendix material excluded, or an
   alternative route misranked as required (or vice versa)?
6. **Source fidelity.** Spot-check transcriptions symbol-by-symbol;
   verify source regions, labels, passages, assumption sourcing
   (author-stated vs inferred).
7. **Support honesty.** Challenge every `SUPPORTED`: does the step truly
   fit the 7-rule algebra, or hide integrals/tensors/geometry? Confirm
   every `UNSUPPORTED` reason. Flag anything dropped to look complete.

## Output contract

Exactly one JSON document, no prose outside it:

```json
{
  "schema_version": "review_v2",
  "prompt_version": "review_v2",
  "candidate_digest": "<sha256 of exact candidate bytes reviewed>",
  "verdict": "ACCEPT|ACCEPT_WITH_CORRECTIONS|REJECT",
  "idea_findings": [{
    "finding_id": "f_001",
    "kind": "missing_idea|bad_split|bad_merge|scope",
    "description": "precise statement with paper locations",
    "proposed_change": "add/merge/split/scope text or null"
  }],
  "derivation_findings": [{
    "finding_id": "f_010",
    "severity": "omission|mistranscription|assumption|dependency|support|coverage",
    "unit_id_or_page": "D1 or p.5",
    "description": "",
    "proposed_correction": "concrete fix or null"
  }],
  "unsupported_confirmations": ["it_005"],
  "notes": "regions you could not verify; low counts prove nothing"
}
```

You flag; the human adjudicates. Findings never auto-modify anything.
