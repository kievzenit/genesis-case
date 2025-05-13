package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kievzenit/genesis-case/internal/api/handlers"
	"github.com/kievzenit/genesis-case/internal/services"
)

func RegisterRoutes(weatherService services.WeatherService) *gin.Engine {
	r := gin.Default()

	r.GET("weather", handlers.GetWeatherForCityHandler(weatherService))

	return r
}
