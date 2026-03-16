package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/nirantaraai/mcp-helm-io/internal/core/ports"
)

// MCPServer implements the MCPPort interface
type MCPServer struct {
	server  *mcp.Server
	handler ports.ToolHandlerPort
}

// NewMCPServer creates a new MCP server
func NewMCPServer(service ports.ServicePort) *MCPServer {
	return &MCPServer{
		handler: NewToolHandler(service),
	}
}

// Start starts the MCP server
func (s *MCPServer) Start(ctx context.Context) error {
	slog.Info("starting MCP server")

	// Create MCP server
	server := mcp.NewServer("helm-mcp-server", "1.0.0", nil)

	s.server = server

	// Register tools
	if err := s.RegisterTools(); err != nil {
		return fmt.Errorf("failed to register tools: %w", err)
	}

	// Start server with stdio transport
	transport := mcp.NewStdioTransport()
	if err := server.Run(ctx, transport); err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	return nil
}

// Stop stops the MCP server
func (s *MCPServer) Stop(ctx context.Context) error {
	slog.Info("stopping MCP server")
	// MCP server cleanup if needed
	return nil
}

// RegisterTools registers all available tools with the MCP server
func (s *MCPServer) RegisterTools() error {
	slog.Info("registering MCP tools")

	// Create tool handlers that wrap our existing handlers
	deployTool := mcp.NewServerTool(
		"deploy_chart",
		"Deploy a Helm chart to Kubernetes",
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]interface{}]) (*mcp.CallToolResultFor[interface{}], error) {
			result, err := s.handler.HandleDeployChart(ctx, params.Arguments)
			if err != nil {
				return nil, err
			}
			resultJSON, _ := json.Marshal(result)
			return &mcp.CallToolResultFor[interface{}]{Content: []mcp.Content{&mcp.TextContent{Text: string(resultJSON)}}}, nil
		},
	)

	upgradeTool := mcp.NewServerTool(
		"upgrade_chart",
		"Upgrade an existing Helm release",
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]interface{}]) (*mcp.CallToolResultFor[interface{}], error) {
			result, err := s.handler.HandleUpgradeChart(ctx, params.Arguments)
			if err != nil {
				return nil, err
			}
			resultJSON, _ := json.Marshal(result)
			return &mcp.CallToolResultFor[interface{}]{Content: []mcp.Content{&mcp.TextContent{Text: string(resultJSON)}}}, nil
		},
	)

	rollbackTool := mcp.NewServerTool(
		"rollback_release",
		"Rollback a Helm release to a previous revision",
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]interface{}]) (*mcp.CallToolResultFor[interface{}], error) {
			result, err := s.handler.HandleRollbackRelease(ctx, params.Arguments)
			if err != nil {
				return nil, err
			}
			resultJSON, _ := json.Marshal(result)
			return &mcp.CallToolResultFor[interface{}]{Content: []mcp.Content{&mcp.TextContent{Text: string(resultJSON)}}}, nil
		},
	)

	uninstallTool := mcp.NewServerTool(
		"uninstall_release",
		"Uninstall a Helm release",
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]interface{}]) (*mcp.CallToolResultFor[interface{}], error) {
			result, err := s.handler.HandleUninstallRelease(ctx, params.Arguments)
			if err != nil {
				return nil, err
			}
			resultJSON, _ := json.Marshal(result)
			return &mcp.CallToolResultFor[interface{}]{Content: []mcp.Content{&mcp.TextContent{Text: string(resultJSON)}}}, nil
		},
	)

	listTool := mcp.NewServerTool(
		"list_releases",
		"List all Helm releases in a namespace",
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]interface{}]) (*mcp.CallToolResultFor[interface{}], error) {
			result, err := s.handler.HandleListReleases(ctx, params.Arguments)
			if err != nil {
				return nil, err
			}
			resultJSON, _ := json.Marshal(result)
			return &mcp.CallToolResultFor[interface{}]{Content: []mcp.Content{&mcp.TextContent{Text: string(resultJSON)}}}, nil
		},
	)

	getStatusTool := mcp.NewServerTool(
		"get_release_status",
		"Get the status of a Helm release",
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]interface{}]) (*mcp.CallToolResultFor[interface{}], error) {
			result, err := s.handler.HandleGetReleaseStatus(ctx, params.Arguments)
			if err != nil {
				return nil, err
			}
			resultJSON, _ := json.Marshal(result)
			return &mcp.CallToolResultFor[interface{}]{Content: []mcp.Content{&mcp.TextContent{Text: string(resultJSON)}}}, nil
		},
	)

	searchTool := mcp.NewServerTool(
		"search_charts",
		"Search for Helm charts",
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]interface{}]) (*mcp.CallToolResultFor[interface{}], error) {
			result, err := s.handler.HandleSearchCharts(ctx, params.Arguments)
			if err != nil {
				return nil, err
			}
			resultJSON, _ := json.Marshal(result)
			return &mcp.CallToolResultFor[interface{}]{Content: []mcp.Content{&mcp.TextContent{Text: string(resultJSON)}}}, nil
		},
	)

	// Add all tools to the server
	s.server.AddTools(
		deployTool,
		upgradeTool,
		rollbackTool,
		uninstallTool,
		listTool,
		getStatusTool,
		searchTool,
	)

	return nil
}

// GetServerInfo returns information about the MCP server
func (s *MCPServer) GetServerInfo() ports.ServerInfo {
	return ports.ServerInfo{
		Name:        "helm-mcp-server",
		Version:     "1.0.0",
		Description: "MCP server for Helm operations",
		Tools: []ports.ToolInfo{
			{Name: "deploy_chart", Description: "Deploy a Helm chart"},
			{Name: "upgrade_chart", Description: "Upgrade a Helm release"},
			{Name: "rollback_release", Description: "Rollback a Helm release"},
			{Name: "uninstall_release", Description: "Uninstall a Helm release"},
			{Name: "list_releases", Description: "List Helm releases"},
			{Name: "get_release_status", Description: "Get release status"},
			{Name: "search_charts", Description: "Search for charts"},
		},
	}
}

// Made with Bob
