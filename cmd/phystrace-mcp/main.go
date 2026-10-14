// Command phystrace-mcp is the PhysTrace MCP server entry point (P0).
//
// It serves the MCP JSON-RPC protocol over stdio for agent-agnostic access
// to PhysTrace review and verification workflows. P0 establishes transport
// and initialization only: no domain tools are registered yet. Tool
// behavior (case_overview, inspect_evidence, run_verification,
// repository_status) arrives in later phases; nothing here fabricates
// tool results.
//
// Transport discipline: standard output carries MCP protocol frames
// exclusively. All logging, diagnostics, and recovered panic details go
// to standard error.
package main

import (
	"context"
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
	if err := run(stderr); err != nil {
		stderr.Printf("server exited: %v", err)
		os.Exit(1)
	}
}

func run(stderr *log.Logger) (err error) {
	defer func() {
		// Recovered panic details go to stderr, never stdout.
		if r := recover(); r != nil {
			stderr.Printf("recovered panic: %v", r)
			err = fmt.Errorf("phystrace-mcp panic: %v", r)
		}
	}()
	server := mcp.NewServer(&mcp.Implementation{Name: serverName, Version: serverVersion}, nil)
	// P0 registers no tools: initialization + transport smoke only.
	return server.Run(context.Background(), &mcp.StdioTransport{})
}
