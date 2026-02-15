MCP Prompts Guide
-----

This document provides detailed documentation for the 5 MCP prompts in the Parseable MCP server.

# What are MCP Prompts?

MCP Prompts are pre-built, multi-step workflows that guide AI agents through common tasks. 
Instead of calling individual tools manually, prompts provide a structured approach to accomplishing complex operations.

When an AI agent receives a prompt, it gets detailed instructions on:
- Which tools to call
- In what order
- What to analyze
- What to report back

When using prompt, we can use natural language to ask questions related to our data:
- "Check the health of the otellogs stream"
- "Find errors in network_logstream from yesterday"
- "Compare otellogs and monitor_logstream"
- "Investigate the severity_text field"
- "Find anomalies in the last 24 hours"

---
# Available Prompts
Agents can use the following pre-built prompts to guide their analysis:

## 1. analyze-errors

**Purpose:** Find and analyze error logs in a data stream.

**When to use:**
- Investigating system errors or failures
- Understanding error patterns over time
- Troubleshooting production issues

**Required Arguments:**
- `streamName` - Name of the data stream to analyze
- `startTime` - Start time in ISO 8601 format (e.g., "2026-02-13T00:00:00Z")
- `endTime` - End time in ISO 8601 format

**Optional Arguments:**
- `errorField` - Specific field to check for errors (default: "body")

**What it does:**
1. Gets the stream schema to understand available fields
2. Queries for errors using pattern matching (error, exception, failed)
3. Analyzes error distribution over time
4. Identifies common error patterns and messages
5. Provides recommendations for investigation

**Usage:**
```
"Find errors in the otellogs stream from yesterday"
"Analyze errors in network_logstream from the last 24 hours"
```

---
## 2. stream-health-check

**Purpose:** Perform a comprehensive health assessment of a data stream.

**When to use:**
- Regular stream monitoring
- Checking if data ingestion is working
- Verifying stream configuration
- Identifying potential issues

**Required Arguments:**
- `streamName` - Name of the data stream to check

**What it does:**
1. Gets stream information (creation time, event timestamps)
2. Gets stream statistics (ingestion rates, storage)
3. Gets stream schema
4. Analyzes:
   - Stream activity status (active/inactive)
   - Data freshness (time since last event)
   - Ingestion rates (events per day/hour)
   - Storage efficiency
5. Reports warnings or concerns
6. Provides optimization recommendations

**Usage:**
```
"Check the health of the otellogs stream"
"Is the network_logstream receiving data?"
```

---
## 3. investigate-field

**Purpose:** Deep dive into a specific field's values and patterns.

**When to use:**
- Understanding field distributions
- Data quality analysis
- Finding most common values
- Checking for nulls or empty values

**Required Arguments:**
- `streamName` - Name of the data stream
- `fieldName` - Name of the field to investigate
- `startTime` - Start time in ISO 8601 format
- `endTime` - End time in ISO 8601 format

**What it does:**
1. Verifies the field exists in the schema
2. Counts total occurrences and unique values
3. Gets top 20 most common values with counts
4. Checks for null/empty values
5. Calculates percentages
6. Provides data quality assessment
7. Recommends filters or queries

**Usage:**
```
"Investigate the severity_text field in otellogs"
"Show me the distribution of the status field"
"What are the most common values for service_name?"
```

---
## 4. compare-streams

**Purpose:** Compare metrics across multiple data streams.

**When to use:**
- Comparing production vs staging streams
- Understanding which streams are most active
- Capacity planning
- Identifying underutilized streams

**Required Arguments:**
- `stream1` - First data stream name
- `stream2` - Second data stream name

**Optional Arguments:**
- `stream3` - Optional third data stream name

**What it does:**
1. Gets info and stats for each stream
2. Compares:
   - Total events (lifetime count)
   - Storage usage
   - Ingestion rates (events per day)
   - Data freshness (last event time)
   - Activity levels
3. Creates comparison table
4. Identifies which stream is most active
5. Highlights concerning trends
6. Provides recommendations

