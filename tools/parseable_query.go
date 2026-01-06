package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterQueryDataStreamTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"query_data_stream",
		mcp.WithDescription("Execute a SQL query against a data stream in Parseable. All fields are required. All times must be in ISO 8601 format."),
		mcp.WithString("query", mcp.Required(), mcp.Description("SQL query to run")),
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
		payload := map[string]string{
			"query":      query,
			"streamName": streamName,
			"startTime":  startTime,
			"endTime":    endTime,
		}
		jsonPayload, _ := json.Marshal(payload)
		url := ParseableBaseURL + parseableSQLPath
		httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return mcp.NewToolResultError("failed to create request"), nil
		}
		httpReq.Header.Set("Content-Type", "application/json")
		addBasicAuth(httpReq)
		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			return mcp.NewToolResultError("query failed: " + err.Error()), nil
		}
		defer resp.Body.Close()
		var result interface{}
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &result); err != nil {
			return mcp.NewToolResultError("failed to parse parseable response"), nil
		}
		return mcp.NewToolResultStructured(map[string]interface{}{"result": result}, "Query successful"), nil
	})
}
