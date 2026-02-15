package prompts

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterParseablePrompts registers all prompts with the MCP server
func RegisterParseablePrompts(mcpServer *server.MCPServer) {
	registerAnalyzeErrorsPrompt(mcpServer)
	registerStreamHealthCheckPrompt(mcpServer)
	registerInvestigateFieldPrompt(mcpServer)
	registerCompareStreamsPrompt(mcpServer)
	registerFindAnomaliesPrompt(mcpServer)
}

func registerAnalyzeErrorsPrompt(mcpServer *server.MCPServer) {
	mcpServer.AddPrompt(mcp.NewPrompt(
		"analyze-errors",
		mcp.WithPromptDescription("Analyze error logs in a data stream over a time range. "+
			"Gets schema, queries for errors, and provides a summary with patterns and recommendations."),
		mcp.WithArgument("streamName", mcp.RequiredArgument(), mcp.ArgumentDescription("Name of the data stream to analyze")),
		mcp.WithArgument("startTime", mcp.RequiredArgument(), mcp.ArgumentDescription("Start time in ISO 8601 format (e.g., '2026-02-01T00:00:00Z')")),
		mcp.WithArgument("endTime", mcp.RequiredArgument(), mcp.ArgumentDescription("End time in ISO 8601 format (e.g., '2026-02-13T23:59:59Z')")),
		mcp.WithArgument("errorField", mcp.ArgumentDescription("Optional: specific field to check for errors (default: body)")),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := request.Params.Arguments
		streamName := args["streamName"]
		startTime := args["startTime"]
		endTime := args["endTime"]
		errorField := args["errorField"]
		if errorField == "" {
			errorField = "body"
		}

		promptText := `You are analyzing error logs in the Parseable data stream "` + streamName + `" from ` + startTime + ` to ` + endTime + `.

Follow these steps:

1. First, call get_data_stream_schema with streamName="` + streamName + `" to understand available fields.

2. Query for errors using query_data_stream with:
   - query: "SELECT * FROM ` + streamName + ` WHERE ` + errorField + ` ILIKE '%error%' OR ` + errorField + ` ILIKE '%exception%' OR ` + errorField + ` ILIKE '%failed%' ORDER BY p_timestamp DESC LIMIT 100"
   - streamName: "` + streamName + `"
   - startTime: "` + startTime + `"
   - endTime: "` + endTime + `"

3. Analyze the results and provide:
   - Total number of errors found
   - Common error patterns or messages
   - Time distribution of errors (any spikes?)
   - Affected components or services (if identifiable)
   - Recommended next steps for investigation

If no errors are found, confirm this and suggest checking different time ranges or search terms.`

		return mcp.NewGetPromptResult(
			"Analyze errors in "+streamName,
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleUser, mcp.TextContent{
					Type: "text",
					Text: promptText,
				}),
			},
		), nil
	})
}

func registerStreamHealthCheckPrompt(mcpServer *server.MCPServer) {
	mcpServer.AddPrompt(mcp.NewPrompt(
		"stream-health-check",
		mcp.WithPromptDescription("Perform a health check on a data stream. "+
			"Analyzes ingestion rates, storage usage, and data freshness to identify potential issues."),
		mcp.WithArgument("streamName", mcp.RequiredArgument(), mcp.ArgumentDescription("Name of the data stream to check")),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := request.Params.Arguments
		streamName := args["streamName"]

		promptText := `You are performing a health check on the Parseable data stream "` + streamName + `".

Follow these steps:

1. Get stream information using get_data_stream_info with streamName="` + streamName + `"
   - Check createdAt, firstEventAt, and latestEventAt timestamps
   - Verify if the stream is receiving recent data

2. Get stream statistics using get_data_stream_stats with streamName="` + streamName + `"
   - Review ingestion counts and sizes
   - Compare lifetime vs current period metrics
   - Check storage utilization

3. Get stream schema using get_data_stream_schema with streamName="` + streamName + `"
   - Verify schema structure is as expected
   - Note important fields for querying

4. Provide a health summary including:
   - Stream status (active/inactive based on latestEventAt)
   - Data ingestion rate (events per day/hour)
   - Storage efficiency (compression ratio if detectable)
   - Data freshness (time since last event)
   - Any warnings or concerns (e.g., no recent data, rapid storage growth)
   - Recommendations for optimization or investigation

Present the findings in a clear, actionable format.`

		return mcp.NewGetPromptResult(
			"Health check for "+streamName,
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleUser, mcp.TextContent{
					Type: "text",
					Text: promptText,
				}),
			},
		), nil
	})
}

