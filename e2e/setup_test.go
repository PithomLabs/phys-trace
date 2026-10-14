// Package e2e drives scripts/setup-antigravity.sh against disposable
// workspaces: hermetic fixture remotes for failure paths, plus one
// real-clone integration run against the pinned tags (network).
package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func scriptPath(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller unavailable")
	}
	p := filepath.Join(filepath.Dir(file), "..", "scripts", "setup-antigravity.sh")
	abs, err := filepath.Abs(p)
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
	return strings.TrimSpace(string(out))
}

// mkfixtureRemote creates a minimal module repo with tag v0.1.0 and
// returns its directory and the tag SHA.
func mkfixtureRemote(t *testing.T, parent, name string) (string, string) {
	t.Helper()
	dir := filepath.Join(parent, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	init := exec.Command("git", "init", "-q")
	init.Dir = dir
	if out, err := init.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v\n%s", err, out)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module github.com/PithomLabs/"+name+"\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# "+name+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	add := exec.Command("git", "add", ".")
	add.Dir = dir
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v\n%s", err, out)
	}
	commit := exec.Command("git", "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-qm", "fixture")
	commit.Dir = dir
	if out, err := commit.CombinedOutput(); err != nil {
		t.Fatalf("git commit: %v\n%s", err, out)
	}
	tag := exec.Command("git", "tag", "v0.1.0")
	tag.Dir = dir
	if out, err := tag.CombinedOutput(); err != nil {
		t.Fatalf("git tag: %v\n%s", err, out)
	}
	sha := gitOut(t, dir, "rev-list", "-n", "1", "v0.1.0")
	return dir, sha
}

// fixtureEnv builds env overrides pointing the script at local fixture
// remotes. Returns env additions and the expected SHA map.
func fixtureEnv(t *testing.T) ([]string, map[string]string) {
	t.Helper()
	parent := t.TempDir()
	shas := map[string]string{}
	for _, name := range []string{"phys-artifact", "phys-math", "phys-lib"} {
		_, sha := mkfixtureRemote(t, parent, name)
		shas[name] = sha
	}
	env := []string{
		"PHYSTRACE_REMOTE_BASE=file://" + parent,
		"PHYSTRACE_EXPECT_PHYS_ARTIFACT=" + shas["phys-artifact"],
		"PHYSTRACE_EXPECT_PHYS_MATH=" + shas["phys-math"],
		"PHYSTRACE_EXPECT_PHYS_LIB=" + shas["phys-lib"],
	}
	return env, shas
}

func runScript(t *testing.T, extraEnv []string, args ...string) (int, string) {
	t.Helper()
	cmd := exec.Command("bash", append([]string{scriptPath(t)}, args...)...)
	cmd.Env = append(os.Environ(), extraEnv...)
	out, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		} else {
			t.Fatalf("script exec: %v", err)
		}
	}
	return code, string(out)
}

func readLock(t *testing.T, ws string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(ws, "workspace.lock.json"))
	if err != nil {
		t.Fatal(err)
	}
	var lock map[string]any
	if err := json.Unmarshal(raw, &lock); err != nil {
		t.Fatal(err)
	}
	return lock
}

