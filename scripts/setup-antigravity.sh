#!/usr/bin/env bash
#
# scripts/setup-antigravity.sh — prepare a pinned local PhysTrace demo workspace.
#
# Usage: setup-antigravity.sh [--workspace DIR]
#
# Creates <workspace>/ with the pinned dependency checkouts, a detached
# worktree of the invoking PhysTrace implementation, a four-member
# go.work, the built MCP binary, workspace-local Antigravity config, and
# a lock record. Prints next steps; never launches Antigravity.
#
# Safety: never overwrites existing repositories, never discards
# uncommitted work, never touches global user configuration. Re-running
# against the same workspace is idempotent.
#
# Test hooks (environment, all optional):
#   PHYSTRACE_REMOTE_BASE   remote root (default https://github.com/PithomLabs)
#   PHYSTRACE_EXPECT_PHYS_ARTIFACT / _PHYS_MATH / _PHYS_LIB
#                           expected tag SHAs (default the verified v0.1.0 SHAs)
#   PHYSTRACE_SKIP_BUILD=1  skip the binary build
#   PHYSTRACE_SKIP_SMOKE=1  skip smoke checks
#
set -euo pipefail

# The demo workspace must resolve in workspace mode regardless of the
# caller's Go environment: an inherited GOWORK=off (or a stray GOFLAGS)
# would blind `go list`/`go build` to the generated go.work.
unset GOWORK

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REMOTE_BASE="${PHYSTRACE_REMOTE_BASE:-https://github.com/PithomLabs}"
EXPECT_PHYS_ARTIFACT="${PHYSTRACE_EXPECT_PHYS_ARTIFACT:-36584f6a802f0641fd03a459176e07c020f671ad}"
EXPECT_PHYS_MATH="${PHYSTRACE_EXPECT_PHYS_MATH:-2053e2baea954d92e3441ab99dc0973faf6fe29c}"
EXPECT_PHYS_LIB="${PHYSTRACE_EXPECT_PHYS_LIB:-38ba25a0ec311cf02ed4ec0be7237b871f9bbff0}"
GO_MIN_MAJOR=1
GO_MIN_MINOR=25

die() { echo "setup-antigravity: error: $*" >&2; exit 1; }
note() { echo "setup-antigravity: $*"; }

# --- Worktree ownership --------------------------------------------------
# CREATED_WORKTREE records a phys-trace worktree this run created.
# SETUP_OK is set only after fully successful completion. On failure or
# interruption, the trap unregisters exactly the worktree this run owns;
# a pre-existing worktree (CREATED_WORKTREE empty) is never touched, and
# a successful run keeps its worktree. Plain `worktree remove` (no
# --force) refuses when dirty, so user content is never discarded by
# cleanup; leftovers print a manual recovery command instead.
CREATED_WORKTREE=""
SETUP_OK=""
cleanup_owned_worktree() {
	[ -n "${SRC_ROOT:-}" ] || return 0
	[ -n "$CREATED_WORKTREE" ] || return 0
	if git -C "$SRC_ROOT" worktree list --porcelain 2>/dev/null | grep -Fxq "worktree $CREATED_WORKTREE"; then
		if git -C "$SRC_ROOT" worktree remove "$CREATED_WORKTREE" 2>/dev/null; then
			note "removed owned worktree $CREATED_WORKTREE"
		else
			note "leaving worktree $CREATED_WORKTREE registered; remove manually: git -C '$SRC_ROOT' worktree remove '$CREATED_WORKTREE'"
		fi
	fi
}
on_exit() {
	if [ -z "$SETUP_OK" ]; then
		cleanup_owned_worktree
	fi
}
trap 'on_exit' EXIT
trap 'on_exit; exit 130' INT
trap 'on_exit; exit 143' TERM

WORKSPACE=""
while [ $# -gt 0 ]; do
	case "$1" in
		--workspace) WORKSPACE="${2:-}"; shift 2 ;;
		--workspace=*) WORKSPACE="${1#--workspace=}"; shift ;;
		-h|--help) sed -n '2,20p' "$0"; exit 0 ;;
		*) die "unknown argument: $1 (try --help)" ;;
	esac
done
[ -n "$WORKSPACE" ] || WORKSPACE="$PWD/phystrace-demo"

