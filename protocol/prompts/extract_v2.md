# PhysTrace Idea-First Extraction Prompt — v2 (canonical for M1)

Paste into Gemini Web with the paper PDF. M1 paper: Perez & Ribisi,
arXiv:2101.11923. Starting proposal: K1 single Maxwell derivation;
K2 scalar extension as candidate. You propose X — the human adjudicates it.

---

You produce a scientific-REPRESENTATION record, not an interpretation,
summary, or essay. Do not judge the paper's physics. Do not add content
the paper does not state. Do not split one argument into many units to
look thorough, and do not merge distinct contributions to look simple.

## Step 1 — propose core ideas FIRST

Inspect the paper in declared scope (all pages incl. figures, tables,
footnotes, appendices). Propose the SMALLEST defensible set of distinct
core scientific ideas preserving its essential contribution. For each:

- `claim`: the scientific contribution in one sentence.
- `distinctness`: why it is independently meaningful, not a step/lemma/
  equation of another idea.
- `source_regions`: sections/pages/equation ranges supporting it.
- `assumptions`: author-stated vs inferred (flag inferred explicitly).
- Exactly one derivation unit containing its necessary equations,
  intermediate steps, lemmas, and load-bearing dependencies.

Identify supporting material, alternative proof routes, and exclusions
WITHOUT elevating them to ideas. Name current-compiler limits honestly
without dropping or falsifying mathematics.

## Step 2 — output contract

Exactly one JSON document, no prose outside it, no Go code. `null`/empty
where unknown — never invent.

```json
{
  "schema_version": "extraction_v2",
  "paper": {"arxiv_id": "2101.11923", "arxiv_version": null,
    "doi": null, "source_url": null, "pdf_sha256_hex": null,
    "page_count": 0,
    "declared_scope": {"included_pages": [[1, 0]],
      "excluded": [{"range": [0, 0], "reason": ""}]}},
  "source_coverage": [{"page": 1,
    "status": "INSPECTED_WITH_CONTENT",
    "note": "optional region-level remark"}],
  "core_ideas": [{
    "idea_id": "K1",
    "claim": "",
    "distinctness": "",
    "source_regions": ["Sec. IV, pp.5-8"],
    "scope_assumptions": ["it_a1"],
    "derivation_unit_id": "D1"
  }],
  "derivation_units": [{
    "unit_id": "D1",
    "idea_id": "K1",
    "intended_conclusion": "",
    "assumptions": ["it_a1"],
    "definitions": ["it_d1"],
    "steps": [{"step_id": "s1", "items": ["it_001"],
      "depends_on": [], "required": true,
      "status": "SUPPORTED|UNSUPPORTED|AMBIGUOUS|UNASSESSED",
      "reason": null}],
    "source_refs": [{"items": ["it_001"], "pages": [5],
      "labels": ["(24)"], "passage": "short quote"}],
    "execution_obligations": ["s4"],
    "limitations": ["integral steps unexecutable under 7-rule baseline"]
  }],
  "items": [{
    "item_id": "it_001", "unit_id": "D1",
    "kind": "equation|definition|assumption|convention|condition|claim|lemma",
    "label": "(24) or unlabeled-N", "latex": null, "pages": [7],
    "passage": "short verbatim quote", "inferred": false,
    "capture_status": "CAPTURED|AMBIGUOUS",
    "compiler_status": "SUPPORTED|UNSUPPORTED|NOT_APPLICABLE",
    "test_status": "NOT_RUN",
    "math_shape": "add|mul|pow|substitute|apply|differentiate|select_branch|canon|other_unsupported",
    "unsupported_reason": null, "confidence": 0.0
  }],
  "relations": [{"from": "it_001", "to": "it_002",
    "rel": "derives|assumes|defines|constrains", "required": true}],
  "ambiguities": [{"item_id": "it_001", "note": ""}],
  "model_provenance": {"model": "<name/version>",
    "prompt_version": "extract_v2", "extraction_version": "rev0"}
}
```

## Rules

- Idea count X = unit count. Every accepted idea has exactly one unit;
  every unit belongs to exactly one idea. No unit exists because an
  equation was extracted.
- `math_shape: other_unsupported` + reason for vector identities,
  integration, tensor ops — never force into a supported shape.
- `test_status` is always `NOT_RUN` at extraction.
- `pdf_sha256_hex` and unverifiable counts: `null` — the human fills them.
- A complete ledger with a thin inventory is a finding about the paper,
  not a failure: report it, do not pad it.
