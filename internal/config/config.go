package config

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

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

type AppConfig struct {
	Listen string `yaml:"listen"`
	Auth   *struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"auth"`
}

type ProviderConfig struct {
	Domain   string   `yaml:"domain"`
	Provider string   `yaml:"provider"`
	Config   *RawYaml `yaml:"config"`
}

type Config struct {
	App       *AppConfig        `yaml:"app"`
	Providers []*ProviderConfig `yaml:"providers"`
}

func LoadConfigFromBytes(b []byte) (*Config, error) {
	return LoadConfigFromReader(bytes.NewReader(b))
}

func LoadConfigFromString(b string) (*Config, error) {
	return LoadConfigFromBytes([]byte(b))
}

func LoadConfigFromFile(p string) (*Config, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LoadConfigFromReader(f)
}

func LoadConfigFromReader(r io.Reader) (*Config, error) {
	c := &Config{}
	err := yaml.NewDecoder(r).Decode(c)
	if err != nil {
		return nil, err
	}
	if err := c.CheckValid(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) CheckValid() error {
	if c.App == nil {
		return fmt.Errorf("app config missing")
	}
	if c.App.Listen == "" {
		return fmt.Errorf("app listen config missing")
	}

	m := make(map[string]struct{})
	for _, cfg := range c.Providers {
		if err := cfg.CheckValid(); err != nil {
			return err
		}
		if _, ok := m[cfg.Domain]; ok {
			return fmt.Errorf("duplicate domain %v in provider %v", cfg.Domain, cfg.Provider)
		}
		m[cfg.Domain] = struct{}{}
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
