package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// workspace is an absolute, cleaned demo-workspace root. All tool file
// access stays inside it; anything else is rejected, never resolved.
type workspace struct {
	root string
}

// resolveWorkspace validates a workspace root for tool use.
func resolveWorkspace(path string) (workspace, error) {
	if path == "" {
		return workspace{}, fmt.Errorf("phystrace-mcp: workspace path is empty")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return workspace{}, fmt.Errorf("phystrace-mcp: workspace path: %w", err)
	}
	abs = filepath.Clean(abs)
	st, err := os.Stat(abs)
	if err != nil {
		return workspace{}, fmt.Errorf("phystrace-mcp: workspace: %w", err)
	}
	if !st.IsDir() {
		return workspace{}, fmt.Errorf("phystrace-mcp: workspace is not a directory: %s", abs)
	}
	return workspace{root: abs}, nil
}

// workspaceFile resolves a repo-relative path inside the workspace,
// rejecting absolute paths and any traversal outside the root.
func (ws workspace) workspaceFile(elem ...string) (string, error) {
	joined := filepath.Join(append([]string{ws.root}, elem...)...)
	clean := filepath.Clean(joined)
	if clean != ws.root && !strings.HasPrefix(clean, ws.root+string(os.PathSeparator)) {
		return "", fmt.Errorf("phystrace-mcp: path escapes workspace: %q", filepath.Join(elem...))
	}
	return clean, nil
}

// readWorkspaceFile reads a bounded repo-relative file.
func (ws workspace) readWorkspaceFile(elem ...string) ([]byte, error) {
	p, err := ws.workspaceFile(elem...)
	if err != nil {
		return nil, err
	}
	const maxEvidenceBytes = 1 << 20
	raw, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("phystrace-mcp: read %s: %w", filepath.Join(elem...), err)
	}
	if len(raw) > maxEvidenceBytes {
		return nil, fmt.Errorf("phystrace-mcp: %s exceeds evidence size bound", filepath.Join(elem...))
	}
	return raw, nil
}

// Pinned demo inventory, verified by TestDerivationInventory in
// phys-lib/bohm/minisuperspace/protocol at phys-lib v0.1.0. These are
// declared values with stated provenance, not computed claims; the P1
// golden test pins them so drift forces a conscious update.
const (
	inventoryUnits     = 12
	inventoryArtifacts = 14
	inventoryVerified  = "TestDerivationInventory (phys-lib/bohm/minisuperspace/protocol)"
)

var inventoryAntiGoldens = []string{
	"TestOperatorCompositionVsSquare",
	"TestWValue",
	"TestHistoricalWrongPath",
}

// manifestDoc is the minimal manifest.json projection the tools read.
type manifestDoc struct {
	Name   string `json:"name"`
	Source struct {
		ArxivID   string `json:"arxiv_id"`
		ScopeNote string `json:"scope_note"`
	} `json:"source"`
	Regime      string          `json:"regime"`
	Premises    []manifestEntry `json:"premises"`
	Limitations []manifestEntry `json:"limitations"`
}

type manifestEntry struct {
	ID        string `json:"id"`
	Statement string `json:"statement"`
}

// manifestRelPath is the manifest location inside a prepared workspace.
var manifestRelPath = []string{"phys-lib", "bohm", "minisuperspace", "manifest.json"}

func readManifest(ws workspace) (manifestDoc, error) {
	var doc manifestDoc
	raw, err := ws.readWorkspaceFile(manifestRelPath...)
	if err != nil {
		return manifestDoc{}, err
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return manifestDoc{}, fmt.Errorf("phystrace-mcp: manifest undecodable: %w", err)
	}
	if doc.Name == "" {
		return manifestDoc{}, fmt.Errorf("phystrace-mcp: manifest has no name")
	}
	return doc, nil
}

// CaseOverview is the case_overview tool output.
type CaseOverview struct {
	Paper           string   `json:"paper"`
	Scope           string   `json:"scope"`
	Regime          string   `json:"regime"`
	DerivationUnits int      `json:"derivation_units"`
	Artifacts       int      `json:"artifacts"`
	AntiGoldens     []string `json:"anti_goldens"`
	Limitations     []string `json:"limitations"`
	VerifiedBy      string   `json:"verified_by"`
}

func caseOverview(ws workspace) (CaseOverview, error) {
	doc, err := readManifest(ws)
	if err != nil {
		return CaseOverview{}, err
	}
	limits := make([]string, 0, len(doc.Limitations))
	for _, l := range doc.Limitations {
		limits = append(limits, l.ID)
	}
	return CaseOverview{
		Paper:           doc.Name,
		Scope:           doc.Source.ScopeNote,
		Regime:          doc.Regime,
		DerivationUnits: inventoryUnits,
		Artifacts:       inventoryArtifacts,
		AntiGoldens:     append([]string{}, inventoryAntiGoldens...),
		Limitations:     limits,
		VerifiedBy:      inventoryVerified,
	}, nil
}

