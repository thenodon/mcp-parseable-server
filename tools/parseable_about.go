package tools

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetAboutTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"get_about",
		mcp.WithDescription(`Get configuration and version information about the Parseable instance.
Use this to understand the Parseable deployment, available features, and version compatibility.
Calls /api/v1/about.

Returns a JSON object with deployment and configuration details:

Deployment Info:
- version: semantic version of Parseable (e.g., "1.2.0")
- uiVersion: version of the web UI
- commit: git commit hash of the Parseable build
- deploymentId: unique identifier for this Parseable instance
- mode: deployment mode ("Standalone" or "Cluster")

Update Information:
- updateAvailable: boolean indicating if a newer version is available
- latestVersion: the latest available version of Parseable

Configuration & Features:
- llmActive: boolean indicating if LLM support is enabled
- llmProvider: name of the LLM provider if configured (e.g., "openai", "anthropic", or null)
- oidcActive: boolean indicating if OpenID Connect authentication is enabled
- analytics: boolean indicating if analytics collection is enabled
- hotTier: boolean indicating if hot tier (fast storage) is enabled

Storage & Infrastructure:
- store: storage backend type (e.g., "local", "s3", "gcs")
- staging: staging/cache path for data processing
- grpcPort: port number for gRPC API connections

Licensing:
- license: license type or status of Parseable

Use this tool to check Parseable capabilities, version information, and configuration state.
`),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		about, err := getParseableAbout()
		if err != nil {
			slog.Error("failed to get response", "tool", "get_about", "error", err)
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultJSON(about)
	})
}
