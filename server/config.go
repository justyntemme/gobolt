package server

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig parses command-line arguments and reads configuration from a YAML file.
func LoadConfig() (*ServerConfig, error) {
	// Define default values
	defaultBaseDir := "./content"
	defaultPort := ":80" // Port must have :
	defaultHostname := "localhost"

	// Command-line flags
	baseDir := flag.String("baseDir", defaultBaseDir, "Base directory for content")
	port := flag.String("port", defaultPort, "Port to listen on")

	configFile := flag.String("config", "./config.yaml", "Path to YAML configuration file")
	hostname := flag.String("hostname", defaultHostname, "hostname to listen on")
	flag.Parse()

	// Start with default config
	config := &ServerConfig{
		BaseDir:  *baseDir,
		Port:     *port,
		Hostname: *hostname,
	}

	// Load and merge configuration from YAML file, if it exists
	if _, err := os.Stat(*configFile); err == nil {
		file, err := os.Open(*configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		fileConfig := ServerConfig{}
		if err := decoder.Decode(&fileConfig); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}

		// Override defaults with YAML config
		if fileConfig.BaseDir != "" {
			config.BaseDir = fileConfig.BaseDir
		}
		if fileConfig.Port != "" {
			config.Port = fileConfig.Port
		}

	}

	return config, nil
}
