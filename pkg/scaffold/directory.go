package scaffold

import (
	"io/fs"
)

type Directory struct {
	Source fs.FS
	Target string
	Data   any
}
