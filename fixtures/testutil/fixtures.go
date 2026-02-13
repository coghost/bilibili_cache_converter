package testutil

import (
	"path/filepath"
	"runtime"
)

// GetProjectRoot returns the absolute path to the project root
func GetProjectRoot() string {
	_, b, _, _ := runtime.Caller(0) //nolint:dogsled
	return filepath.Dir(filepath.Dir(filepath.Dir(b)))
}
