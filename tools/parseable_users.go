package tools

import (
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterGetUsersTool(mcpServer *server.MCPServer) {
	mcpServer.AddTool(mcp.NewTool(
		"get_users",
		mcp.WithDescription(`Get all configured users in the Parseable instance with their authentication methods and role assignments.
Use this to understand user access, authentication configuration, and role-based permissions.
Calls /api/v1/users.

Returns a JSON object with a 'users' array where each element is a user object with the following structure, plus 'count' (number of users):

User Information:
- id: unique user identifier
- username: the user's login name
- method: authentication method ("native" for local auth, "oidc" for OpenID Connect)
- email: user's email address (null if not configured)
- picture: user's profile picture URL (null if not set)

Role Assignments:
- roles: object mapping role names to arrays of privilege grants
  Each privilege grant contains:
  - privilege: the permission level ("admin", "editor", "reader", "writer", or "ingestor")
  - resource: object specifying the resource this privilege applies to
    - stream: the data stream name this privilege grants access to
- group_roles: object containing roles inherited from group membership
- user_groups: array of groups this user belongs to

Use this tool to:
- Verify user access permissions before executing operations
- Understand which users have access to specific data streams
- Check authentication methods configured for users
- Audit user-role-stream relationships
`),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		users, err := getParseableUsers()
		if err != nil {
			slog.Error("failed to get users", "error", err)
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultJSON(map[string]interface{}{
			"users": users,
			"count": len(users),
		})
	})
}
