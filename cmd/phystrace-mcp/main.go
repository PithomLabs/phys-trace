// Command phystrace-mcp is the PhysTrace MCP server entry point.
//
// It serves the MCP JSON-RPC protocol over stdio for agent-agnostic access
// to PhysTrace review and verification workflows. P1 registers the three
// read-only inspection tools (case_overview, inspect_evidence,
// repository_status); run_verification arrives in P2. Nothing here
// fabricates tool results.
//
// Transport discipline: standard output carries MCP protocol frames
// exclusively. All logging, diagnostics, and recovered panic details go
// to standard error.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server identity reported during MCP initialization.
const (
	serverName    = "phystrace-mcp"
	serverVersion = "0.1.0"
)

func main() {
	// Stderr-only logger: nothing in this process may write to stdout
	// except the SDK's protocol transport.
	stderr := log.New(os.Stderr, "phystrace-mcp: ", log.LstdFlags)
	var workspacePath string
	flag.StringVar(&workspacePath, "workspace", "", "demo-workspace root (default: current directory)")
	flag.Parse()
	if workspacePath == "" {
		if env := os.Getenv("PHYSTRACE_WORKSPACE"); env != "" {
			workspacePath = env
		} else {
			workspacePath = "."
		}
	}
	if err := run(stderr, workspacePath); err != nil {
		stderr.Printf("server exited: %v", err)
		os.Exit(1)
	}
}

func run(stderr *log.Logger, workspacePath string) (err error) {
	defer func() {
		// Recovered panic details go to stderr, never stdout.
		if r := recover(); r != nil {
			stderr.Printf("recovered panic: %v", r)
			err = fmt.Errorf("phystrace-mcp panic: %v", r)
		}
	}()
	ws, err := resolveWorkspace(workspacePath)
	if err != nil {
		return fmt.Errorf("phystrace-mcp: %w", err)
	}
	stderr.Printf("workspace: %s", ws.root)
	server := newServer(ws)
	// P1 registers the three read-only inspection tools. run_verification
	// arrives in P2; the tool list is pinned by test.
	return server.Run(context.Background(), &mcp.StdioTransport{})
}
