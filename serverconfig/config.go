package serverconfig

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config is loaded from worldfiles/server.yaml.
// See that file for an example of this configuration.
type Config struct {
	ContentPath string `yaml:"contentPath"`
	Log         struct {
		File  string
		Level string
	}
	ServerPort int `yaml:"serverPort"`
	WebPort    int `yaml:"webPort"`
	DB         struct {
		UseSSH bool `yaml:"useSSH"`
		SSH    struct {
			User    string
			Host    string
			Port    int
			KeyFile string `yaml:"keyfile"`
		}
		User     string
		Password string
		Host     string
		Port     int
		Name     string
	}
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read server config: %w", err)
	}
	var cfg Config
	// parse the configuration file
	if err := yaml.UnmarshalStrict(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &cfg, nil
}
