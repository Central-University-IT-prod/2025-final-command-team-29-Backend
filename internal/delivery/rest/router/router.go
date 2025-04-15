package router

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"

	"backend/internal/delivery/rest/handlers"
	"backend/internal/infrastructure/logger"
	"backend/internal/infrastructure/repository"
	usecase "backend/internal/use_case"
	jwtconfig "backend/pkg/token_jwt"
)

func RegisterRouter(e *echo.Echo, db *gorm.DB, s3Client *minio.Client, s3Bucket string, log *logger.Logger, secret string) {
	clientRepo := repository.NewClientRepository(db)
	partnerRepo := repository.NewPartnerRepository(db)
	promoRepo := repository.NewPromoRepository(db)
	statRepo := repository.NewStatRepository(db)
	activationRepo := repository.NewActivationRepository(db)

	clientService := usecase.NewClientService(clientRepo)
	partnerService := usecase.NewPartnerService(partnerRepo)
	promoService := usecase.NewPromoService(promoRepo, s3Client, s3Bucket)
	statService := usecase.NewStatService(statRepo)
	activationService := usecase.NewActivationService(activationRepo, promoRepo)

	authHandler := handlers.NewAuthHandler(clientService, partnerService)
	partnerHandler := handlers.NewPartnerHandler(partnerService)
	statHandler := handlers.NewStatsrHandler(statService)

	activationHandler := handlers.NewActivationHandler(activationService)

	// programmatically set swagger info

	api := e.Group("/api/v1")
	api.GET("/ping", handlers.Ping)
	client := api.Group("/clients")
	authClient := client.Group("/auth")
	{
		authClient.POST("/sign-up", authHandler.SignUpClient)
		authClient.POST("/sign-in", authHandler.SignInClient)
	}

	partner := api.Group("/partners")
	{
		partner.GET("", partnerHandler.GetPartnersList)
	}
	authPartner := partner.Group("/auth")
	{
		authPartner.POST("/sign-up", authHandler.SignUpPartner)
		authPartner.POST("/sign-in", authHandler.SignInPartner)
	}

	jwtware := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtconfig.Claims)
		},
		SigningKey: []byte(secret),
	}
	partnersProtected := partner.Group("/")
	partnersProtected.Use(echojwt.WithConfig(jwtware))
	promoHandler := handlers.NewPromoHandler(promoService)
	{
		partnersProtected.GET(":partnerID/promos", promoHandler.GetPartnerPromos)
		partnersProtected.GET(":partnerID/promos/:promoID", promoHandler.GetPromo)
		partnersProtected.GET(":partnerID/promos/:promoID/img", promoHandler.GetPromoImage)

	}

	promosProtected := api.Group("/promos")
	{
		promosProtected.GET("/moderation", promoHandler.GetNotApprovedPromo)
		promosProtected.POST("/moderation", promoHandler.ModeratePromo)
	}
	promosProtected.Use(echojwt.WithConfig(jwtware))

	{
		promosProtected.GET("/:promoID/:clientID", promoHandler.GetPromoForPartnerWithClientData)
		promosProtected.GET("/:promoID", promoHandler.GetPromoForPartner)
		promosProtected.GET("", promoHandler.GetParnterPromosProtected)
		promosProtected.POST("", promoHandler.CreatePromo)
		promosProtected.DELETE("/:promoID", promoHandler.DeletePromo)
		promosProtected.POST("/activate", activationHandler.Activate)
		promosProtected.POST("/:promoID/img", promoHandler.UploadImage)
	}

	partnersStat := partner.Group("/:partnerID/stat")
	partnersStat.Use(echojwt.WithConfig(jwtware))

	{
		partnersStat.GET("", statHandler.GetStatList)
		partnersStat.GET("/:promoID", statHandler.GetStatByPromo)
		partnersStat.GET("/export", statHandler.StatExport)

	}
	log.Info("Routes successfully registered")
}
