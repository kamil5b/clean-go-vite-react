package web

import (
	"net/http/httputil"
	"net/url"
	"os"

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
	// This would serve embedded assets from the frontend/dist directory
	// For now, we're keeping it simple - actual implementation depends on
	// how you want to embed the frontend (using go:embed)

	e.Static("/", "frontend/dist")

	// SPA routing: serve index.html for any unmatched routes
	e.GET("*", func(c echo.Context) error {
		return c.File("frontend/dist/index.html")
	})
}