func TestSetupFreshWorkspaceMechanics(t *testing.T) {
	env, shas := fixtureEnv(t)
	ws, err := os.MkdirTemp("", "demo-*")
	if err != nil {
		t.Fatal(err)
	}
	// RemoveAll registered first so the worktree cleanup below runs first.
	t.Cleanup(func() { os.RemoveAll(ws) })
	target := filepath.Join(ws, "phystrace-demo")
	code, out := runScript(t, env, "--workspace", target)
	if code != 0 {
		t.Fatalf("setup exit %d:\n%s", code, out)
	}
	// Detached worktree registered in the invoking repo: remove it first.
	invoking, err := filepath.Abs(filepath.Join(filepath.Dir(scriptPath(t)), ".."))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cmd := exec.Command("git", "-C", invoking, "worktree", "remove", "--force", filepath.Join(target, "phys-trace"))
		cmd.CombinedOutput()
	})
	lock := readLock(t, target)
	baselines := lock["baselines"].(map[string]any)
	for _, name := range []string{"phys-artifact", "phys-math", "phys-lib"} {
		entry := baselines[name].(map[string]any)
		if entry["tag"] != "v0.1.0" || entry["resolved_sha"] != shas[name] || entry["expected_sha"] != shas[name] {
			t.Fatalf("baseline wrong for %s: %v", name, entry)
		}
		if entry["clean"] != true {
			t.Fatalf("baseline not clean for %s", name)
		}
	}
	adapter := lock["adapter"].(map[string]any)
	if adapter["source"] != "local-detached-worktree" {
		t.Fatalf("adapter source wrong: %v", adapter["source"])
	}
	head := gitOut(t, invoking, "rev-parse", "HEAD")
	if adapter["implementation_sha"] != head {
		t.Fatalf("implementation SHA %v != invoking HEAD %s", adapter["implementation_sha"], head)
	}
	wantWork := "go 1.25.7\n\nuse (\n\t./phys-artifact\n\t./phys-math\n\t./phys-lib\n\t./phys-trace\n)\n"
	raw, err := os.ReadFile(filepath.Join(target, "go.work"))
	if err != nil || string(raw) != wantWork {
		t.Fatalf("go.work wrong: %q, err=%v", raw, err)
	}
	bin := filepath.Join(target, "bin", "phystrace-mcp")
	if fi, err := os.Stat(bin); err != nil || fi.Mode()&0o111 == 0 {
		t.Fatalf("binary missing or not executable: %v", err)
	}
	cfgRaw, err := os.ReadFile(filepath.Join(target, ".agents", "mcp_config.json"))
	if err != nil {
		t.Fatal(err)
	}
	var cfg map[string]any
	if err := json.Unmarshal(cfgRaw, &cfg); err != nil {
		t.Fatal(err)
	}
	entry := cfg["mcpServers"].(map[string]any)["phystrace"].(map[string]any)
	if entry["command"] != bin {
		t.Fatalf("config command %v != binary %s", entry["command"], bin)
	}
	// Idempotent second run: exit 0, no backup created, lock identical.
	lockBefore, _ := os.ReadFile(filepath.Join(target, "workspace.lock.json"))
	code2, out2 := runScript(t, env, "--workspace", target)
	if code2 != 0 {
		t.Fatalf("re-run exit %d:\n%s", code2, out2)
	}
	lockAfter, _ := os.ReadFile(filepath.Join(target, "workspace.lock.json"))
	if string(lockBefore) != string(lockAfter) {
		t.Fatal("re-run changed the lock record")
	}
	matches, _ := filepath.Glob(filepath.Join(target, ".agents", "*.bak.*"))
	if len(matches) != 0 {
		t.Fatalf("re-run created backups: %v", matches)
	}
}

