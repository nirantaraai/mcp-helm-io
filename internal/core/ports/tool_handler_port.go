package ports

import "context"

// ToolHandlerPort defines the interface for MCP tool handlers
type ToolHandlerPort interface {
	// HandleDeployChart handles the deploy_chart tool call
	HandleDeployChart(ctx context.Context, args map[string]interface{}) (interface{}, error)

	// HandleUpgradeChart handles the upgrade_chart tool call
	HandleUpgradeChart(ctx context.Context, args map[string]interface{}) (interface{}, error)

	// HandleRollbackRelease handles the rollback_release tool call
	HandleRollbackRelease(ctx context.Context, args map[string]interface{}) (interface{}, error)

	// HandleUninstallRelease handles the uninstall_release tool call
	HandleUninstallRelease(ctx context.Context, args map[string]interface{}) (interface{}, error)

	// HandleListReleases handles the list_releases tool call
	HandleListReleases(ctx context.Context, args map[string]interface{}) (interface{}, error)

	// HandleGetReleaseStatus handles the get_release_status tool call
	HandleGetReleaseStatus(ctx context.Context, args map[string]interface{}) (interface{}, error)

	// HandleSearchCharts handles the search_charts tool call
	HandleSearchCharts(ctx context.Context, args map[string]interface{}) (interface{}, error)
}

// Made with Bob
