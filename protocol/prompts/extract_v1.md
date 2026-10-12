# PhysTrace Candidate-Extraction Prompt — v1

Paste this prompt into Gemini Web together with the uploaded paper PDF.
Target paper (M1): Perez & Ribisi, arXiv:2101.11923.

---

You are producing a candidate scientific-REPRESENTATION record, not an
interpretation, summary, or essay. Do not judge whether the paper's physics
is correct. Do not add equations, assumptions, or conclusions the paper does
not state.

## Output contract

Output exactly one JSON document conforming to the schema below. No prose
outside the JSON. No Go code. If a field is unknown, use `null` — never
invent values.

```json
{
  "schema_version": "extraction_v1",
  "paper": {
    "arxiv_id": "2101.11923",
    "arxiv_version": null,
    "doi": null,
    "source_url": null,
    "pdf_sha256_hex": null,
    "page_count": 0,
    "declared_scope": {
      "included_pages": [[1, 0]],
      "excluded": [{"range": [0, 0], "reason": "references"}]
    }
  },
  "source_coverage": [{"page": 1, "status": "INSPECTED_WITH_CONTENT"}],
  "items": [{
    "item_id": "it_001",
    "kind": "equation|definition|assumption|convention|condition|claim|dependency",
    "label": "(3.2) or unlabeled-N",
    "latex": null,
    "pages": [1],
    "passage": "short verbatim quote",
    "inferred": false,
    "capture_status": "CAPTURED|AMBIGUOUS",
    "compiler_status": "SUPPORTED|UNSUPPORTED|NOT_APPLICABLE",
    "test_status": "NOT_RUN",
    "math_shape": "add|mul|pow|substitute|apply|differentiate|select_branch|canon|other_unsupported",
    "unsupported_reason": "trig|integration|general_solve|sparse_tensor|unbound_symbol|null",
    "confidence": 0.0
  }],
  "relations": [{"from": "it_001", "to": "it_002",
    "rel": "derives|assumes|defines|constrains", "required": true}],
  "ambiguities": [{"item_id": "it_001", "note": "missing domain"}],
  "model_provenance": {"model": "<model name/version>",
    "prompt_version": "extract_v1", "extraction_version": "rev0"}
}
```

## Field rules

- `declared_scope`: list every page you inspected in `included_pages`;
  record skipped regions (references, appendices out of scope) in `excluded`
  with reasons. Never silently drop pages.
- `source_coverage`: one entry per in-scope page. Status is one of
  `INSPECTED_WITH_CONTENT`, `INSPECTED_NO_FORMAL_CONTENT`, `UNREADABLE`,
  `UNASSESSED`. Sweep body text, figures, tables, footnotes, appendices.
- `items`: every numbered equation plus every definition, stated assumption,
  convention, boundary/initial condition, regime statement, and conclusion
  relevant to the derivation. `label` is the paper's own number or
  `unlabeled-N`. `passage` is a short verbatim quote locating the item.
- `inferred`: `true` only for assumptions YOU supply that the authors do not
  state. Author-stated assumptions are `false`. Never present an inference
  as an author claim.
- `math_shape`: your best judgment of the operation shape; use
  `other_unsupported` with a reason whenever the mathematics involves
  vector calculus identities, surface/volume integration, tensor
  manipulation, or anything beyond plain algebra — do NOT force it into a
  supported shape.
- `relations`: derivation dependencies between items; `required:true` when
  the argument breaks without that link.
- `capture_status`: `CAPTURED` (fully transcribed with refs) or `AMBIGUOUS`
  (unclear transcription, missing definition, or uncertain scope).
- `test_status`: always `NOT_RUN` at extraction time.
- `pdf_sha256_hex`, counts you cannot verify: `null`/`0` — the human fills
  these in; do not fabricate them.

## What success looks like

A complete, honest inventory: everything found is transcribed with its
source location; everything uncertain is flagged; unsupported mathematics
is preserved with its reason, not hidden; no page in scope is unaccounted.
