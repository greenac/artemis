package tools

import (
	"os"
	"path"
)

type File struct {
	Path    string
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

func (f *File) GetNewTotalPath() string {
	return path.Join(f.NewPath, f.NewName)
}

func (f *File) RenamePath() string {
	return path.Join(f.Path, f.NewName)
}
