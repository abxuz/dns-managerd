package cmd

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/abxuz/dns-manager/internal/app"
	"github.com/abxuz/dns-manager/internal/model"
)

func NewCmd() *cobra.Command {
	flag := &model.Flag{}

	c := &cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: "dns manager system",
		Args:  cobra.OnlyValidArgs,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

			app.App.SetFlag(flag)
			err := app.App.Run(ctx)
			stop()

			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		},
	}

	flags := c.Flags()
	flags.StringVar(&flag.Storage, "storage", "local", "storage type: local or etcd")

	flags.StringVar(&flag.LocalConfig, "local-config", "config.yaml", "local config file path [required for local storage]")
	c.MarkFlagFilename("local-config", "yaml", "yml")

	flags.StringArrayVar(&flag.EtcdEndpoints, "etcd-endpoints", nil, "etcd endpoints [required for etcd storage]")
	flags.StringVar(&flag.EtcdConfigKey, "etcd-config-key", "", "etcd config key [required for etcd storage]")
	flags.StringVar(&flag.EtcdCA, "etcd-ca", "ca.crt", "etcd ca cert path")
	flags.StringVar(&flag.EtcdCert, "etcd-cert", "client.crt", "etcd client cert path")
	flags.StringVar(&flag.EtcdKey, "etcd-key", "client.key", "etcd client key path")
	flags.StringVar(&flag.EtcdUsername, "etcd-username", "", "etcd username")
	flags.StringVar(&flag.EtcdPassword, "etcd-password", "", "etcd password")
	c.MarkFlagFilename("etcd-ca", "crt", "pem")
	c.MarkFlagFilename("etcd-cert", "crt", "pem")
	c.MarkFlagFilename("etcd-key", "crt", "pem")
	return c
}
