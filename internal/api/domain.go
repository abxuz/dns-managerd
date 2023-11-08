package api

import (
	"sort"

	"github.com/abxuz/dns-manager/internal/model"
	"github.com/abxuz/dns-manager/internal/service"
	"github.com/gin-gonic/gin"
)

var Domain = &aDomain{}

type aDomain struct {
}

func (a *aDomain) List(c *gin.Context) {
	list := make([]*model.Domain, 0)

	cfg := service.Config.GetCachedConfig()
	for _, p := range cfg.Providers {
		list = append(list, &model.Domain{
			Domain:   p.Domain,
			Provider: p.Provider,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Domain < list[j].Domain
	})
	c.Set("data", list)
}
