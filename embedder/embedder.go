package embedder

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Handler returns an http.Handler that serves the frontend
func Handler() http.Handler {
	isDev := os.Getenv("DEV_MODE") == "true"

	if isDev {
		return devProxyHandler()
	}
	return staticAssetsHandler()
}

// devProxyHandler returns a handler that proxies to Vite dev server
func devProxyHandler() http.Handler {
	viteURL, err := url.Parse("http://localhost:5173")
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(viteURL)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
}

// staticAssetsHandler returns a handler that serves embedded frontend assets
func staticAssetsHandler() http.Handler {
	distDir := findDistDirectory()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Serve static assets
		if strings.HasPrefix(path, "/assets/") || path == "/vite.svg" {
			filePath := filepath.Join(distDir, strings.TrimPrefix(path, "/"))
			if _, err := os.Stat(filePath); err == nil {
				http.ServeFile(w, r, filePath)
				return
			}
		}

		// SPA routing: serve index.html for all other routes
		indexPath := filepath.Join(distDir, "index.html")
		http.ServeFile(w, r, indexPath)
	})
}

// findDistDirectory locates the frontend dist directory
func findDistDirectory() string {
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
			return path
		}
	}

	// Default fallback
	return "frontend/dist"
}
