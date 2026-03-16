package infrastructure

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	// Server configuration
	ServerPort int
	ServerHost string

	// Kubernetes configuration
	KubeConfig    string
	InCluster     bool
	KubeContext   string
	KubeNamespace string

	// Helm configuration
	HelmDriver     string
	HelmMaxHistory int

	// Logging configuration
	LogLevel  string
	LogFormat string

	// MCP configuration
	MCPTransport string // stdio, http
	MCPPort      int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		ServerPort:     getEnvAsInt("SERVER_PORT", 8080),
		ServerHost:     getEnv("SERVER_HOST", "0.0.0.0"),
		KubeConfig:     getEnv("KUBECONFIG", ""),
		InCluster:      getEnvAsBool("IN_CLUSTER", false),
		KubeContext:    getEnv("KUBE_CONTEXT", ""),
		KubeNamespace:  getEnv("KUBE_NAMESPACE", "default"),
		HelmDriver:     getEnv("HELM_DRIVER", "secret"),
		HelmMaxHistory: getEnvAsInt("HELM_MAX_HISTORY", 10),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		LogFormat:      getEnv("LOG_FORMAT", "json"),
		MCPTransport:   getEnv("MCP_TRANSPORT", "stdio"),
		MCPPort:        getEnvAsInt("MCP_PORT", 9090),
	}

	// If KUBECONFIG is not set, try default location
	if config.KubeConfig == "" && !config.InCluster {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			config.KubeConfig = fmt.Sprintf("%s/.kube/config", homeDir)
		}
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if !c.InCluster && c.KubeConfig == "" {
		return fmt.Errorf("either IN_CLUSTER must be true or KUBECONFIG must be set")
	}

	if c.ServerPort < 1 || c.ServerPort > 65535 {
		return fmt.Errorf("invalid server port: %d", c.ServerPort)
	}

	if c.MCPPort < 1 || c.MCPPort > 65535 {
		return fmt.Errorf("invalid MCP port: %d", c.MCPPort)
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	validTransports := map[string]bool{
		"stdio": true,
		"http":  true,
	}
	if !validTransports[c.MCPTransport] {
		return fmt.Errorf("invalid MCP transport: %s", c.MCPTransport)
	}

	return nil
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// Made with Bob
