package cmd

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xbugio/dns-manager/internal/app"
	"github.com/xbugio/dns-manager/internal/config"
)

type Cmd struct {
	cobra.Command

	htmlFs fs.FS
	config string
}

func NewCmd(htmlFs fs.FS) *Cmd {
	c := &Cmd{
		Command: cobra.Command{
			Use:   filepath.Base(os.Args[0]),
			Short: "dns manager system",
			Args:  cobra.OnlyValidArgs,
		},
		htmlFs: htmlFs,
	}
	c.Command.Run = c.Run
	c.Flags().StringVarP(&c.config, "config", "c", "config.yaml", "config file path")
	c.MarkFlagFilename("config")
	return c
}

func (c *Cmd) Run(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfigFromFile(c.config)
	if err != nil {
		c.PrintErrln(err)
		return
	}

	htmlFs, err := fs.Sub(c.htmlFs, "html")
	if err != nil {
		c.PrintErrln(err)
		return
	}

	if err := app.New(cfg, htmlFs).Run(); err != nil {
		c.PrintErrln(err)
	}
}
