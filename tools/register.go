package tools

import "github.com/mark3labs/mcp-go/server"

func RegisterParseableTools(mcpServer *server.MCPServer) {
	RegisterQueryDataStreamTool(mcpServer)
	RegisterListDataStreamsTool(mcpServer)
	RegisterGetDataStreamSchemaTool(mcpServer)
	RegisterGetDataStreamStatsTool(mcpServer)
	RegisterGetDataStreamInfoTool(mcpServer)
	RegisterGetAboutTool(mcpServer)
	RegisterGetRolesTool(mcpServer)
}
