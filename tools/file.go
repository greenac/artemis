package tools

import (
	"os"
	"github.com/greenac/artemis/movie"
	"strings"
	"errors"
)

type File struct {
	Info os.FileInfo
}

func (f *File)Name() *string {
	n := f.Info.Name()
	return &n
}

func (f *File)IsDir() bool {
	return f.Info.IsDir()
}

func (f *File)IsMovie() bool {
	mt, err := f.MovieType()
	if err != nil {
		if err.Error() == "NotMovie" {
			return false
		}

		panic(err)
	}

	return mt != nil
}

func (f *File)MovieType() (*movie.MovieType, error) {
	if f.IsDir() {
			return nil, errors.New("NotMovie")
	}

	parts := strings.Split(*f.Name(), ".")
	if len(parts) == 1 {
		return nil, errors.New("NotMovie")
	}

	var movType *movie.MovieType = nil
	mt := parts[len(parts) - 1]
	for _, t := range *movie.MovieTypes() {
		if mt == string(t) {
			movType = &t
			break
		}
	}

	if 	movType == nil {
		return nil, errors.New("NotMovie")
	}

	return movType, nil
}