func TestSetupTagMismatchAborts(t *testing.T) {
	env, _ := fixtureEnv(t)
	// Corrupt one expectation: verification must fail, changing nothing.
	env = append(env, "PHYSTRACE_EXPECT_PHYS_MATH=0000000000000000000000000000000000000000")
	ws, err := os.MkdirTemp("", "demo-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(ws) })
	code, out := runScript(t, env, "--workspace", filepath.Join(ws, "phystrace-demo"))
	if code == 0 {
		t.Fatal("mismatched SHA accepted")
	}
	if !strings.Contains(out, "mismatch") && !strings.Contains(out, "refusing") {
		t.Fatalf("failure message unhelpful:\n%s", out)
	}
}

func TestSetupDirtyExistingAborts(t *testing.T) {
	env, _ := fixtureEnv(t)
	ws, err := os.MkdirTemp("", "demo-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(ws) })
	target := filepath.Join(ws, "phystrace-demo")
	// Establish a good workspace first (full mechanics run).
	if code, out := runScript(t, env, "--workspace", target); code != 0 {
		t.Fatalf("initial setup exit %d:\n%s", code, out)
	}
	invoking, err := filepath.Abs(filepath.Join(filepath.Dir(scriptPath(t)), ".."))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cmd := exec.Command("git", "-C", invoking, "worktree", "remove", "--force", filepath.Join(target, "phys-trace"))
		cmd.CombinedOutput()
	})
	// Dirty one member checkout; the re-run must refuse and preserve it.
	dirty := filepath.Join(target, "phys-math", "uncommitted.txt")
	if err := os.WriteFile(dirty, []byte("work"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, out := runScript(t, env, "--workspace", target)
	if code == 0 {
		t.Fatal("dirty checkout accepted")
	}
	if !strings.Contains(out, "dirty") {
		t.Fatalf("failure message unhelpful:\n%s", out)
	}
	if _, err := os.Stat(dirty); err != nil {
		t.Fatal("user file was discarded")
	}
}

func TestSetupRejectsDangerousAndNested(t *testing.T) {
	env, _ := fixtureEnv(t)
	if code, _ := runScript(t, env, "--workspace", "/"); code == 0 {
		t.Fatal("root destination accepted")
	}
	repoRoot, err := filepath.Abs(filepath.Join(filepath.Dir(scriptPath(t)), ".."))
	if err != nil {
		t.Fatal(err)
	}
	if code, out := runScript(t, env, "--workspace", filepath.Join(repoRoot, "nested-demo")); code == 0 {
		t.Fatal("nested destination accepted")
	} else if !strings.Contains(out, "inside the invoking checkout") {
		t.Fatalf("failure message unhelpful:\n%s", out)
	}
	if code, _ := runScript(t, env, "--bogus-flag"); code == 0 {
		t.Fatal("unknown flag accepted")
	}
}

// worktreeRegistered reports whether path is a registered worktree of repo.
func worktreeRegistered(t *testing.T, repo, path string) bool {
	t.Helper()
	cmd := exec.Command("git", "-C", repo, "worktree", "list", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		if line == "worktree "+path {
			return true
		}
	}
	return false
}

func invokingRepo(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join(filepath.Dir(scriptPath(t)), ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

// TestSetupAdoptsPartialWorkspace reproduces the interrupted-setup case:
// a workspace with only script-created entries and no lock record must
// be adopted and completed, with everything re-validated.
func TestSetupAdoptsPartialWorkspace(t *testing.T) {
	env, _ := fixtureEnv(t)
	ws, err := os.MkdirTemp("", "demo-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(ws) })
	target := filepath.Join(ws, "phystrace-demo")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	// Simulate interruption after two clones: real checkouts, no lock.
	base := strings.TrimPrefix(env[0], "PHYSTRACE_REMOTE_BASE=file://")
	for _, name := range []string{"phys-artifact", "phys-math"} {
		cmd := exec.Command("git", "clone", "-q", "--branch", "v0.1.0", base+"/"+name, filepath.Join(target, name))
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("clone: %v\n%s", err, out)
		}
	}
	code, out := runScript(t, env, "--workspace", target)
	if code != 0 {
		t.Fatalf("partial workspace not adopted:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(target, "workspace.lock.json")); err != nil {
		t.Fatal("adopted run produced no lock record")
	}
	invoking := invokingRepo(t)
	t.Cleanup(func() {
		cmd := exec.Command("git", "-C", invoking, "worktree", "remove", "--force", filepath.Join(target, "phys-trace"))
		cmd.CombinedOutput()
	})
	if !worktreeRegistered(t, invoking, filepath.Join(target, "phys-trace")) {
		t.Fatal("adopted run created no worktree")
	}
}

// TestSetupRefusesForeignPartial proves unrelated files block adoption
// with an exact recovery command, and are preserved untouched.
func TestSetupRefusesForeignPartial(t *testing.T) {
	env, _ := fixtureEnv(t)
	ws, err := os.MkdirTemp("", "demo-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(ws) })
	target := filepath.Join(ws, "phystrace-demo")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	mine := filepath.Join(target, "my-notes.txt")
	if err := os.WriteFile(mine, []byte("user data"), 0o644); err != nil {
		t.Fatal(err)
	}
	code, out := runScript(t, env, "--workspace", target)
	if code == 0 {
		t.Fatal("foreign partial workspace adopted")
	}
	if !strings.Contains(out, "rm -rf") || !strings.Contains(out, target) {
		t.Fatalf("no exact recovery command:\n%s", out)
	}
	if raw, err := os.ReadFile(mine); err != nil || string(raw) != "user data" {
		t.Fatal("user file disturbed")
	}
}

// TestSetupCleansOwnedWorktreeOnFailure forces a failure after worktree
// creation (wrong go.work content) and requires the trap to unregister
// exactly the owned worktree while leaving everything else intact.
func TestSetupCleansOwnedWorktreeOnFailure(t *testing.T) {
	env, _ := fixtureEnv(t)
	ws, err := os.MkdirTemp("", "demo-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(ws) })
	target := filepath.Join(ws, "phystrace-demo")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	wrong := filepath.Join(target, "go.work")
	if err := os.WriteFile(wrong, []byte("go 1.25\n\nuse ./something-else\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	invoking := invokingRepo(t)
	headBefore, err := exec.Command("git", "-C", invoking, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	code, out := runScript(t, env, "--workspace", target)
	if code == 0 {
		t.Fatal("wrong go.work accepted")
	}
	if worktreeRegistered(t, invoking, filepath.Join(target, "phys-trace")) {
		t.Fatal("failed run left its worktree registered")
	}
	headAfter, err := exec.Command("git", "-C", invoking, "rev-parse", "HEAD").Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(headBefore) != string(headAfter) {
		t.Fatal("invoking checkout moved")
	}
	if raw, err := os.ReadFile(wrong); err != nil || !strings.Contains(string(raw), "something-else") {
		t.Fatal("pre-existing go.work disturbed")
	}
	if !strings.Contains(out, "go.work") {
		t.Fatalf("failure message unhelpful:\n%s", out)
	}
}

// TestSetupRealClones exercises actual pinned clones over the network:
// small repos, full setup including build and smoke. Requires network;
// the test fails honestly (no skip) when clones break.
func TestSetupRealClones(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git unavailable")
	}
	ws, err := os.MkdirTemp("", "demo-real-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(ws) })
	target := filepath.Join(ws, "phystrace-demo")
	code, out := runScript(t, nil, "--workspace", target)
	if code != 0 {
		t.Fatalf("setup exit %d:\n%s", code, out)
	}
	invoking, _ := filepath.Abs(filepath.Join(filepath.Dir(scriptPath(t)), ".."))
	t.Cleanup(func() {
		cmd := exec.Command("git", "-C", invoking, "worktree", "remove", "--force", filepath.Join(target, "phys-trace"))
		cmd.CombinedOutput()
	})
	lock := readLock(t, target)
	baselines := lock["baselines"].(map[string]any)
	for name, want := range map[string]string{
		"phys-artifact": "36584f6a802f0641fd03a459176e07c020f671ad",
		"phys-math":     "2053e2baea954d92e3441ab99dc0973faf6fe29c",
		"phys-lib":      "38ba25a0ec311cf02ed4ec0be7237b871f9bbff0",
	} {
		entry := baselines[name].(map[string]any)
		if entry["resolved_sha"] != want {
			t.Fatalf("%s resolved to %v, want %s", name, entry["resolved_sha"], want)
		}
	}
	if fi, err := os.Stat(filepath.Join(target, "bin", "phystrace-mcp")); err != nil || fi.Mode()&0o111 == 0 {
		t.Fatalf("binary missing: %v", err)
	}
}
