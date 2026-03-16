package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nirantaraai/mcp-helm-io/internal/core/domain"
	"github.com/nirantaraai/mcp-helm-io/internal/core/ports"
)

// GetReleaseStatusUseCase handles getting the status of a Helm release
type GetReleaseStatusUseCase struct {
	helmPort ports.HelmPort
	logger   *slog.Logger
}

// NewGetReleaseStatusUseCase creates a new GetReleaseStatusUseCase
func NewGetReleaseStatusUseCase(helmPort ports.HelmPort, logger *slog.Logger) *GetReleaseStatusUseCase {
	return &GetReleaseStatusUseCase{
		helmPort: helmPort,
		logger:   logger,
	}
}

// Execute gets the status of a Helm release
func (uc *GetReleaseStatusUseCase) Execute(ctx context.Context, name, namespace string) (*domain.HelmRelease, error) {
	uc.logger.Info("getting release status",
		slog.String("release", name),
		slog.String("namespace", namespace),
	)

	// Validate inputs
	if err := uc.validateInputs(name, namespace); err != nil {
		uc.logger.Error("invalid inputs", slog.String("error", err.Error()))
		return nil, fmt.Errorf("invalid inputs: %w", err)
	}

	// Get release status using the Helm port
	release, err := uc.helmPort.GetReleaseStatus(ctx, name, namespace)
	if err != nil {
		uc.logger.Error("failed to get release status",
			slog.String("release", name),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to get release status: %w", err)
	}

	uc.logger.Info("release status retrieved successfully",
		slog.String("release", release.Name),
		slog.String("namespace", release.Namespace),
		slog.String("status", release.Status.String()),
	)

	return release, nil
}

// validateInputs validates the inputs
func (uc *GetReleaseStatusUseCase) validateInputs(name, namespace string) error {
	if name == "" {
		return fmt.Errorf("release name is required")
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	return nil
}

// Made with Bob
