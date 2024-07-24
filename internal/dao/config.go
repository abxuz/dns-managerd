package dao

import (
	"os"

	"github.com/abxuz/dns-manager/internal/model"
	"gopkg.in/yaml.v3"
)

var Config = &dConfig{}

type dConfig struct {
	cfg *model.Config
}

func (d *dConfig) Init(p string) error {
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(&d.cfg)
	if err != nil {
		return err
	}
	return d.cfg.CheckValid()
}

func (d *dConfig) Cfg() *model.Config {
	return d.cfg
}
