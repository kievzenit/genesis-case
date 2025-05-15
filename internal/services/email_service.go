package services

import (
	"net/smtp"

	"github.com/google/uuid"
	"github.com/kievzenit/genesis-case/internal/config"
)

type WeatherData struct {
	City        string
	Temp        float64
	Humidity    float64
	Description string
}

type EmailService interface {
	SendConfirmationEmail(email string, token uuid.UUID) error
	SendWeatherReport(email string, city string, weatherData WeatherData) error
}

func NewEmailService(cfg *config.EmailServiceConfig) EmailService {
	return &emailService{
		smtpClient: nil,
	}
}

type emailService struct {
	smtpClient *smtp.Client
}

func (e *emailService) SendConfirmationEmail(email string, token uuid.UUID) error {
	panic("unimplemented")
}

func (e *emailService) SendWeatherReport(email string, city string, weatherData WeatherData) error {
	panic("unimplemented")
}
