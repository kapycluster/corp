package views

import (
	"embed"
	"io/fs"
)

//go:embed static/* node_modules/preline/dist/preline.js
var static embed.FS

func Style() embed.FS {
	return static
}

func Preline() fs.FS {
	f, _ := fs.Sub(static, "node_modules/preline/dist")
	return f
}
