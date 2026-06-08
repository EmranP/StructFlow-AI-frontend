package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EmranP/Design-Struct-Project-AI/backend/configs"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/ai/gemini"
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
	aiService, err := gemini.New(
		context.Background(),
		cfg.AIKey,
	)
	if err != nil {
		log.Fatal(err)
	}

	generator := generationService.New(
		generationRepo,
		templateRepo,
		projectRepo,
		aiService,
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

	auth := api.Group("/auth")
	project := api.Group("/project", jwtMiddleware.Protected)
	generations := api.Group(
		"/gen",
		jwtMiddleware.Protected,
	)

	// Auth Route
	auth.Post(
		"/register",
		authHandler.Register,
	)
	auth.Post(
		"/login",
		authHandler.Login,
	)
	auth.Post(
		"/verify-email",
		authHandler.VerifyEmail,
	)
	auth.Get(
		"/me",
		jwtMiddleware.Protected,
		authHandler.Me,
	)
	auth.Get(
		"/refresh",
		authHandler.Refresh,
	)
	auth.Post(
		"/logout",
		jwtMiddleware.Protected,
		authHandler.Logout,
	)

	// Project Route
	project.Post(
		"/new",
		projectHandler.Create,
	)
	project.Get(
		"/all",
		projectHandler.GetAll,
	)
	project.Delete(
		"/remove/all",
		projectHandler.DeleteAll,
	)

	project.Get(
		"/:id",
		projectHandler.GetById,
	)
	project.Patch(
		"/edit/:id",
		projectHandler.Edit,
	)
	project.Delete(
		"/remove/:id",
		projectHandler.Delete,
	)
	// Gen
	project.Post(
		"/gen/:id",
		generationHandler.Create,
	)
	generations.Get(
		"/all/:id",
		generationHandler.GetAll,
	)
	generations.Get(
		"/:id",
		generationHandler.GetByID,
	)
	// Temp
	generations.Get(
		"/:id/templates",
		generationHandler.GetTemplates,
	)
	generations.Get(
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
