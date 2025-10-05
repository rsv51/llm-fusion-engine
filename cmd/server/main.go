package main

import (
	"fmt"
	"llm-fusion-engine/internal/api/admin"
	"llm-fusion-engine/internal/api/v1"
	"llm-fusion-engine/internal/database"
	"llm-fusion-engine/internal/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Initializing LLM Fusion Engine...")

	// 1. Initialize Database
	db, err := database.InitDatabase("fusion.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 2. Initialize Services
	keyManager := services.NewKeyManager(db)
	providerRouter := services.NewProviderRouter(db, keyManager)
	healthChecker := services.NewHealthChecker(db)
	// TODO: Initialize ProviderFactory
	// providerFactory := services.NewProviderFactory()
	multiProviderService := services.NewMultiProviderService(providerRouter, nil) // Pass nil for factory for now

	// 3. Initialize Handlers
	chatHandler := v1.NewChatHandler(multiProviderService)
	authHandler := admin.NewAuthHandler(db)
	groupHandler := admin.NewGroupHandler(db)
	statsHandler := admin.NewStatsHandler(db)
	keyHandler := admin.NewKeyHandler(db)
	logHandler := admin.NewLogHandler(db)
	exportHandler := admin.NewExportHandler(db)
	importHandler := admin.NewImportHandler(db)
	providerHandler := admin.NewProviderHandler(db)
	modelHandler := admin.NewModelHandler(db)
	modelMappingHandler := admin.NewModelMappingHandler(db)
	healthHandler := admin.NewHealthHandler(db, healthChecker)

	// 4. Setup Router
	router := gin.Default()

	// Simple auth middleware
	authMiddleware := func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}
		// In a real app, you'd validate the token here.
		c.Next()
	}
	
	// Serve static files from web/dist
	router.Static("/assets", "./web/dist/assets")
	router.StaticFile("/", "./web/dist/index.html")
	router.StaticFile("/favicon.ico", "./web/dist/favicon.ico")
	
	// API v1 for proxying
	v1Group := router.Group("/v1")
	{
		v1Group.POST("/chat/completions", chatHandler.ChatCompletions)
	}

	// Auth API (no auth required)
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", authHandler.Login)
	}

	// Admin API for management
	adminGroup := router.Group("/api/admin")
	adminGroup.Use(authMiddleware) // Protect all admin routes
	{
		// Statistics
		adminGroup.GET("/stats", statsHandler.GetStats)
		
		// Groups
		adminGroup.POST("/groups", groupHandler.CreateGroup)
		adminGroup.GET("/groups", groupHandler.GetGroups)
		adminGroup.GET("/groups/:id", groupHandler.GetGroup)
		adminGroup.PUT("/groups/:id", groupHandler.UpdateGroup)
		adminGroup.DELETE("/groups/:id", groupHandler.DeleteGroup)
		
		// Model Mappings
		adminGroup.POST("/model-mappings", modelMappingHandler.CreateModelMapping)
		adminGroup.GET("/model-mappings", modelMappingHandler.GetModelMappings)
		adminGroup.GET("/model-mappings/:id", modelMappingHandler.GetModelMapping)
		adminGroup.PUT("/model-mappings/:id", modelMappingHandler.UpdateModelMapping)
		adminGroup.DELETE("/model-mappings/:id", modelMappingHandler.DeleteModelMapping)

		// Keys
		adminGroup.POST("/keys", keyHandler.CreateKey)
		adminGroup.GET("/keys", keyHandler.GetKeys)
		adminGroup.GET("/keys/:id", keyHandler.GetKey)
		adminGroup.PUT("/keys/:id", keyHandler.UpdateKey)
		adminGroup.DELETE("/keys/:id", keyHandler.DeleteKey)
		
		// Logs
		adminGroup.GET("/logs", logHandler.GetLogs)
		adminGroup.GET("/logs/:id", logHandler.GetLog)
		adminGroup.DELETE("/logs", logHandler.DeleteLogs)
		
		// Export
		adminGroup.GET("/export/all", exportHandler.ExportAll)

		// Import
		adminGroup.POST("/import/all", importHandler.ImportAll)

		// Providers
		adminGroup.POST("/providers", providerHandler.CreateProvider)
		adminGroup.GET("/providers", providerHandler.GetProviders)
		adminGroup.GET("/providers/:id", providerHandler.GetProvider)
		adminGroup.PUT("/providers/:id", providerHandler.UpdateProvider)
		adminGroup.DELETE("/providers/:id", providerHandler.DeleteProvider)

		// Models
		adminGroup.POST("/models", modelHandler.CreateModel)
		adminGroup.GET("/models", modelHandler.GetModels)
		adminGroup.GET("/models/:id", modelHandler.GetModel)
		adminGroup.PUT("/models/:id", modelHandler.UpdateModel)
		adminGroup.DELETE("/models/:id", modelHandler.DeleteModel)
		adminGroup.POST("/models/:id/clone", modelHandler.CloneModel)

		// Health Checks
		adminGroup.POST("/health/providers/:id", healthHandler.CheckProviderHealth)
		adminGroup.POST("/health/providers", healthHandler.CheckAllProvidersHealth)

		// User account management
		adminGroup.PUT("/account/profile", authHandler.UpdateProfile)
	}
	
	// NoRoute handler for SPA routing
	router.NoRoute(func(c *gin.Context) {
		// Only serve index.html for GET requests that are not API calls
		if c.Request.Method == "GET" && !strings.HasPrefix(c.Request.URL.Path, "/api/") && !strings.HasPrefix(c.Request.URL.Path, "/v1/") {
			c.File("./web/dist/index.html")
			return
		}
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	// 5. Start Server
	fmt.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}