package jobs

import (
	"context"

	"github.com/kievzenit/genesis-case/internal/database"
	"github.com/kievzenit/genesis-case/internal/database/repositories"
	"github.com/kievzenit/genesis-case/internal/models"
	"github.com/kievzenit/genesis-case/internal/services"
)

type SendWeatherReportJob struct {
	weatherService         services.WeatherService
	emailService           services.EmailService
	subscriptionRepository repositories.SubscriptionRepository
}

func NewSendWeatherReportJob(
	weatherService services.WeatherService,
	emailService services.EmailService,
	database database.Database,
) *SendWeatherReportJob {
	return &SendWeatherReportJob{
		weatherService:         weatherService,
		emailService:           emailService,
		subscriptionRepository: repositories.NewSubscriptionRepository(database),
	}
}

func (j *SendWeatherReportJob) Run(frequency models.Frequency) {
	ctx := context.Background()

	subscriptions, err := j.subscriptionRepository.GetSubscriptionsByFrequencyContext(ctx, frequency)
	if err != nil {
		return
	}

	weathers := make(map[string]services.CurrentWeatherResponse)
	for _, subscription := range subscriptions {
		weather, ok := weathers[subscription.City]
		if !ok {
			weather, err = j.weatherService.GetCurrentWeatherForCity(subscription.City)
			if err != nil {
				continue
			}
		}

		err = j.emailService.SendWeatherReport(subscription.Email, subscription.City, services.WeatherData{
			Temp:        weather.Temperature,
			Humidity:    weather.Humidity,
			Description: weather.Condition,
		})
		if err != nil {
			continue
		}
	}
}
