//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/nirantaraai/go-mage-shared/dockermagex"
	"github.com/nirantaraai/go-mage-shared/golang"
)

// Build namespace for build-related targets
type Build mg.Namespace

// All builds the binary for the current platform
func (Build) All() error {
	fmt.Println("Building mcp-helm-server...")
	return golang.RunBuild(golang.BuildOptions{
		Binary:         "mcp-helm-server",
		Version:        "dev",
		OS:             "linux",
		Arch:           "amd64",
		Debug:          false,
		Packages:       []string{"./cmd/server"},
		DestinationDir: "bin",
	})
}

// Linux builds for Linux amd64
func (Build) Linux() error {
	fmt.Println("Building mcp-helm-server for Linux amd64...")
	return golang.RunBuild(golang.BuildOptions{
		Binary:         "mcp-helm-server-linux-amd64",
		Version:        "dev",
		OS:             "linux",
		Arch:           "amd64",
		Debug:          false,
		Packages:       []string{"./cmd/server"},
		DestinationDir: "bin",
	})
}

// Darwin builds for macOS arm64
func (Build) Darwin() error {
	fmt.Println("Building mcp-helm-server for macOS arm64...")
	return golang.RunBuild(golang.BuildOptions{
		Binary:         "mcp-helm-server-darwin-arm64",
		Version:        "dev",
		OS:             "darwin",
		Arch:           "arm64",
		Debug:          false,
		Packages:       []string{"./cmd/server"},
		DestinationDir: "bin",
	})
}

// DarwinAmd builds for macOS amd64
func (Build) DarwinAmd() error {
	fmt.Println("Building mcp-helm-server for macOS amd64...")
	return golang.RunBuild(golang.BuildOptions{
		Binary:         "mcp-helm-server-darwin-amd64",
		Version:        "dev",
		OS:             "darwin",
		Arch:           "amd64",
		Debug:          false,
		Packages:       []string{"./cmd/server"},
		DestinationDir: "bin",
	})
}

// Multi builds for multiple platforms
func (Build) Multi() error {
	mg.Deps(Build.Linux, Build.Darwin, Build.DarwinAmd)
	return nil
}

// Docker namespace for Docker-related targets
type Docker mg.Namespace

// Build builds the Docker image
func (Docker) Build() error {
	fmt.Println("Building Docker image...")
	if err := dockermagex.LoadConfig("docker-build.yaml"); err != nil {
		return err
	}
	return dockermagex.Build()
}

// Push pushes the Docker image
func (Docker) Push() error {
	fmt.Println("Pushing Docker image...")
	if err := dockermagex.LoadConfig("docker-build.yaml"); err != nil {
		return err
	}
	return dockermagex.Push()
}

// BuildAndPush builds and pushes the Docker image
func (d Docker) BuildAndPush() error {
	if err := d.Build(); err != nil {
		return err
	}
	return d.Push()
}

// Test namespace for testing targets
type Test mg.Namespace

// Unit runs unit tests
func (Test) Unit() error {
	fmt.Println("Running unit tests...")
	return golang.RunTests("-v")
}

// Coverage runs tests with coverage report
func (Test) Coverage() error {
	fmt.Println("Running tests with coverage...")
	return golang.RunTestsWithCoverage()
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	dirs := []string{"bin", "dist"}
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", dir, err)
		}
	}

	files := []string{"coverage.out", "mcp-helm-server"}
	for _, file := range files {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", file, err)
		}
	}

	fmt.Println("Clean complete!")
	return nil
}

// Fmt formats the code
func Fmt() error {
	fmt.Println("Formatting code...")
	return golang.RunFormat()
}

// Vet runs go vet
func Vet() error {
	fmt.Println("Running go vet...")
	return golang.RunVet()
}

// Lint runs formatting and vetting
func Lint() error {
	mg.Deps(Fmt, Vet)
	return nil
}

// Mod runs go mod tidy
func Mod() error {
	fmt.Println("Running go mod tidy...")
	return golang.Run()
}

// ModTasks runs go mod tidy and verify
func ModTasks() error {
	fmt.Println("Running go mod maintenance...")
	return golang.RunModTasks()
}

// Dev runs the server in development mode
func Dev() error {
	fmt.Println("Running server in development mode...")
	return golang.Run()
}

// Made with Bob
