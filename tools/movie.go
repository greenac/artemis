package tools

import (
	"github.com/greenac/artemis/movie"
	"errors"
	"strings"
	"regexp"
	"fmt"
)

func FormatMovieName(f *File) (*[]byte, error) {
	nn := make([]byte, len(*f.Name()))
	copy(nn, *f.Name())
	fmt.Println("old name:", string(nn))
	name := strings.ToLower(string(nn))
	ext := ""
	if IsMovie(f) {
		parts := strings.Split(name, ".")
		ext = parts[len(parts) - 1]
		name = strings.Join(parts[:len(parts) - 1], ".")
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
		if err.Error() == "NotMovie" {
			return false
		}

		panic(err)
	}

	return mt != nil
}

func MovieType(f *File) (*movie.MovieType, error) {
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

func MovieFiles(fh *FileHandler) *[]File {
	movieFiles := make([]File, 0)
	for _, f := range *fh.Files {
		if IsMovie(&f) {
			movieFiles = append(movieFiles, f)
		}
	}

	return &movieFiles
}

func MovieFileNames(fh *FileHandler) *[][]byte {
	mFiles := MovieFiles(fh)
	names := make([][]byte, len(*mFiles))
	for i, f := range *mFiles {
		names[i] = []byte(*f.Name())
	}

	return &names
}
