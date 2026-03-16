package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nirantaraai/mcp-helm-io/internal/core/domain"
	"github.com/nirantaraai/mcp-helm-io/internal/core/ports"
)

// DeployChartUseCase handles the deployment of Helm charts
type DeployChartUseCase struct {
	helmPort ports.HelmPort
	logger   *slog.Logger
}

// NewDeployChartUseCase creates a new DeployChartUseCase
func NewDeployChartUseCase(helmPort ports.HelmPort, logger *slog.Logger) *DeployChartUseCase {
	return &DeployChartUseCase{
		helmPort: helmPort,
		logger:   logger,
	}
}

// Execute deploys a Helm chart
func (uc *DeployChartUseCase) Execute(ctx context.Context, cmd ports.DeployChartCommand) (*domain.HelmRelease, error) {
	uc.logger.Info("deploying chart",
		slog.String("chart", cmd.Chart),
		slog.String("release", cmd.ReleaseName),
		slog.String("namespace", cmd.Namespace),
	)

	// Validate command
	if err := uc.validateCommand(cmd); err != nil {
		uc.logger.Error("invalid deploy command", slog.String("error", err.Error()))
		return nil, fmt.Errorf("invalid deploy command: %w", err)
	}

	// Deploy the chart using the Helm port
	release, err := uc.helmPort.DeployChart(ctx, cmd)
	if err != nil {
		uc.logger.Error("failed to deploy chart",
			slog.String("chart", cmd.Chart),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to deploy chart: %w", err)
	}

	uc.logger.Info("chart deployed successfully",
		slog.String("release", release.Name),
		slog.String("namespace", release.Namespace),
		slog.String("status", release.Status.String()),
	)

	return release, nil
}

// validateCommand validates the deploy command
func (uc *DeployChartUseCase) validateCommand(cmd ports.DeployChartCommand) error {
	if cmd.Chart == "" {
		return fmt.Errorf("chart name is required")
	}
	if cmd.ReleaseName == "" {
		return fmt.Errorf("release name is required")
	}
	if cmd.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	return nil
}

// Made with Bob