# --- 1. Prerequisites -------------------------------------------------
command -v git >/dev/null || die "git not found in PATH"
command -v go >/dev/null || die "go toolchain not found in PATH"
command -v python3 >/dev/null || die "python3 not found in PATH (needed for config/lock JSON)"
GO_VER="$(go version | awk '{print $3}')" || die "cannot determine go version"
GO_VER="${GO_VER#go}"
GO_MAJOR="${GO_VER%%.*}"; GO_REST="${GO_VER#*.}"; GO_MINOR="${GO_REST%%.*}"
if [ "$GO_MAJOR" -lt "$GO_MIN_MAJOR" ] || { [ "$GO_MAJOR" -eq "$GO_MIN_MAJOR" ] && [ "$GO_MINOR" -lt "$GO_MIN_MINOR" ]; }; then
	die "go >= $GO_MIN_MAJOR.$GO_MIN_MINOR required, found $GO_VER"
fi

# --- 2. Destination ----------------------------------------------------
case "$WORKSPACE" in
	""|"/") die "refusing dangerous workspace destination: '$WORKSPACE'" ;;
esac
# Resolve lexically; realpath may be unavailable for missing paths.
WS_ABS="$(cd "$(dirname "$WORKSPACE")" 2>/dev/null && pwd)/$(basename "$WORKSPACE")" || die "cannot resolve workspace parent"
if [ -e "$WS_ABS" ] && [ ! -d "$WS_ABS" ]; then
	die "workspace path exists and is not a directory: $WS_ABS"
