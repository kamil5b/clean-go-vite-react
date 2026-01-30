package embedder

import (
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

// RegisterHandlers sets up frontend serving
func RegisterHandlers(e *echo.Echo) {
	isDev := os.Getenv("DEV_MODE") == "true"

	if isDev {
		// In development, proxy to Vite dev server
		setupDevProxy(e)
	} else {
		// In production, serve embedded assets
		setupStaticAssets(e)
	}
}

// setupDevProxy configures proxy to Vite dev server
func setupDevProxy(e *echo.Echo) {
	viteURL, err := url.Parse("http://localhost:5173")
	if err != nil {
		e.Logger.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(viteURL)

	// Fallback to index.html for SPA routing
	e.Any("/*", func(c echo.Context) error {
		proxy.ServeHTTP(c.Response(), c.Request())
		return nil
	})
}

// setupStaticAssets configures serving of embedded frontend assets
func setupStaticAssets(e *echo.Echo) {
	// Determine the frontend dist directory
	distDir := "frontend/dist"

	// Try to find the dist directory in multiple locations
	possiblePaths := []string{
		"frontend/dist",
		"./frontend/dist",
		"../frontend/dist",
	}

	// Get the directory of the executable
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)
		possiblePaths = append(possiblePaths,
			filepath.Join(execDir, "frontend/dist"),
			filepath.Join(execDir, "../frontend/dist"),
		)
	}

	// Find which path exists
	for _, path := range possiblePaths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			distDir = path
			break
		}
	}

	// Serve static files with proper file serving first
	e.GET("/assets/*", func(c echo.Context) error {
		return c.File(filepath.Join(distDir, c.Request().URL.Path[1:]))
	})

	e.GET("/vite.svg", func(c echo.Context) error {
		return c.File(filepath.Join(distDir, "vite.svg"))
	})

	// SPA routing: serve index.html for any unmatched routes
	e.GET("/*", func(c echo.Context) error {
		return c.File(filepath.Join(distDir, "index.html"))
	})
}