// EvidenceDoc is one inspect_evidence document: curated text with exact
// source locations, or raw file content for the "file" selector.
type EvidenceDoc struct {
	Selector string   `json:"selector"`
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Sources  []string `json:"sources"`
}

// InspectEvidenceIn is the inspect_evidence input. Path is honored only
// for selector "file" and is confined to the workspace.
type InspectEvidenceIn struct {
	Selector string `json:"selector" jsonschema:"one of: manifest, inventory, q-relation, anti-goldens, covariant-fragment, limitations, file"`
	Path     string `json:"path,omitempty" jsonschema:"repo-relative file path, only for selector file"`
}

func inspectEvidence(ws workspace, in InspectEvidenceIn) (EvidenceDoc, error) {
	switch in.Selector {
	case "manifest":
		doc, err := readManifest(ws)
		if err != nil {
			return EvidenceDoc{}, err
		}
		raw, err := ws.readWorkspaceFile(manifestRelPath...)
		if err != nil {
			return EvidenceDoc{}, err
		}
		return EvidenceDoc{Selector: "manifest", Title: doc.Name, Body: string(raw),
			Sources: []string{filepath.Join(manifestRelPath...)}}, nil
	case "inventory":
		return EvidenceDoc{Selector: "inventory", Title: "Derivation-unit inventory",
			Body: "12 mathematical units across 14 replayed derivation artifacts " +
				"(13 Derivs appends plus separately stored FRWDerivRef; the " +
				"three-step Ada composition counts as one operator unit). " +
				"Anti-goldens: TestOperatorCompositionVsSquare, TestWValue, " +
				"TestHistoricalWrongPath.",
			Sources: []string{
				"phys-lib/bohm/minisuperspace/fixture.go",
				"phys-lib/bohm/minisuperspace/protocol/derivation_test.go:TestDerivationInventory",
			}}, nil
	case "q-relation":
		return EvidenceDoc{Selector: "q-relation", Title: "Q = a^3 rhoQ defining relation",
			Body: "QRef is Mul(a^3, RhoQRef); QRelRef is a stored relation/1 " +
				"EQ(Q, a^3 rhoQ). Fixture.Verify decodes and validates the " +
				"stored relation against the actual Q and density artifacts " +
				"(byte-identical left side; canonical right side with `a` " +
				"resolved through the amplitude declaring function). Pressure " +
				"pQ differentiates Q, never rhoQ directly.",
			Sources: []string{
				"phys-lib/bohm/minisuperspace/fixture.go:checkQRelation",
				"phys-lib/bohm/minisuperspace/protocol/derivation_test.go:TestQRelation",
			}}, nil
	case "anti-goldens":
		return EvidenceDoc{Selector: "anti-goldens", Title: "Rejection rows",
			Body: "TestOperatorCompositionVsSquare rejects the a^2 d^2_a shape " +
				"for (a d_a)^2; TestWValue rejects a collapsed W shape; " +
				"TestHistoricalWrongPath builds pressure from rhoQ directly " +
				"and requires it to differ from the Q-derived pressure.",
			Sources: []string{
				"phys-lib/bohm/minisuperspace/protocol/derivation_test.go",
			}}, nil
	case "covariant-fragment":
		return EvidenceDoc{Selector: "covariant-fragment", Title: "FRW covariant fragment",
			Body: "Flat-FRW-chart covector fragment built in the fixture store " +
				"(no side stores): general inverse, vanishing Levi-Civita " +
				"coefficients, covariant_evaluation/1 certificate. Verified " +
				"by Fixture.Verify via verifyCovariant (resolution, replay, " +
				"recomputation, tensor structure, Kronecker goldens).",
			Sources: []string{
				"phys-lib/bohm/minisuperspace/fixture.go:buildFRWCovariant",
				"phys-lib/bohm/minisuperspace/fixture.go:verifyCovariant",
				"phys-lib/bohm/minisuperspace/frw_covariant_test.go",
			}}, nil
	case "limitations":
		doc, err := readManifest(ws)
		if err != nil {
			return EvidenceDoc{}, err
		}
		var b strings.Builder
		for _, l := range doc.Limitations {
			b.WriteString("- " + l.ID + ": " + l.Statement + "\n")
		}
		return EvidenceDoc{Selector: "limitations", Title: "Unsupported scope",
			Body:    b.String(),
			Sources: []string{filepath.Join(manifestRelPath...)}}, nil
	case "file":
		if in.Path == "" {
			return EvidenceDoc{}, fmt.Errorf("phystrace-mcp: file selector needs a path")
		}
		raw, err := ws.readWorkspaceFile(in.Path)
		if err != nil {
			return EvidenceDoc{}, err
		}
		return EvidenceDoc{Selector: "file", Title: in.Path, Body: string(raw),
			Sources: []string{in.Path}}, nil
	default:
		return EvidenceDoc{}, fmt.Errorf("phystrace-mcp: unknown evidence selector %q", in.Selector)
	}
}

