package models

import (
	"errors"
	"github.com/greenac/artemis/logger"
	"os"
	"path"
	"strings"
)

type MovieExt string

var movieExts = [12]MovieExt{
	"mp4",
	"wmv",
	"avi",
	"mpg",
	"mpeg",
	"mov",
	"asf",
	"mkv",
	"flv",
	"m4v",
	"rmvb",
	"si",
}

var movHash *map[MovieExt]int

func MovieExtsHash() *map[MovieExt]int {
	if movHash == nil {
		mh := make(map[MovieExt]int, len(movieExts))
		for _, ext := range movieExts {
			mh[ext] = 0
		}

		movHash = &mh
	}

	return movHash
}

type File struct {
	BasePath string
	Info    os.FileInfo
	NewName string
	NewBasePath string
}

func (f *File) Name() string {
	return  f.Info.Name()
}

func (f *File) Path() string {
	return path.Join(f.BasePath, f.Name())
}

func (f *File) NewPath() string {
	return path.Join(f.NewBasePath, f.NewName)
}

func (f *File) IsDir() bool {
	return f.Info.IsDir()
}

func (f *File) IsMovie() bool {
	mt, err := f.MovieType()
	if err != nil {
		return false
	}

	return mt != nil
}

func (f *File) GetNewTotalPath() string {
	var nn string
	if f.NewName == "" {
		nn = f.Info.Name()
	} else {
		nn = f.NewName
	}

	return path.Join(f.NewBasePath, nn)
}

func (f *File) RenamePath() string {
	pts := strings.Split(f.BasePath, "/")
	if len(pts) == 1 {
		return path.Join(f.BasePath, f.NewName)
	}

	return path.Join(strings.Join(pts[:len(pts)-1], "/"), f.NewName)
}

func (f *File) MovieType() (*MovieExt, error) {
	if f.IsDir() {
		return nil, errors.New("NotMovie")
	}

	parts := strings.Split(f.Name(), ".")
	if len(parts) == 1 {
		return nil, errors.New("NotMovie")
	}

	movExt := MovieExt(strings.ToLower(parts[len(parts)-1]))
	exts := *MovieExtsHash()
	_, has := exts[movExt]
	if !has {
		logger.Warn("`MovieType` Unknown movie type:", movExt)
		return nil, errors.New("UnknownMovieType")
	}

	return &movExt, nil
}
