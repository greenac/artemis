package bin

import (
	"github.com/greenac/artemis/dbinteractors"
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"os"
	"path"
)

func RemoveSymLinks(basePath string) error {
	fh := handlers.FileHandler{BasePath: models.FilePath{Path: basePath}}
	err := fh.SetFiles()
	if err != nil { return err }

	for _, df := range fh.Files {
		if !df.IsDir() {
			continue
		}

		mfh := handlers.FileHandler{BasePath: models.FilePath{Path: path.Join(df.Path())}}
		err := mfh.SetFiles()
		if err != nil { return err }

		for _, f := range mfh.Files {
			if !f.IsMovie() {
				continue
			}

			info, err := os.Lstat(f.Path())
			if err != nil || info == nil {
				logger.Error("Getting file info for:", f.Path(), err)
				continue
			}

			if !(info.Mode() & os.ModeSymlink == os.ModeSymlink) {
				continue
			}

			mid := models.MovieIdentifier(f.Path())
			m, err := dbinteractors.GetMovieByIdentifier(mid)
			if err != nil { continue }

			for _, aid := range m.ActorIds {
				act, err := dbinteractors.GetActorById(aid)
				if err != nil { continue }
				up := act.RemoveMovie(m.Id)
				if up {
					logger.Log("Removing movie:", m.Name, "for actor:", act.FullNameNoUnderscores())
					_ = act.Save()
				}
			}

			_ = dbinteractors.DeleteMovie(m.Id)
			_ = os.Remove(f.Path())

			logger.Log("Deleting & removing sym linked movie:", m.Name)
		}
	}

	return nil
}
