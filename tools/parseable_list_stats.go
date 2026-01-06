package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetDataStreamStatsTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"get_data_stream_stats",
		mcp.WithDescription(`Get stats for a Parseable data stream by name. Calls /api/v1/logstream/<streamName>/stats.

Returned fields:
- stream: data stream name
- time: stats timestamp (ISO 8601)
- ingestion: object with ingestion stats
    - count: number of ingested records
    - size: total bytes ingested
    - format: data format (e.g. json)
    - lifetime_count: cumulative ingested records
    - lifetime_size: cumulative ingested bytes
    - deleted_count: number of deleted records
    - deleted_size: bytes deleted
- storage: object with storage stats
    - size: total bytes stored
    - format: storage format (e.g. parquet)
    - lifetime_size: cumulative stored bytes
    - deleted_size: bytes deleted from storage
`),
		mcp.WithString("streamName", mcp.Required(), mcp.Description("Name of the data stream (e.g. otellogs)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		streamName := mcp.ParseString(req, "streamName", "")
		if streamName == "" {
			return mcp.NewToolResultError("missing required field: streamName"), nil
		}
		stats, err := getParseableStats(streamName)
		if err != nil {
			return mcp.NewToolResultError("failed to get stats: " + err.Error()), nil
		}
		fieldDescriptions := map[string]interface{}{
			"stream": "Data stream name",
			"time":   "Stats timestamp (ISO 8601)",
			"ingestion": map[string]string{
				"count":          "Number of ingested records",
				"size":           "Total bytes ingested",
				"format":         "Data format (e.g. json)",
				"lifetime_count": "Cumulative ingested records",
				"lifetime_size":  "Cumulative ingested bytes",
				"deleted_count":  "Number of deleted records",
				"deleted_size":   "Bytes deleted",
			},
			"storage": map[string]string{
				"size":          "Total bytes stored",
				"format":        "Storage format (e.g. parquet)",
				"lifetime_size": "Cumulative stored bytes",
				"deleted_size":  "Bytes deleted from storage",
			},
		}
		return mcp.NewToolResultStructured(map[string]interface{}{
			"stats":             stats,
			"fieldDescriptions": fieldDescriptions,
		}, "Stats returned"), nil
	})
}
