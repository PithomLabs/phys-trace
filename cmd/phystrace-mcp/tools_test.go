package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// minimalManifest is a hermetic fixture manifest for tool tests.
const minimalManifest = `{
  "name": "Fixture Demo",
  "source": {"arxiv_id": "0000.00000", "scope_note": "test scope"},
  "regime": "test regime",
  "premises": [{"id": "p1", "statement": "s1"}],
  "limitations": [
    {"id": "no_x", "statement": "x unsupported"},
    {"id": "no_y", "statement": "y unsupported"}
  ]
}`

// mkworkspace builds a temp workspace with the manifest at its demo path.
func mkworkspace(t *testing.T, manifest string) workspace {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "phys-lib", "bohm", "minisuperspace")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if manifest != "" {
		if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(manifest), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	ws, err := resolveWorkspace(root)
	if err != nil {
		t.Fatal(err)
	}
	return ws
}

func TestResolveWorkspaceRejects(t *testing.T) {
	if _, err := resolveWorkspace(""); err == nil {
		t.Fatal("empty path accepted")
	}
	if _, err := resolveWorkspace(filepath.Join("no-such-dir-xyz", "sub")); err == nil {
		t.Fatal("missing dir accepted")
	}
	f := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := resolveWorkspace(f); err == nil {
		t.Fatal("regular file accepted as workspace")
	}
}

func TestCaseOverviewGolden(t *testing.T) {
	ws := mkworkspace(t, minimalManifest)
	out, err := caseOverview(ws)
	if err != nil {
		t.Fatal(err)
	}
	if out.Paper != "Fixture Demo" || out.Scope != "test scope" || out.Regime != "test regime" {
		t.Fatalf("manifest fields not echoed: %+v", out)
	}
	if out.DerivationUnits != 12 || out.Artifacts != 14 {
		t.Fatalf("inventory counts drifted: %+v", out)
	}
	if len(out.AntiGoldens) != 3 || out.AntiGoldens[2] != "TestHistoricalWrongPath" {
		t.Fatalf("anti-golden list wrong: %+v", out.AntiGoldens)
	}
	if len(out.Limitations) != 2 || out.Limitations[0] != "no_x" {
		t.Fatalf("limitations wrong: %+v", out.Limitations)
	}
	if out.VerifiedBy == "" {
		t.Fatal("provenance missing")
	}
}

func TestCaseOverviewMissingManifest(t *testing.T) {
	ws := mkworkspace(t, "")
	if _, err := caseOverview(ws); err == nil {
		t.Fatal("missing manifest accepted")
	}
}

func TestInspectEvidenceSelectors(t *testing.T) {
	ws := mkworkspace(t, minimalManifest)
	for _, sel := range []string{"manifest", "inventory", "q-relation", "anti-goldens", "covariant-fragment", "limitations"} {
		doc, err := inspectEvidence(ws, InspectEvidenceIn{Selector: sel})
		if err != nil {
			t.Fatalf("selector %q failed: %v", sel, err)
		}
		if doc.Body == "" || len(doc.Sources) == 0 {
			t.Fatalf("selector %q returned empty evidence", sel)
		}
	}
	if _, err := inspectEvidence(ws, InspectEvidenceIn{Selector: "prove-riemann-hypothesis"}); err == nil {
		t.Fatal("unknown selector accepted")
	}
	if _, err := inspectEvidence(ws, InspectEvidenceIn{Selector: "file"}); err == nil {
		t.Fatal("file selector without path accepted")
	}
}

func TestInspectEvidenceFileGuard(t *testing.T) {
	ws := mkworkspace(t, minimalManifest)
	rel := filepath.Join("phys-lib", "bohm", "minisuperspace", "manifest.json")
	doc, err := inspectEvidence(ws, InspectEvidenceIn{Selector: "file", Path: rel})
	if err != nil {
		t.Fatalf("legitimate file read failed: %v", err)
	}
	if !strings.Contains(doc.Body, "Fixture Demo") {
		t.Fatal("file content wrong")
	}
	for _, evil := range []string{"../escape", "../../etc/passwd", "/etc/passwd", "phys-lib/../../x"} {
		if _, err := inspectEvidence(ws, InspectEvidenceIn{Selector: "file", Path: evil}); err == nil {
			t.Fatalf("traversal path accepted: %q", evil)
		}
	}
	if _, err := inspectEvidence(ws, InspectEvidenceIn{Selector: "file", Path: "no-such-file.json"}); err == nil {
		t.Fatal("missing file accepted")
	}
}

// gitRepo builds a one-commit git repo for status tests; skips when git
// is unavailable.
func gitRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git unavailable")
	}
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q")
	run("config", "user.name", "t")
	run("config", "user.email", "t@t")
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", "f.txt")
	run("commit", "-qm", "init")
	return dir
}

func TestRepositoryStatusCleanAndDirty(t *testing.T) {
	member := gitRepo(t)
	root := t.TempDir()
	if err := os.Rename(member, filepath.Join(root, "phys-math")); err != nil {
		t.Fatal(err)
	}
	ws, err := resolveWorkspace(root)
	if err != nil {
		t.Fatal(err)
	}
	out, err := gitStatus(ws)
	if err != nil {
		t.Fatal(err)
	}
	var math *RepoStatus
	present := 0
	for i := range out.Repos {
		if out.Repos[i].Present {
			present++
			if out.Repos[i].Name == "phys-math" {
				math = &out.Repos[i]
			}
		}
	}
	if present != 1 || math == nil {
		t.Fatalf("membership wrong: %+v", out.Repos)
	}
	if !math.Clean || math.HEAD == "" || len(math.HEAD) != 40 || math.Branch == "" {
		t.Fatalf("clean repo misreported: %+v", math)
	}
	if math.Tags != "" {
		t.Fatalf("unexpected tags: %q", math.Tags)
	}
	// Dirty file must be reported, never cleaned.
	if err := os.WriteFile(filepath.Join(root, "phys-math", "dirty.txt"), []byte("y"), 0o644); err != nil {
		t.Fatal(err)
	}
	out2, err := gitStatus(ws)
	if err != nil {
		t.Fatal(err)
	}
	for i := range out2.Repos {
		if out2.Repos[i].Name == "phys-math" {
			if out2.Repos[i].Clean {
				t.Fatal("dirty worktree reported clean")
			}
			if !strings.Contains(out2.Repos[i].Status, "dirty.txt") {
				t.Fatalf("dirty file not listed: %q", out2.Repos[i].Status)
			}
		}
	}
	// Absent members report absence, not failure.
	for i := range out2.Repos {
		if out2.Repos[i].Name == "phys-lib" && (out2.Repos[i].Present || out2.Repos[i].Error == "") {
			t.Fatalf("absent member misreported: %+v", out2.Repos[i])
		}
	}
}

// TestCaseOverviewRoundTrip exercises flags → server → tool end to end:
// the binary serves case_overview against the given workspace.
func TestCaseOverviewRoundTrip(t *testing.T) {
	ws := mkworkspace(t, minimalManifest)
	bin := buildServerBinary(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client := mcp.NewClient(&mcp.Implementation{Name: "p1-rt", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, &mcp.CommandTransport{Command: exec.Command(bin, "--workspace", ws.root)}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer cs.Close()
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "case_overview"})
	if err != nil {
		t.Fatalf("case_overview: %v", err)
	}
	if res.IsError {
		t.Fatalf("case_overview returned error: %v", res.Content)
	}
}
