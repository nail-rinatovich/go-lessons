package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// corsMiddleware allows the React dev server (Vite, port 5173) to call the API
// directly. With the Vite proxy or the production build this is not needed,
// but it shows how CORS works.
func corsMiddleware(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Origin") == allowedOrigin {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type")
		}
		// Preflight request: answer immediately
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// serveFrontend serves the React production build and falls back to
// index.html for any non-API route (so client-side routing works).
func serveFrontend(r *gin.Engine, distDir string) {
	index := filepath.Join(distDir, "index.html")
	if _, err := os.Stat(index); err != nil {
		log.Printf("frontend build not found in %s (run `npm run build` in ../frontend)", distDir)
		return
	}
	r.Static("/assets", filepath.Join(distDir, "assets"))
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.File(index)
	})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	db, err := openDB(getEnv("DB_PATH", "./tasks.db"))
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	r := gin.Default()
	r.Use(corsMiddleware(getEnv("CORS_ORIGIN", "http://localhost:5173")))

	api := r.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	handler := &TaskHandler{store: &TaskStore{db: db}}
	handler.registerRoutes(api)

	serveFrontend(r, getEnv("FRONTEND_DIST", "../frontend/dist"))

	addr := ":" + getEnv("PORT", "8080")
	log.Printf("Server starting on http://localhost%s", addr)
	log.Fatal(r.Run(addr))
}
