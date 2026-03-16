package ports

import (
	"context"

	"github.com/nirantaraai/mcp-helm-io/internal/core/domain"
)

// HelmPort defines the interface for Helm operations
type HelmPort interface {
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

// DeployChartCommand contains parameters for deploying a chart
type DeployChartCommand struct {
	Chart           string
	ReleaseName     string
	Namespace       string
	Version         string
	Values          map[string]interface{}
	CreateNamespace bool
	Wait            bool
	Timeout         int // seconds
}

// UpgradeChartCommand contains parameters for upgrading a release
type UpgradeChartCommand struct {
	Chart       string
	ReleaseName string
	Namespace   string
	Version     string
	Values      map[string]interface{}
	ResetValues bool
	ReuseValues bool
	Wait        bool
	Timeout     int // seconds
}

// RollbackReleaseCommand contains parameters for rolling back a release
type RollbackReleaseCommand struct {
	ReleaseName string
	Namespace   string
	Revision    int
	Wait        bool
	Timeout     int // seconds
}

// UninstallReleaseCommand contains parameters for uninstalling a release
type UninstallReleaseCommand struct {
	ReleaseName string
	Namespace   string
	KeepHistory bool
	Timeout     int // seconds
}

// ChartInfo represents information about a Helm chart
type ChartInfo struct {
	Name        string
	Version     string
	AppVersion  string
	Description string
	Repository  string
}

// Made with Bob
