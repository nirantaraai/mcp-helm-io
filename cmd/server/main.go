package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nirantaraai/mcp-helm-io/internal/adapters/helm"
	"github.com/nirantaraai/mcp-helm-io/internal/adapters/mcp"
	"github.com/nirantaraai/mcp-helm-io/internal/infrastructure"
)

func main() {
	// Load configuration
	config, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize logger
	logger := infrastructure.NewLogger(config)
	logger.Info("Starting MCP Helm Server",
		"version", "1.0.0",
		"transport", config.MCPTransport,
	)

	// Initialize Kubernetes client
	_, err = infrastructure.NewKubernetesClient(config)
	if err != nil {
		logger.Error("Failed to create Kubernetes client", "error", err)
		os.Exit(1)
	}
	logger.Info("Kubernetes client initialized successfully")

	// Initialize Helm adapter (implements ServicePort interface)
	helmService := helm.NewHelmAdapter()
	logger.Info("Helm service initialized")

	// Initialize MCP server
	mcpServer := mcp.NewMCPServer(helmService)
	logger.Info("MCP server initialized")

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		logger.Info("Starting MCP server", "transport", config.MCPTransport)
		if err := mcpServer.Start(ctx); err != nil {
			errChan <- fmt.Errorf("MCP server error: %w", err)
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-sigChan:
		logger.Info("Received shutdown signal")
	case err := <-errChan:
		logger.Error("Server error", "error", err)
	}

	// Graceful shutdown
	logger.Info("Shutting down server...")
	if err := mcpServer.Stop(ctx); err != nil {
		logger.Error("Error during shutdown", "error", err)
	}

	logger.Info("Server stopped successfully")
}

// Made with Bob
