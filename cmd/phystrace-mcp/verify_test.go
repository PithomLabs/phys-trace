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

func TestProfileMappingsExact(t *testing.T) {
	specs := profileSpecs()
	if len(specs) != 3 {
		t.Fatalf("allowlist holds %d profiles, want 3", len(specs))
	}
	demo := specs[ProfileDemoArc]
	if strings.Join(demo.argv, " ") != "go test -count=1 -run TestDerivationInventory|TestHistoricalWrongPath|TestOperatorCompositionVsSquare|TestWValue|TestQRelation|TestPressureFromQ|TestSubstitutionPoint|TestVerifyCoverageAudit ./bohm/minisuperspace/protocol/" {
		t.Fatalf("demo-arc argv drifted: %q", strings.Join(demo.argv, " "))
	}
	if strings.Join(demo.dir, "/") != "phys-lib" {
		t.Fatalf("demo-arc dir drifted: %q", strings.Join(demo.dir, "/"))
	}
	cov := specs[ProfileCoverageAudit]
	if strings.Join(cov.argv, " ") != "go test -count=1 -run TestVerifyCoverageAudit|TestDerivationInventory ./bohm/minisuperspace/protocol/" {
		t.Fatalf("coverage-audit argv drifted: %q", strings.Join(cov.argv, " "))
	}
	full := specs[ProfileFullSuite]
	if strings.Join(full.argv, " ") != "go test -count=1 ./..." || strings.Join(full.dir, "/") != "phys-lib" {
		t.Fatalf("full-suite mapping drifted: %q in %q", strings.Join(full.argv, " "), strings.Join(full.dir, "/"))
	}
	for name, spec := range specs {
		if len(spec.argv) < 2 || spec.timeout <= 0 {
			t.Fatalf("profile %q has degenerate spec", name)
		}
		// The executable is a bare name resolved via PATH: no caller
		// path can reach the command, and no shell interprets argv
		// (exec runs the vector directly).
		if spec.argv[0] != filepath.Base(spec.argv[0]) {
			t.Fatalf("profile %q executable is a path, not a bare name: %q", name, spec.argv[0])
		}
	}
}

func TestRunVerificationRejectsUnknown(t *testing.T) {
	ws := mkworkspace(t, minimalManifest)
	ctx := context.Background()
	for _, bad := range []string{"", "demo-arc; rm -rf /", "../../../bin/sh", "FULL-SUITE", "demo_arc", "/abs/path"} {
		if _, err := runVerification(ctx, ws, bad); err == nil {
			t.Fatalf("profile accepted: %q", bad)
		}
	}
}

func TestExecutorStartupFailure(t *testing.T) {
	ctx := context.Background()
	spec := profileSpec{dir: []string{"."}, argv: []string{"/nonexistent-binary-xyz-123"}, timeout: 10 * time.Second}
	if _, err := executeSpec(ctx, t.TempDir(), spec); err == nil {
		t.Fatal("missing executable reported success")
	}
}

func TestExecutorExitStatusFaithful(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go toolchain unavailable")
	}
	if _, err := exec.LookPath("false"); err != nil {
		t.Skip("false unavailable")
	}
	ctx := context.Background()
	dir := t.TempDir()
	// Positive control.
	okSpec := profileSpec{dir: []string{"."}, argv: []string{"go", "version"}, timeout: 30 * time.Second}
	okRes, err := executeSpec(ctx, dir, okSpec)
	if err != nil {
		t.Fatalf("go version failed: %v", err)
	}
	if okRes.ExitCode != 0 || okRes.TimedOut || !strings.Contains(okRes.Stdout, "go version") {
		t.Fatalf("positive control wrong: %+v", okRes)
	}
	// Failing command: reported code must equal the directly observed code.
	failSpec := profileSpec{dir: []string{"."}, argv: []string{"false"}, timeout: 30 * time.Second}
	failRes, err := executeSpec(ctx, dir, failSpec)
	if err != nil {
		t.Fatalf("failing command returned error instead of result: %v", err)
	}
	direct := exec.Command("false")
	directErr := direct.Run()
	directCode := -99
	if exitErr, ok := directErr.(*exec.ExitError); ok {
		directCode = exitErr.ExitCode()
	}
	if failRes.ExitCode != directCode || failRes.TimedOut {
		t.Fatalf("exit code %d != direct %d (timed out: %v)", failRes.ExitCode, directCode, failRes.TimedOut)
	}
}

