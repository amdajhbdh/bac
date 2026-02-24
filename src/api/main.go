package main

import (
	"log"
	"os"

	"github.com/bac-unified/api/internal/handlers"
	"github.com/bac-unified/api/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize services
	authService := services.NewAuthService()
	solverService := services.NewSolverService()
	submissionService := services.NewSubmissionService()
	predictionService := services.NewPredictionService()

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
		c.JSON(200, gin.H{"status": "ok", "service": "bac-unified-api"})
	})

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(authHandler.AuthMiddleware())
		{
			// Questions
			protected.GET("/questions", handlers.ListQuestions)
			protected.GET("/questions/:id", handlers.GetQuestion)
			protected.POST("/questions", handlers.CreateQuestion)
			protected.PUT("/questions/:id", handlers.UpdateQuestion)
			protected.DELETE("/questions/:id", handlers.DeleteQuestion)

			// Submission (OCR + processing)
			protected.POST("/submit/image", submissionHandler.SubmitImage)
			protected.POST("/submit/pdf", submissionHandler.SubmitPDF)
			protected.POST("/submit/url", submissionHandler.SubmitURL)
			protected.GET("/submit/status/:job_id", submissionHandler.GetStatus)

			// Solver
			protected.POST("/solve", solverHandler.Solve)
			protected.GET("/solve/:id/steps", solverHandler.GetSteps)
			protected.GET("/solve/:id/animation", solverHandler.GetAnimation)

			// Predictions
			protected.GET("/predictions", predictionHandler.ListPredictions)
			protected.GET("/predictions/:id", predictionHandler.GetPrediction)
			protected.GET("/predictions/subject/:subject", predictionHandler.GetBySubject)

			// User
			protected.GET("/user/me", handlers.GetProfile)
			protected.PUT("/user/me", handlers.UpdateProfile)
			protected.GET("/user/progress", handlers.GetProgress)
			protected.GET("/user/stats", handlers.GetStats)
			protected.GET("/user/badges", handlers.GetBadges)

			// Practice
			protected.POST("/practice/start", handlers.StartPractice)
			protected.POST("/practice/answer", handlers.SubmitAnswer)
			protected.POST("/practice/:session/end", handlers.EndPractice)
		}

		// Public routes
		v1.GET("/predictions/latest", predictionHandler.GetLatest)
		v1.GET("/leaderboard", handlers.GetLeaderboard)
		v1.GET("/subjects", handlers.ListSubjects)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("BAC Unified API starting on port %s", port)
	r.Run(":" + port)
}
