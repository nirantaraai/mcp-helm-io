package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nirantaraai/mcp-helm-io/internal/core/domain"
	"github.com/nirantaraai/mcp-helm-io/internal/core/ports"
)

// RollbackReleaseUseCase handles rolling back Helm releases
type RollbackReleaseUseCase struct {
	helmPort ports.HelmPort
	logger   *slog.Logger
}

// NewRollbackReleaseUseCase creates a new RollbackReleaseUseCase
func NewRollbackReleaseUseCase(helmPort ports.HelmPort, logger *slog.Logger) *RollbackReleaseUseCase {
	return &RollbackReleaseUseCase{
		helmPort: helmPort,
		logger:   logger,
	}
}

// Execute rolls back a Helm release
func (uc *RollbackReleaseUseCase) Execute(ctx context.Context, cmd ports.RollbackReleaseCommand) (*domain.HelmRelease, error) {
	uc.logger.Info("rolling back release",
		slog.String("release", cmd.ReleaseName),
		slog.String("namespace", cmd.Namespace),
		slog.Int("revision", cmd.Revision),
	)

	// Validate command
	if err := uc.validateCommand(cmd); err != nil {
		uc.logger.Error("invalid rollback command", slog.String("error", err.Error()))
		return nil, fmt.Errorf("invalid rollback command: %w", err)
	}

	// Rollback the release using the Helm port
	release, err := uc.helmPort.RollbackRelease(ctx, cmd)
	if err != nil {
		uc.logger.Error("failed to rollback release",
			slog.String("release", cmd.ReleaseName),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to rollback release: %w", err)
	}

	uc.logger.Info("release rolled back successfully",
		slog.String("release", release.Name),
		slog.String("namespace", release.Namespace),
		slog.Int("revision", release.Revision),
	)

	return release, nil
}

// validateCommand validates the rollback command
func (uc *RollbackReleaseUseCase) validateCommand(cmd ports.RollbackReleaseCommand) error {
	if cmd.ReleaseName == "" {
		return fmt.Errorf("release name is required")
	}
	if cmd.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if cmd.Revision < 0 {
		return fmt.Errorf("revision must be non-negative")
	}
	return nil
}

// Made with Bob
