package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kievzenit/genesis-case/internal/api/handlers"
	"github.com/kievzenit/genesis-case/internal/config"
	"github.com/kievzenit/genesis-case/internal/database"
	"github.com/kievzenit/genesis-case/internal/services"
)

func RegisterRoutes(
	weatherService services.WeatherService,
	emailService services.EmailService,
	database database.Database,
	txManager *database.TransactionManger,
	corsConfig *config.CORSConfig,
) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     corsConfig.AllowOrigins,
		AllowMethods:     corsConfig.AllowMethods,
		AllowHeaders:     corsConfig.AllowHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
	}))

	r.GET("weather", handlers.GetWeatherForCityHandler(weatherService))

	r.POST("subscribe", handlers.SubscribeForWeatherHandler(
		weatherService,
		emailService,
		database,
		txManager,
	))
	r.GET("confirm/:token", handlers.ConfirmSubscriptionHandler(database))
	r.GET("unsubscribe/:token", handlers.UnsubscribeHandler(database))

	return r
}
