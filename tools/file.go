package tools

import (
	"os"
	"path"
	"strings"
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
	pts := strings.Split(f.Path, "/")
	if len(pts) == 1 {
		return path.Join(f.Path, f.NewName)
	}

	return path.Join(strings.Join(pts[:len(pts)-1], "/"), f.NewName)
}
