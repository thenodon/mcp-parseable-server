package tools

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetDataStreamStatsTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"get_data_stream_stats",
		mcp.WithDescription(`Get comprehensive statistics for a Parseable data stream, including ingestion and storage metrics.
Use this to monitor stream health, data ingestion rates, and storage usage. Calls /api/v1/logstream/<streamName>/stats.

Returns a JSON object with the following structure:

- stream: data stream name
- time: stats timestamp (ISO 8601 format)
- ingestion: object with ingestion metrics
    - count: number of records ingested in the current period
    - size: total bytes ingested in the current period
    - format: data format (e.g., "json")
    - lifetime_count: cumulative total records ever ingested
    - lifetime_size: cumulative total bytes ever ingested
    - deleted_count: number of deleted records
    - deleted_size: bytes freed from deleted records
- storage: object with storage metrics
    - size: total bytes currently stored
    - format: storage format (e.g., "parquet")
    - lifetime_size: cumulative total bytes ever stored
    - deleted_size: bytes freed from storage deletions

All metrics are calculated at the time of the API call and reflect the current state of the stream.
`),
		mcp.WithString("streamName", mcp.Required(), mcp.Description("Name of the data stream to get stats for. Example: 'otellogs' or 'monitor_logstream'. Stream must exist in Parseable.")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		streamName := mcp.ParseString(req, "streamName", "")
		if streamName == "" {
			slog.Warn("called with missing parameter", "parameter", "streamName", "tool", "get_data_stream_stats")
			return mcp.NewToolResultError("missing required field: streamName"), nil
		}

		stats, err := getParseableStats(streamName)
		if err != nil {
			slog.Error("failed to get response", "streamName", streamName, "error", err, "tool", "get_data_stream_stats")
			return mcp.NewToolResultError("failed to get stats: " + err.Error()), nil
		}

		return mcp.NewToolResultJSON(stats)
	})
}
