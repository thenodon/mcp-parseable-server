package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetDataStreamInfoTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"get_data_stream_info",
		mcp.WithDescription(`Get info for a Parseable data stream by name. Calls /api/v1/logstream/<streamName>/info.

Returned fields:
- createdAt: when the data stream was created (ISO 8601)
- firstEventAt: timestamp of the first event (ISO 8601)
- latestEventAt: timestamp of the latest event (ISO 8601)
- streamType: type of data stream (e.g. UserDefined)
- logSource: array of log source objects
    - log_source_format: format of the log source (e.g. otel-logs)
    - fields: list of field names in the log source
- telemetryType: type of telemetry (e.g. logs, metrics, traces)
`),
		mcp.WithString("streamName", mcp.Required(), mcp.Description("Name of the data stream (e.g. otellogs)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		streamName := mcp.ParseString(req, "streamName", "")
		if streamName == "" {
			return mcp.NewToolResultError("missing required field: streamName"), nil
		}
		info, err := getParseableInfo(streamName)
		if err != nil {
			return mcp.NewToolResultError("failed to get info: " + err.Error()), nil
		}
		// Default: return as text
		var lines []string
		for k, v := range info {
			lines = append(lines, k+": "+fmt.Sprintf("%v", v))
		}
		return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
		// Optionally, for structured output:
		// return mcp.NewToolResultStructured(map[string]interface{}{"info": info}, "Info returned"), nil
	})
}
