package application

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/minio-go/v7"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"

	_ "backend/docs"
	"backend/internal/delivery/rest/router"
	"backend/internal/infrastructure"
	"backend/internal/infrastructure/database"
	"backend/internal/infrastructure/logger"
	"backend/internal/infrastructure/s3"
	"backend/internal/utils"
	tokenjwt "backend/pkg/token_jwt"
)

type Application struct {
	e        *echo.Echo
	Address  string
	Db       *gorm.DB
	S3Client *minio.Client
	S3Bucket string
	Secret   string
	Logger   *logger.Logger
}

func NewApplication(config *infrastructure.Config) *Application {
	e := echo.New()
	l := logger.NewLogger() // Создаём логгер

	db, err := database.NewPostgresDB(config, l)
	if err != nil {
		l.Errorf("Failed to connect to database: %s", err.Error())
		return nil
	}

	// Инициализация MinIO
	s3Client, s3Bucket, err := s3.NewMinioClient(config, l)
	if err != nil {
		l.Errorf("Failed to initialize MinIO: %s", err.Error())
		return nil
	}

	// Инициализация JWT
	tokenjwt.InitJWTKey(config.Other.JWTKey)

	return &Application{
		e:        e,
		Address:  config.Server.Address,
		Db:       db,
		S3Client: s3Client,
		S3Bucket: s3Bucket,
		Secret:   config.Other.JWTKey,
		Logger:   l, // Добавляем логгер
	}
}

// RunServer запускает сервер
func (a *Application) RunServer() error {
	e := initServer(a)
	a.Logger.Info("Starting server on " + a.Address)

	if err := e.Start(a.Address); err != nil {
		if err != http.ErrServerClosed {
			a.Logger.Errorf("Failed to start server: %s", err.Error())
		}
	}
	return nil
}

func initServer(a *Application) *echo.Echo {
	e := a.e
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = utils.NewValidator()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// Теперь логгер передаётся корректно
	router.RegisterRouter(e, a.Db, a.S3Client, a.S3Bucket, a.Logger, a.Secret)

	return e
}
