package tools

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterListDataStreamsTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"list_data_streams",
		mcp.WithDescription("List all available data streams in Parseable. A stream is equivalent to a table in the SQL query."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		streams, err := listParseableStreams()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		//return mcp.NewToolResultStructured(map[string]interface{}{"streams": streams}, "Streams listed"), nil
		return mcp.NewToolResultText(strings.Join(streams, "\n")), nil

	})
}
