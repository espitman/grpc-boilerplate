package gutil

import (
	"os"
)

func CreateDir(name string) {
	_ = os.Mkdir(name, 0755)
}
