# Demo Checklist (local rehearsal)

## Preflight

- [ ] Branch `phystrace/k1-mcp-bootstrap`, worktree clean
      (`git status --porcelain` empty).
- [ ] Deps at pinned tags (see `workspace.lock.json` after setup).
- [ ] `bin/phystrace-mcp` present and executable.
- [ ] `.agents/mcp_config.json` valid JSON, `command` absolute.
- [ ] Tool registry is exactly the four tools
      (`TestServerToolListPinsFourTools`).

## Setup

- [ ] Fresh disposable workspace via `scripts/setup-antigravity.sh`.
- [ ] Four-module layout, lock record SHAs match pins.
- [ ] Build + smoke checks pass as printed by the script.

## Review

- [ ] Workspace opened in Antigravity 2.0.2 (version from the installer
  desktop entry; no CLI version oracle — `--version` launches the GUI);
  `phystrace` server listed. **Manual connection check required before
  declaring demo readiness: automated tests cover the config file, not
  the live client handshake.**
- [ ] `prompts/review-durr-minisuperspace.md` submitted unchanged.
- [ ] Agent ran an approved profile; real exit code captured.
- [ ] Anti-goldens named with their rejection behavior.

## Evidence

- [ ] Review output saved; test commands + exit codes recorded.
- [ ] Baseline/implementation SHAs, worktree before/after status saved.
- [ ] Unresolved scientific claims listed as unresolved.

## Failure conditions (stop, report, never auto-repair)

Tag/SHA mismatch, dirty repo, missing source evidence, MCP init
failure, approved-test failure, claim beyond evidence. A completed
local rehearsal is not independent judge reproducibility.
