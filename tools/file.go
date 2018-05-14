package tools

import (
	"os"
)

type File struct {
	Info    os.FileInfo
	NewName string
}

func (f *File) Name() *string {
	n := f.Info.Name()
	return &n
}

func (f *File) IsDir() bool {
	return f.Info.IsDir()
}
