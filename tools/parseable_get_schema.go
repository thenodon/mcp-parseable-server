package tools

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetDataStreamSchemaTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"get_data_stream_schema",
		mcp.WithDescription(`Get the complete field schema for a Parseable data stream.
Use this to discover field names, data types, and structure before constructing queries.
Calls /api/v1/logstream/<streamName>/schema.

Returns a JSON object with a 'fields' array containing field definitions for each available field in the stream.
Each field includes:
- name: the field name (string)
- data_type: the data type of the field (e.g., "String", "i64", "f64", "bool", "DateTime")

Use this tool to understand what fields are available for filtering, grouping, or selecting in query_data_stream operations.
`),
		mcp.WithString("streamName", mcp.Required(), mcp.Description("Name of the data stream to get the schema for. Example: 'otellogs' or 'monitor_logstream'. Stream must exist in Parseable.")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		stream := mcp.ParseString(req, "streamName", "")
		if stream == "" {
			slog.Warn("Missing parameter", "tool", "get_data_stream_schema", "parameter", "streamName")
			return mcp.NewToolResultError("missing required field: streamName"), nil
		}

		schema, err := getParseableSchema(stream)
		if err != nil {
			slog.Error("failed to get response", "tool", "get_data_stream_schema", "streamName", stream, "error", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultJSON(schema)
	})
}
