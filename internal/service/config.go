package service

import (
	"context"

	"github.com/abxuz/dns-manager/internal/model"
	"gopkg.in/yaml.v3"
)

var Config = &sConfig{}

type ConfigStorage interface {
	GetConfig(ctx context.Context, k string) ([]byte, error)
	SetConfig(ctx context.Context, k string, v []byte) error
}

type sConfig struct {
	path      string
	storage   ConfigStorage
	cachedCfg *model.Config
}

func (s *sConfig) SetPath(p string) {
	s.path = p
}

func (s *sConfig) SetStorage(storage ConfigStorage) {
	s.storage = storage
}

func (s *sConfig) GetConfig(ctx context.Context) (*model.Config, error) {
	data, err := s.storage.GetConfig(ctx, s.path)
	if err != nil {
		return nil, err
	}

	var cfg *model.Config
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

func (s *sConfig) SetCachedConfig(cfg *model.Config) {
	s.cachedCfg = cfg
}

func (s *sConfig) GetCachedConfig() *model.Config {
	return s.cachedCfg
}
