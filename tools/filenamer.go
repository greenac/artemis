package tools

import "strings"

type FileNamer struct {
	Path FilePath
	File File
}

func (fn *FileNamer) fileName() *[]byte {
	if !fn.Path.PathDefined() {
		panic("File Path Not Set")
	}

	parts := strings.Split(string(fn.Path.Path), "/")
	name := parts[len(parts)-1]
	rv := []byte(name)
	return &rv
}
