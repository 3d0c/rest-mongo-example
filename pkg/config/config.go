package config

import (
	"sync"
)

var (
	instance *Config
	once     sync.Once
)

// Server config secon
type Server struct {
	APIVersion string
	Address    string
	JWTSecret  string
	Static     string
}

// Logger config section
type Logger struct {
	Level       string
	AddCaller   bool
	OutputPaths []string
}

// Database config section
type Database struct {
	URI     string
	Name    string
	Timeout int
}

// Config is a complete config file
type Config struct {
	Server
	Logger
	Database
}

// TheConfig config singleton
func TheConfig() *Config {
	once.Do(func() {
		instance = new(Config)
	})

	return instance
}
