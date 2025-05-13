package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	*ServerConfig
	*WeatherServiceConfig
	*DatabaseConfig
	*CacheConfig
}

type ServerConfig struct {
	Address      string
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type WeatherServiceConfig struct {
	ApiKey      string
	HttpTimeout int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type CacheConfig struct {
	Host string
	Port int
}

func LoadConfig() (*Config, error) {
	config := getDefaultConfig()

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

	if cacheHost := os.Getenv("WAPP_CACHE_HOST"); cacheHost != "" {
		config.CacheConfig.Host = cacheHost
	}
	if cachePort := os.Getenv("WAPP_CACHE_PORT"); cachePort != "" {
		p, err := strconv.Atoi(cachePort)
		if err != nil {
			return nil, fmt.Errorf("malformed environment variable WAPP_CACHE_PORT: %w", err)
		}
		config.CacheConfig.Port = p
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
		WeatherServiceConfig: &WeatherServiceConfig{
			ApiKey:      "",
			HttpTimeout: 3,
		},
		DatabaseConfig: &DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "user",
			Password: "password",
		},
		CacheConfig: &CacheConfig{
			Host: "localhost",
			Port: 6379,
		},
	}
}
