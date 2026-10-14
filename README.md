# PhysTrace Local MCP Demo

One-command local setup plus an adversarial review of the existing
Dürr–Struyve minisuperspace formalization. Demonstrates: supported
derivations become runnable, verifiable test suites — nothing more.

## What this is not

No full Quantum Einstein Equations, no general Wheeler–DeWitt solving,
no time evolution, no general curvature support, no Bianchi proof, no
empirical claims. See `phys-lib/bohm/minisuperspace/manifest.json`.

## Prerequisites

Git, Go >= 1.25, python3, network access to github.com.

## Setup

```bash
scripts/setup-antigravity.sh --workspace /path/to/phystrace-demo
```

The script clones `phys-artifact`, `phys-math`, `phys-lib` at pinned
`v0.1.0` tags (tag + commit-SHA enforced, mismatch aborts), sources the
adapter from a detached worktree of this checkout at its exact commit,
writes a four-member `go.work` (frozen Phys excluded), builds
`bin/phystrace-mcp`, runs smoke checks, writes
`.agents/mcp_config.json` (absolute binary path, unrelated servers
preserved, backup on change), and records `workspace.lock.json`.
Re-running is idempotent. It never overwrites repos, discards work, or
touches global user config. It does not launch Antigravity.

Test hooks: `PHYSTRACE_REMOTE_BASE`, `PHYSTRACE_EXPECT_PHYS_*`,
`PHYSTRACE_SKIP_BUILD`, `PHYSTRACE_SKIP_SMOKE` (see script header).

## Open in Antigravity

Open the workspace directory. Confirm the `phystrace` server appears in
MCP settings (restart the client if it was already running). Then paste
`prompts/review-durr-minisuperspace.md` as a new agent task.

## MCP server

```bash
bin/phystrace-mcp --workspace /path/to/phystrace-demo
# or: PHYSTRACE_WORKSPACE=/path/to/phystrace-demo bin/phystrace-mcp
```

Stdio transport; stdout carries protocol frames only, logs go to
stderr. Startup stdout purity is pinned by an automated test; purity
later in an active session rests on construction (single stderr logger,
SDK-owned protocol writes) and is not independently test-pinned. Four
tools, exactly: `case_overview`, `inspect_evidence`
(selectors: manifest, inventory, q-relation, anti-goldens,
covariant-fragment, limitations, file), `run_verification` (profiles
below), `repository_status`. No shell, no writes, no other tools.
Workspace path checks reject lexical escapes; symlink escapes are
rejected after symlink resolution, with the residual limitation that
resolution and file access are not atomic (time-of-check/time-of-use),
so the guard protects the review workflow rather than serving as an
OS-level sandbox.

## Verification profiles (`run_verification` only)

| Profile | Command (in phys-lib) | Timeout |
|---|---|---|
| `demo-arc` | 8-test filter, protocol package | 120s |
| `coverage-audit` | audit + inventory guard | 180s |
| `full-suite` | `go test -count=1 ./...` | 600s |

Output capped at 32KB per stream (truncation flagged); real exit codes
reported; client cancellation propagates. Running tests may use caches
and temp dirs — this tool executes, it does not merely observe.

## Lock record and troubleshooting

`workspace.lock.json` records baselines (remote/tag/expected/resolved
SHA), the adapter worktree revision, and the tool list. On failure:
tag mismatch → wrong code, re-clone; dirty tree → commit/clean and
re-run; config invalid JSON → fix or remove `.agents/mcp_config.json`
(a backup is made before changes); binary missing → re-run setup
without `PHYSTRACE_SKIP_BUILD`.

## Distribution status

This checkout's adapter exists on a local branch only. A fresh,
independent workspace cannot fetch it yet — local bootstrap mechanics
are proven, judge-side obtainability is not. See `plan1.24.2.md`.
