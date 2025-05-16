package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BaseURL string
	*ServerConfig
	*JobsConfig
	*WeatherServiceConfig
	*EmailServiceConfig
	*DatabaseConfig
	*CORSConfig
}

type ServerConfig struct {
	Address      string
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type JobsConfig struct {
	EmailConfirmationInterval int
}

type WeatherServiceConfig struct {
	ApiKey      string
	HttpTimeout int
}

type EmailServiceConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	SSL      bool
}

type DatabaseConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	DatabaseName    string
	ApplyMigrations bool
}

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
}

func LoadConfig() (*Config, error) {
	config := getDefaultConfig()

	baseURL := os.Getenv("WAPP_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("missing required environment variable: WAPP_BASE_URL")
	}
	config.BaseURL = baseURL

	if addr := os.Getenv("WAPP_SERVER_ADDRESS"); addr != "" {
		config.ServerConfig.Address = addr
	}
	if port := os.Getenv("WAPP_SERVER_PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("malformed environment variable WAPP_SERVER_PORT: %w", err)
		}
		config.ServerConfig.Port = p
	}
	if readTimeout := os.Getenv("WAPP_SERVER_READ_TIMEOUT"); readTimeout != "" {
		rt, err := strconv.Atoi(readTimeout)
		if err != nil {
			return nil, fmt.Errorf("malformed environment variable WAPP_SERVER_READ_TIMEOUT: %w", err)
		}
		config.ServerConfig.ReadTimeout = rt
	}
	if writeTimeout := os.Getenv("WAPP_SERVER_WRITE_TIMEOUT"); writeTimeout != "" {
		wrt, err := strconv.Atoi(writeTimeout)
		if err != nil {
			return nil, fmt.Errorf("malformed environment variable WAPP_SERVER_WRITE_TIMEOUT: %w", err)
		}
		config.ServerConfig.WriteTimeout = wrt
	}

	if emailConfirmationInterval := os.Getenv("WAPP_EMAIL_CONFIRMATION_INTERVAL"); emailConfirmationInterval != "" {
		eci, err := strconv.Atoi(emailConfirmationInterval)
		if err != nil {
			return nil, fmt.Errorf("malformed environment variable WAPP_EMAIL_CONFIRMATION_INTERVAL: %w", err)
		}
		config.JobsConfig.EmailConfirmationInterval = eci
	}

	apiKey := os.Getenv("WAPP_WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("missing required environment variable: WAPP_WEATHER_API_KEY")
	}
	if weatherApiHttpTimeout := os.Getenv("WAPP_WEATHER_API_HTTP_TIMEOUT"); weatherApiHttpTimeout != "" {
		httpTimeout, err := strconv.Atoi(weatherApiHttpTimeout)
		if err != nil {
			return nil, fmt.Errorf("malformed environment variable WAPP_WEATHER_API_HTTP_TIMEOUT: %w", err)
		}
		config.WeatherServiceConfig.HttpTimeout = httpTimeout
	}

	config.WeatherServiceConfig.ApiKey = apiKey

	if emailHost := os.Getenv("WAPP_EMAIL_HOST"); emailHost != "" {
		config.EmailServiceConfig.Host = emailHost
	}
	if emailPort := os.Getenv("WAPP_EMAIL_PORT"); emailPort != "" {
		port, err := strconv.Atoi(emailPort)
		if err != nil {
			return nil, fmt.Errorf("malformed environment variable WAPP_EMAIL_PORT: %w", err)
		}
		config.EmailServiceConfig.Port = port
	}
	emailUsername := os.Getenv("WAPP_EMAIL_USERNAME")
	if emailUsername == "" {
		return nil, fmt.Errorf("missing required environment variable: WAPP_EMAIL_USERNAME")
	}
	config.EmailServiceConfig.Username = emailUsername
	emailPassword := os.Getenv("WAPP_EMAIL_PASSWORD")
	if emailPassword == "" {
		return nil, fmt.Errorf("missing required environment variable: WAPP_EMAIL_PASSWORD")
	}
	config.EmailServiceConfig.Password = emailPassword
	emailFrom := os.Getenv("WAPP_EMAIL_FROM")
	if emailFrom == "" {
		return nil, fmt.Errorf("missing required environment variable: WAPP_EMAIL_FROM")
	}
	config.EmailServiceConfig.From = emailFrom
	if emailSSL := os.Getenv("WAPP_EMAIL_SSL"); emailSSL != "" {
		emailSSLBool, err := strconv.ParseBool(emailSSL)
		if err != nil {
			return nil, fmt.Errorf("malformed environment variable WAPP_EMAIL_SSL: %w", err)
		}
		config.EmailServiceConfig.SSL = emailSSLBool
	}

	if dbHost := os.Getenv("WAPP_DB_HOST"); dbHost != "" {
		config.DatabaseConfig.Host = dbHost
	}
	if dbPort := os.Getenv("WAPP_DB_PORT"); dbPort != "" {
		p, err := strconv.Atoi(dbPort)
		if err != nil {
			return nil, fmt.Errorf("malformed environment variable WAPP_DB_PORT: %w", err)
		}
		config.DatabaseConfig.Port = p
	}
	if dbUser := os.Getenv("WAPP_DB_USER"); dbUser != "" {
		config.DatabaseConfig.Username = dbUser
	}
	if dbPass := os.Getenv("WAPP_DB_PASS"); dbPass != "" {
		config.DatabaseConfig.Password = dbPass
	}
	if dbName := os.Getenv("WAPP_DB_NAME"); dbName != "" {
		config.DatabaseConfig.DatabaseName = dbName
	}
	if applyMigrations := os.Getenv("WAPP_DB_APPLY_MIGRATIONS"); applyMigrations != "" {
		applyMigrationsBool, err := strconv.ParseBool(applyMigrations)
		if err != nil {
			return nil, fmt.Errorf("malformed environment variable WAPP_DB_APPLY_MIGRATIONS: %w", err)
		}
		config.DatabaseConfig.ApplyMigrations = applyMigrationsBool
	}

	if allowOrigins := os.Getenv("WAPP_CORS_ALLOW_ORIGINS"); allowOrigins != "" {
		config.CORSConfig.AllowOrigins = strings.Split(allowOrigins, ",")
	}
	if allowMethods := os.Getenv("WAPP_CORS_ALLOW_METHODS"); allowMethods != "" {
		config.CORSConfig.AllowMethods = strings.Split(allowMethods, ",")
	}
	if allowHeaders := os.Getenv("WAPP_CORS_ALLOW_HEADERS"); allowHeaders != "" {
		config.CORSConfig.AllowHeaders = strings.Split(allowHeaders, ",")
	}
	if allowCredentials := os.Getenv("WAPP_CORS_ALLOW_CREDENTIALS"); allowCredentials != "" {
		allowCredentialsBool, err := strconv.ParseBool(allowCredentials)
		if err != nil {
			return nil, fmt.Errorf("malformed environment variable WAPP_CORS_ALLOW_CREDENTIALS: %w", err)
		}
		config.CORSConfig.AllowCredentials = allowCredentialsBool
	}

	return config, nil
}

func getDefaultConfig() *Config {
	return &Config{
		ServerConfig: &ServerConfig{
			Address:      "localhost",
			Port:         8080,
			ReadTimeout:  10,
			WriteTimeout: 10,
		},
		JobsConfig: &JobsConfig{
			EmailConfirmationInterval: 1,
		},
		WeatherServiceConfig: &WeatherServiceConfig{
			ApiKey:      "",
			HttpTimeout: 3,
		},
		EmailServiceConfig: &EmailServiceConfig{
			Host: "smtp.gmail.com",
			Port: 587,
			SSL:  true,
		},
		DatabaseConfig: &DatabaseConfig{
			Host:            "localhost",
			Port:            5432,
			Username:        "user",
			Password:        "password",
			DatabaseName:    "dbname",
			ApplyMigrations: false,
		},
		CORSConfig: &CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
			AllowCredentials: false,
		},
	}
}
