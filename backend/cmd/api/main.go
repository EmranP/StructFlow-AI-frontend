package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EmranP/Design-Struct-Project-AI/backend/configs"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/ai"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/ai/claude"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/ai/gemini"
	aiHandle "github.com/EmranP/Design-Struct-Project-AI/backend/internal/ai/handler"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/ai/openai"
	authHandle "github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/handler"
	authMiddleware "github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/middleware"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/password"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/token"
	authUseCase "github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/usecase"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/container"
	generationHandler "github.com/EmranP/Design-Struct-Project-AI/backend/internal/generation/handler"
	generationService "github.com/EmranP/Design-Struct-Project-AI/backend/internal/generation/service"
	generationUseCase "github.com/EmranP/Design-Struct-Project-AI/backend/internal/generation/usecase"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/generation/zip"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/infrastructure/database"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/infrastructure/email"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/infrastructure/logger"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/infrastructure/postgres"
	projectHandle "github.com/EmranP/Design-Struct-Project-AI/backend/internal/project/handler"
	projectUseCase "github.com/EmranP/Design-Struct-Project-AI/backend/internal/project/usecase"
	httpUtils "github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/http"
	sharedMiddleware "github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/middleware"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/validator"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal(err)
	}
	logg, err := logger.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(
		fiber.Config{
			ErrorHandler: httpUtils.ErrorHandler,
		},
	)

	app.Use(
		sharedMiddleware.CORS(cfg),
	)

	api := app.Group("/api")

	userRepo := postgres.NewUserRepository(db)
	verifyRepo := postgres.NewVerificationRepository(db)
	sessionRepo := postgres.NewSessionRepository(db)
	projectRepo := postgres.NewProjectRepository(db)
	generationRepo := postgres.NewGenerationRepository(db)
	templateRepo := postgres.NewGeneratedTemplateRepository(db)

	c := &container.Container{
		DB:                          db,
		Logger:                      logg,
		UserRepository:              userRepo,
		ProjectRepository:           projectRepo,
		GenerationRepository:        generationRepo,
		GeneratedTemplateRepository: templateRepo,
	}
	// Service
	passwordService := password.New()
	tokenService := token.New(
		cfg.JWTSecret,
	)

	jwtMiddleware := authMiddleware.
		NewJWTMiddleware(
			tokenService,
		)

	v := validator.New()

	geminiProvider, err := gemini.New(
		context.Background(),
		cfg.AIGeminiKey,
	)
	if err != nil {
		log.Fatal(err)
	}
	claudeProvider := claude.New(
		cfg.AIClaudeKey,
	)
	gptProvider := openai.New(
		cfg.AIChatGPTKey,
	)
	aiManager := ai.NewManager(
		geminiProvider,
		claudeProvider,
		gptProvider,
	)

	generator := generationService.New(
		generationRepo,
		templateRepo,
		projectRepo,
	)
	zipService := zip.New()
	emailService := email.NewResend(
		cfg.SmtpHost,
		cfg.SmtpPort,
		cfg.SmtpEmail,
		cfg.SmtpPassword,
	)
	// UseCase
	sessionUC := authUseCase.NewSession(
		sessionRepo,
		userRepo,
		passwordService,
		tokenService,
	)
	authUC := authUseCase.New(
		userRepo,
		verifyRepo,
		sessionUC,
		passwordService,
		tokenService,
		emailService,
	)
	generationUC := generationUseCase.New(
		projectRepo,
		generationRepo,
		templateRepo,
		generator,
		zipService,
		aiManager,
	)
	authHandler := authHandle.New(
		authUC,
		sessionUC,
		v,
	)
	projectUC := projectUseCase.New(
		projectRepo,
	)
	projectHandler := projectHandle.New(
		projectUC,
		v,
	)
	generationHandler := generationHandler.New(
		generationUC,
	)
	aiHandler := aiHandle.NewAI(aiManager)

	authRoute := api.Group("/auth")
	aiRoute := api.Group("/ai", jwtMiddleware.Protected)
	projectRoute := api.Group("/project", jwtMiddleware.Protected)
	generationsRoute := api.Group(
		"/gen",
		jwtMiddleware.Protected,
	)

	// Auth Route
	authRoute.Post(
		"/register",
		authHandler.Register,
	)
	authRoute.Post(
		"/login",
		authHandler.Login,
	)
	authRoute.Post(
		"/verify-email",
		authHandler.VerifyEmail,
	)
	authRoute.Get(
		"/me",
		jwtMiddleware.Protected,
		authHandler.Me,
	)
	authRoute.Get(
		"/refresh",
		authHandler.Refresh,
	)
	authRoute.Post(
		"/logout",
		jwtMiddleware.Protected,
		authHandler.Logout,
	)

	// AI
	aiRoute.Get(
		"/models",
		aiHandler.Models,
	)

	// Project Route
	projectRoute.Post(
		"/new",
		projectHandler.Create,
	)
	projectRoute.Get(
		"/all",
		projectHandler.GetAll,
	)
	projectRoute.Delete(
		"/remove/all",
		projectHandler.DeleteAll,
	)

	projectRoute.Get(
		"/:id",
		projectHandler.GetById,
	)
	projectRoute.Patch(
		"/edit/:id",
		projectHandler.Edit,
	)
	projectRoute.Delete(
		"/remove/:id",
		projectHandler.Delete,
	)
	// Gen
	projectRoute.Post(
		"/gen/:id",
		generationHandler.Create,
	)
	generationsRoute.Get(
		"/all/:id",
		generationHandler.GetAll,
	)
	generationsRoute.Get(
		"/:id",
		generationHandler.GetByID,
	)
	// Temp
	generationsRoute.Get(
		"/:id/templates",
		generationHandler.GetTemplates,
	)
	generationsRoute.Get(
		"/download/:id",
		generationHandler.Download,
	)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "API endpoint not found",
			"path":    c.Path(),
		})
	})

	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Fatal(err)
		}
	}()

	gracefulShutdown(app, c)
}

func gracefulShutdown(
	app *fiber.App,
	c *container.Container,
) {
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	c.Logger.Info("shutdown started")

	c.DB.Close()

	if err := app.ShutdownWithContext(context.Background()); err != nil {
		c.Logger.Error("shutdown error")
	}

	c.Logger.Info("shutdown completed")
}
