Testing MCP server
-------------
This guide shows you how to test the MCP tools and MCP prompts in your Parseable server.

# Testing Tools with mcp-cli

We use `mcp-cli` from https://github.com/philschmid/mcp-cli to test the MCP tools.

```bash
# Install mcp-cli (philschmid/mcp-cli)
# Follow the install instructions in the repo:
# https://github.com/philschmid/mcp-cli
```

You can use `mcp-cli` to test the MCP tools:

```bash
# Test get_roles tool
mcp-cli call parseable get_roles

# Test get_data_streams tool
mcp-cli call parseable get_data_streams

# Test query_data_stream tool
mcp-cli call parseable query_data_stream \
'{"streamName": "a_logstream", "query": "select hostname, count(*) as count from a_logstream group by hostname limit 5", "startTime": "2026-02-12T17:00:00Z", "endTime":"2026-02-12T19:00:00Z"}'
  
# Test get_data_stream_schema tool
mcp-cli call parseable get_data_stream_schema \
  '{"streamName": "a_logstream"}'
```

Example output from `mcp-cli call parseable get_roles`:
```json
{
  "content": [
    {
      "type": "text",
      "text": "{\"admins\":[{\"privilege\":\"admin\"}],\"network_role\":[{\"privilege\":\"reader\",\"resource\":{\"stream\":\"network_logstream\"}}]}"
    }
  ],
  "structuredContent": {
    "admins": [
      {
        "privilege": "admin"
      }
    ],
    "network_role": [
      {
        "privilege": "reader",
        "resource": {
          "stream": "network_logstream"
        }
      }
    ]
  }
}
```
---
# Testing Prompts

The following sections show how to test MCP prompts. When testing prompts with the test script, the server runs 
in stdio mode, which allows you to see the full prompt instructions and responses directly in the terminal.

## 1. Build the server
```bash
cd /home/andersh/github/mcp-parseable-server
go build -o mcp-parseable-server ./cmd/mcp_parseable_server
```

## 2. Run the test script
```bash
./examples/test-prompts-stdio.sh a_logstream 2026-02-01T00:00:00Z 2026-02-14T23:59:59Z
```

That's it! All 5 prompts will be tested in stdio mode.

## Test Script Usage

```bash
./examples/test-prompts-stdio.sh [stream-name] [start-time] [end-time]
```

**Arguments:**
- `stream-name` (optional): Name of your data stream (default: "otellogs")
- `start-time` (optional): ISO 8601 start time (default: "2026-02-13T00:00:00Z")
- `end-time` (optional): ISO 8601 end time (default: "2026-02-14T00:00:00Z")

**Examples:**
```bash
# Use default stream and dates
./examples/test-prompts-stdio.sh

# Test specific stream
./examples/test-prompts-stdio.sh pmeta

# Test with custom date range
./examples/test-prompts-stdio.sh otellogs 2026-02-01T00:00:00Z 2026-02-14T23:59:59Z
```

## Expected Output

```
=== MCP Prompts Test Suite (stdio mode) ===
Stream: network_logstream
Time Range: 2026-02-01T00:00:00Z to 2026-02-14T23:59:59Z

[1/7] Initializing MCP connection...
  ✓ Server initialized

[2/7] Listing available prompts...
  ✓ analyze-errors - Analyze error logs in a data stream...
  ✓ stream-health-check - Perform a health check...
  ✓ investigate-field - Investigate a specific field...
  ✓ compare-streams - Compare metrics across streams...
  ✓ find-anomalies - Find anomalies and unusual patterns...

[3/7] Testing stream-health-check prompt...
You are performing a health check on the Parseable data stream "network_logstream".

Follow these steps:
1. Get stream information using get_data_stream_info...
...

=== All tests complete ===
```

## Why stdio Mode for Testing?

### HTTP Mode with SSE Transport

The MCP server runs in HTTP mode using the StreamableHTTPServer, which uses Server-Sent Events (SSE) for transport. This requires session management, making simple curl testing complex.

**How StreamableHTTPServer works:**
1. Client opens SSE connection: `GET /mcp`
2. Server sends session ID via SSE stream
3. Client includes session ID in subsequent POST requests
4. Server responds via SSE events

This is handled automatically by MCP clients (Claude Desktop, etc.) but is complex to implement manually with curl/scripts.

Note: `mcp-cli` does not support the MCP prompts feature. Use Claude Desktop or other MCP clients that support prompts  
or use the stdio test script for testing.

### stdio Mode Advantages

- ✅ No session management needed
- ✅ Works with standard bash/jq
- ✅ Immediate results
- ✅ Perfect for CI/CD and testing

## Understanding Prompt Responses

**Important:** Prompts return **instructions** for AI agents, not actual data.

Each prompt generates a detailed workflow that an AI agent (like Claude) would follow to complete the task. 
The agent then uses the MCP tools to execute those instructions.

Example response structure:
```json
{
  "result": {
    "description": "Health check for otellogs",
    "messages": [{
      "role": "user",
      "content": {
        "type": "text",
        "text": "You are performing a health check...\n\nFollow these steps:\n1. Get stream information...\n2. Get stream statistics..."
      }
    }]
  }
}
```

## Using Prompts with AI Agents

### Claude Desktop

Add to your Claude Desktop config:
- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`
- Linux: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "parseable": {
      "command": "/path/to/mcp-parseable-server",
      "args": ["--mode=stdio"],
      "env": {
        "PARSEABLE_URL": "http://localhost:8000",
        "PARSEABLE_USER": "admin",
        "PARSEABLE_PASS": "admin",
        "LOG_LEVEL": "info"
      }
    }
  }
}
```

### Other MCP Clients

Any MCP-compatible client can use these prompts. The server works in both stdio and HTTP modes. 
For HTTP mode, clients must properly handle SSE and session management.


## Manual Testing (Advanced)

You can test manually by sending JSON-RPC via stdin:

```bash
# List all prompts
{
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}'
  echo '{"jsonrpc":"2.0","id":2,"method":"prompts/list","params":{}}'
} | ./mcp-parseable-server --mode=stdio
```

Test a specific prompt:
```bash
{
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}'
  echo '{"jsonrpc":"2.0","id":2,"method":"prompts/get","params":{"name":"stream-health-check","arguments":{"streamName":"otellogs"}}}'
} | ./mcp-parseable-server --mode=stdio | tail -1 | jq -r '.result.messages[0].content.text'
```

