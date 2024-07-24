package cmd

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/abxuz/b-tools/bhttp"
	"github.com/abxuz/dns-manager/assets"
	"github.com/abxuz/dns-manager/internal/api"
	"github.com/abxuz/dns-manager/internal/dao"
	"github.com/abxuz/dns-manager/internal/middleware"
	"github.com/abxuz/dns-manager/internal/service"

	// providers initial
	_ "github.com/abxuz/dns-manager/provider/aliyun"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func NewCmd() *cobra.Command {
	var config string
	c := &cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: "dns manager system",
		Args:  cobra.OnlyValidArgs,
		Run: func(cmd *cobra.Command, args []string) {
			Must(dao.Config.Init(config))

			cfg := dao.Config.Cfg()
			for _, c := range cfg.Providers {
				Must(service.Provider.SetProvider(c))
			}

			gin.SetMode(gin.ReleaseMode)

			g := gin.New()
			g.Use(gin.Recovery())
			if cfg.App.Auth != nil {
				g.Use(middleware.BasicAuth(cfg.App.Auth.Username, cfg.App.Auth.Password))
			}

			v1 := g.Group("/api/v1/")
			{
				v1.GET("/domain", api.Domain.List)

				g := v1.Group("/domain/:domain")
				{
					g.GET("/record", api.Record.List)
					g.GET("/record/:id", api.Record.Get)
					g.POST("/record", api.Record.Add)
					g.DELETE("/record/:id", api.Record.Delete)
					g.PATCH("/record", api.Record.Update)
				}
			}

			fileserver := http.FileServer(&bhttp.NoAutoIndexFileSystem{
				FileSystem: http.FS(assets.HtmlFs()),
			})
			g.NoRoute(gin.WrapH(fileserver))

			server := &http.Server{
				Addr:    cfg.App.Listen,
				Handler: g,
			}
			Must(server.ListenAndServe())
		},
	}

	flags := c.Flags()
	flags.StringVarP(&config, "config", "c", "config.yaml", "config file path")
	c.MarkFlagFilename("config", "yaml", "yml")

	return c
}
