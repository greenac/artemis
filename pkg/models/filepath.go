package models

import (
	"os"
	"path"
)

type FilePath struct {
	Path string
}

func (fp *FilePath) PathDefined() bool {
	return fp.Path != ""
}

func (fp *FilePath) PathAsBytes() *[]byte {
	p := []byte(fp.Path)
	return &p
}

func (fp *FilePath) PathAsString() string {
	return fp.Path
}

func (fp *FilePath) IsDir() (bool, error) {
	fi, err := os.Stat(fp.Path)
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}

func (fp *FilePath) FileName() string {
	p := fp.PathAsString()
	if p == "" {
		return p
	}

	chrs := []rune(p)

	if chrs[len(chrs)-1] == '/' {
		p = string(chrs[:len(chrs)-1])
	}

	_, file := path.Split(p)

	return file
}
