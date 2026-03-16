package mcp

import (
	"context"
	"fmt"

	"github.com/nirantaraai/mcp-helm-io/internal/core/ports"
)

// ToolHandler implements the ToolHandlerPort interface
type ToolHandler struct {
	service ports.ServicePort
}

// NewToolHandler creates a new ToolHandler
func NewToolHandler(service ports.ServicePort) *ToolHandler {
	return &ToolHandler{
		service: service,
	}
}

// HandleDeployChart handles the deploy_chart tool call
func (h *ToolHandler) HandleDeployChart(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cmd := ports.DeployChartCommand{
		Chart:           getString(args, "chart"),
		ReleaseName:     getString(args, "release_name"),
		Namespace:       getString(args, "namespace"),
		Version:         getString(args, "version"),
		Values:          getMap(args, "values"),
		CreateNamespace: getBool(args, "create_namespace"),
		Wait:            getBool(args, "wait"),
		Timeout:         getInt(args, "timeout"),
	}

	release, err := h.service.DeployChart(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"release_name": release.Name,
		"namespace":    release.Namespace,
		"status":       release.Status.String(),
		"revision":     release.Revision,
		"chart":        release.Chart,
		"version":      release.Version,
	}, nil
}

// HandleUpgradeChart handles the upgrade_chart tool call
func (h *ToolHandler) HandleUpgradeChart(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cmd := ports.UpgradeChartCommand{
		Chart:       getString(args, "chart"),
		ReleaseName: getString(args, "release_name"),
		Namespace:   getString(args, "namespace"),
		Version:     getString(args, "version"),
		Values:      getMap(args, "values"),
		ResetValues: getBool(args, "reset_values"),
		ReuseValues: getBool(args, "reuse_values"),
		Wait:        getBool(args, "wait"),
		Timeout:     getInt(args, "timeout"),
	}

	release, err := h.service.UpgradeChart(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"release_name": release.Name,
		"namespace":    release.Namespace,
		"status":       release.Status.String(),
		"revision":     release.Revision,
	}, nil
}

// HandleRollbackRelease handles the rollback_release tool call
func (h *ToolHandler) HandleRollbackRelease(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cmd := ports.RollbackReleaseCommand{
		ReleaseName: getString(args, "release_name"),
		Namespace:   getString(args, "namespace"),
		Revision:    getInt(args, "revision"),
		Wait:        getBool(args, "wait"),
		Timeout:     getInt(args, "timeout"),
	}

	release, err := h.service.RollbackRelease(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"release_name": release.Name,
		"namespace":    release.Namespace,
		"status":       release.Status.String(),
		"revision":     release.Revision,
	}, nil
}

// HandleUninstallRelease handles the uninstall_release tool call
func (h *ToolHandler) HandleUninstallRelease(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	cmd := ports.UninstallReleaseCommand{
		ReleaseName: getString(args, "release_name"),
		Namespace:   getString(args, "namespace"),
		KeepHistory: getBool(args, "keep_history"),
		Timeout:     getInt(args, "timeout"),
	}

	err := h.service.UninstallRelease(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": fmt.Sprintf("Release %s uninstalled successfully", cmd.ReleaseName),
	}, nil
}

// HandleListReleases handles the list_releases tool call
func (h *ToolHandler) HandleListReleases(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	namespace := getString(args, "namespace")

	releases, err := h.service.ListReleases(ctx, namespace)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, rel := range releases {
		result = append(result, map[string]interface{}{
			"name":      rel.Name,
			"namespace": rel.Namespace,
			"status":    rel.Status.String(),
			"revision":  rel.Revision,
			"chart":     rel.Chart,
			"version":   rel.Version,
			"updated":   rel.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

// HandleGetReleaseStatus handles the get_release_status tool call
func (h *ToolHandler) HandleGetReleaseStatus(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	name := getString(args, "release_name")
	namespace := getString(args, "namespace")

	release, err := h.service.GetReleaseStatus(ctx, name, namespace)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"release_name": release.Name,
		"namespace":    release.Namespace,
		"status":       release.Status.String(),
		"revision":     release.Revision,
		"chart":        release.Chart,
		"version":      release.Version,
		"updated":      release.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// HandleSearchCharts handles the search_charts tool call
func (h *ToolHandler) HandleSearchCharts(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	keyword := getString(args, "keyword")

	charts, err := h.service.SearchCharts(ctx, keyword)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, chart := range charts {
		result = append(result, map[string]interface{}{
			"name":        chart.Name,
			"version":     chart.Version,
			"app_version": chart.AppVersion,
			"description": chart.Description,
			"repository":  chart.Repository,
		})
	}

	return result, nil
}

// Helper functions

func getString(args map[string]interface{}, key string) string {
	if val, ok := args[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getInt(args map[string]interface{}, key string) int {
	if val, ok := args[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case float64:
			return int(v)
		}
	}
	return 0
}

func getBool(args map[string]interface{}, key string) bool {
	if val, ok := args[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

func getMap(args map[string]interface{}, key string) map[string]interface{} {
	if val, ok := args[key]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			return m
		}
	}
	return make(map[string]interface{})
}

// Made with Bob
