package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kievzenit/genesis-case/internal/api/handlers"
	"github.com/kievzenit/genesis-case/internal/database"
	"github.com/kievzenit/genesis-case/internal/services"
)

func RegisterRoutes(
	weatherService services.WeatherService,
	emailService services.EmailService,
	database database.Database,
	txManager *database.TransactionManger,
) *gin.Engine {
	r := gin.Default()

	r.GET("weather", handlers.GetWeatherForCityHandler(weatherService))

	r.POST("subscribe", handlers.SubscribeForWeatherHandler(
		weatherService,
		emailService,
		database,
		txManager,
	))

	return r
}
