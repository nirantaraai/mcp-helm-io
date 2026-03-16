package ports

import "context"

// MCPPort defines the interface for MCP server operations
type MCPPort interface {
	// Start starts the MCP server
	Start(ctx context.Context) error

	// Stop stops the MCP server
	Stop(ctx context.Context) error

	// RegisterTools registers all available tools with the MCP server
	RegisterTools() error

	// GetServerInfo returns information about the MCP server
	GetServerInfo() ServerInfo
}

// ServerInfo contains information about the MCP server
type ServerInfo struct {
	Name        string
	Version     string
	Description string
	Tools       []ToolInfo
}

// ToolInfo contains information about an MCP tool
type ToolInfo struct {
	Name        string
	Description string
	InputSchema map[string]interface{}
}

// Made with Bob
