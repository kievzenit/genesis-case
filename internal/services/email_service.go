package services

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/google/uuid"
	"github.com/kievzenit/genesis-case/internal/config"
	"github.com/kievzenit/genesis-case/internal/models"
	"github.com/kievzenit/genesis-case/internal/utils"
	"gopkg.in/gomail.v2"
)

type WeatherData struct {
	City        string
	Temp        float64
	Humidity    float64
	Description string
}

type EmailService interface {
	SendConfirmationEmail(email string, city string, frequency models.Frequency, token uuid.UUID) error
	SendWeatherReport(email string, city string, frequency models.Frequency, weatherData WeatherData) error
}

func NewEmailService(baseURL string, cfg *config.EmailServiceConfig) EmailService {
	return &emailService{
		from:    cfg.From,
		baseURL: baseURL,
		dialer: gomail.NewDialer(
			cfg.Host,
			cfg.Port,
			cfg.Username,
			cfg.Password,
		),
	}
}

type emailService struct {
	from    string
	baseURL string
	dialer  *gomail.Dialer
}

func convertFrequencyToReportPeriod(frequency models.Frequency) string {
	switch frequency {
	case models.Daily:
		return "24 hours"
	case models.Hourly:
		return "1 hour"
	default:
		panic("unknown frequency")
	}
}

func (e *emailService) SendConfirmationEmail(
	email string,
	city string,
	frequency models.Frequency,
	token uuid.UUID,
) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", e.from)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "Weather subscription confirmation")

	emailBodyTemplate, err := template.ParseFiles("templates/email/subscription_confirmation_email.html")
	if err != nil {
		return err
	}

	emailBodyTemplate = emailBodyTemplate.Funcs(template.FuncMap{
		"upperFirstLetter": utils.UpperFirstLetter,
	})

	var emailBodyBuf bytes.Buffer
	err = emailBodyTemplate.Execute(&emailBodyBuf, struct {
		CustomerEmail    string
		City             string
		Frequency        string
		Date             string
		ReportPeriod     string
		ConfirmationLink string
	}{
		CustomerEmail:    email,
		City:             city,
		Frequency:        string(frequency),
		Date:             time.Now().Format("January 2, 2006"),
		ReportPeriod:     convertFrequencyToReportPeriod(frequency),
		ConfirmationLink: fmt.Sprintf("http://%s/confirm/%s", e.baseURL, token.String()),
	})
	if err != nil {
		return err
	}

	msg.SetBody("text/html", emailBodyBuf.String())

	return e.dialer.DialAndSend(msg)
}

func (e *emailService) SendWeatherReport(
	email string,
	city string,
	frequency models.Frequency,
	weatherData WeatherData,
) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", e.from)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", fmt.Sprintf("Weather report for %s", city))

	emailBodyTemplate, err := template.ParseFiles("templates/email/weather_report_email.html")
	if err != nil {
		return err
	}

	emailBodyTemplate = emailBodyTemplate.Funcs(template.FuncMap{
		"upperFirstLetter": utils.UpperFirstLetter,
	})

	var emailBodyBuf bytes.Buffer
	err = emailBodyTemplate.Execute(&emailBodyBuf, struct {
		Frequency       string
		Date            string
		City            string
		FullDate        string
		Time            string
		Description     string
		Temperature     string
		Humidity        string
		UnsubscribeLink string
		CustomerEmail   string
	}{
		Frequency:       string(frequency),
		Date:            time.Now().Format("January 2, 2006"),
		City:            weatherData.City,
		FullDate:        time.Now().Format("Monday, January 2, 2006"),
		Time:            time.Now().Format("15:04"),
		Description:     weatherData.Description,
		Temperature:     fmt.Sprintf("%.2f", weatherData.Temp),
		Humidity:        fmt.Sprintf("%.2f", weatherData.Humidity),
		UnsubscribeLink: fmt.Sprintf("http://%s/unsubscribe/%s", e.baseURL, email),
		CustomerEmail:   email,
	})
	if err != nil {
		return err
	}

	msg.SetBody("text/html", emailBodyBuf.String())

	return e.dialer.DialAndSend(msg)
}
