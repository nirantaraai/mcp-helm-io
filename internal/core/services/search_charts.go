package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nirantaraai/mcp-helm-io/internal/core/ports"
)

// SearchChartsUseCase handles searching for Helm charts
type SearchChartsUseCase struct {
	helmPort ports.HelmPort
	logger   *slog.Logger
}

// NewSearchChartsUseCase creates a new SearchChartsUseCase
func NewSearchChartsUseCase(helmPort ports.HelmPort, logger *slog.Logger) *SearchChartsUseCase {
	return &SearchChartsUseCase{
		helmPort: helmPort,
		logger:   logger,
	}
}

// Execute searches for Helm charts
func (uc *SearchChartsUseCase) Execute(ctx context.Context, keyword string) ([]ports.ChartInfo, error) {
	uc.logger.Info("searching charts", slog.String("keyword", keyword))

	// Validate keyword
	if keyword == "" {
		return nil, fmt.Errorf("search keyword is required")
	}

	// Search charts using the Helm port
	charts, err := uc.helmPort.SearchCharts(ctx, keyword)
	if err != nil {
		uc.logger.Error("failed to search charts",
			slog.String("keyword", keyword),
			slog.String("error", err.Error()),
		)
		return nil, fmt.Errorf("failed to search charts: %w", err)
	}

	uc.logger.Info("charts search completed",
		slog.String("keyword", keyword),
		slog.Int("count", len(charts)),
	)

	return charts, nil
}

// Made with Bob
