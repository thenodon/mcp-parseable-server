package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterQueryDataStreamTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"query_data_stream",
		mcp.WithDescription("Execute a SQL query against a data stream in Parseable. All fields are required. All times must be in ISO 8601 format."),
		mcp.WithString("query", mcp.Required(), mcp.Description("SQL query to run, but the FROM must always be set to streamName")),
		mcp.WithString("streamName", mcp.Required(), mcp.Description("Name of the data stream (table)")),
		mcp.WithString("startTime", mcp.Required(), mcp.Description("Query start time in ISO 8601 (format: yyyy-MM-ddTHH:mm:ss+hh:mm)")),
		mcp.WithString("endTime", mcp.Required(), mcp.Description("Query end time in ISO 8601 (format: yyyy-MM-ddTHH:mm:ss+hh:mm)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := mcp.ParseString(req, "query", "")
		streamName := mcp.ParseString(req, "streamName", "")
		startTime := mcp.ParseString(req, "startTime", "")
		endTime := mcp.ParseString(req, "endTime", "")
		if query == "" || streamName == "" || startTime == "" || endTime == "" {
			return mcp.NewToolResultError("missing required fields: query, streamName, startTime, and endTime are required"), nil
		}
		queryResult, err := doParseableQuery(query, streamName, startTime, endTime)
		if err != nil {
			return mcp.NewToolResultError("query failed: " + err.Error()), nil
		}

		b, err := json.MarshalIndent(queryResult, "", "  ")
		if err != nil {
			return nil, err
		}

		// Optional: a one-liner summary that sets expectations.
		text := fmt.Sprintf(
			"Returned %d rows as JSON (array of objects). Use keys exactly as shown.\n```json\n%s\n```",
			len(queryResult),
			string(b),
		)

		return mcp.NewToolResultText(text), nil
		// Optionally, for structured output:
		// return mcp.NewToolResultStructured(map[string]interface{}{"result": result}, "Query successful"), nil
	})
}
