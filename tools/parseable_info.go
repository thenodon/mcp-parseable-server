package tools

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetDataStreamInfoTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"get_data_stream_info",
		mcp.WithDescription(`Get comprehensive metadata information for a Parseable data stream.
Use this to understand stream composition, available fields, and data ingestion timeline.
Calls /api/v1/logstream/<streamName>/info.

Returns a JSON object with the following structure:

- createdAt: ISO 8601 timestamp when the data stream was created
- firstEventAt: ISO 8601 timestamp of the first event (null if stream has no events)
- latestEventAt: ISO 8601 timestamp of the most recent event (null if stream has no events)
- streamType: classification of the stream (e.g., "UserDefined", "System")
- logSource: array of log source objects describing data sources
    - log_source_format: format of the ingested data (e.g., "otel-logs", "json", "logfmt")
    - fields: array of field names available in this data source
- telemetryType: category of telemetry data (e.g., "logs", "metrics", "traces")

Use this tool before querying a stream to understand its fields and structure.
`),
		mcp.WithString("streamName", mcp.Required(), mcp.Description("Name of the data stream to get info for. Example: 'otellogs' or 'monitor_logstream'. Stream must exist in Parseable.")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		streamName := mcp.ParseString(req, "streamName", "")
		if streamName == "" {
			slog.Warn("called with missing parameter", "parameter", "streamName", "tool", "get_data_stream_info")
			return mcp.NewToolResultError("missing required field: streamName"), nil
		}

		info, err := getParseableInfo(streamName)
		if err != nil {
			slog.Error("failed to get response", "streamName", streamName, "error", err, "tool", "get_data_stream_info")
			return mcp.NewToolResultError("failed to get info: " + err.Error()), nil
		}

		return mcp.NewToolResultJSON(info)
	})
}
