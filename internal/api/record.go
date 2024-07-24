package api

import (
	"fmt"

	"github.com/abxuz/dns-manager/internal/service"
	"github.com/abxuz/dns-manager/provider"
	"github.com/gin-gonic/gin"
)

var Record = &aRecord{}

type aRecord struct {
	baseApi
}

func (a *aRecord) List(c *gin.Context) {
	domain := c.Param("domain")

	page := provider.Page{}
	if err := c.ShouldBindQuery(&page); err != nil {
		a.Error(c, err.Error())
		return
	}

	provider, ok := service.Provider.GetProvider(domain)
	if !ok {
		a.Error(c, fmt.Sprintf("domain %v not configured", domain))
		return
	}

	list, total, err := provider.ListRecords(page)
	if err != nil {
		a.Error(c, err.Error())
		return
	}
	a.Success(c, gin.H{"list": list, "total": total})
}

func (a *aRecord) Get(c *gin.Context) {
	domain := c.Param("domain")
	id := c.Param("id")

	provider, ok := service.Provider.GetProvider(domain)
	if !ok {
		a.Error(c, fmt.Sprintf("domain %v not configured", domain))
		return
	}

	data, err := provider.GetRecord(id)
	if err != nil {
		a.Error(c, err.Error())
		return
	}
	a.Success(c, data)
}

func (a *aRecord) Add(c *gin.Context) {
	domain := c.Param("domain")

	data := &provider.Record{}
	if err := c.ShouldBindJSON(data); err != nil {
		a.Error(c, err.Error())
		return
	}

	provider, ok := service.Provider.GetProvider(domain)
	if !ok {
		a.Error(c, fmt.Sprintf("domain %v not configured", domain))
		return
	}

	if err := provider.AddRecord(data); err != nil {
		a.Error(c, err.Error())
		return
	}
	a.Success(c, nil)
}

func (a *aRecord) Delete(c *gin.Context) {
	domain := c.Param("domain")
	id := c.Param("id")

	provider, ok := service.Provider.GetProvider(domain)
	if !ok {
		a.Error(c, fmt.Sprintf("domain %v not configured", domain))
		return
	}

	if err := provider.DeleteRecord(id); err != nil {
		a.Error(c, err.Error())
		return
	}
	a.Success(c, nil)
}

func (a *aRecord) Update(c *gin.Context) {
	domain := c.Param("domain")

	data := &provider.Record{}
	if err := c.ShouldBindJSON(data); err != nil {
		a.Error(c, err.Error())
		return
	}

	provider, ok := service.Provider.GetProvider(domain)
	if !ok {
		a.Error(c, fmt.Sprintf("domain %v not configured", domain))
		return
	}

	if err := provider.ModifyRecord(data); err != nil {
		a.Error(c, err.Error())
		return
	}
	a.Success(c, nil)
}
