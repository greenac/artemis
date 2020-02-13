package handlers

import (
	"github.com/greenac/artemis/models"
	"github.com/greenac/artemis/utils"
)

type MovieFileHandler struct {
	FileHandler
}

func (mfh *MovieFileHandler) MovieFiles() *[]models.File {
	movieFiles := make([]models.File, 0)
	for _, f := range *mfh.Files {
		if utils.IsMovie(&f) {
			movieFiles = append(movieFiles, f)
		}
	}

	return &movieFiles
}

func (mfh *MovieFileHandler) MovieFileNames() *[][]byte {
	mFiles := mfh.MovieFiles()
	names := make([][]byte, len(*mFiles))
	for i, f := range *mFiles {
		names[i] = []byte(*f.Name())
	}

	return &names
}
