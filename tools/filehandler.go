package tools

import (
	"io/ioutil"
	"github.com/greenac/artemis/logger"
)

type FileHandler struct {
	Files *[]File
	BasePath FilePath
}

func (fh *FileHandler)SetFiles() error {
	if !fh.BasePath.PathDefined() {
		panic("File Handler Base Path Not Set")
	}

	fi, err := ioutil.ReadDir(string(*fh.BasePath.Path))
	if err != nil {
		logger.Error("Failed to set file handler file names with error:", err)
		return err
	}

	files := make([]File, len(fi))
	for i, f := range fi {
		files[i] = File{Info: f}
	}

	fh.Files = &files
	return nil
}

func (fh *FileHandler)FileNames() *[][]byte {
	names := make([][]byte, len(*fh.Files))
	for i, f := range *fh.Files {
		names[i] = []byte(*f.Name())
	}

	return &names
}

func (fh *FileHandler)MovieFiles() *[]File {
	movieFiles := make([]File, 0)
	for _, f := range *fh.Files {
		if f.IsMovie() {
			movieFiles = append(movieFiles, f)
		}
	}

	return &movieFiles
}

func (fh *FileHandler)DirFiles() *[]File {
	dFiles := make([]File, 0)
	for _, f := range *fh.Files {
		if f.IsDir() {
			dFiles = append(dFiles, f)
		}
	}

	return &dFiles
}

func (fh *FileHandler)MovieFileNames() *[][]byte {
	mFiles := fh.MovieFiles()
	names := make([][]byte, len(*mFiles))
	for i, f := range *mFiles {
		names[i] = []byte(*f.Name())
	}

	return &names
}

func (fh *FileHandler)DirFileNames() *[][]byte {
	dFiles := fh.DirFiles()
	names := make([][]byte, len(*dFiles))
	for i, f := range *dFiles {
		names[i] = []byte(*f.Name())
	}

	return &names
}