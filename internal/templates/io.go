package templates

import (
	"path/filepath"
	"runtime"
)

var (
	UpTemplatePath   = "up.tpl"
	DownTemplatePath = "down.tpl"
)

func GetUpTemplatePath() string {
	_, file, _, _ := runtime.Caller(0)

	return filepath.Join(file, "..", UpTemplatePath)
}

func GetDownTemplatePath() string {
	_, file, _, _ := runtime.Caller(0)

	return filepath.Join(file, "..", DownTemplatePath)
}
