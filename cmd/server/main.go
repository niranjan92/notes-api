package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	_ "github.com/lib/pq"
	"github.com/qiangxue/go-rest-api/internal/auth"
	"github.com/qiangxue/go-rest-api/internal/config"
	"github.com/qiangxue/go-rest-api/internal/errors"
	"github.com/qiangxue/go-rest-api/internal/healthcheck"
	"github.com/qiangxue/go-rest-api/internal/notes"
	"github.com/qiangxue/go-rest-api/pkg/accesslog"
	"github.com/qiangxue/go-rest-api/pkg/dbcontext"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()
	// create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

	// load application configurations
	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	// connect to the database
	db, err := dbx.MustOpen("postgres", cfg.DSN)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
	db.QueryLogFunc = logDBQuery(logger)
	db.ExecLogFunc = logDBExec(logger)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
	}()

	dsn := "host=localhost user=postgres password=postgres dbname=go_restful port=5433 sslmode=disable TimeZone=Asia/Kolkata"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// build HTTP server
	address := fmt.Sprintf(":%v", cfg.ServerPort)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, gormDB, dbcontext.New(db), cfg),
	}

	// start the HTTP server with graceful shutdown
	go routing.GracefulShutdown(hs, 10*time.Second, logger.Infof)
	logger.Infof("server %v is running at %v", Version, address)
	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
}

// buildHandler sets up the HTTP routing and builds an HTTP handler.
func buildHandler(logger log.Logger, gormDB *gorm.DB, db *dbcontext.DB, cfg *config.Config) http.Handler {
	router := routing.New()

	router.Use(
		accesslog.Handler(logger),
		errors.Handler(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.AllowAll),
	)

	healthcheck.RegisterHandlers(router, Version)

	// rg := router.Group("/api/v1")
	rg := router.Group("/api")

	authHandler := auth.Handler(cfg.JWTSigningKey)
	rateLimiter := auth.RateLimiter()

	notes.RegisterHandlers(rg.Group(""),
		notes.NewService(notes.NewRepository(gormDB, db, logger), logger),
		authHandler, rateLimiter, logger)

	auth.RegisterHandlers(rg.Group(""),
		auth.NewService(auth.NewRepository(db, logger), cfg.JWTSigningKey, cfg.JWTExpiration, logger),
		logger,
	)

	return router
}

// logDBQuery returns a logging function that can be used to log SQL queries.
func logDBQuery(logger log.Logger) dbx.QueryLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, rows *sql.Rows, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB query successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB query error: %v", err)
		}
	}
}

// logDBExec returns a logging function that can be used to log SQL executions.
func logDBExec(logger log.Logger) dbx.ExecLogFunc {
	return func(ctx context.Context, t time.Duration, sql string, result sql.Result, err error) {
		if err == nil {
			logger.With(ctx, "duration", t.Milliseconds(), "sql", sql).Info("DB execution successful")
		} else {
			logger.With(ctx, "sql", sql).Errorf("DB execution error: %v", err)
		}
	}
}
