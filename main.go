package main

import (
	"embed"

	"github.com/xbugio/dns-manager/cmd"
)

//go:embed html
var htmlFs embed.FS

func main() {
	cmd.NewCmd(&htmlFs).Execute()
}
