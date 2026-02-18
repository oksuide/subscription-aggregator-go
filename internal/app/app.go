package app

import (
	"fmt"
	"net/http"

	"subscription_aggregator_go/internal/config"
	"subscription_aggregator_go/internal/log"
	"subscription_aggregator_go/internal/storage"
	"subscription_aggregator_go/internal/subscription"

	_ "subscription_aggregator_go/docs"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type App struct {
	Config *config.Config
	Logger *zap.Logger
	DB     *sqlx.DB
	Server *http.Server
}

func New() (*App, error) {

	// Init conifg
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Init logger
	logger, _ := NewLogger(cfg)
	defer logger.Sync()
	logger.Info("App initialized", zap.String("env", cfg.Env))

	// Init storage
	dbConn, err := storage.Connect(&cfg.DBConfig)
	if err != nil {
		return nil, err
	}

	// Migrations
	if err := storage.RunMigrations(dbConn, cfg.DBConfig.MigrationsPath); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	// Init router
	router := buildRouter(logger, dbConn)

	address := fmt.Sprintf(":%v", cfg.HTTPServerConfig.ServerPort)
	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	return &App{
		Config: cfg,
		Logger: logger,
		DB:     dbConn,
		Server: srv,
	}, nil
}

func NewLogger(env *config.Config) (*zap.Logger, error) {
	var cfg zap.Config
	if env.Env == "prod" {
		cfg = zap.NewProductionConfig()
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.DisableStacktrace = true
	}

	return cfg.Build()
}

func buildRouter(logger *zap.Logger, dbConn *sqlx.DB) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(log.LoggerMiddleware(logger))

	subscriptionRepo := subscription.NewRepository(dbConn, logger)
	subscriptionService := subscription.NewService(subscriptionRepo, logger)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	public := router.Group("/api")
	subscription.RegisterRoutes(public, subscriptionService, logger)

	return router
}
