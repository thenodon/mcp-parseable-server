package tools

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetRolesTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"get_roles",
		mcp.WithDescription(`Get role-based access control (RBAC) information for the Parseable instance.
Use this to understand user roles, permissions, and dataset access controls.
Calls /api/v1/role.

Returns a JSON object where each key is a role name and the value is an array of privileges assigned to that role.

Available Privilege Types:
- admin: Full system access with all privileges (no resource restrictions)
- editor: Limited administrative privileges for cluster features (no resource restrictions)
- reader: Read-only access to specific datasets (requires at least one dataset resource)
- writer: Read and write access to specific datasets (requires at least one dataset resource)
- ingestor: Write-only access to specific datasets for data ingestion (requires at least one dataset resource)

Resource Assignment:
- admin and editor roles apply globally across the entire Parseable instance
- reader, writer, and ingestor roles are always associated with specific dataset resources
- Each role entry includes the list of datasets (resources) the role has access to

Use this tool to understand access controls before querying or ingesting data.
For detailed RBAC documentation, see: https://www.parseable.com/docs/user-guide/rbac
`),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		slog.Info("fetching roles information")
		roles, err := getParseableRoles()
		if err != nil {
			slog.Error("failed to get roles", "error", err)
			return mcp.NewToolResultError(err.Error()), nil
		}
		slog.Info("successfully retrieved roles", "count", len(roles))
		return mcp.NewToolResultJSON(roles)
	})
}
