package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nirantaraai/mcp-helm-io/internal/core/ports"
)

// UninstallReleaseUseCase handles uninstalling Helm releases
type UninstallReleaseUseCase struct {
	helmPort ports.HelmPort
	logger   *slog.Logger
}

// NewUninstallReleaseUseCase creates a new UninstallReleaseUseCase
func NewUninstallReleaseUseCase(helmPort ports.HelmPort, logger *slog.Logger) *UninstallReleaseUseCase {
	return &UninstallReleaseUseCase{
		helmPort: helmPort,
		logger:   logger,
	}
}

// Execute uninstalls a Helm release
func (uc *UninstallReleaseUseCase) Execute(ctx context.Context, cmd ports.UninstallReleaseCommand) error {
	uc.logger.Info("uninstalling release",
		slog.String("release", cmd.ReleaseName),
		slog.String("namespace", cmd.Namespace),
	)

	// Validate command
	if err := uc.validateCommand(cmd); err != nil {
		uc.logger.Error("invalid uninstall command", slog.String("error", err.Error()))
		return fmt.Errorf("invalid uninstall command: %w", err)
	}

	// Uninstall the release using the Helm port
	if err := uc.helmPort.UninstallRelease(ctx, cmd); err != nil {
		uc.logger.Error("failed to uninstall release",
			slog.String("release", cmd.ReleaseName),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to uninstall release: %w", err)
	}

	uc.logger.Info("release uninstalled successfully",
		slog.String("release", cmd.ReleaseName),
		slog.String("namespace", cmd.Namespace),
	)

	return nil
}

// validateCommand validates the uninstall command
func (uc *UninstallReleaseUseCase) validateCommand(cmd ports.UninstallReleaseCommand) error {
	if cmd.ReleaseName == "" {
		return fmt.Errorf("release name is required")
	}
	if cmd.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	return nil
}

// Made with Bob