func registerInvestigateFieldPrompt(mcpServer *server.MCPServer) {
	mcpServer.AddPrompt(mcp.NewPrompt(
		"investigate-field",
		mcp.WithPromptDescription("Investigate a specific field in a data stream. "+
			"Analyzes field values, distributions, and patterns over a time range."),
		mcp.WithArgument("streamName", mcp.RequiredArgument(), mcp.ArgumentDescription("Name of the data stream")),
		mcp.WithArgument("fieldName", mcp.RequiredArgument(), mcp.ArgumentDescription("Name of the field to investigate")),
		mcp.WithArgument("startTime", mcp.RequiredArgument(), mcp.ArgumentDescription("Start time in ISO 8601 format")),
		mcp.WithArgument("endTime", mcp.RequiredArgument(), mcp.ArgumentDescription("End time in ISO 8601 format")),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := request.Params.Arguments
		streamName := args["streamName"]
		fieldName := args["fieldName"]
		startTime := args["startTime"]
		endTime := args["endTime"]

		promptText := `You are investigating the field "` + fieldName + `" in the Parseable data stream "` + streamName + `" from ` + startTime + ` to ` + endTime + `.

Follow these steps:

1. Verify the field exists using get_data_stream_schema with streamName="` + streamName + `"
   - Confirm field name and data type
   - If the field name include space, dot or special characters, escape them with backticks

2. Get field statistics using query_data_stream:
   - query: "SELECT COUNT(*) as total_count, COUNT(DISTINCT ` + fieldName + `) as unique_values FROM ` + streamName + `"
   - This shows how many total records and unique values exist

3. Get top values using query_data_stream:
   - query: "SELECT ` + fieldName + `, COUNT(*) as count FROM ` + streamName + ` GROUP BY ` + fieldName + ` ORDER BY count DESC LIMIT 20"
   - This shows the most common values

4. Check for nulls/empty values using query_data_stream:
   - query: "SELECT COUNT(*) FROM ` + streamName + ` WHERE ` + fieldName + ` IS NULL OR ` + fieldName + ` = ''"

5. Analyze and report:
   - Total occurrences and unique values
   - Top 10 most common values with counts
   - Percentage of null/empty values
   - Any unusual patterns or outliers
   - Data quality assessment
   - Recommendations for queries or filters

Present findings in a structured, easy-to-read format.`

		return mcp.NewGetPromptResult(
			"Investigate field "+fieldName+" in "+streamName,
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleUser, mcp.TextContent{
					Type: "text",
					Text: promptText,
				}),
			},
		), nil
	})
}

