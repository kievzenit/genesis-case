package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kievzenit/genesis-case/internal/services"
)

func GetWeatherForCityHandler(weatherService services.WeatherService) gin.HandlerFunc {
	return func(c *gin.Context) {
		city := c.Query("city")
		if city == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		weatherResponse, err := weatherService.GetCurrentWeatherForCity(city)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"temperature": weatherResponse.Temperature,
				"humidity":    weatherResponse.Humidity,
				"description": weatherResponse.Condition,
			})
			return
		}

		if errors.Is(err, services.ErrCityNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
		}

		c.AbortWithError(http.StatusInternalServerError, err)
	}
}
