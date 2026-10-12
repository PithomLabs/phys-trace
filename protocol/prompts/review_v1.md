# PhysTrace Independent Adversarial-Review Prompt — v1

Run this in a FRESH Gemini Web conversation. Supply TWO inputs: the same
paper PDF (identical bytes as extraction) and the candidate JSON produced
by `extract_v1`. This session must not rely on the extraction conversation.

---

You are an adversarial reviewer of a scientific-representation record, not
of the paper's physics. The candidate JSON claims to inventory a paper's
formal content. Your job is to break that claim: find what is missing,
mistranscribed, misclassified, or inconsistent.

You have the original paper. The candidate is guilty until proven
innocent on each point below. Do not summarize the paper. Do not offer
physical interpretations.

## Review procedure

1. **Page sweep.** Walk every in-scope page (body, figures, tables,
   footnotes, appendices). For each numbered equation, definition,
   assumption, convention, condition, and conclusion, check it appears in
   the candidate with correct `label`, `pages`, and `passage`. List every
   in-scope formal item the candidate OMITS, with page + label.
2. **Transcription check.** For each captured equation, compare `latex`
   against the paper symbol-by-symbol (signs, indices, factors of c,
   derivatives vs integrals, boundary terms). List mistranscriptions.
3. **Assumption audit.** List every assumption the paper states that the
   candidate misses. For each candidate item with `inferred:true`, judge
   whether the inference is justified or smuggled in — flag unjustified
   ones. Flag any author-stated assumption mislabeled as inferred.
4. **Dependency check.** Verify each `relations` entry against the paper's
   argument order. Flag reversed, missing, or spurious links, and any
   `required` flag you judge wrong.
5. **Support-classification check.** For each item marked `SUPPORTED`,
   challenge it: does the mathematics really fit plain algebra
   (add/mul/pow/substitute/apply/differentiate/select-branch/canonicalize),
   or does it hide vector identities, integration, or tensor operations
   that require `UNSUPPORTED`? For each `UNSUPPORTED`, confirm the reason
   code fits.
6. **Coverage-ledger check.** Confirm every in-scope page has a coverage
   entry with a plausible status. Flag pages marked `INSPECTED_*` whose
   content suggests they were not truly inspected.

## Output contract

Output exactly one JSON document, no prose outside it:

```json
{
  "schema_version": "review_v1",
  "prompt_version": "review_v1",
  "candidate_digest": "<sha256 of the candidate JSON you reviewed, computed over its exact bytes>",
  "verdict": "ACCEPT|ACCEPT_WITH_CORRECTIONS|REJECT",
  "findings": [{
    "finding_id": "f_001",
    "severity": "omission|mistranscription|assumption|dependency|support|coverage",
    "item_id_or_page": "it_003 or p.5",
    "description": "precise statement with paper location",
    "proposed_correction": "concrete fix or null if only flagging"
  }],
  "unsupported_confirmations": ["it_007"],
  "notes": "anything the human must decide that you could not"
}
```

Severity guide: `omission` = in-scope content absent from candidate;
`mistranscription` = wrong symbols/content; `assumption` = missing or
mishandled assumption; `dependency` = wrong link; `support` = wrong
SUPPORTED/UNSUPPORTED call; `coverage` = ledger gap.

Your findings are flags for human adjudication. You do not modify the
representation. A low finding count is not proof of completeness — state
in `notes` any region you could not verify.
