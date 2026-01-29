package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetAboutTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"get_about",
		mcp.WithDescription(`Get information about the Parseable instance. Calls /api/v1/about.

Returned fields:
- version: version of Parseable
- uiVersion: the UI version of Parseable
- commit: the git commit hash
- deploymentId: the deployment ID of Parseable
- updateAvailable: if updates of Parseable is available
- latestVersion: the latest version of Parseable
- llmActive: if the Parseable is configured with LLM support
- llmProvider: what LLM provider is used
- oidcActive: if Parseable is configured with OpenID Connect support
- license: the license of Parseable
- mode: if Parseable is running as Standalone or Cluster mode
- staging: the staging path
- hotTier: if hot tier is enabled or disabled
- grpcPort: the grpc port of Parseable
- store: the storage type used for Parseable like local or object store
- analytics: if analytics is enabled or disabled
`),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		about, err := getParseableAbout()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var lines []string
		for k, v := range about {
			lines = append(lines, k+": "+fmt.Sprintf("%v", v))
		}
		return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
		// Optionally, for structured output:
		// return mcp.NewToolResultStructured(map[string]interface{}{"info": info}, "Info returned"), nil
	})
}
