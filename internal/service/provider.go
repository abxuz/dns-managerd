package service

import (
	"fmt"

	"github.com/abxuz/dns-manager/internal/model"
	"github.com/abxuz/dns-manager/provider"
)

var Provider = &sProvider{
	providers: make(map[string]provider.Provider),
}

type sProvider struct {
	providers map[string]provider.Provider
}

func (s *sProvider) SetProvider(cfg *model.ProviderConfig) error {
	factory, ok := provider.GetProviderFactory(cfg.Provider)
	if !ok {
		return fmt.Errorf("provider %v not found", cfg.Provider)
	}

	p, err := factory.NewProvider(cfg.Domain, cfg.Config)
	if err != nil {
		return fmt.Errorf("failed to create provider, error: %v", err.Error())
	}

	s.providers[cfg.Domain] = p
	return nil
}

func (s *sProvider) GetProvider(domain string) (provider.Provider, bool) {
	p, ok := s.providers[domain]
	return p, ok
}
