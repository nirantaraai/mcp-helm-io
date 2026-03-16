package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nirantaraai/mcp-helm-io/internal/core/domain"
	"github.com/nirantaraai/mcp-helm-io/internal/core/ports"
)

// UpgradeChartUseCase handles the upgrade of Helm releases
type UpgradeChartUseCase struct {
	helmPort ports.HelmPort
	logger   *slog.Logger
}

// NewUpgradeChartUseCase creates a new UpgradeChartUseCase
func NewUpgradeChartUseCase(helmPort ports.HelmPort, logger *slog.Logger) *UpgradeChartUseCase {
	return &UpgradeChartUseCase{
		helmPort: helmPort,
		logger:   logger,
	}
}

// Execute upgrades a Helm release
func (uc *UpgradeChartUseCase) Execute(ctx context.Context, cmd ports.UpgradeChartCommand) (*domain.HelmRelease, error) {
	uc.logger.Info("upgrading release",
		slog.String("release", cmd.ReleaseName),
		slog.String("chart", cmd.Chart),
		slog.String("namespace", cmd.Namespace),
	)

	// Validate command
	if err := uc.validateCommand(cmd); err != nil {
		uc.logger.Error("invalid upgrade command", slog.String("error", err.Error()))
		return nil, fmt.Errorf("invalid upgrade command: %w", err)
	}

	// Upgrade the release using the Helm port
	release, err := uc.helmPort.UpgradeChart(ctx, cmd)
	if err != nil {
		uc.logger.Error("failed to upgrade release",
			slog.String("release", cmd.ReleaseName),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to upgrade release: %w", err)
	}

	uc.logger.Info("release upgraded successfully",
		slog.String("release", release.Name),
		slog.String("namespace", release.Namespace),
		slog.Int("revision", release.Revision),
		slog.String("status", release.Status.String()),
	)

	return release, nil
}

// validateCommand validates the upgrade command
func (uc *UpgradeChartUseCase) validateCommand(cmd ports.UpgradeChartCommand) error {
	if cmd.ReleaseName == "" {
		return fmt.Errorf("release name is required")
	}
	if cmd.Chart == "" {
		return fmt.Errorf("chart name is required")
	}
	if cmd.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	return nil
}

// Made with Bob