func registerCompareStreamsPrompt(mcpServer *server.MCPServer) {
	mcpServer.AddPrompt(mcp.NewPrompt(
		"compare-streams",
		mcp.WithPromptDescription("Compare metrics across multiple data streams. "+
			"Useful for understanding relative activity, storage usage, and data patterns."),
		mcp.WithArgument("stream1", mcp.RequiredArgument(), mcp.ArgumentDescription("First data stream name")),
		mcp.WithArgument("stream2", mcp.RequiredArgument(), mcp.ArgumentDescription("Second data stream name")),
		mcp.WithArgument("stream3", mcp.ArgumentDescription("Optional third data stream name")),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := request.Params.Arguments
		stream1 := args["stream1"]
		stream2 := args["stream2"]
		stream3 := args["stream3"]

		streams := `"` + stream1 + `" and "` + stream2 + `"`
		if stream3 != "" {
			streams = `"` + stream1 + `", "` + stream2 + `", and "` + stream3 + `"`
		}

		promptText := `You are comparing the following Parseable data streams: ` + streams + `.

Follow these steps for each stream:

1. Get stream info using get_data_stream_info
   - Note creation time and event timestamps
   - Check stream type and telemetry type

2. Get stream stats using get_data_stream_stats
   - Record ingestion counts and sizes
   - Note storage metrics

3. Compare and analyze:
   - Ingestion rates (events per day)
   - Storage usage and growth
   - Data freshness (time since last event)
   - Ingestion format and storage format
   - Activity level (active vs dormant)

4. Present a comparison table with:
   - Stream names
   - Total events (lifetime_count)
   - Storage size (human-readable)
   - Average event size
   - Events per day (calculated)
   - Last event time
   - Status (active/inactive)

5. Provide insights:
   - Which stream is most active?
   - Which uses most storage?
   - Any concerning trends (e.g., stream stopped receiving data)?
   - Recommendations for stream management

Format the comparison clearly for easy decision-making.`

		return mcp.NewGetPromptResult(
			"Compare data streams",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleUser, mcp.TextContent{
					Type: "text",
					Text: promptText,
				}),
			},
		), nil
	})
}

func registerFindAnomaliesPrompt(mcpServer *server.MCPServer) {
	mcpServer.AddPrompt(mcp.NewPrompt(
		"find-anomalies",
		mcp.WithPromptDescription("Find anomalies and unusual patterns in a data stream over time. "+
			"Looks for spikes, drops, and irregular patterns in event volumes."),
		mcp.WithArgument("streamName", mcp.RequiredArgument(), mcp.ArgumentDescription("Name of the data stream to analyze")),
		mcp.WithArgument("startTime", mcp.RequiredArgument(), mcp.ArgumentDescription("Start time in ISO 8601 format")),
		mcp.WithArgument("endTime", mcp.RequiredArgument(), mcp.ArgumentDescription("End time in ISO 8601 format")),
		mcp.WithArgument("groupBy", mcp.ArgumentDescription("Time grouping: 'hour' or 'day' (default: hour)")),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := request.Params.Arguments
		streamName := args["streamName"]
		startTime := args["startTime"]
		endTime := args["endTime"]
		groupBy := args["groupBy"]
		if groupBy == "" {
			groupBy = "hour"
		}

		// Determine the date_trunc interval
		interval := "hour"
		if groupBy == "day" {
			interval = "day"
		}

		promptText := `You are analyzing the Parseable data stream "` + streamName + `" for anomalies from ` + startTime + ` to ` + endTime + `.

Follow these steps:

1. Get event counts over time using query_data_stream:
   - query: "SELECT date_trunc('` + interval + `', p_timestamp) as time_bucket, COUNT(*) as event_count FROM ` + streamName + ` GROUP BY time_bucket ORDER BY time_bucket"
   - streamName: "` + streamName + `"
   - startTime: "` + startTime + `"
   - endTime: "` + endTime + `"

2. Analyze the time series data:
   - Calculate average event count per ` + interval + `
   - Identify spikes (periods with significantly higher counts)
   - Identify drops (periods with significantly lower counts or zero events)
   - Look for irregular patterns or gaps

3. For any anomalies found, investigate further:
   - Query sample events from anomalous periods
   - Check if specific fields or patterns are different

4. Report findings:
   - Summary of normal activity baseline
   - List of anomalies with timestamps and severity
   - Potential causes or correlations
   - Recommended actions (e.g., investigate specific time windows, check for system issues)

Present a clear timeline of activity with highlighted anomalies.`

		return mcp.NewGetPromptResult(
			"Find anomalies in "+streamName,
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(mcp.RoleUser, mcp.TextContent{
					Type: "text",
					Text: promptText,
				}),
			},
		), nil
	})
}
