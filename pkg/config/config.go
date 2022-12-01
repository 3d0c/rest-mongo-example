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

// Docview config section
type Docview struct {
	Path string
}

// SAP config section
type SAP struct {
	Auth         string
	DocList      string
	DocGet       string
	UserTest     string
	ValidateUser bool
}

// Config is a complete config file
type Config struct {
	Server
	Logger
	Database
	SAP
	Docview
}

// TheConfig config singleton
func TheConfig() *Config {
	once.Do(func() {
		instance = new(Config)
	})

	return instance
}