// RepoStatus is one member repository's reported state. Dirty worktrees
// are reported, never cleaned.
type RepoStatus struct {
	Name    string `json:"name"`
	Present bool   `json:"present"`
	Remote  string `json:"remote,omitempty"`
	Branch  string `json:"branch,omitempty"`
	HEAD    string `json:"head,omitempty"`
	Tags    string `json:"tags,omitempty"`
	Clean   bool   `json:"clean"`
	Status  string `json:"status,omitempty"`
	Error   string `json:"error,omitempty"`
}

// RepositoryStatusOut is the repository_status tool output.
type RepositoryStatusOut struct {
	Workspace string       `json:"workspace"`
	Repos     []RepoStatus `json:"repos"`
}

// workspaceMembers is the fixed demo-workspace membership. Frozen phys
// is deliberately absent: nothing the demo builds depends on it.
var workspaceMembers = []string{"phys-artifact", "phys-math", "phys-lib", "phys-trace"}

// gitTimeout bounds every git subprocess; output is truncated to a fixed cap.
const (
	gitTimeoutSeconds = 15
	maxStatusLines    = 20
)

func gitStatus(ws workspace) (RepositoryStatusOut, error) {
	out := RepositoryStatusOut{Workspace: ws.root}
	for _, name := range workspaceMembers {
		out.Repos = append(out.Repos, repoStatus(ws, name))
	}
	return out, nil
}

func gitRun(dir string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gitTimeoutSeconds*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	raw, err := cmd.Output()
	if err != nil {
		return "", err
	}
	const maxGitBytes = 8192
	if len(raw) > maxGitBytes {
		raw = raw[:maxGitBytes]
	}
	return string(raw), nil
}

func repoStatus(ws workspace, name string) RepoStatus {
	st := RepoStatus{Name: name}
	dir, err := ws.workspaceFile(name)
	if err != nil {
		st.Error = err.Error()
		return st
	}
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		st.Error = fmt.Sprintf("phystrace-mcp: member %s absent", name)
		return st
	}
	st.Present = true
	// Absent origin is reported, not fatal: the remaining fields are
	// still meaningful for a local checkout.
	if remote, err := gitRun(dir, "remote", "get-url", "origin"); err == nil {
		st.Remote = strings.TrimSpace(remote)
	}
	if branch, err := gitRun(dir, "rev-parse", "--abbrev-ref", "HEAD"); err == nil {
		st.Branch = strings.TrimSpace(branch)
	} else {
		st.Error = "phystrace-mcp: HEAD unresolvable"
		return st
	}
	if head, err := gitRun(dir, "rev-parse", "HEAD"); err == nil {
		st.HEAD = strings.TrimSpace(head)
	}
	if tags, err := gitRun(dir, "tag", "--points-at", "HEAD"); err == nil {
		st.Tags = strings.TrimSpace(tags)
	}
	porcelain, err := gitRun(dir, "status", "--porcelain")
	if err != nil {
		st.Error = "phystrace-mcp: status failed"
		return st
	}
	lines := strings.Split(strings.TrimRight(porcelain, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		st.Clean = true
		return st
	}
	if len(lines) > maxStatusLines {
		lines = append(lines[:maxStatusLines], "... (truncated)")
	}
	st.Status = strings.Join(lines, "\n")
	return st
}

// newServer builds the P1 server: exactly case_overview,
// inspect_evidence, repository_status. run_verification arrives in P2;
// the tool list is pinned by test so placeholders cannot sneak in.
func newServer(ws workspace) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: serverName, Version: serverVersion}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "case_overview",
		Description: "Summarize the Durr-Struyve minisuperspace case: paper, scope, derivation counts, limitations, and what verified them.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, CaseOverview, error) {
		out, err := caseOverview(ws)
		if err != nil {
			return nil, CaseOverview{}, err
		}
		return &mcp.CallToolResult{}, out, nil
	})
	mcp.AddTool(server, &mcp.Tool{
		Name:        "inspect_evidence",
		Description: "Inspect curated demo evidence (manifest, inventory, q-relation, anti-goldens, covariant-fragment, limitations) or read a workspace-confined file via selector file.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, in InspectEvidenceIn) (*mcp.CallToolResult, EvidenceDoc, error) {
		out, err := inspectEvidence(ws, in)
		if err != nil {
			return nil, EvidenceDoc{}, err
		}
		return &mcp.CallToolResult{}, out, nil
	})
	mcp.AddTool(server, &mcp.Tool{
		Name:        "repository_status",
		Description: "Report each workspace member's remote, branch, HEAD, tags, and cleanliness. Read-only: dirty trees are reported, never cleaned.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, RepositoryStatusOut, error) {
		out, err := gitStatus(ws)
		if err != nil {
			return nil, RepositoryStatusOut{}, err
		}
		return &mcp.CallToolResult{}, out, nil
	})
	return server
}
