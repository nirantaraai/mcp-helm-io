# Hexagonal Architecture Documentation

## Overview

This project implements **Hexagonal Architecture** (also known as Ports and Adapters pattern), which provides a clean separation between business logic and external concerns.

## Core Principles

### 1. Dependency Rule

Dependencies point **inward** toward the core:

```
External World → Adapters → Ports (Interfaces) → Services → Domain
```

- **Domain** has no dependencies
- **Services** depend only on Domain and Ports (interfaces)
- **Adapters** implement Ports and depend on external libraries
- **Infrastructure** provides cross-cutting concerns

### 2. Ports (Interfaces)

Ports define contracts between layers:

```go
// internal/core/ports/helm_port.go
type HelmPort interface {
    DeployChart(ctx context.Context, cmd DeployChartCommand) (*domain.HelmRelease, error)
    UpgradeChart(ctx context.Context, cmd UpgradeChartCommand) (*domain.HelmRelease, error)
    // ... other methods
}
```

**Benefits:**
- Services depend on abstractions, not implementations
- Easy to mock for testing
- Implementations can be swapped without changing business logic

### 3. Services (Use Cases)

Services contain pure business logic:

```go
// internal/core/services/deploy_chart.go
type DeployChartUseCase struct {
    helmPort ports.HelmPort  // Depends on interface, not implementation
    logger   *slog.Logger
}

func (uc *DeployChartUseCase) Execute(ctx context.Context, cmd ports.DeployChartCommand) (*domain.HelmRelease, error) {
    // Business logic here
    // Validation
    // Orchestration
    // Error handling
}
```

**Characteristics:**
- No external dependencies (only interfaces)
- Testable in isolation
- Single responsibility
- Framework-agnostic

### 4. Adapters (Implementations)

Adapters implement ports and handle external concerns:

```go
// internal/adapters/helm/helm_adapter.go
type HelmAdapter struct {
    settings *cli.EnvSettings
    logger   *slog.Logger
}

// Implements ports.HelmPort interface
func (h *HelmAdapter) DeployChart(ctx context.Context, cmd ports.DeployChartCommand) (*domain.HelmRelease, error) {
    // Helm SDK specific implementation
}
```

**Types of Adapters:**
- **Primary (Driving)**: MCP Server, HTTP API, CLI - drive the application
- **Secondary (Driven)**: Helm Adapter, Database - driven by the application

## Layer Breakdown

### Domain Layer (`internal/core/domain/`)

Pure business entities with no dependencies:

```go
type HelmRelease struct {
    Name      string
    Namespace string
    Chart     string
    Version   string
    Values    map[string]interface{}
    Status    ReleaseStatus
    Revision  int
    UpdatedAt time.Time
    CreatedAt time.Time
}
```

**Rules:**
- No external imports (except standard library)
- Business rules and invariants
- Value objects and entities

### Ports Layer (`internal/core/ports/`)

Interface definitions:

```go
// Repository port (driven/secondary)
type HelmPort interface {
    DeployChart(ctx context.Context, cmd DeployChartCommand) (*domain.HelmRelease, error)
}

// Service port (driving/primary)
type MCPPort interface {
    Start(ctx context.Context) error
    RegisterTools() error
}
```

**Rules:**
- Only interfaces
- No implementations
- Define contracts

### Services Layer (`internal/core/services/`)

Business logic implementation:

```go
type DeployChartUseCase struct {
    helmPort ports.HelmPort  // Dependency injection via interface
    logger   *slog.Logger
}

func NewDeployChartUseCase(helmPort ports.HelmPort, logger *slog.Logger) *DeployChartUseCase {
    return &DeployChartUseCase{
        helmPort: helmPort,
        logger:   logger,
    }
}
```

**Rules:**
- Depend only on ports (interfaces)
- Contain business logic
- Orchestrate domain objects
- Framework-agnostic

### Adapters Layer (`internal/adapters/`)

External integrations:

#### Primary Adapters (Driving)
- **MCP Server** (`internal/adapters/mcp/`): AI agent interface
- **HTTP API** (`internal/adapters/http/`): REST endpoints
- **CLI** (`internal/adapters/cli/`): Command-line interface

#### Secondary Adapters (Driven)
- **Helm Adapter** (`internal/adapters/helm/`): Helm SDK integration

```go
// Helm adapter implements HelmPort
type HelmAdapter struct {
    settings *cli.EnvSettings
    logger   *slog.Logger
}

func NewHelmAdapter(logger *slog.Logger) *HelmAdapter {
    return &HelmAdapter{
        settings: cli.New(),
        logger:   logger,
    }
}

// Implements ports.HelmPort
func (h *HelmAdapter) DeployChart(ctx context.Context, cmd ports.DeployChartCommand) (*domain.HelmRelease, error) {
    // Use Helm SDK
    actionConfig, err := h.getActionConfig(cmd.Namespace)
    // ... implementation
}
```

### Infrastructure Layer (`internal/infrastructure/`)

Cross-cutting concerns:

- **Config** (`config.go`): Configuration management
- **Logger** (`logger.go`): Structured logging
- **Kubernetes Client** (`kubernetes_client.go`): K8s connectivity

## Dependency Injection

Dependencies are injected through constructors:

