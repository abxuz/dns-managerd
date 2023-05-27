package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/xbugio/dns-manager/fs"
	"github.com/xbugio/dns-manager/provider"
)

func (app *App) serveApi() error {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.RecoveryWithWriter(io.Discard))
	router.Use(app.MiddlewareBasicAuth)

	v1 := router.Group("/api/v1/")
	v1.Use(app.MiddlewareApiResponse)
	{
		v1.GET("/domain", app.ApiListDomain)

		g := v1.Group("/domain/:domain")
		{
			g.POST("/record", app.ApiAddRecord)
			g.DELETE("/record/:id", app.ApiDeleteRecord)
			g.PATCH("/record", app.ApiModifyRecord)
			g.GET("/record", app.ApiListRecord)
			g.GET("/record/:id", app.ApiGetRecord)
		}
	}

	router.NoRoute(app.FileServer())
	server := &http.Server{
		Addr:     app.cfg.App.Listen,
		Handler:  router,
		ErrorLog: log.New(io.Discard, "", log.LstdFlags),
	}
	return server.ListenAndServe()
}

func (app *App) MiddlewareBasicAuth(c *gin.Context) {
	if app.cfg.App.Auth == nil {
		c.Next()
		return
	}

	auth := app.cfg.App.Auth
	username, password, ok := c.Request.BasicAuth()
	authOk := ok && (username == auth.Username) && (password == auth.Password)
	if authOk {
		c.Next()
		return
	}

	c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
	c.AbortWithStatus(http.StatusUnauthorized)
}

func (app *App) MiddlewareApiResponse(c *gin.Context) {
	c.Next()

	err := c.Errors.Last()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"errno":  1,
			"errmsg": err.Error(),
		})
		return
	}

	obj := gin.H{
		"errno": 0,
	}
	data, exists := c.Get("data")
	if exists {
		obj["data"] = data
	}
	c.JSON(http.StatusOK, obj)
}

func (app *App) FileServer() gin.HandlerFunc {
	fileserver := http.FileServer(&fs.NoAutoIndexFileSystem{
		FileSystem: http.FS(app.htmlFs),
	})
	return gin.WrapH(fileserver)
}

func (app *App) ApiListDomain(c *gin.Context) {
	type Data struct {
		Domain   string `json:"domain"`
		Provider string `json:"provider"`
	}

	list := make([]*Data, 0)
	for domain, provider := range app.providers {
		list = append(list, &Data{
			Domain:   domain,
			Provider: provider.Name,
		})
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Domain < list[j].Domain
	})
	c.Set("data", list)
}

func (app *App) ApiAddRecord(c *gin.Context) {
	domain := c.Param("domain")

	data := &provider.Record{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.Error(err)
		return
	}

	provider, ok := app.providers[domain]
	if !ok {
		c.Error(fmt.Errorf("domain %v not configured", domain))
		return
	}

	if err := provider.AddRecord(data); err != nil {
		c.Error(err)
		return
	}
}

func (app *App) ApiDeleteRecord(c *gin.Context) {
	domain := c.Param("domain")
	id := c.Param("id")

	provider, ok := app.providers[domain]
	if !ok {
		c.Error(fmt.Errorf("domain %v not configured", domain))
		return
	}

	if err := provider.DeleteRecord(id); err != nil {
		c.Error(err)
		return
	}
}

func (app *App) ApiModifyRecord(c *gin.Context) {
	domain := c.Param("domain")

	data := &provider.Record{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.Error(err)
		return
	}

	provider, ok := app.providers[domain]
	if !ok {
		c.Error(fmt.Errorf("domain %v not configured", domain))
		return
	}

	if err := provider.ModifyRecord(data); err != nil {
		c.Error(err)
		return
	}
}

func (app *App) ApiListRecord(c *gin.Context) {
	domain := c.Param("domain")

	page := provider.Page{}
	if err := c.ShouldBindQuery(&page); err != nil {
		c.Error(err)
		return
	}

	provider, ok := app.providers[domain]
	if !ok {
		c.Error(fmt.Errorf("domain %v not configured", domain))
		return
	}

	list, total, err := provider.ListRecords(page)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", gin.H{"list": list, "total": total})
}

func (app *App) ApiGetRecord(c *gin.Context) {
	domain := c.Param("domain")
	id := c.Param("id")

	provider, ok := app.providers[domain]
	if !ok {
		c.Error(fmt.Errorf("domain %v not configured", domain))
		return
	}

	data, err := provider.GetRecord(id)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("data", data)
}
