package handlers

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"path"
)

type DirectoryMover struct {
	ToDirPath models.FilePath
}

func (dm *DirectoryMover) MoveDir(dir models.File) error {
	if !dir.IsDir() {
		return nil
	}

	fh := FileHandler{BasePath: models.FilePath{Path: dir.Path()}}

	err := fh.SetFiles()
	if err != nil {
		logger.Error("DirectoryMover::MoveDir failed to read files in dir:", dir.Path, err)
		return err
	}

	ex, err := fh.DoesFileExistAtPath(dir.NewPath)
	if err != nil {
		logger.Warn("DirectoryMover::MoveDir failed to move directory to:", dir.GetNewTotalPath(), err)
		return err
	}

	if ex {
		for _, f := range fh.Files {
			f.NewPath = dir.NewPath
			if f.IsMovie() {
				fh.Rename(f.Path, )
				err := fm.Move()
			}
		}
	} else {

	}
}

