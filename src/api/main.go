package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bac-unified/api/internal/handlers"
	"github.com/bac-unified/api/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           BAC Unified API
// @version         1.0
// @description     AI-powered exam preparation platform for Mauritania's BAC C students
// @termsOfService  http://swagger.io/terms/

// @contact.name   BAC Unified Team
// @contact.url    https://bac-unified.mr
// @contact.email  support@bac-unified.mr

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Get database URL from environment or use default
	dbURL := os.Getenv("NEON_DB_URL")
	if dbURL == "" {
		dbURL = "postgresql://neondb_owner:npg_ubkCLmerS03Z@ep-fragrant-violet-ai2ew4vx-pooler.c-4.us-east-1.aws.neon.tech/neondb?sslmode=require"
	}

	// Initialize services
	db, err := services.NewDB(dbURL)
	if err != nil {
		log.Printf("Warning: DB connection failed: %v", err)
	}

	// Initialize enhanced DB service with cron jobs
	dbService, err := services.NewDBService(dbURL)
	if err != nil {
		log.Printf("Warning: Enhanced DB service failed: %v", err)
	} else {
		defer dbService.Close()
		log.Println("Enhanced DB service started with cron jobs")
	}

	authService := services.NewAuthService(db)
	solverService := services.NewSolverService(db)
	submissionService := services.NewSubmissionService(db)
	predictionService := services.NewPredictionService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	solverHandler := handlers.NewSolverHandler(solverService)
	submissionHandler := handlers.NewSubmissionHandler(submissionService)
	predictionHandler := handlers.NewPredictionHandler(predictionService)

	// Setup router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		status := "ok"
		if db == nil {
			status = "db_disconnected"
		}
		c.JSON(200, gin.H{"status": status, "service": "bac-unified-api"})
	})

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 - All routes public (no auth)
	v1 := r.Group("/api/v1")
	{
		// Auth routes (optional)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Questions - public
		v1.GET("/questions", handlers.ListQuestions)
		v1.GET("/questions/:id", handlers.GetQuestion)
		v1.POST("/questions", handlers.CreateQuestion)
		v1.PUT("/questions/:id", handlers.UpdateQuestion)
		v1.DELETE("/questions/:id", handlers.DeleteQuestion)

		// Submission (OCR + processing)
		v1.POST("/submit/image", submissionHandler.SubmitImage)
		v1.POST("/submit/pdf", submissionHandler.SubmitPDF)
		v1.POST("/submit/url", submissionHandler.SubmitURL)
		v1.GET("/submit/status/:job_id", submissionHandler.GetStatus)

		// Solver
		v1.POST("/solve", solverHandler.Solve)
		v1.GET("/solve/:id/steps", solverHandler.GetSteps)
		v1.GET("/solve/:id/animation", solverHandler.GetAnimation)

		// Predictions - public
		v1.GET("/predictions", predictionHandler.ListPredictions)
		v1.GET("/predictions/:id", predictionHandler.GetPrediction)
		v1.GET("/predictions/subject/:subject", predictionHandler.GetBySubject)
		v1.GET("/predictions/latest", predictionHandler.GetLatest)

		// User (simplified - no auth required)
		v1.GET("/user/me", handlers.GetProfile)
		v1.PUT("/user/me", handlers.UpdateProfile)
		v1.GET("/user/progress", handlers.GetProgress)
		v1.GET("/user/stats", handlers.GetStats)
		v1.GET("/user/badges", handlers.GetBadges)

		// Practice
		v1.POST("/practice/start", handlers.StartPractice)
		v1.POST("/practice/answer", handlers.SubmitAnswer)
		v1.POST("/practice/:session/end", handlers.EndPractice)

		// Public
		v1.GET("/leaderboard", handlers.GetLeaderboard)
		v1.GET("/subjects", handlers.ListSubjects)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("BAC Unified API starting on port %s", port)
	log.Printf("Database: %s", maskPassword(dbURL))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		if dbService != nil {
			dbService.Close()
		}
		os.Exit(0)
	}()

	r.Run(":" + port)
}

func maskPassword(url string) string {
	// Simple mask for logging
	for i := 0; i < len(url); i++ {
		if url[i] == ':' && i+1 < len(url) && url[i+1] == '/' {
			// Found scheme, look for @ for password
			continue
		}
		if url[i] == '@' {
			return url[:i+1] + "***@***"
		}
	}
	return url
}
