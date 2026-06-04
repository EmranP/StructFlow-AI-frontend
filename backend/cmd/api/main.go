package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EmranP/Design-Struct-Project-AI/backend/configs"
	authHandle "github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/handler"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/middleware"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/password"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/token"
	authusecase "github.com/EmranP/Design-Struct-Project-AI/backend/internal/auth/usecase"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/container"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/infrastructure/database"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/infrastructure/logger"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/infrastructure/postgres"
	projectHandle "github.com/EmranP/Design-Struct-Project-AI/backend/internal/project/handler"
	projectusecase "github.com/EmranP/Design-Struct-Project-AI/backend/internal/project/usecase"
	httputils "github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/http"
	"github.com/EmranP/Design-Struct-Project-AI/backend/internal/shared/validator"
	userhandler "github.com/EmranP/Design-Struct-Project-AI/backend/internal/user/handler"
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
			ErrorHandler: httputils.ErrorHandler,
		},
	)
	api := app.Group("/api")

	userRepo := postgres.NewUserRepository(db)
	projectRepo := postgres.NewProjectRepository(db)

	c := &container.Container{
		DB:                db,
		Logger:            logg,
		UserRepository:    userRepo,
		ProjectRepository: projectRepo,
	}

	passwordService := password.New()

	tokenService := token.New(
		cfg.JWTSecret,
	)
	jwtMiddleware := middleware.
		NewJWTMiddleware(
			tokenService,
		)

	v := validator.New()

	authUC := authusecase.New(
		userRepo,
		passwordService,
		tokenService,
	)
	authHandler := authHandle.New(
		authUC,
		v,
	)

	auth := api.Group("/auth")

	auth.Post(
		"/register",
		authHandler.Register,
	)

	auth.Post(
		"/login",
		authHandler.Login,
	)

	api.Get(
		"/me",
		jwtMiddleware.Protected,
		userhandler.Me,
	)

	projectUC := projectusecase.New(
		projectRepo,
	)
	projectHandler := projectHandle.New(
		projectUC,
		v,
	)

	project := api.Group("/project", jwtMiddleware.Protected)

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
