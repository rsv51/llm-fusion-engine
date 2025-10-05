package main

import (
	"fmt"
	"llm-fusion-engine/internal/api/admin"
	"llm-fusion-engine/internal/api/v1"
	"llm-fusion-engine/internal/database"
	"llm-fusion-engine/internal/services"
	"log"

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
	// TODO: Initialize ProviderFactory
	// providerFactory := services.NewProviderFactory()
	multiProviderService := services.NewMultiProviderService(providerRouter, nil) // Pass nil for factory for now

	// 3. Initialize Handlers
	chatHandler := v1.NewChatHandler(multiProviderService)
	groupHandler := admin.NewGroupHandler(db)
	statsHandler := admin.NewStatsHandler(db)

	// 4. Setup Router
	router := gin.Default()
	
	// Serve static files from web/dist
	router.Static("/assets", "./web/dist/assets")
	router.StaticFile("/", "./web/dist/index.html")
	router.StaticFile("/favicon.ico", "./web/dist/favicon.ico")
	
	// API v1 for proxying
	v1Group := router.Group("/v1")
	{
		v1Group.POST("/chat/completions", chatHandler.ChatCompletions)
	}

	// Admin API for management
	adminGroup := router.Group("/api/admin")
	{
		adminGroup.GET("/stats", statsHandler.GetStats)
		adminGroup.POST("/groups", groupHandler.CreateGroup)
		adminGroup.GET("/groups", groupHandler.GetGroups)
		adminGroup.GET("/groups/:id", groupHandler.GetGroup)
		adminGroup.PUT("/groups/:id", groupHandler.UpdateGroup)
		adminGroup.DELETE("/groups/:id", groupHandler.DeleteGroup)
	}
	
	// NoRoute handler for SPA routing
	router.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})

	// 5. Start Server
	fmt.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}