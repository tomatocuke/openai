package util

import (
	"path/filepath"
	"runtime"
)

var (
	rootDir string
)

func GetRootDir() string {
	if rootDir == "" {
		_, file, _, _ := runtime.Caller(0)
		rootDir = filepath.Dir(filepath.Dir(file))
	}

	return rootDir
}
