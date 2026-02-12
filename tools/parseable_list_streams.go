package tools

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterListDataStreamsTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"list_data_streams",
		mcp.WithDescription("List all available data streams in Parseable. "+
			"Use this tool to discover which data streams are available before executing queries. "+
			"Each stream is a table-like collection of data and must be referenced by exact name in query_data_stream operations. "+
			"Returns a JSON object with a 'streams' array containing stream names as strings. "+
			"All returned streams are accessible and queryable by the current user."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		slog.Info("listing all data streams")
		streams, err := listParseableStreams()
		if err != nil {
			slog.Error("failed to get response", "error", err, "tool", "list_data_streams")
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultJSON(map[string]interface{}{"streams": streams})
	})
}
