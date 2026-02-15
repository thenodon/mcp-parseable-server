#!/bin/bash
# Test MCP Prompts using stdio mode
# Usage: ./test-prompts-stdio.sh [stream-name] [start-time] [end-time]

set -e

STREAM_NAME="${1:-otellogs}"
START_TIME="${2:-2026-02-13T00:00:00Z}"
END_TIME="${3:-2026-02-14T00:00:00Z}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== MCP Prompts Test Suite (stdio mode) ===${NC}"
echo "Stream: $STREAM_NAME"
echo "Time Range: $START_TIME to $END_TIME"
echo ""

# Helper function to send JSON-RPC via stdio
send_rpc() {
    local method="$1"
    local params="$2"

    echo "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"$method\",\"params\":$params}"
}

# Test 1: Initialize
echo -e "${YELLOW}[1/7] Initializing MCP connection...${NC}"
send_rpc "initialize" '{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}' | \
    ./mcp-parseable-server --mode=stdio | jq -r '.result.serverInfo.name' && echo -e "${GREEN}  ✓ Server initialized${NC}" || echo -e "${RED}  ✗ Failed${NC}"
echo ""

# Test 2: List prompts
echo -e "${YELLOW}[2/7] Listing available prompts...${NC}"
{
    send_rpc "initialize" '{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}'
    send_rpc "prompts/list" '{}'
} | ./mcp-parseable-server --mode=stdio | tail -1 | jq -r '.result.prompts[]? | "  ✓ \(.name) - \(.description)"'
echo ""

# Test 3: stream-health-check
echo -e "${YELLOW}[3/7] Testing stream-health-check prompt...${NC}"
{
    send_rpc "initialize" '{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}'
    send_rpc "prompts/get" "{\"name\":\"stream-health-check\",\"arguments\":{\"streamName\":\"$STREAM_NAME\"}}"
} | ./mcp-parseable-server --mode=stdio | tail -1 | jq -r '.result.messages[0].content.text' | head -10
echo ""

# Test 4: analyze-errors
echo -e "${YELLOW}[4/7] Testing analyze-errors prompt...${NC}"
{
    send_rpc "initialize" '{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}'
    send_rpc "prompts/get" "{\"name\":\"analyze-errors\",\"arguments\":{\"streamName\":\"$STREAM_NAME\",\"startTime\":\"$START_TIME\",\"endTime\":\"$END_TIME\"}}"
} | ./mcp-parseable-server --mode=stdio | tail -1 | jq -r '.result.messages[0].content.text' | head -10
echo ""

# Test 5: investigate-field
echo -e "${YELLOW}[5/7] Testing investigate-field prompt...${NC}"
{
    send_rpc "initialize" '{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}'
    send_rpc "prompts/get" "{\"name\":\"investigate-field\",\"arguments\":{\"streamName\":\"$STREAM_NAME\",\"fieldName\":\"severity_text\",\"startTime\":\"$START_TIME\",\"endTime\":\"$END_TIME\"}}"
} | ./mcp-parseable-server --mode=stdio | tail -1 | jq -r '.result.messages[0].content.text' | head -10
echo ""

# Test 6: compare-streams
echo -e "${YELLOW}[6/7] Testing compare-streams prompt...${NC}"
{
    send_rpc "initialize" '{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}'
    send_rpc "prompts/get" "{\"name\":\"compare-streams\",\"arguments\":{\"stream1\":\"otellogs\",\"stream2\":\"network_logstream\"}}"
} | ./mcp-parseable-server --mode=stdio | tail -1 | jq -r '.result.messages[0].content.text' | head -10
echo ""

# Test 7: find-anomalies
echo -e "${YELLOW}[7/7] Testing find-anomalies prompt...${NC}"
{
    send_rpc "initialize" '{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}'
    send_rpc "prompts/get" "{\"name\":\"find-anomalies\",\"arguments\":{\"streamName\":\"$STREAM_NAME\",\"startTime\":\"$START_TIME\",\"endTime\":\"$END_TIME\",\"groupBy\":\"hour\"}}"
} | ./mcp-parseable-server --mode=stdio | tail -1 | jq -r '.result.messages[0].content.text' | head -10
echo ""

echo -e "${GREEN}=== All tests complete ===${NC}"
echo ""
echo "Note: The prompts return instructions for AI agents to follow."
echo "For HTTP mode testing, use an MCP-compatible client like Claude Desktop or mcp-cli."

