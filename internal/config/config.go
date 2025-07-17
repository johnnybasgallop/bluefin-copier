package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Subscriber is just a placeholder here; define its fields when you need them
type Subscriber struct {
	ID        string            `yaml:"id"`
	LotRatio  float64           `yaml:"lotRatio"`
	SymbolMap map[string]string `yaml:"symbolMap"`
}

type Config struct {
	MasterID    string       `yaml:"masterID"`
	Subscribers []Subscriber `yaml:"subscribers"`
}

// Load reads and parses your YAML file into a Config
func Load(path string) (*Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
