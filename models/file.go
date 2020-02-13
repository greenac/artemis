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

func (f *File) IsMovie() bool {
	mt, err := f.MovieType()
	if err != nil {
		return false
	}

	return mt != nil
}

func (f *File) MovieType() (*MovieExt, error) {
	if f.IsDir() {
		return nil, errors.New("NotMovie")
	}

	parts := strings.Split(*f.Name(), ".")
	if len(parts) == 1 {
		return nil, errors.New("NotMovie")
	}

	movExt := MovieExt(strings.ToLower(parts[len(parts)-1]))
	exts := *MovieExtsHash()
	_, has := exts[movExt]
	if !has {
		logger.Error("`MovieType` Unknown movie type:", movExt)
		return nil, errors.New("UnknownMovieType")
	}

	return &movExt, nil
}
