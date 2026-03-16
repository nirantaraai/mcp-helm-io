package helm

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nirantaraai/mcp-helm-io/internal/core/domain"
	"github.com/nirantaraai/mcp-helm-io/internal/core/ports"
	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart"
	"helm.sh/helm/v4/pkg/chart/loader"
	"helm.sh/helm/v4/pkg/cli"
	"helm.sh/helm/v4/pkg/release"
	"helm.sh/helm/v4/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// HelmAdapter implements the HelmPort interface using Helm Go SDK v4
type HelmAdapter struct {
	settings *cli.EnvSettings
}

// NewHelmAdapter creates a new HelmAdapter
func NewHelmAdapter() *HelmAdapter {
	settings := cli.New()
	return &HelmAdapter{
		settings: settings,
	}
}

// DeployChart deploys a new Helm chart
func (h *HelmAdapter) DeployChart(ctx context.Context, cmd ports.DeployChartCommand) (*domain.HelmRelease, error) {
	slog.Info("deploying chart via Helm adapter",
		slog.String("chart", cmd.Chart),
		slog.String("release", cmd.ReleaseName),
	)

	// Create action configuration
	actionConfig, err := h.getActionConfig(cmd.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create action config: %w", err)
	}

	// Create install action
	install := action.NewInstall(actionConfig)
	install.ReleaseName = cmd.ReleaseName
	install.Namespace = cmd.Namespace
	install.CreateNamespace = cmd.CreateNamespace
	install.WaitForJobs = cmd.Wait
	install.Timeout = time.Duration(cmd.Timeout) * time.Second
	install.Version = cmd.Version

	// Load chart
	chartPath, err := install.ChartPathOptions.LocateChart(cmd.Chart, h.settings)
	if err != nil {
		return nil, fmt.Errorf("failed to locate chart: %w", err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %w", err)
	}

	// Install the chart
	rel, err := install.RunWithContext(ctx, chart, cmd.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to install chart: %w", err)
	}

	return h.convertToHelmRelease(rel), nil
}

// UpgradeChart upgrades an existing Helm release
func (h *HelmAdapter) UpgradeChart(ctx context.Context, cmd ports.UpgradeChartCommand) (*domain.HelmRelease, error) {
	slog.Info("upgrading chart via Helm adapter",
		slog.String("chart", cmd.Chart),
		slog.String("release", cmd.ReleaseName),
	)

	// Create action configuration
	actionConfig, err := h.getActionConfig(cmd.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create action config: %w", err)
	}

	// Create upgrade action
	upgrade := action.NewUpgrade(actionConfig)
	upgrade.Namespace = cmd.Namespace
	upgrade.WaitForJobs = cmd.Wait
	upgrade.Timeout = time.Duration(cmd.Timeout) * time.Second
	upgrade.Version = cmd.Version
	upgrade.ResetValues = cmd.ResetValues
	upgrade.ReuseValues = cmd.ReuseValues

	// Load chart
	chartPath, err := upgrade.ChartPathOptions.LocateChart(cmd.Chart, h.settings)
	if err != nil {
		return nil, fmt.Errorf("failed to locate chart: %w", err)
	}

	chart, err := loader.Load(chartPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart: %w", err)
	}

	// Upgrade the release
	rel, err := upgrade.RunWithContext(ctx, cmd.ReleaseName, chart, cmd.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade release: %w", err)
	}

	return h.convertToHelmRelease(rel), nil
}

// RollbackRelease rolls back a release to a previous revision
func (h *HelmAdapter) RollbackRelease(ctx context.Context, cmd ports.RollbackReleaseCommand) (*domain.HelmRelease, error) {
	slog.Info("rolling back release via Helm adapter",
		slog.String("release", cmd.ReleaseName),
		slog.Int("revision", cmd.Revision),
	)

	// Create action configuration
	actionConfig, err := h.getActionConfig(cmd.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create action config: %w", err)
	}

	// Create rollback action
	rollback := action.NewRollback(actionConfig)
	rollback.WaitForJobs = cmd.Wait
	rollback.Timeout = time.Duration(cmd.Timeout) * time.Second
	rollback.Version = cmd.Revision

	// Rollback the release
	if err := rollback.Run(cmd.ReleaseName); err != nil {
		return nil, fmt.Errorf("failed to rollback release: %w", err)
	}

	// Get the updated release status
	return h.GetReleaseStatus(ctx, cmd.ReleaseName, cmd.Namespace)
}

// UninstallRelease uninstalls a Helm release
func (h *HelmAdapter) UninstallRelease(ctx context.Context, cmd ports.UninstallReleaseCommand) error {
	slog.Info("uninstalling release via Helm adapter",
		slog.String("release", cmd.ReleaseName),
	)

	// Create action configuration
	actionConfig, err := h.getActionConfig(cmd.Namespace)
	if err != nil {
		return fmt.Errorf("failed to create action config: %w", err)
	}

	// Create uninstall action
	uninstall := action.NewUninstall(actionConfig)
	uninstall.KeepHistory = cmd.KeepHistory
	uninstall.Timeout = time.Duration(cmd.Timeout) * time.Second

	// Uninstall the release
	_, err = uninstall.Run(cmd.ReleaseName)
	if err != nil {
		return fmt.Errorf("failed to uninstall release: %w", err)
	}

	return nil
}

// ListReleases lists all Helm releases in a namespace
func (h *HelmAdapter) ListReleases(ctx context.Context, namespace string) ([]*domain.HelmRelease, error) {
	slog.Info("listing releases via Helm adapter", slog.String("namespace", namespace))

	// Create action configuration
	actionConfig, err := h.getActionConfig(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create action config: %w", err)
	}

	// Create list action
	list := action.NewList(actionConfig)
	list.AllNamespaces = namespace == ""

	// List releases
	releases, err := list.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	// Convert to domain releases
	var domainReleases []*domain.HelmRelease
	for _, rel := range releases {
		domainReleases = append(domainReleases, h.convertToHelmRelease(rel))
	}

	return domainReleases, nil
}

// GetReleaseStatus gets the status of a specific release
func (h *HelmAdapter) GetReleaseStatus(ctx context.Context, name, namespace string) (*domain.HelmRelease, error) {
	slog.Info("getting release status via Helm adapter",
		slog.String("release", name),
		slog.String("namespace", namespace),
	)

	// Create action configuration
	actionConfig, err := h.getActionConfig(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create action config: %w", err)
	}

	// Create status action
	status := action.NewStatus(actionConfig)

	// Get release status
	rel, err := status.Run(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get release status: %w", err)
	}

	return h.convertToHelmRelease(rel), nil
}

// SearchCharts searches for charts in configured repositories
func (h *HelmAdapter) SearchCharts(ctx context.Context, keyword string) ([]ports.ChartInfo, error) {
	slog.Info("searching charts via Helm adapter", slog.String("keyword", keyword))

	// For simplicity, return mock data
	// In a real implementation, you would use helm search or query chart repositories
	charts := []ports.ChartInfo{
		{
			Name:        fmt.Sprintf("%s-chart", keyword),
			Version:     "1.0.0",
			AppVersion:  "1.0.0",
			Description: fmt.Sprintf("Chart matching keyword: %s", keyword),
			Repository:  "stable",
		},
	}

	return charts, nil
}

// getActionConfig creates an action configuration for the given namespace
func (h *HelmAdapter) getActionConfig(namespace string) (*action.Configuration, error) {
	actionConfig := new(action.Configuration)

	// Create REST client getter
	configFlags := &genericclioptions.ConfigFlags{
		Namespace: &namespace,
	}

	if err := actionConfig.Init(configFlags, namespace, driver.SecretsDriverName); err != nil {
		return nil, fmt.Errorf("failed to initialize action config: %w", err)
	}

	return actionConfig, nil
}

// convertToHelmRelease converts a Helm release to domain HelmRelease
func (h *HelmAdapter) convertToHelmRelease(rel release.Releaser) *domain.HelmRelease {
	// Convert Releaser to Accessor to access release information
	accessor, err := release.NewAccessor(rel)
	if err != nil {
		slog.Error("failed to create release accessor", slog.Any("error", err))
		return &domain.HelmRelease{
			Name:   "unknown",
			Status: domain.StatusUnknown,
		}
	}

	chartObj := accessor.Chart()
	chartAccessor, err := chart.NewAccessor(chartObj)
	if err != nil {
		slog.Error("failed to create chart accessor", slog.Any("error", err))
		return &domain.HelmRelease{
			Name:      accessor.Name(),
			Namespace: accessor.Namespace(),
			Status:    h.convertStatus(accessor.Status()),
			Revision:  accessor.Version(),
			UpdatedAt: accessor.DeployedAt(),
			CreatedAt: accessor.DeployedAt(),
		}
	}

	chartMetadata := chartAccessor.MetadataAsMap()
	chartName := chartAccessor.Name()
	chartVersion := ""
	if v, ok := chartMetadata["version"].(string); ok {
		chartVersion = v
	}

	return &domain.HelmRelease{
		Name:      accessor.Name(),
		Namespace: accessor.Namespace(),
		Chart:     chartName,
		Version:   chartVersion,
		Values:    chartAccessor.Values(),
		Status:    h.convertStatus(accessor.Status()),
		Revision:  accessor.Version(),
		UpdatedAt: accessor.DeployedAt(),
		CreatedAt: accessor.DeployedAt(), // Helm v4 doesn't expose FirstDeployed separately
	}
}

// convertStatus converts Helm release status string to domain status
func (h *HelmAdapter) convertStatus(status string) domain.ReleaseStatus {
	switch status {
	case "deployed":
		return domain.StatusDeployed
	case "uninstalled":
		return domain.StatusUninstalled
	case "superseded":
		return domain.StatusSuperseded
	case "failed":
		return domain.StatusFailed
	case "pending-install":
		return domain.StatusPendingInstall
	case "pending-upgrade":
		return domain.StatusPendingUpgrade
	case "pending-rollback":
		return domain.StatusPendingRollback
	default:
		return domain.StatusUnknown
	}
}

// Made with Bob
