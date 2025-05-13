package config

type Config struct {
	*ServerConfig
	*DatabaseConfig
	*CacheConfig
}

type ServerConfig struct {
	Address string
	Port int
}

type DatabaseConfig struct {
	Host string
	Port int
	Username string
	Password string
}

type CacheConfig struct {
	Host string
	Port int
}

func LoadConfig() (*Config, error) {
	return nil, nil
}