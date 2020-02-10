package models

import (
	"errors"
	"fmt"
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"regexp"
	"strings"
)

func FormatMovieName(f *File) (*[]byte, error) {
	nn := make([]byte, len(*f.Name()))
	copy(nn, *f.Name())
	fmt.Println("old name:", string(nn))
	name := strings.ToLower(string(nn))
	ext := ""
	if IsMovie(f) {
		parts := strings.Split(name, ".")
		ext = parts[len(parts)-1]
		name = strings.Join(parts[:len(parts)-1], ".")
	}

	re, err := regexp.Compile(`[-\s\t!@#$%^&*()[\]<>,.?~]`)
	if err != nil {
		fmt.Println("Cannot format name compiling:", err)
		return nil, err
	}

	rs := re.ReplaceAll(nn, []byte{'_'})
	fmt.Println("matched:", string(rs))

	newName := append([]byte(string(rs)), []byte(string(ext))...)
	fmt.Println("new file name:", newName)
	return &newName, nil
}

func IsMovie(f *File) bool {
	mt, err := MovieType(f)
	if err != nil {
		return false
	}

	return mt != nil
}

func MovieType(f *File) (*MovieExt, error) {
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

func MovieFiles(fh *handlers.FileHandler) *[]File {
	movieFiles := make([]File, 0)
	for _, f := range *fh.Files {
		if IsMovie(&f) {
			movieFiles = append(movieFiles, f)
		}
	}

	return &movieFiles
}

func MovieFileNames(fh *handlers.FileHandler) *[][]byte {
	mFiles := MovieFiles(fh)
	names := make([][]byte, len(*mFiles))
	for i, f := range *mFiles {
		names[i] = []byte(*f.Name())
	}

	return &names
}
