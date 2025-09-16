package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all configuration for our application
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Weather   WeatherConfig   `mapstructure:"weather"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Telemetry TelemetryConfig `mapstructure:"telemetry"`
	ServiceB  ServiceBConfig  `mapstructure:"service_b"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// WeatherConfig holds weather API configuration
type WeatherConfig struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

// DatabaseConfig holds database configuration (for future use)
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// TelemetryConfig holds telemetry configuration
type TelemetryConfig struct {
	ZipkinEndpoint string `mapstructure:"zipkin_endpoint"`
}

// ServiceBConfig holds Service B configuration
type ServiceBConfig struct {
	URL string `mapstructure:"url"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*Config, error) {
	// Load .env file first
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables and defaults")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.cep-temperatura")

	// Set default values
	setDefaults()

	// Enable reading from environment variables
	viper.AutomaticEnv()

	// Bind environment variables
	bindEnvVars()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		log.Println("No config file found, using defaults and environment variables")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("weather.base_url", "http://api.weatherapi.com/v1")
	viper.SetDefault("weather.api_key", "")
	viper.SetDefault("telemetry.zipkin_endpoint", "http://localhost:9411/api/v2/spans")
	viper.SetDefault("service_b.url", "http://localhost:8081")
}

// bindEnvVars binds environment variables to configuration keys
func bindEnvVars() {
	// Server configuration
	viper.BindEnv("server.port", "PORT")
	viper.BindEnv("server.host", "HOST")

	// Weather API configuration
	viper.BindEnv("weather.api_key", "WEATHER_API_KEY")
	viper.BindEnv("weather.base_url", "WEATHER_BASE_URL")

	// Telemetry configuration
	viper.BindEnv("telemetry.zipkin_endpoint", "ZIPKIN_ENDPOINT")

	// Service B configuration
	viper.BindEnv("service_b.url", "SERVICE_B_URL")
}

// GetServerAddress returns the server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Weather.APIKey == "" {
		return fmt.Errorf("weather API key is required")
	}

	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	return nil
}
