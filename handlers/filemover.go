package handlers

import (
	"github.com/greenac/artemis/logger"
	"os"
)

type FileMover struct {
	FromPath FilePath
	ToPath   FilePath
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
