package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds the configuration for the API Gateway
type Config struct {
	// Server configuration
	ServerPort     string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration

	// User service configuration
	UserServiceAddress string
	UserServiceTimeout time.Duration

	// CORS configuration
	CORSAllowOrigins     string
	CORSAllowMethods     string
	CORSAllowHeaders     string
	CORSAllowCredentials bool
}

// Load loads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		ReadTimeout:         getDurationEnv("READ_TIMEOUT", 15*time.Second),
		WriteTimeout:        getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:         getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
		UserServiceAddress:  getEnv("USER_SERVICE_ADDRESS", "localhost:50051"),
		UserServiceTimeout:  getDurationEnv("USER_SERVICE_TIMEOUT", 5*time.Second),
		CORSAllowOrigins:    getEnv("CORS_ALLOW_ORIGINS", "*"),
		CORSAllowMethods:    getEnv("CORS_ALLOW_METHODS", "GET, POST, PUT, PATCH, DELETE, OPTIONS"),
		CORSAllowHeaders:    getEnv("CORS_ALLOW_HEADERS", "Content-Type, Authorization"),
		CORSAllowCredentials: getBoolEnv("CORS_ALLOW_CREDENTIALS", false),
	}
}

// String returns a string representation of the config (for logging)
func (c *Config) String() string {
	return fmt.Sprintf(
		"ServerPort=%s, ReadTimeout=%v, WriteTimeout=%v, IdleTimeout=%v, UserServiceAddress=%s, UserServiceTimeout=%v",
		c.ServerPort, c.ReadTimeout, c.WriteTimeout, c.IdleTimeout, c.UserServiceAddress, c.UserServiceTimeout,
	)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv gets a duration environment variable in seconds or returns a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if seconds, err := strconv.Atoi(value); err == nil {
			return time.Duration(seconds) * time.Second
		}
	}
	return defaultValue
}

// getBoolEnv gets a boolean environment variable or returns a default value
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
