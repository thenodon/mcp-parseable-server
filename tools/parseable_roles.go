package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetRolesTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"get_roles",
		mcp.WithDescription(`Get information about the Parseable roles. Roles is used for handling RBAC permissions/privilege and define access to datasets. Calls /api/v1/roles.

Data is returned as a dictionary with the role name and a list of privilege. The privilege can be of the following:
- admin - have all privileges
- editor - have limited privileges like cluster features 
- reader - allow read from datasets
- writer - allow read and write from datasets 
- ingestor - allow write from datasets

For reader, writer and ingestor role there is always at least one resource connected to the role. This resources is typical a dataset.
For full description of roles and RBAC use https://www.parseable.com/docs/user-guide/rbac
`),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		about, err := getParseableRoles()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		var lines []string
		for k, v := range about {
			lines = append(lines, k+": "+fmt.Sprintf("%v", v))
		}
		return mcp.NewToolResultText(strings.Join(lines, "\n")), nil
		// Optionally, for structured output:
		// return mcp.NewToolResultStructured(map[string]interface{}{"info": info}, "Info returned"), nil
	})
}
