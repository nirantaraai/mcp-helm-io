package ports

import (
	"context"

	"github.com/nirantaraai/mcp-helm-io/internal/core/domain"
)

// ServicePort defines the interface for all Helm service operations
type ServicePort interface {
	// DeployChart deploys a new Helm chart
	DeployChart(ctx context.Context, cmd DeployChartCommand) (*domain.HelmRelease, error)

	// UpgradeChart upgrades an existing Helm release
	UpgradeChart(ctx context.Context, cmd UpgradeChartCommand) (*domain.HelmRelease, error)

	// RollbackRelease rolls back a release to a previous revision
	RollbackRelease(ctx context.Context, cmd RollbackReleaseCommand) (*domain.HelmRelease, error)

	// UninstallRelease uninstalls a Helm release
	UninstallRelease(ctx context.Context, cmd UninstallReleaseCommand) error

	// ListReleases lists all Helm releases in a namespace
	ListReleases(ctx context.Context, namespace string) ([]*domain.HelmRelease, error)

	// GetReleaseStatus gets the status of a specific release
	GetReleaseStatus(ctx context.Context, name, namespace string) (*domain.HelmRelease, error)

	// SearchCharts searches for charts in configured repositories
	SearchCharts(ctx context.Context, keyword string) ([]ChartInfo, error)
}

// Made with Bob
