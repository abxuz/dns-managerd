package app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/abxuz/dns-manager/assets"
	"github.com/abxuz/dns-manager/fs"
	"github.com/abxuz/dns-manager/internal/api"
	"github.com/abxuz/dns-manager/internal/dao"
	"github.com/abxuz/dns-manager/internal/middleware"
	"github.com/abxuz/dns-manager/internal/model"
	"github.com/abxuz/dns-manager/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	// providers initial
	_ "github.com/abxuz/dns-manager/provider/aliyun"

	etcd "go.etcd.io/etcd/client/v3"
)

var App = app{}

type app struct {
	flag *model.Flag
}

func (a *app) SetFlag(f *model.Flag) {
	a.flag = f
}

func (a *app) Init(ctx context.Context) error {
	switch a.flag.Storage {
	case "local":
		service.Config.SetPath(a.flag.LocalConfig)
		service.Config.SetStorage(dao.Local)
	case "etcd":
		etcdCfg := etcd.Config{
			Endpoints:   a.flag.EtcdEndpoints,
			TLS:         &tls.Config{},
			Context:     ctx,
			DialOptions: []grpc.DialOption{grpc.WithBlock()},
			Username:    a.flag.EtcdUsername,
			Password:    a.flag.EtcdPassword,
			Logger:      zap.NewNop(),
		}

		if a.flag.EtcdCA != "" {
			data, err := os.ReadFile(a.flag.EtcdCA)
			if err != nil {
				return err
			}
			etcdCfg.TLS.RootCAs = x509.NewCertPool()
			etcdCfg.TLS.RootCAs.AppendCertsFromPEM(data)
		}

		if a.flag.EtcdCert != "" {
			cert, err := tls.LoadX509KeyPair(a.flag.EtcdCert, a.flag.EtcdKey)
			if err != nil {
				return err
			}
			etcdCfg.TLS.Certificates = []tls.Certificate{cert}
		}

		client, err := etcd.New(etcdCfg)
		if err != nil {
			return err
		}
		// important, 目前的需求只有读完config就关，dao.Etcd就不可再调用了。
		// 后续有变化的话要注意这个
		defer client.Close()

		dao.Etcd.SetClient(client)
		service.Config.SetPath(a.flag.EtcdConfigKey)
		service.Config.SetStorage(dao.Etcd)
	}

	cfg, err := service.Config.GetConfig(ctx)
	if err != nil {
		return err
	}
	service.Config.SetCachedConfig(cfg)

	for _, c := range cfg.Providers {
		err = service.Provider.SetProvider(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *app) Run(ctx context.Context) error {
	if err := a.Init(ctx); err != nil {
		return err
	}

	cfg := service.Config.GetCachedConfig()

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.RecoveryWithWriter(io.Discard))
	router.Use(middleware.Api.BasicAuth)

	v1 := router.Group("/api/v1/")
	v1.Use(middleware.Api.ApiResponse)

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

	fileserver := http.FileServer(&fs.NoAutoIndexFileSystem{
		FileSystem: http.FS(assets.HtmlFs()),
	})
	router.NoRoute(gin.WrapH(fileserver))

	server := &http.Server{
		Addr:     cfg.App.Listen,
		Handler:  router,
		ErrorLog: log.New(io.Discard, "", log.LstdFlags),
	}
	return server.ListenAndServe()
}
