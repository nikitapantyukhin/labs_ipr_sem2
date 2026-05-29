package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sport_platform/application/controllers"
	"sport_platform/application/models/claims"
	"sport_platform/internal/configuration"
	"sport_platform/internal/env_loader"
	"sport_platform/internal/jwt"
	"sport_platform/internal/middleware"
	"sport_platform/internal/minio_config"
	"sport_platform/internal/password"
	"sport_platform/internal/service_wrapper"
	"sport_platform/internal/sqlc/db"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Application struct {
	engine        *gin.Engine
	wrapper       *service_wrapper.Wrapper
	configuration configuration.IConfiguration
}

var appStartTime = time.Now()

func (appl *Application) GetEnv() error {
	appEnvLoader := env_loader.CreateLoaderFromEnv()

	var dbConfig db.Config
	if err := appEnvLoader.LoadDataIntoStruct(&dbConfig); err != nil {
		return err
	}

	var passwordConfig password.PasswordConfig
	if err := appEnvLoader.LoadDataIntoStruct(&passwordConfig); err != nil {
		return err
	}

	var jwtConfig jwt.JwtConfig
	if err := appEnvLoader.LoadDataIntoStruct(&jwtConfig); err != nil {
		return err
	}

	var minioConfig minio_config.MinioConfig
	if err := appEnvLoader.LoadDataIntoStruct(&minioConfig); err != nil {
		return err
	}

	appl.configuration.
		AddConfiguration(&dbConfig).
		AddConfiguration(&passwordConfig).
		AddConfiguration(&jwtConfig).
		AddConfiguration(&minioConfig)

	return nil
}

func (appl *Application) ConstructClients() error {
	dbConfig, dbConfigGetterError := appl.configuration.Get(&db.Config{})
	if dbConfigGetterError != nil {
		return dbConfigGetterError
	}

	dbClient, dbConnectionError := db.CreateConnection(dbConfig.(*db.Config), context.Background())
	if dbConnectionError != nil {
		return dbConnectionError
	}
	appl.wrapper.Db = dbClient

	minioConfig, minioConfigGetterError := appl.configuration.Get(&minio_config.MinioConfig{})
	if minioConfigGetterError != nil {
		return minioConfigGetterError
	}

	minioClient, minioConnectionError := minio.New(minioConfig.(*minio_config.MinioConfig).Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(
			minioConfig.(*minio_config.MinioConfig).AccessKeyID,
			minioConfig.(*minio_config.MinioConfig).SecretKey,
			"",
		),
		Secure: minioConfig.(*minio_config.MinioConfig).UseSSL,
	})
	if minioConnectionError != nil {
		return minioConnectionError
	}
	if err := minio_config.InitBuckets(context.Background(), minioClient, minioConfig.(*minio_config.MinioConfig).BucketName); err != nil {
		return err
	}
	appl.wrapper.Minio = minioClient

	passwordConfig, passwordConfigGetterError := appl.configuration.Get(&password.PasswordConfig{})
	if passwordConfigGetterError != nil {
		return passwordConfigGetterError
	}
	appl.wrapper.PasswordHandler = password.CreateHandler(passwordConfig.(*password.PasswordConfig))

	jwtConfig, jwtConfigGetterError := appl.configuration.Get(&jwt.JwtConfig{})
	if jwtConfigGetterError != nil {
		return jwtConfigGetterError
	}
	appl.wrapper.JwtHandler = jwt.CreateHandler[claims.UserClaims](jwtConfig.(*jwt.JwtConfig))

	return nil
}

func (appl *Application) Configure(engine *gin.Engine) error {
	envGetterError := appl.GetEnv()
	if envGetterError != nil {
		return envGetterError
	}

	clientsConstructionError := appl.ConstructClients()
	if clientsConstructionError != nil {
		return clientsConstructionError
	}

	if demoSeedEnabled() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := appl.SeedDemoData(ctx); err != nil {
			return err
		}
	}

	engine.Use(corsMiddleware())
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	engine.GET("/metrics", metricsHandler)

	engine.Use(middleware.AuthMiddleware(appl.wrapper))

	controllers.UserController(engine, appl.wrapper)

	controllers.ClubController(engine, appl.wrapper)

	controllers.WorkoutController(engine, appl.wrapper)

	controllers.ClubJoinRequestController(engine, appl.wrapper)
	return nil
}

func (appl *Application) Run() {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	port := getHTTPPort()
	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: engine,
	}

	if err := appl.Configure(engine); err != nil {
		panic(err)
	}

	go func() {
		fmt.Printf("Server started on http://localhost:%s\n", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Received a shutdown call. Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := appl.wrapper.Close(); err != nil {
		panic(err)
	}

	fmt.Println("Closed all connections. Shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func CreateApplication() *Application {
	engine := gin.New()

	return &Application{
		engine:        engine,
		wrapper:       &service_wrapper.Wrapper{},
		configuration: configuration.CreateConfiguration(),
	}
}

func getHTTPPort() string {
	port := strings.TrimSpace(os.Getenv("APP_PORT"))
	if port == "" {
		return "8080"
	}
	return port
}

func demoSeedEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("DEMO_SEED")))
	return value == "true" || value == "1" || value == "yes"
}

func corsMiddleware() gin.HandlerFunc {
	allowedOrigins := map[string]bool{
		"http://localhost:3000": true,
		"http://127.0.0.1:3000": true,
	}

	if rawOrigins := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")); rawOrigins != "" {
		allowedOrigins = make(map[string]bool)
		for _, origin := range strings.Split(rawOrigins, ",") {
			trimmedOrigin := strings.TrimSpace(origin)
			if trimmedOrigin != "" {
				allowedOrigins[trimmedOrigin] = true
			}
		}
	}

	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if origin != "" && (allowedOrigins["*"] || allowedOrigins[origin]) {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Vary", "Origin")
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
			ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		}

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func metricsHandler(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	ctx.String(http.StatusOK, strings.Join([]string{
		"# HELP sport_platform_backend_up Backend availability marker.",
		"# TYPE sport_platform_backend_up gauge",
		"sport_platform_backend_up 1",
		"# HELP sport_platform_backend_uptime_seconds Backend process uptime in seconds.",
		"# TYPE sport_platform_backend_uptime_seconds gauge",
		fmt.Sprintf("sport_platform_backend_uptime_seconds %.0f", time.Since(appStartTime).Seconds()),
		"",
	}, "\n"))
}
