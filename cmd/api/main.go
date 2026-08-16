// Package swaggerdocs содержит "теневые" (никогда не вызываемые) функции
// с swag-аннотациями — источник для генерации docs/swagger/openapi.json.
//
// У каждой CQRS-команды/запроса — свой задокументированный путь вида
// /api/commands/{name} или /api/queries/{name}. Эти пути РЕАЛЬНЫ (см.
// CommandByNameEndpoint/QueryByNameEndpoint в internal/interfaces/http) —
// тело запроса в них это сразу "data", без конверта. Классический вызов
// через POST /api/commands с конвертом {"command"/"query", "data"}
// продолжает работать параллельно — это просто два способа обратиться к
// одному и тому же диспатчеру.
//
// ВАЖНО: @BasePath намеренно пустой (не "/api"), потому что /auth/* живут
// на корне (см. router.go), а не под /api — префикс "/api" прописан явно
// в @Router каждой команды/запроса.
//package swaggerdocs

// @title Agro Platform API
// @version 1.0
// @description CQRS-style API. У каждой команды/запроса — свой путь вида
// @description /api/commands/{name} или /api/queries/{name}, тело = data.
// @host api.lab.note
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT из POST /auth/login. Формат: "Bearer <token>"
// func GeneralInfo() {}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samurenkoroma/agro-platform/internal/bootstrap"
	configs "github.com/samurenkoroma/agro-platform/internal/shared/config"
	"github.com/samurenkoroma/agro-platform/pkg/db"
	"github.com/samurenkoroma/agro-platform/pkg/logger"
)

func main() {
	conf := configs.LoadConfig()

	log := logger.New(conf.Logger.Level, nil)

	ctx := logger.WithContext(context.Background(), log)

	pool, err := db.NewPool(conf.Db)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		fmt.Fprintf(os.Stderr, "unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()
	log.Info("database connected")

	app, err := bootstrap.Build(ctx, pool, conf)
	if err != nil {
		log.Error("failed to build app", "error", err)
		os.Exit(1)
	}

	addr := conf.Server.ApiPort
	srv := &http.Server{
		Addr:         addr,
		Handler:      app.HTTPHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("server started", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	log.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	log.Info("server stopped")
}
