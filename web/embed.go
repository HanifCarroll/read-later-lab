package web

import (
	"embed"
	"io/fs"
)

//go:embed build
var embedded embed.FS

var Assets = mustSub(embedded, "build")

func mustSub(source fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(source, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
