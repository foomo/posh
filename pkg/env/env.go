package env

import (
	"os"
	"path"
)

const projectRoot = "PROJECT_ROOT"

func ProjectRoot() string {
	return os.Getenv(projectRoot)
}

func SetProjectRoot(v string) error {
	return os.Setenv(projectRoot, v)
}

func Path(elem ...string) string {
	return path.Join(ProjectRoot(), path.Join(elem...))
}