fi
# Refuse destinations nested inside the invoking PhysTrace checkout:
# sibling clones and worktrees belong beside it, not inside it.
SRC_ROOT="$(git -C "$SCRIPT_DIR/.." rev-parse --show-toplevel 2>/dev/null)" || die "cannot locate invoking PhysTrace checkout"
case "$WS_ABS" in
	"$SRC_ROOT"|"$SRC_ROOT"/*) die "workspace must not live inside the invoking checkout ($SRC_ROOT); pass --workspace outside it" ;;
esac
if [ -d "$WS_ABS" ] && [ -n "$(ls -A "$WS_ABS")" ] && [ ! -f "$WS_ABS/workspace.lock.json" ]; then
	# Adoptable partial workspace: only directories and files the script
	# itself creates. Anything else means user data lives here; refuse
	# with the exact recovery command rather than adopting or deleting.
	# Adopted contents are fully re-validated below (remotes, tags, SHAs,
	# cleanliness), so adoption relaxes no trust.
	ADOPTABLE=1
	while IFS= read -r entry; do
		case "$(basename "$entry")" in
			phys-artifact|phys-math|phys-lib|phys-trace|go.work|bin|.agents) ;;
			*) ADOPTABLE=0; break ;;
		esac
	done < <(find "$WS_ABS" -mindepth 1 -maxdepth 1)
	if [ "$ADOPTABLE" != "1" ]; then
		die "destination $WS_ABS holds files this setup did not create; refusing to adopt it (nothing was changed). Recovery: inspect, back up what matters, then run: rm -rf '$WS_ABS'"
	fi
	note "adopting partial workspace (contents will be re-validated)"
fi
mkdir -p "$WS_ABS"
note "workspace: $WS_ABS"

# --- 3. Pinned dependency checkouts ------------------------------------
clone_or_verify() {
	local name="$1" expect_sha="$2"
	local dest="$WS_ABS/$name" remote="$REMOTE_BASE/$name"
	if [ ! -e "$dest" ]; then
		note "cloning $name at v0.1.0"
		git clone --branch v0.1.0 "$remote" "$dest" >&2 || die "clone of $remote failed"
	fi
	local actual_remote
	actual_remote="$(git -C "$dest" remote get-url origin 2>/dev/null)" || die "$name: no origin remote"
	# Normalize file:// URLs: `git clone /path` and `git clone file:///path`
	# denote the same repository but record different strings.
	actual_remote="${actual_remote#file://}"
	local want_remote="${remote#file://}"
	[ "$actual_remote" = "$want_remote" ] || die "$name: origin is $actual_remote, want $remote (not adopting a foreign checkout)"
	local resolved
	if ! resolved="$(git -C "$dest" rev-list -n 1 v0.1.0 2>/dev/null)"; then
		note "$name: tag v0.1.0 missing locally, fetching tags"
		git -C "$dest" fetch --tags origin >&2 || die "$name: tag fetch failed"
		resolved="$(git -C "$dest" rev-list -n 1 v0.1.0 2>/dev/null)" || die "$name: tag v0.1.0 unavailable after fetch"
	fi
	[ "$resolved" = "$expect_sha" ] || die "$name: v0.1.0 resolves to $resolved, want $expect_sha (refusing mismatched code)"
	if [ -n "$(git -C "$dest" status --porcelain)" ]; then
		die "$name: working tree is dirty; commit, stash, or clean it, then re-run (nothing was changed)"
	fi
	note "$name: $resolved (clean)"
}
clone_or_verify phys-artifact "$EXPECT_PHYS_ARTIFACT"
clone_or_verify phys-math "$EXPECT_PHYS_MATH"
clone_or_verify phys-lib "$EXPECT_PHYS_LIB"

# --- 4. PhysTrace implementation worktree -------------------------------
# The v0.1.0 tag holds docs/protocol only, so the adapter is sourced from
# the invoking checkout at its exact commit via a detached worktree: the
# main checkout is never disturbed and the revision is recorded exactly.
IMPL_SHA="$(git -C "$SRC_ROOT" rev-parse HEAD)" || die "cannot resolve implementation commit"
IMPL_BRANCH="$(git -C "$SRC_ROOT" rev-parse --abbrev-ref HEAD)"
if [ -e "$WS_ABS/phys-trace" ]; then
	REG="$(git -C "$WS_ABS/phys-trace" rev-parse HEAD 2>/dev/null)" || die "existing $WS_ABS/phys-trace is not a usable checkout (remove it or point elsewhere)"
	[ "$REG" = "$IMPL_SHA" ] || die "existing phys-trace worktree is at $REG, want implementation $IMPL_SHA (remove it or reconcile manually)"
	note "phys-trace worktree already at $IMPL_SHA"
else
	note "creating detached phys-trace worktree at $IMPL_SHA"
	git -C "$SRC_ROOT" worktree add --detach "$WS_ABS/phys-trace" "$IMPL_SHA" >&2 || die "worktree creation failed"
	CREATED_WORKTREE="$WS_ABS/phys-trace"
fi

# --- 5. Workspace file ---------------------------------------------------
GOWORK_WANT="$(printf 'go 1.25.7\n\nuse (\n\t./phys-artifact\n\t./phys-math\n\t./phys-lib\n\t./phys-trace\n)\n')"
if [ -f "$WS_ABS/go.work" ]; then
	if [ "$(cat "$WS_ABS/go.work")" != "$GOWORK_WANT" ]; then
		die "$WS_ABS/go.work exists with unexpected content; refusing to overwrite (inspect or remove it, then re-run)"
	fi
	note "go.work already correct"
else
	printf '%s\n' "$GOWORK_WANT" > "$WS_ABS/go.work"
	note "go.work written"
fi

# --- 6. Build ------------------------------------------------------------
BIN="$WS_ABS/bin/phystrace-mcp"
if [ "${PHYSTRACE_SKIP_BUILD:-0}" = "1" ]; then
	note "build skipped (PHYSTRACE_SKIP_BUILD=1)"
else
	note "building bin/phystrace-mcp"
	# Standalone module build: the adapter imports only the MCP SDK and
	# the standard library (thin-adapter boundary), so it resolves
	# without the workspace. A workspace-mode build would force module
	# graph materialization across phys-lib's unpublished v0.0.0
	# requirement, which no network can satisfy; that is a property of
	# the deferred version scheme (see phys-lib/go.mod M-4 note), not of
	# this workspace.
	(cd "$WS_ABS/phys-trace" && GOWORK=off go build -o "$BIN" ./cmd/phystrace-mcp) >&2 || die "adapter build failed"
	[ -x "$BIN" ] || die "built binary missing: $BIN"
fi

# --- 7. Smoke checks ------------------------------------------------------
if [ "${PHYSTRACE_SKIP_SMOKE:-0}" = "1" ]; then
	note "smoke checks skipped (PHYSTRACE_SKIP_SMOKE=1)"
else
	note "smoke: module graph (workspace membership)"
	(cd "$WS_ABS" && go list -m >/dev/null) >&2 || die "module graph incoherent (go list -m failed)"
	note "smoke: stdio startup (stdin EOF must exit 0 with empty stdout)"
	SMOKE_TMP="$(mktemp)"
	"$BIN" --workspace "$WS_ABS" </dev/null >"$SMOKE_TMP" 2>/dev/null &
	SMOKE_PID=$!
	SMOKE_WAIT=0
	while kill -0 "$SMOKE_PID" 2>/dev/null; do
		sleep 1
		SMOKE_WAIT=$((SMOKE_WAIT + 1))
		if [ "$SMOKE_WAIT" -ge 30 ]; then
			kill -9 "$SMOKE_PID" 2>/dev/null || true
			rm -f "$SMOKE_TMP"
			die "binary did not exit 30s after stdin EOF"
		fi
	done
	wait "$SMOKE_PID" || die "binary failed clean stdio startup"
	[ -s "$SMOKE_TMP" ] && { rm -f "$SMOKE_TMP"; die "binary wrote to stdout outside protocol framing"; }
	rm -f "$SMOKE_TMP"
	note "smoke checks passed"
fi

# --- 8. Antigravity config -------------------------------------------------
note "writing Antigravity config"
CONFIG_SUMMARY="$(python3 "$SCRIPT_DIR/antigravity_config.py" --workspace "$WS_ABS" --binary "$BIN")" || die "config generation failed"
note "config: $CONFIG_SUMMARY"

# --- 9. Lock record ----------------------------------------------------------
python3 - "$WS_ABS" "$REMOTE_BASE" "$EXPECT_PHYS_ARTIFACT" "$EXPECT_PHYS_MATH" "$EXPECT_PHYS_LIB" "$IMPL_SHA" "$IMPL_BRANCH" <<'PYEOF' || die "lock record failed"
import json, subprocess, sys
ws, base, a_sha, m_sha, l_sha, impl_sha, impl_branch = sys.argv[1:8]
def rev(d, *args):
    return subprocess.run(["git", "-C", d, *args], capture_output=True, text=True, check=True).stdout.strip()
def lock():
    deps = {}
    for name, sha in (("phys-artifact", a_sha), ("phys-math", m_sha), ("phys-lib", l_sha)):
        d = "%s/%s" % (ws, name)
        deps[name] = {
            "remote": "%s/%s" % (base, name),
            "tag": "v0.1.0", "expected_sha": sha,
            "resolved_sha": rev(d, "rev-list", "-n", "1", "v0.1.0"),
            "head": rev(d, "rev-parse", "HEAD"),
            "clean": rev(d, "status", "--porcelain") == "",
        }
    return {
        "baselines": deps,
        "adapter": {
            "repo": "%s/phys-trace" % base,
            "baseline_tag": "v0.1.0",
            "baseline_sha": rev("%s/phys-trace" % ws, "rev-list", "-n", "1", "v0.1.0"),
            "implementation_ref": "worktree@%s" % impl_sha,
            "implementation_sha": impl_sha,
            "implementation_branch": impl_branch,
            "worktree": "clean",
            "source": "local-detached-worktree",
            "tools": ["case_overview", "inspect_evidence", "repository_status", "run_verification"],
        },
    }
with open("%s/workspace.lock.json" % ws, "w") as f:
    json.dump(lock(), f, indent=2, sort_keys=True)
    f.write("\n")
PYEOF
note "lock record written"
SETUP_OK=1

# --- 10. Next steps ------------------------------------------------------------
cat <<EOF
setup-antigravity: done.
  workspace: $WS_ABS
  binary:    $BIN
  config:    $WS_ABS/.agents/mcp_config.json
  lock:      $WS_ABS/workspace.lock.json
Next:
  1. Open Antigravity at this workspace: $WS_ABS
  2. Confirm the phystrace server appears in MCP settings (restart the
     client if it was already running).
  3. Paste prompts/review-durr-minisuperspace.md as a new agent task
     (see the phys-trace repository prompts/ directory once P4 lands;
     until then use TOE-V4/review drafts).
  4. Ask the agent for repository_status first: dirty trees are reported,
     never cleaned.
EOF
