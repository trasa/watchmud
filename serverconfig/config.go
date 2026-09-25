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
	Telnet     struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	}

	Mongo MongoConfig `yaml:"mongo"`
}

// MongoConfig says where characters are persisted. An empty Uri means the
// in-memory store, which is a real choice: it is what you want for a throwaway
// server, and it is what you get if you haven't started docker-compose yet.
// Nothing falls back the other way -- a Uri that is set and wrong fails
// startup, because silently forgetting every character is worse than not
// starting.
type MongoConfig struct {
	Uri      string `yaml:"uri"`
	Database string `yaml:"database"`
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
	// verify contents
	if len(cfg.Telnet.Host) == 0 || cfg.Telnet.Port == 0 {
		return nil, fmt.Errorf("telnet host and port must be configured")
	}
	return &cfg, nil
}
