package di

import (
	"context"
	"testing"

	"github.com/kamil5b/clean-go-vite-react/internal/platform"
)

func TestNewContainer_CreatesAllDependencies(t *testing.T) {
	// Arrange
	cfg := &platform.Config{
		Server: platform.ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
	}

	// Act
	container := NewContainer(cfg)

	// Assert
	if container == nil {
		t.Fatal("expected container to be created, got nil")
	}
	if container.Config == nil {
		t.Error("expected Config to be set")
	}
	if container.Echo == nil {
		t.Error("expected Echo to be set")
	}
	if container.Services == nil {
		t.Error("expected Services to be set")
	}
}

func TestNewContainer_Services_MessageServiceExists(t *testing.T) {
	// Arrange
	cfg := &platform.Config{}

	// Act
	container := NewContainer(cfg)

	// Assert
	if container.Services.Message == nil {
		t.Fatal("expected MessageService to be set")
	}

	// Verify the service works
	msg, err := container.Services.Message.GetMessage(context.Background())
	if err != nil {
		t.Errorf("MessageService.GetMessage failed: %v", err)
	}
	if msg != "Hello, from the golang World!" {
		t.Errorf("expected message, got %q", msg)
	}
}

func TestNewContainer_RoutesAreRegistered(t *testing.T) {
	// Arrange
	cfg := &platform.Config{}

	// Act
	container := NewContainer(cfg)

	// Assert
	if container.Echo == nil {
		t.Fatal("expected Echo to be set")
	}

	// Check that routes are registered by inspecting the echo instance
	routes := container.Echo.Routes()
	if len(routes) == 0 {
		t.Error("expected routes to be registered")
	}

	// Verify health endpoint exists
	healthRouteFound := false
	messageRouteFound := false
	for _, route := range routes {
		if route.Path == "/api/health" {
			healthRouteFound = true
		}
		if route.Path == "/api/message" {
			messageRouteFound = true
		}
	}

	if !healthRouteFound {
		t.Error("expected /api/health route to be registered")
	}
	if !messageRouteFound {
		t.Error("expected /api/message route to be registered")
	}
}

func TestNewContainer_ConfigIsStored(t *testing.T) {
	// Arrange
	expectedPort := 9000
	cfg := &platform.Config{
		Server: platform.ServerConfig{
			Port: expectedPort,
		},
	}

	// Act
	container := NewContainer(cfg)

	// Assert
	if container.Config.Server.Port != expectedPort {
		t.Errorf("expected port %d, got %d", expectedPort, container.Config.Server.Port)
	}
}

func TestNewContainer_MultipleInstances(t *testing.T) {
	// Arrange
	cfg := &platform.Config{}

	// Act
	container1 := NewContainer(cfg)
	container2 := NewContainer(cfg)

	// Assert - Each container should have its own instances
	if container1.Echo == container2.Echo {
		t.Error("expected separate Echo instances")
	}

	// Services are created fresh for each container
	if container1.Services == container2.Services {
		t.Error("expected separate Services instances")
	}
}

func TestNewContainer_ServicesDependencies(t *testing.T) {
	// Arrange
	cfg := &platform.Config{}

	// Act
	container := NewContainer(cfg)

	// Assert
	// Services should be properly initialized and wired
	services := container.Services

	if services == nil {
		t.Fatal("Services should not be nil")
	}

	// Test that all services are initialized
	if services.Message == nil {
		t.Error("MessageService should not be nil")
	}
	if services.Email == nil {
		t.Error("EmailService should not be nil")
	}
	if services.Health == nil {
		t.Error("HealthService should not be nil")
	}
}

func TestNewContainer_EchoMiddleware(t *testing.T) {
	// Arrange
	cfg := &platform.Config{}

	// Act
	container := NewContainer(cfg)

	// Assert
	if container.Echo == nil {
		t.Fatal("expected Echo to be set")
	}

	// Echo should be configured (middleware can be checked but is implementation detail)
	if container.Echo.Logger == nil {
		t.Error("expected Echo logger to be set")
	}
}

func TestContainer_CanBeUsedAsGlobalRegistry(t *testing.T) {
	// Arrange
	cfg := &platform.Config{
		Server: platform.ServerConfig{
			Port: 8080,
		},
	}

	// Act
	container := NewContainer(cfg)

	// Assert - verify we can access all parts easily
	if container.Config.Server.Port != 8080 {
		t.Error("cannot access Config from container")
	}

	if container.Echo == nil {
		t.Error("cannot access Echo from container")
	}

	if container.Services == nil {
		t.Error("cannot access Services from container")
	}
}