```go
// main.go
func main() {
    // 1. Initialize infrastructure
    config := infrastructure.LoadConfig()
    logger := infrastructure.NewLogger(config)
    
    // 2. Initialize adapters (implementations)
    helmAdapter := helm.NewHelmAdapter(logger)
    
    // 3. Initialize services (inject interfaces)
    deployService := services.NewDeployChartUseCase(helmAdapter, logger)
    
    // 4. Initialize primary adapters
    mcpServer := mcp.NewMCPServer(deployService, logger)
    
    // 5. Start server
    mcpServer.Start(ctx)
}
```

## Testing Strategy

### Unit Testing Services

Services are easy to test because they depend on interfaces:

```go
// Mock implementation
type MockHelmPort struct {
    mock.Mock
}

func (m *MockHelmPort) DeployChart(ctx context.Context, cmd ports.DeployChartCommand) (*domain.HelmRelease, error) {
    args := m.Called(ctx, cmd)
    return args.Get(0).(*domain.HelmRelease), args.Error(1)
}

// Test
func TestDeployChartUseCase(t *testing.T) {
    mockHelm := new(MockHelmPort)
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    
    useCase := services.NewDeployChartUseCase(mockHelm, logger)
    
    mockHelm.On("DeployChart", mock.Anything, mock.Anything).
        Return(&domain.HelmRelease{Name: "test"}, nil)
    
    result, err := useCase.Execute(ctx, cmd)
    
    assert.NoError(t, err)
    assert.Equal(t, "test", result.Name)
    mockHelm.AssertExpectations(t)
}
```

### Integration Testing Adapters

Test adapters with real external systems:

```go
func TestHelmAdapter_DeployChart(t *testing.T) {
    // Use real Helm SDK
    adapter := helm.NewHelmAdapter(logger)
    
    result, err := adapter.DeployChart(ctx, cmd)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

## Benefits

### 1. Testability
- Mock interfaces easily
- Test business logic in isolation
- No need for complex test setups

### 2. Flexibility
- Swap implementations without changing business logic
- Add new adapters (gRPC, GraphQL) easily
- Support multiple transports simultaneously

### 3. Maintainability
- Clear separation of concerns
- Changes in one layer don't affect others
- Easy to understand and navigate

### 4. Independence
- Business logic independent of frameworks
- Can change databases, APIs, UIs without touching core
- Technology decisions deferred to adapters

## SOLID Principles Applied

### Single Responsibility Principle (SRP)
Each service has one reason to change:
- `DeployChartUseCase` - only changes if deploy logic changes
- `HelmAdapter` - only changes if Helm SDK changes

### Open/Closed Principle (OCP)
Open for extension, closed for modification:
- Add new adapters without modifying services
- Add new services without modifying adapters

### Liskov Substitution Principle (LSP)
Interfaces are substitutable:
- Any `HelmPort` implementation works with services
- Mock or real implementations are interchangeable

### Interface Segregation Principle (ISP)
Small, focused interfaces:
- `HelmPort` - only Helm operations
- `MCPPort` - only MCP operations

### Dependency Inversion Principle (DIP)
Depend on abstractions:
- Services depend on `ports.HelmPort` (interface)
- Not on `helm.HelmAdapter` (concrete implementation)

## Common Patterns

### Command Pattern
Commands encapsulate request parameters:

```go
type DeployChartCommand struct {
    Chart       string
    ReleaseName string
    Namespace   string
    Values      map[string]interface{}
}
```

### Repository Pattern
Adapters act as repositories:

```go
type HelmPort interface {
    DeployChart(...)
    GetReleaseStatus(...)
    ListReleases(...)
}
```

### Factory Pattern
Constructors create properly initialized objects:

```go
func NewDeployChartUseCase(helmPort ports.HelmPort, logger *slog.Logger) *DeployChartUseCase {
    return &DeployChartUseCase{
        helmPort: helmPort,
        logger:   logger,
    }
}
```

## Anti-Patterns to Avoid

### ❌ Services Depending on Concrete Implementations

```go
// BAD
type DeployChartUseCase struct {
    helmAdapter *helm.HelmAdapter  // Concrete dependency
}
```

```go
// GOOD
type DeployChartUseCase struct {
    helmPort ports.HelmPort  // Interface dependency
}
```

### ❌ Domain Depending on External Libraries

```go
// BAD
import "helm.sh/helm/v4/pkg/release"

type HelmRelease struct {
    Release *release.Release  // External dependency
}
```

```go
// GOOD
type HelmRelease struct {
    Name      string
    Namespace string
    // Pure domain fields
}
```

### ❌ Business Logic in Adapters

```go
// BAD
func (h *HelmAdapter) DeployChart(...) {
    // Validation logic here
    // Business rules here
}
```

```go
// GOOD
func (uc *DeployChartUseCase) Execute(...) {
    // Validation logic here
    // Business rules here
    return uc.helmPort.DeployChart(...)
}
```

## Conclusion

Hexagonal Architecture provides:
- **Clean separation** between business logic and infrastructure
- **High testability** through dependency injection
- **Flexibility** to change implementations
- **Maintainability** through clear boundaries

This architecture ensures the codebase remains clean, testable, and adaptable to changing requirements.