package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Verification profiles. Each maps to a fixed executable, argument
// vector, working directory, and timeout below. There is no other way
// to invoke execution: unknown profiles are rejected, and no
// caller-supplied value can change the approved command.
const (
	ProfileDemoArc       = "demo-arc"
	ProfileCoverageAudit = "coverage-audit"
	ProfileFullSuite     = "full-suite"
)

// demoArcFilter is the rehearsed 8-test filter from the demo script.
const demoArcFilter = "TestDerivationInventory|TestHistoricalWrongPath|TestOperatorCompositionVsSquare|TestWValue|TestQRelation|TestPressureFromQ|TestSubstitutionPoint|TestVerifyCoverageAudit"

// coverageAuditFilter runs the coverage audit plus the inventory guard.
const coverageAuditFilter = "TestVerifyCoverageAudit|TestDerivationInventory"

type profileSpec struct {
	// dir is the workspace-relative working directory (a fixed member).
	dir []string
	// argv is the fixed executable + arguments. No shell is involved.
	argv []string
	// timeout bounds the whole execution.
	timeout time.Duration
}

func profileSpecs() map[string]profileSpec {
	return map[string]profileSpec{
		ProfileDemoArc: {
			dir:     []string{"phys-lib"},
			argv:    []string{"go", "test", "-count=1", "-run", demoArcFilter, "./bohm/minisuperspace/protocol/"},
			timeout: 120 * time.Second,
		},
		ProfileCoverageAudit: {
			dir:     []string{"phys-lib"},
			argv:    []string{"go", "test", "-count=1", "-run", coverageAuditFilter, "./bohm/minisuperspace/protocol/"},
			timeout: 180 * time.Second,
		},
		ProfileFullSuite: {
			dir:     []string{"phys-lib"},
			argv:    []string{"go", "test", "-count=1", "./..."},
			timeout: 600 * time.Second,
		},
	}
}

// maxVerifyBytes caps each captured stream; output beyond the cap is
// dropped and reported via Truncated, never silently.
const maxVerifyBytes = 32 << 10

// VerificationResult is the run_verification tool output. Executing
// tests can have filesystem side effects (caches, test temp dirs), so
// this tool is an allowlisted verification executor, not a purely
// observational operation.
type VerificationResult struct {
	Profile   string `json:"profile"`
	Command   string `json:"command"`
	Dir       string `json:"dir"`
	ExitCode  int    `json:"exit_code"`
	TimedOut  bool   `json:"timed_out"`
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	Truncated bool   `json:"truncated"`
	Error     string `json:"error,omitempty"`
}

// RunVerificationIn is the run_verification tool input: a profile name
// only. Commands, paths, executables, and directories are never accepted
// from the caller.
type RunVerificationIn struct {
	Profile string `json:"profile" jsonschema:"one of: demo-arc, coverage-audit, full-suite"`
}

// cappedWriter keeps the tail of a stream up to a byte cap.
type cappedWriter struct {
	buf       bytes.Buffer
	cap       int
	truncated bool
}

func (w *cappedWriter) Write(p []byte) (int, error) {
	if w.buf.Len()+len(p) > w.cap {
		w.truncated = true
		keep := w.cap
		overflow := w.buf.Len() + len(p) - w.cap
		if overflow >= w.buf.Len() {
			w.buf.Reset()
			w.buf.Write(p[len(p)-keep:])
			return len(p), nil
		}
		w.buf.Next(overflow)
		w.buf.Write(p)
		return len(p), nil
	}
	return w.buf.Write(p)
}

// executeSpec runs one fixed spec. Startup failures are returned as
// errors; completed processes (including test failures) are reported via
// the result with their real exit status.
func executeSpec(ctx context.Context, dir string, spec profileSpec) (VerificationResult, error) {
	res := VerificationResult{
		Command: strings.Join(spec.argv, " "),
		Dir:     dir,
	}
	bin, err := exec.LookPath(spec.argv[0])
	if err != nil {
		return VerificationResult{}, fmt.Errorf("phystrace-mcp: executable %q unavailable: %w", spec.argv[0], err)
	}
	ctx, cancel := context.WithTimeout(ctx, spec.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, spec.argv[1:]...)
	cmd.Dir = dir
	// Inherit the environment: the Go toolchain requires PATH, HOME,
	// and module cache locations. Nothing is added or scrubbed; the
	// executable, arguments, and directory are fixed by the allowlist.
	cmd.Env = os.Environ()
	var stdout, stderr cappedWriter
	stdout.cap, stderr.cap = maxVerifyBytes, maxVerifyBytes
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err = cmd.Run()
	res.Stdout, res.Stderr = stdout.buf.String(), stderr.buf.String()
	res.Truncated = stdout.truncated || stderr.truncated
	if err == nil {
		return res, nil
	}
	if ctx.Err() == context.DeadlineExceeded {
		res.TimedOut = true
		res.ExitCode = -1
		return res, nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		res.ExitCode = exitErr.ExitCode()
		return res, nil
	}
	return VerificationResult{}, fmt.Errorf("phystrace-mcp: startup failure: %w", err)
}

// runVerification maps a profile name to its fixed spec and executes it.
// The working directory is resolved inside the workspace; profiles are
// the only accepted input.
func runVerification(ctx context.Context, ws workspace, profile string) (VerificationResult, error) {
	spec, ok := profileSpecs()[profile]
	if !ok {
		return VerificationResult{}, fmt.Errorf("phystrace-mcp: unknown verification profile %q", profile)
	}
	dir, err := ws.workspaceFile(spec.dir...)
	if err != nil {
		return VerificationResult{}, err
	}
	res, err := executeSpec(ctx, dir, spec)
	if err != nil {
		return VerificationResult{}, err
	}
	res.Profile = profile
	res.Dir = filepath.Join(spec.dir...)
	return res, nil
}

func registerVerificationTool(server *mcp.Server, ws workspace) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "run_verification",
		Description: "Execute one allowlisted verification profile (demo-arc, coverage-audit, full-suite) in the phys-lib module and report its real exit status. No other commands can be run through this tool.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, in RunVerificationIn) (*mcp.CallToolResult, VerificationResult, error) {
		out, err := runVerification(ctx, ws, in.Profile)
		if err != nil {
			return nil, VerificationResult{}, err
		}
		return &mcp.CallToolResult{}, out, nil
	})
}
