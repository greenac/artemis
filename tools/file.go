package tools

import (
	"os"
)

type File struct {
	Path string
	Info    os.FileInfo
	NewName string
	NewPath string
}

func (f *File) Name() *string {
	n := f.Info.Name()
	return &n
}

func (f *File) IsDir() bool {
	return f.Info.IsDir()
}
