package embed

import (
	"embed"
	"io/fs"
	"path/filepath"
)

//go:embed scaffold
var scaffold embed.FS

func Scaffold(path string) (fs.FS, error) {
	return fs.Sub(scaffold, filepath.Join("scaffold", path))
}
