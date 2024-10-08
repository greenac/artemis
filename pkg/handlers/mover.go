package handlers

import (
	"github.com/greenac/artemis/pkg/logger"
	"github.com/greenac/artemis/pkg/models"
	"os"
)

type FileMover struct {
	FromPath models.FilePath
	ToPath   models.FilePath
}

func (fm *FileMover) checkPaths() bool {
	return fm.FromPath.PathDefined() && fm.ToPath.PathDefined()
}

func (fm *FileMover) Move() error {
	if !fm.checkPaths() {
		// TODO: throw correct error here
		panic("File mover paths not instantiated")
	}

	err := os.Rename(fm.FromPath.PathAsString(), fm.ToPath.PathAsString())
	if err != nil {
		if os.IsNotExist(err) {
			logger.Error("Failed to move file. File at path:", fm.FromPath, "does not exist")
		}

		return err
	}

	return nil
}
