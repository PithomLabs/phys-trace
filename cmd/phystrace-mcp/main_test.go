package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// buildServerBinary compiles the adapter into a temp dir for subprocess tests.
func buildServerBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "phystrace-mcp")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir, _ = os.Getwd()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build server: %v\n%s", err, out)
	}
	return bin
}

// TestServerToolListPinsFourTools connects over stdio and asserts the
// exact P2 tool surface. No other tools may be registered without a
// scope decision.
func TestServerToolListPinsFourTools(t *testing.T) {
	bin := buildServerBinary(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client := mcp.NewClient(&mcp.Implementation{Name: "p1-test", Version: "0.0.0"}, nil)
	cs, err := client.Connect(ctx, &mcp.CommandTransport{Command: exec.Command(bin)}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer cs.Close()
	res, err := cs.ListTools(ctx, &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("tools/list: %v", err)
	}
	got := map[string]bool{}
	for _, tool := range res.Tools {
		got[tool.Name] = true
	}
	for _, want := range []string{"case_overview", "inspect_evidence", "repository_status", "run_verification"} {
		if !got[want] {
			t.Fatalf("tool %q missing from registry", want)
		}
	}
	if len(res.Tools) != 4 {
		t.Fatalf("registry holds %d tools, want exactly 4", len(res.Tools))
	}
}

// TestStdoutCleanAtStartup spawns the server with stdout captured and
// stdin held open: no banner, log line, or diagnostic may appear on
// stdout before any session exists. Mid-session stdout purity rests on
// construction (single stderr logger; the SDK owns protocol writes) and
// is a code-review invariant, not covered by this spawn test.
func TestStdoutCleanAtStartup(t *testing.T) {
	bin := buildServerBinary(t)
	cmd := exec.Command(bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	if stdout.Len() != 0 {
		stdin.Close()
		cmd.Wait()
		t.Fatalf("server wrote %d unexpected stdout bytes at startup: %q", stdout.Len(), stdout.String())
	}
	// Closing stdin ends the stdio session; the server must exit cleanly.
	stdin.Close()
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("server exit after stdin EOF: %v (stderr: %q)", err, stderr.String())
		}
	case <-time.After(15 * time.Second):
		cmd.Process.Kill()
		t.Fatal("server did not exit after stdin EOF")
	}
	if stdout.Len() != 0 {
		t.Fatalf("server wrote %d stdout bytes over its lifetime: %q", stdout.Len(), stdout.String())
	}
}
