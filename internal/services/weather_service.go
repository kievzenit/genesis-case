package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/kievzenit/genesis-case/internal/config"
)

type WeatherService interface {
	GetCurrentWeatherForCity(city string) (CurrentWeatherResponse, error)
}

type weatherService struct {
	cfg *config.WeatherServiceConfig
}

func NewWeatherService(cfg *config.WeatherServiceConfig) WeatherService {
	return &weatherService{
		cfg: cfg,
	}
}

type CurrentWeatherResponse struct {
	Temperature float64
	Humidity    float64
	Condition   string
}

type currentWeatherApiResponse struct {
	Current weatherApiResponse `json:"current"`
}

type weatherApiResponse struct {
	TempCelsius float64                            `json:"temp_c"`
	Humidity    float64                            `json:"humidity"`
	Condition   currentWeatherApiResponseCondition `json:"condition"`
}

type currentWeatherApiResponseCondition struct {
	Text string `json:"text"`
}

type weatherApiErrorResponse struct {
	Error weatherApiInnerErrorResponse `json:"error"`
}

type weatherApiInnerErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var ErrCityNotFound = errors.New("city not found")

const cityNotFoundApiErrorCode = 1006

func (ws *weatherService) GetCurrentWeatherForCity(city string) (CurrentWeatherResponse, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(ws.cfg.HttpTimeout) * time.Second,
	}

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", ws.cfg.ApiKey, city)
	resp, err := httpClient.Get(url)
	if err != nil {
		return CurrentWeatherResponse{}, fmt.Errorf("failed to get weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse weatherApiErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return CurrentWeatherResponse{}, fmt.Errorf("failed to decode error response: %w", err)
		}

		if errorResponse.Error.Code == cityNotFoundApiErrorCode {
			return CurrentWeatherResponse{}, fmt.Errorf("location %s not found: %w", city, ErrCityNotFound)
		}

		return CurrentWeatherResponse{}, fmt.Errorf("weather API error: %s", errorResponse.Error.Message)
	}

	var apiResponse currentWeatherApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return CurrentWeatherResponse{}, fmt.Errorf("failed to decode weather response: %w", err)
	}
	return CurrentWeatherResponse{
		Temperature: apiResponse.Current.TempCelsius,
		Humidity:    apiResponse.Current.Humidity,
		Condition:   apiResponse.Current.Condition.Text,
	}, nil
}
