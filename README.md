mcp-parseable-server - A Parseable MCP Server 
----------------------

> This project is currently in early development. Feedback and contributions are welcome.
> 
> Testing has been done using:
> - vscode with Github Copilot as the agent.
> - opencode agent CLI tool - works from v0.2.0

# Overview
This project provides an MCP (Message Context Protocol) server for [Parseable](https://github.com/parseablehq/parseable), enabling AI agents and tools to interact with Parseable data streams (logs, metrics, traces) using natural language and structured tool calls.

# Features
- List available data streams in Parseable
- Query data streams using SQL
- Get schema, stats, and info for any data stream
- Modular MCP tool registration for easy extension
- Supports both HTTP and stdio MCP modes
- Environment variable and flag-based configuration
- The mcp server returns responses in json where the payload is both in text and structured format.

> In Parseable dataset and data stream names are used interchangeably, as Parseable's datasets are essentially 
> named data streams. In all tool description we try to use the term data stream to avoid confusion with the term dataset which can have different 
> meanings in other contexts.

# Testing
To test the server you can use the [mcp-cli](https://github.com/philschmid/mcp-cli)
```shell
mcp-cli  call parseable get_roles
``` 
Returns 
```json
{
  "content": [
    {
      "type": "text",
      "text": "{\"admins\":[{\"privilege\":\"admin\"}],\"network_role\":[{\"privilege\":\"reader\",\"resource\":{\"stream\":\"network_logstream\"}}],\"otel_gateway\":[{\"privilege\":\"editor\"}]}"
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
    ],
    "otel_gateway": [
      {
        "privilege": "editor"
      }
    ]
  }
}

```

# Building

Ensure you have Go 1.20+ installed.

```sh
git clone https://github.com/thenodon/mcp-parseable-server
cd mcp-parseable-server
# Build the MCP server binary
go build -o mcp-parseable-server ./cmd/mcp_parseable_server
```

# Running

## HTTP Mode (default)

```sh
./mcp-parseable-server --listen :9034
```

The MCP server will listen on `http://localhost:9034/mcp` for agent/tool requests.

## Stdio Mode

```sh
./mcp-parseable-server --mode stdio
```

This mode is used for CLI or agent-to-agent workflows.

# Configuration

You can configure the Parseable connection using environment variables or flags:

- `PARSEABLE_URL` or `--parseable-url`- url to the parseable instance (default: http://localhost:8000)
- `PARSEABLE_USERNAME` or `--parseable-user` (default: admin)
- `PARSEABLE_PASSWORD` or `--parseable-pass` (default: admin)
- `LISTEN_ADDR` or `--listen` - the address when running the mcp server in http mode (default: :9034)
- `INSECURE` - set to `true` to skip TLS verification (default: false)`
- `LOG_LEVEL` - set log level. Supported levels are debug, info, warn and error (default: info)

Example:
```sh
PARSEABLE_URL="http://your-parseable-host:8000" PARSEABLE_USER="admin" PARSEABLE_PASS="admin" ./mcp-parseable-server
```

# MCP Tools

## 1. `query_data_stream`
Execute a SQL query against a data stream.
- **Inputs:**
  - `query`: SQL query string
  - `streamName`: Name of the data stream
  - `startTime`: ISO 8601 start time (e.g. 2026-01-01T00:00:00+00:00)
  - `endTime`: ISO 8601 end time
- **Returns:** Query result

## 2. `get_data_streams`
List all available data streams in Parseable.
- **Returns:** Array of stream names

## 3. `get_data_stream_schema`
Get the schema for a specific data stream.
- **Inputs:**
  - `stream`: Name of the data stream
- **Returns:** Schema fields and types

## 4. `get_data_stream_stats`
Get stats for a data stream.
- **Inputs:**
  - `streamName`: Name of the data stream
- **Returns:** Stats object (see tool description for details)

## 5. `get_data_stream_info`
Get info for a data stream.
- **Inputs:**
  - `streamName`: Name of the data stream
- **Returns:** Info object (see tool description for details)

## 6. `get_parseable_about`
Get Parseable about info.
- **Returns:** About object (see tool description for details)

## 7. `get_parseable_roles`
Get Parseable roles.
- **Returns:** Roles object (see tool description for details)

# Example: Querying with curl

1. **List streams:**
```sh
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "callTool",
    "params": { "tool": "list_data_streams", "arguments": {} }
  }'
```

2. **Query a stream:**
```sh
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "callTool",
    "params": {
      "tool": "query_data_stream",
      "arguments": {
        "query": "SELECT COUNT(*) FROM log WHERE body ILIKE '%clamd%'",
        "streamName": "otellogs",
        "startTime": "2026-01-01T00:00:00+00:00",
        "endTime": "2026-01-06T00:00:00+00:00"
      }
    }
  }'
```

3. **Get schema for a stream:**
```sh
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "callTool",
    "params": {
      "tool": "get_data_stream_schema",
      "arguments": {
        "stream": "otellogs"
      }
    }
  }'
```

# Tool Discovery

Agents can discover all available tools and their input/output schemas via the MCP protocol. Each tool description includes details about returned fields and their meanings.

# Extending

To add new tools, create a new file in `tools/`, implement the registration function, and add it to `RegisterParseableTools` in `tools/register.go`.

# License

This work is licensed under the GNU GENERAL PUBLIC LICENSE Version 3.

# Todo 
- No prompts or resources are currently included for agent usage.
- No tools to understand the Parseable setups. It can by using the `about` tool understand if cluster or standalone, but nothing about the configuration.
- No authentication mechanisms implemented. 
- .....