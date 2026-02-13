package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"mcp-pb/tools"
)

var version = "undefined"

func main() {
	mode := flag.String("mode", "http", "server mode: http or stdio")
	parseableURL := flag.String("parseable-url", "", "base URL for Parseable instance (or set PARSEABLE_URL env var)")
	parseableUserFlag := flag.String("parseable-username", "", "Parseable basic auth username (or set PARSEABLE_USER env var)")
	parseablePassFlag := flag.String("parseable-password", "", "Parseable basic auth password (or set PARSEABLE_PASS env var)")
	listenAddr := flag.String("listen", ":9034", "address to listen on")
	logLevel := flag.String("log-level", "info", "log level: debug, info, warn, error (or set LOG_LEVEL env var)")
	versionFlag := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	// Determine log level from environment variable or flag
	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr == "" {
		logLevelStr = *logLevel
	}

	// Parse log level string to slog.Level
	var level slog.Level
	switch logLevelStr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Setup structured logger for stdout
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	if *versionFlag {
		println("mcp-parseable-server " + version)
		os.Exit(0)
	}

	// Prefer environment variables if set, otherwise use flags or defaults
	tools.ParseableBaseURL = os.Getenv("PARSEABLE_URL")
	if tools.ParseableBaseURL == "" {
		if *parseableURL != "" {
			tools.ParseableBaseURL = *parseableURL
		} else {
			tools.ParseableBaseURL = "http://localhost:8000"
		}
	}
	tools.ParseableUser = os.Getenv("PARSEABLE_USERNAME")
	if tools.ParseableUser == "" {
		if *parseableUserFlag != "" {
			tools.ParseableUser = *parseableUserFlag
		} else {
			tools.ParseableUser = "admin"
		}
	}
	tools.ParseablePass = os.Getenv("PARSEABLE_PASSWORD")
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

	mcpServer := server.NewMCPServer("parseable-mcp", version,
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
		slog.Info("MCP server running in stdio mode", "parseable_url", tools.ParseableBaseURL)
		if err := server.ServeStdio(mcpServer); err != nil {
			slog.Error("MCP stdio server failed", "error", err)
			os.Exit(1)
		}
		return
	}

	httpServer := server.NewStreamableHTTPServer(mcpServer)
	slog.Info("MCP server running", "address", *listenAddr, "parseable_url", tools.ParseableBaseURL)
	if err := httpServer.Start(*listenAddr); err != nil {
		slog.Error("MCP server failed", "error", err)
		os.Exit(1)
	}
}
