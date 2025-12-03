package api

import (
	"time"

	"pbkk-quizlit-backend/internal/config"
	"pbkk-quizlit-backend/internal/database"
	"pbkk-quizlit-backend/internal/handlers"
	"pbkk-quizlit-backend/internal/middleware"
	"pbkk-quizlit-backend/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config *config.Config
	router *gin.Engine
}

func NewServer(cfg *config.Config) *Server {
	// Set Gin mode
	if cfg.Port == "8080" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize auth middleware with JWT secret
	if cfg.SupabaseJWTSecret != "" {
		middleware.InitAuth(cfg.SupabaseJWTSecret)
		logrus.Info("✅ Auth middleware initialized")
	} else {
		logrus.Warn("⚠️  No SUPABASE_JWT_SECRET configured, auth will not work")
	}

	// Initialize database connection
	if cfg.DatabaseURL != "" {
		if err := database.Connect(cfg.DatabaseURL); err != nil {
			logrus.WithError(err).Warn("⚠️  Failed to connect to database")
			logrus.Info("   Continuing without database (will use in-memory storage)")
		} else {
			logrus.Info("✅ Database connected successfully")
		}
	} else {
		logrus.Warn("⚠️  No DATABASE_URL configured, using in-memory storage")
	}

	router := gin.Default()

	// Configure CORS - allow all origins for flexibility
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false, // Must be false when AllowOrigins is "*"
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(corsConfig))

	server := &Server{
		config: cfg,
		router: router,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Initialize services
	fileService := services.NewFileService()
	aiService := services.NewAIServiceWithOptions(s.config.OpenAIKey, s.config.EnableRAG)
	quizService := services.NewQuizService()

	// Initialize handlers
	quizHandler := handlers.NewQuizHandler(quizService, aiService, fileService)

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "QuizLit API is running",
		})
	})

	// API routes
	api := s.router.Group("/api/v1")
	{
		// Quiz routes (protected)
		quizzes := api.Group("/quizzes")
		quizzes.Use(middleware.AuthMiddleware()) // Apply auth to all quiz routes
		{
			quizzes.POST("/upload", quizHandler.UploadFileAndGenerateQuiz)
			quizzes.POST("/generate", quizHandler.GenerateQuizFromText)
			quizzes.GET("/", quizHandler.GetAllQuizzes)
			quizzes.GET("/:id", quizHandler.GetQuiz)
			quizzes.PUT("/:id", quizHandler.UpdateQuiz)
			quizzes.DELETE("/:id", quizHandler.DeleteQuiz)

			// Quiz taking endpoints
			quizzes.GET("/take/:id", quizHandler.GetQuizForTaking)
			quizzes.POST("/submit", quizHandler.SubmitQuizAttempt)
			quizzes.GET("/attempt/:id", quizHandler.GetQuizAttempt)
			quizzes.GET("/attempts", quizHandler.ListUserAttempts)
		}
	}
}

func (s *Server) Start() error {
	return s.router.Run(":" + s.config.Port)
}
