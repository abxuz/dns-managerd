package model

import (
	"fmt"

	"github.com/abxuz/b-tools/bset"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App       *AppConfig        `yaml:"app"`
	Providers []*ProviderConfig `yaml:"providers"`
}

type AppConfig struct {
	Listen string         `yaml:"listen"`
	Auth   *AppAuthConfig `yaml:"auth"`
}

type AppAuthConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type ProviderConfig struct {
	Domain   string   `yaml:"domain"`
	Provider string   `yaml:"provider"`
	Config   *RawYaml `yaml:"config"`
}

type RawYaml struct {
	node *yaml.Node
}

func (raw *RawYaml) UnmarshalYAML(value *yaml.Node) error {
	raw.node = value
	return nil
}

func (raw *RawYaml) Unmarshal(v any) error {
	return raw.node.Decode(v)
}

func (c *Config) CheckValid() error {
	if c.App == nil {
		return fmt.Errorf("app config missing")
	}
	if c.App.Listen == "" {
		return fmt.Errorf("app listen config missing")
	}

	domains := bset.NewSetString()
	for _, cfg := range c.Providers {
		if err := cfg.CheckValid(); err != nil {
			return err
		}
		if domains.Has(cfg.Domain) {
			return fmt.Errorf("duplicate domain %v in provider %v", cfg.Domain, cfg.Provider)
		}
		domains.Set(cfg.Domain)
	}

	return nil
}

func (c *ProviderConfig) CheckValid() error {
	if c.Domain == "" {
		return fmt.Errorf("missing domain in provider")
	}
	if c.Provider == "" {
		return fmt.Errorf("missing provider type for domain %v", c.Domain)
	}
	if c.Config == nil {
		return fmt.Errorf("missing provider config for domain %v", c.Domain)
	}
	return nil
}
