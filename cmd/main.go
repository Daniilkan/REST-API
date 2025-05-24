package main

// @title TestRest API
// @version 1.0
// @description This is a REST API for managing people.
// @host localhost:8080
// @BasePath /

import (
	_ "TestRest/docs"
	"TestRest/internal/config"
	"TestRest/internal/handlers"
	"TestRest/pkg/logger"
	"TestRest/pkg/migrations"
	"TestRest/pkg/postgres"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	ctx := context.Background()
	ctx, err := logger.New(ctx)

	cfg, err := config.New()
	if err != nil {
		logger.GetLoggerFromContext(ctx).Fatal(ctx, "failed to load config", zap.Error(err))
	}
	logger.GetLoggerFromContext(ctx).Info(ctx, "Config loaded")

	if err != nil {
		logger.GetLoggerFromContext(ctx).Fatal(ctx, "failed to create logger", zap.Error(err))
		return
	}
	logger.GetLoggerFromContext(ctx).Info(ctx, "Middleware initialized")

	db, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		logger.GetLoggerFromContext(ctx).Fatal(ctx, "failed to connect to database", zap.Error(err))
		return
	}
	if err = db.Ping(ctx); err != nil {
		logger.GetLoggerFromContext(ctx).Fatal(ctx, "failed to ping database", zap.Error(err))
		return
	}
	logger.GetLoggerFromContext(ctx).Info(ctx, "Database connected")
	defer db.Close()

	if err = migrations.RunMigrations(db, "./pkg/migrations", *cfg); err != nil {
		logger.GetLoggerFromContext(ctx).Fatal(ctx, "failed to run migrations", zap.Error(err))
		return
	}
	logger.GetLoggerFromContext(ctx).Info(ctx, "Migrations applied")

	handlers.InitHandlers(db, ctx)
	logger.GetLoggerFromContext(ctx).Info(ctx, "Handlers initialized")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.Middleware(ctx))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	router.Get("/get", handlers.GetInfo)
	router.Delete("/delete", handlers.DeletePerson)
	router.Post("/post", handlers.InsertPerson)
	router.Put("/put", handlers.UpdatePerson)

	if err = http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.RESTHost, cfg.RESTPort), router); err != nil {
		logger.GetLoggerFromContext(ctx).Fatal(ctx, "failed to start server", zap.Error(err))
		return
	}
}
