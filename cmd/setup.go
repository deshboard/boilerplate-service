package main

import (
	"flag"
	"os"

	"github.com/goph/fxt/debug"
	"github.com/goph/fxt/log"
	"github.com/kelseyhightower/envconfig"
)

// NewConfig creates the application Config from flags and the environment.
func NewConfig() (*Config, error) {
	config := new(Config)

	// Load Config from flags first to determine environment prefix
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	config.flags(flags)

	prefix := flags.String("prefix", "", "Environment variable prefix (useful when multiple apps use the same environment)")

	flags.Parse(os.Args[1:])

	// Load Config from environment
	err := envconfig.Process(*prefix, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// NewLoggerConfig creates a logger config constructor.
func NewLoggerConfig(config *Config) (*log.Config, error) {
	c := log.NewConfig()
	f, err := log.ParseFormat(config.LogFormat)
	if err != nil {
		return nil, err
	}

	c.Format = f
	c.Debug = config.Debug
	c.Context = []interface{}{
		"environment", config.Environment,
		"service", ServiceName,
		"tag", LogTag,
	}

	return c, nil
}

// NewDebugConfig creates a debug config.
func NewDebugConfig(config *Config) *debug.Config {
	addr := config.DebugAddr

	// Listen on loopback interface in development mode
	if config.Environment == "development" && addr[0] == ':' {
		addr = "127.0.0.1" + addr
	}

	c := debug.NewConfig(addr)
	c.Debug = config.Debug

	return c
}
