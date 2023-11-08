package assets

import (
	"embed"
	"io/fs"
)

//go:embed html
var fsHtml embed.FS

func HtmlFs() fs.FS {
	f, _ := fs.Sub(fsHtml, "html")
	return f
}
