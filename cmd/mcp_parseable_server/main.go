package main

import (
	"flag"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"mcp-pb/tools"
)

func main() {
	mode := flag.String("mode", "http", "server mode: http or stdio")
	parseableURL := flag.String("parseable-url", "", "base URL for Parseable instance (or set PARSEABLE_URL env var)")
	parseableUserFlag := flag.String("parseable-username", "", "Parseable basic auth username (or set PARSEABLE_USER env var)")
	parseablePassFlag := flag.String("parseable-password", "", "Parseable basic auth password (or set PARSEABLE_PASS env var)")
	listenAddr := flag.String("listen", ":9034", "address to listen on")
	flag.Parse()

	// Prefer environment variables if set, otherwise use flags or defaults
	tools.ParseableBaseURL = os.Getenv("PARSEABLE_URL")
	if tools.ParseableBaseURL == "" {
		if *parseableURL != "" {
			tools.ParseableBaseURL = *parseableURL
		} else {
			tools.ParseableBaseURL = "http://localhost:8000"
		}
	}
	tools.ParseableUser = os.Getenv("PARSEABLE_USER")
	if tools.ParseableUser == "" {
		if *parseableUserFlag != "" {
			tools.ParseableUser = *parseableUserFlag
		} else {
			tools.ParseableUser = "admin"
		}
	}
	tools.ParseablePass = os.Getenv("PARSEABLE_PASS")
	if tools.ParseablePass == "" {
		if *parseablePassFlag != "" {
			tools.ParseablePass = *parseablePassFlag
		} else {
			tools.ParseablePass = "admin"
		}
	}

	listenAddrEnv := os.Getenv("LISTEN_ADDR")
	if listenAddrEnv != "" {
		*listenAddr = listenAddrEnv
	}

	mcpServer := server.NewMCPServer("parseable-mcp", "0.1.0",
		server.WithRecovery(),
		server.WithLogging(),
		server.WithInstructions(`
You are Virtual Assistant, a tool for interacting with Parseable API and documentation in different tasks related to monitoring and observability.

You have many tools to get data from Parseable, but try to specify the query as accurately as possible.

Try not to second guess information - if you don't know something or lack information, it's better to ask.
	`),
	)

	tools.RegisterParseableTools(mcpServer)

	if *mode == "stdio" {
		log.Printf("MCP server running in stdio mode (Parseable at %s)", tools.ParseableBaseURL)
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("MCP stdio server failed: %v", err)
		}
		return
	}

	httpServer := server.NewStreamableHTTPServer(mcpServer)
	log.Printf("MCP server running on %s, Parseable at %s", *listenAddr, tools.ParseableBaseURL)
	if err := httpServer.Start(*listenAddr); err != nil {
		log.Fatalf("MCP server failed: %v", err)
	}
}
