package handlers

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
)


func MoveDir(dir models.File) error {
	if !dir.IsDir() {
		return nil
	}

	fh := FileHandler{BasePath: models.FilePath{Path: dir.Path()}}

	err := fh.SetFiles()
	if err != nil {
		logger.Error("DirectoryMover::MoveDir failed to read files in dir:", dir.Path, err)
		return err
	}

	ex, err := fh.DoesFileExistAtPath(dir.NewPath())
	if err != nil {
		logger.Warn("DirectoryMover::MoveDir failed to move directory to:", dir.GetNewTotalPath(), err)
		return err
	}

	if ex {
		for _, f := range fh.Files {
			f.NewBasePath = dir.NewPath()
			if f.IsMovie() {
				err := fh.Rename(f.Path(), f.GetNewTotalPath())
				if err != nil {
					continue
				}
			}
		}
	} else {
		fh.Rename(dir.Path(), dir.GetNewTotalPath())
	}
}

