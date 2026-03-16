package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nirantaraai/mcp-helm-io/internal/core/domain"
	"github.com/nirantaraai/mcp-helm-io/internal/core/ports"
)

// ListReleasesUseCase handles listing Helm releases
type ListReleasesUseCase struct {
	helmPort ports.HelmPort
	logger   *slog.Logger
}

// NewListReleasesUseCase creates a new ListReleasesUseCase
func NewListReleasesUseCase(helmPort ports.HelmPort, logger *slog.Logger) *ListReleasesUseCase {
	return &ListReleasesUseCase{
		helmPort: helmPort,
		logger:   logger,
	}
}

// Execute lists all Helm releases in a namespace
func (uc *ListReleasesUseCase) Execute(ctx context.Context, namespace string) ([]*domain.HelmRelease, error) {
	uc.logger.Info("listing releases", slog.String("namespace", namespace))

	// Validate namespace
	if namespace == "" {
		namespace = "default"
	}

	// List releases using the Helm port
	releases, err := uc.helmPort.ListReleases(ctx, namespace)
	if err != nil {
		uc.logger.Error("failed to list releases",
			slog.String("namespace", namespace),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to list releases: %w", err)
	}

	uc.logger.Info("releases listed successfully",
		slog.String("namespace", namespace),
		slog.Int("count", len(releases)),
	)

	return releases, nil
}

// Made with Bob
