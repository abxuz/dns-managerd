package app

import (
	"fmt"
	"io/fs"

	"github.com/xbugio/dns-manager/internal/config"
	"github.com/xbugio/dns-manager/provider"
	_ "github.com/xbugio/dns-manager/provider/aliyun"
)

type ProviderWithName struct {
	Name string
	provider.Provider
}

type App struct {
	cfg       *config.Config
	htmlFs    fs.FS
	providers map[string]*ProviderWithName
}

func New(cfg *config.Config, htmlFs fs.FS) *App {
	return &App{
		cfg:       cfg,
		htmlFs:    htmlFs,
		providers: make(map[string]*ProviderWithName),
	}
}

func (app *App) Run() error {
	if err := app.loadProviders(); err != nil {
		return err
	}
	return app.serveApi()
}

func (app *App) loadProviders() error {
	for _, cfg := range app.cfg.Providers {
		factory, ok := provider.GetProviderFactory(cfg.Provider)
		if !ok {
			return fmt.Errorf("provider %v not found", cfg.Provider)
		}
		provider, err := factory.NewProvider(cfg.Domain, cfg.Config)
		if err != nil {
			return fmt.Errorf("failed to create provider, error: %v", err.Error())
		}
		app.providers[cfg.Domain] = &ProviderWithName{
			Name:     cfg.Provider,
			Provider: provider,
		}
	}
	return nil
}
