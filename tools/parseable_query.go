package tools

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterQueryDataStreamTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"query_data_stream",
		mcp.WithDescription("Execute a SQL query against a data stream in Parseable and retrieve rows of data. "+
			"All parameters are required. "+
			"Supported SQL operations: SELECT with column selection, WHERE conditions (but not time-based), GROUP BY, ORDER BY, LIMIT, and aggregate functions (COUNT, SUM, AVG, MIN, MAX). "+
			"Time filtering is handled by startTime and endTime parameters - do not include time conditions in the WHERE clause. "+
			"The FROM clause table name must exactly match the streamName parameter. "+
			"Returns a JSON object with 'rows' (array of data objects) and 'count' (number of rows returned)."),
		mcp.WithString("query", mcp.Required(), mcp.Description("SQL query to execute. FROM clause table must exactly match the streamName parameter. Example: 'SELECT field1, field2 FROM streamName WHERE field1 > 100 ORDER BY timestamp DESC LIMIT 100'")),
		mcp.WithString("streamName", mcp.Required(), mcp.Description("Exact name of the data stream (table) to query. Must match the table name in the FROM clause. Example: 'monitor_logstream'")),
		mcp.WithString("startTime", mcp.Required(), mcp.Description("Query start time in ISO 8601 format with timezone. Examples: '2026-02-12T00:00:00Z' or '2026-02-12T00:00:00+00:00'")),
		mcp.WithString("endTime", mcp.Required(), mcp.Description("Query end time in ISO 8601 format with timezone. Must be after startTime. Examples: '2026-02-12T23:59:59Z' or '2026-02-12T23:59:59+00:00'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := mcp.ParseString(req, "query", "")
		streamName := mcp.ParseString(req, "streamName", "")
		startTime := mcp.ParseString(req, "startTime", "")
		endTime := mcp.ParseString(req, "endTime", "")
		if query == "" || streamName == "" || startTime == "" || endTime == "" {
			slog.Warn("called with missing parameter",
				"query", query,
				"streamName", streamName,
				"startTime", startTime,
				"endTime", endTime)
			return mcp.NewToolResultError("missing required fields: query, streamName, startTime, and endTime are required"), nil
		}

		slog.Debug("executing query_data_stream",
			"streamName", streamName,
			"startTime", startTime,
			"endTime", endTime,
			"query", query)

		queryResult, err := doParseableQuery(query, streamName, startTime, endTime)
		if err != nil {
			slog.Error("failed to get response",
				"streamName", streamName,
				"error", err, "tool", "query_data_stream", "query", query)
			return mcp.NewToolResultError("query failed: " + err.Error()), nil
		}

		slog.Debug("query_data_stream completed successfully",
			"streamName", streamName,
			"rowCount", len(queryResult))

		return mcp.NewToolResultJSON(map[string]interface{}{
			"rows":  queryResult,
			"count": len(queryResult),
		})
	})
}