**Usage:**
```
"Compare otellogs and network_logstream"
"Which stream is more active: prod-logs or staging-logs?"
"Compare the top 3 streams by size"
```

---
## 5. find-anomalies

**Purpose:** Detect unusual patterns, spikes, or drops in event volumes.

**When to use:**
- Detecting system anomalies
- Finding traffic spikes
- Identifying data ingestion issues
- Monitoring for unusual activity

**Required Arguments:**
- `streamName` - Name of the data stream to analyze
- `startTime` - Start time in ISO 8601 format
- `endTime` - End time in ISO 8601 format

**Optional Arguments:**
- `groupBy` - Time grouping: "hour" or "day" (default: "hour")

**What it does:**
1. Queries event counts are grouped by time (hour or day)
2. Calculates average event count (baseline)
3. Identifies spikes (significantly higher counts)
4. Identifies drops (significantly lower counts or gaps)
5. Investigates anomalous periods
6. Provides timeline with highlighted anomalies
7. Suggests potential causes
8. Recommends follow-up actions

**Usage:**
```
"Find anomalies in otellogs from the last week"
"Are there any spikes in network_logstream today?"
"Detect unusual patterns in the last 24 hours"
```

---
# Using Prompts with AI Agents

## In AI Agent Conversations

When using MCP-compatible clients, you can reference prompts using natural language:

**Health Checks:**
- "Check the health of otellogs"
- "Is network_logstream working?"
- "Run a health check on all my streams"

**Error Analysis:**
- "Find errors in the last hour"
- "Analyze errors in otellogs from yesterday"
- "Show me all failed requests"

**Field Investigation:**
- "What are the common values for severity_text?"
- "Investigate the user_id field"
- "Show me the distribution of status codes"

**Stream Comparison:**
- "Compare otellogs and network_logstream"
- "Which stream has more data?"
- "Compare all production streams"

**Anomaly Detection:**
- "Find any unusual activity"
- "Are there spikes in the last 24 hours?"
- "Detect anomalies this week"


## Testing Prompts

To test prompts, see [TESTING.md](TESTING.md) for the complete testing guide using the stdio test script.

---
# Best Practices

## Time Ranges

- Start with shorter time ranges (1-24 hours) for faster results
- Expand to longer ranges if needed
- Use ISO 8601 format: `2026-02-14T00:00:00Z`

## Field Names

- Use `get_data_stream_schema` first to see available fields
- Field names are case-sensitive and fields with special characters like `.` or `-` must be quoted with backticks.
- Common fields: `body`, `severity_text`, `p_timestamp`

## Stream Names

- Use exact stream names (case-sensitive)
- List streams with `get_data_streams` tool if unsure
- Stream names in prompts must match exactly

## Performance

- **analyze-errors**: Limit to 100-1000 results
- **investigate-field**: Works best with shorter time ranges
- **find-anomalies**: Use `groupBy=day` for long time ranges
- **compare-streams**: Compare 2-3 streams at a time

---
# Troubleshooting

**"Stream not found"**
- Verify stream name with `get_data_streams` tool
- Check spelling and case sensitivity

**"No data returned"**
- Verify time range contains data
- Check stream has events in that period
- Use `get_data_stream_info` to see first/last event times

**"Field not found"**
- Use `get_data_stream_schema` to see available fields
- Check field name spelling and case

**Slow responses**
- Reduce time range
- Use hourly grouping instead of minute-level queries
- For anomalies, use `groupBy=day` for long ranges

---
# Next Steps

1. Review [TESTING.md](TESTING.md) to test prompts
2. Configure Claude Desktop or another MCP client
3. Try prompts with natural language
4. Customize prompts in `prompts/prompts.go` for your needs
5. Create new prompts for your specific use cases

---
# Customizing Prompts

Prompts are defined in `prompts/prompts.go`. You can:

- Modify SQL queries in existing prompts
- Change default values (e.g., error patterns, field names)
- Add new arguments
- Create entirely new prompts

Each prompt is a function that returns instructions for the AI agent. The agent then uses the MCP tools to execute those instructions.