func TestExecutorTimeoutAndTruncation(t *testing.T) {
	if _, err := exec.LookPath("sleep"); err != nil {
		t.Skip("sleep unavailable")
	}
	if _, err := exec.LookPath("cat"); err != nil {
		t.Skip("cat unavailable")
	}
	ctx := context.Background()
	dir := t.TempDir()
	slow := profileSpec{dir: []string{"."}, argv: []string{"sleep", "30"}, timeout: time.Second}
	res, err := executeSpec(ctx, dir, slow)
	if err != nil {
		t.Fatalf("slow command returned error instead of result: %v", err)
	}
	if !res.TimedOut || res.ExitCode != -1 {
		t.Fatalf("timeout not reported distinctly: %+v", res)
	}
	// 200KB through a 32KB cap.
	big := filepath.Join(dir, "big.txt")
	blob := make([]byte, 200<<10)
	for i := range blob {
		blob[i] = byte('a' + i%26)
	}
	if err := os.WriteFile(big, blob, 0o644); err != nil {
		t.Fatal(err)
	}
	loud := profileSpec{dir: []string{"."}, argv: []string{"cat", big}, timeout: 30 * time.Second}
	loudRes, err := executeSpec(ctx, dir, loud)
	if err != nil {
		t.Fatal(err)
	}
	if loudRes.ExitCode != 0 || !loudRes.Truncated || len(loudRes.Stdout) != maxVerifyBytes {
		t.Fatalf("truncation wrong: exit=%d truncated=%v len=%d", loudRes.ExitCode, loudRes.Truncated, len(loudRes.Stdout))
	}
}

func TestInspectEvidenceSymlinkEscape(t *testing.T) {
	ws := mkworkspace(t, minimalManifest)
	outside := t.TempDir()
	secret := filepath.Join(outside, "secret.txt")
	if err := os.WriteFile(secret, []byte("secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	linkDir := filepath.Join(ws.root, "phys-lib", "bohm", "linkdir")
	if err := os.Symlink(outside, linkDir); err != nil {
		t.Fatal(err)
	}
	if _, err := inspectEvidence(ws, InspectEvidenceIn{Selector: "file", Path: "phys-lib/bohm/linkdir/secret.txt"}); err == nil {
		t.Fatal("symlink escape to outside file accepted")
	} else if !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("wrong boundary for symlink escape: %v", err)
	}
	// A symlink pointing inside the workspace stays usable.
	inside := filepath.Join(ws.root, "phys-lib", "bohm", "minisuperspace")
	inLink := filepath.Join(ws.root, "phys-lib", "bohm", "inlink")
	if err := os.Symlink(inside, inLink); err != nil {
		t.Fatal(err)
	}
	doc, err := inspectEvidence(ws, InspectEvidenceIn{Selector: "file", Path: "phys-lib/bohm/inlink/manifest.json"})
	if err != nil {
		t.Fatalf("inside-workspace symlink rejected: %v", err)
	}
	if !strings.Contains(doc.Body, "Fixture Demo") {
		t.Fatal("inside symlink content wrong")
	}
}

// fakePhysLibModule builds a minimal hermetic Go module standing in for
// phys-lib: one passing test, so verification profiles exercise real
// `go test` plumbing without depending on the developer checkout.
func fakePhysLibModule(t *testing.T, ws workspace) {
	t.Helper()
	dir := filepath.Join(ws.root, "phys-lib")
	if err := os.MkdirAll(filepath.Join(dir, "bohm", "minisuperspace", "protocol"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/fakelib\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "bohm", "minisuperspace", "protocol", "fake_test.go"),
		[]byte("package protocol\n\nimport \"testing\"\n\nfunc TestFakeOK(t *testing.T) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRunVerificationRoundTrip(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go toolchain unavailable")
	}
	root := t.TempDir()
	ws, err := resolveWorkspace(root)
	if err != nil {
		t.Fatal(err)
	}
	fakePhysLibModule(t, ws)
	bin := buildServerBinary(t)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	client := mcp.NewClient(&mcp.Implementation{Name: "p2-rt", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, &mcp.CommandTransport{Command: exec.Command(bin, "--workspace", root)}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer cs.Close()
	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "run_verification", Arguments: map[string]any{"profile": "demo-arc"}})
	if err != nil {
		t.Fatalf("run_verification call failed: %v", err)
	}
	if res.IsError {
		t.Fatalf("run_verification returned error result: %v", res.Content)
	}
	res2, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "run_verification", Arguments: map[string]any{"profile": "rm -rf /"}})
	if err == nil && !res2.IsError {
		t.Fatal("injection profile accepted")
	}
}
