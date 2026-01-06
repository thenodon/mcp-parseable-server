package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetDataStreamSchemaTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"get_data_stream_schema",
		mcp.WithDescription("Get the schema for a specific data stream in Parseable. The full content of the stream is typically in the 'body' field as a string."),
		mcp.WithString("stream", mcp.Required(), mcp.Description("Data stream name")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stream := mcp.ParseString(req, "stream", "")
		if stream == "" {
			return mcp.NewToolResultError("missing stream in context"), nil
		}
		schema, err := getParseableSchema(stream)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultStructured(map[string]interface{}{"schema": schema}, "Schema returned"), nil
	})
}
