package views

import (
	"embed"
	"io/fs"
)

//go:embed static/*  node_modules/htmx.org/dist/htmx.min.js
var static embed.FS

func Style() embed.FS {
	return static
}

func HTMX() fs.FS {
	f, _ := fs.Sub(static, "node_modules/htmx.org/dist")
	return f
}
