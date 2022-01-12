package config

import (
	"sync"
)

var (
	instance *Config
	once     sync.Once
)

type (
	Server struct {
		APIVersion string
		Address    string
		JWTSecret  string
	}

	Logger struct {
		Level       string
		AddCaller   bool
		OutputPaths []string
	}
)

type Config struct {
	Server
	Logger
}

// TheConfig config singleton
func TheConfig() *Config {
	once.Do(func() {
		instance = new(Config)
	})

	return instance
}
