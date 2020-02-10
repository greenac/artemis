package handlers

import (
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"github.com/greenac/artemis/utils"
	"path"
)

type MovieHandler struct {
	DirPaths      *[]models.FilePath
	Movies        *[]models.Movie
	NewToPath     *models.FilePath
	UnknownMovies []*models.Movie
}

func (mh *MovieHandler) SetMovies() error {
	if mh.DirPaths == nil {
		logger.Error("Cannot fill movies from dirs. DirPaths not initialized")
		return artemiserror.New(artemiserror.ArgsNotInitialized)
	}

	mvs := make([]models.Movie, 0)
	for _, p := range *mh.DirPaths {
		fh := FileHandler{BasePath: p}
		err := fh.SetFiles()
		if err != nil {
			logger.Warn("Could not fill movies from path:", p.PathAsString())
			continue
		}

		for _, f := range *fh.Files {
			if utils.IsMovie(&f) {
				m := models.Movie{File: f}
				m.Path = path.Join(p.Path, *m.Name())
				mvs = append(mvs, m)
			}
		}
	}

	mh.Movies = &mvs
	return nil
}

func (mh *MovieHandler) RenameMovies(mvs []*models.Movie) {
	for _, m := range mvs {
		mh.RenameMovie(m)
	}
}

func (mh *MovieHandler) RenameMovie(m *models.Movie) error {
	if m.Path == "" {
		logger.Warn("`MovieHandler::RenameMovie` movie:", m.Name(), "does not have path set")
		return artemiserror.New(artemiserror.PathNotSet)
	}

	fh := FileHandler{}
	err := fh.Rename(m.Path, m.RenamePath())
	if err != nil {
		logger.Warn("`MovieHandler::RenameMovie` movie:", m.Name(), "failed to be renamed with error:", err)
		return err
	}

	return nil
}

func (mh *MovieHandler) AddUnknownMovie(m *models.Movie) {
	mh.UnknownMovies = append(mh.UnknownMovies, m)
}

func (mh *MovieHandler) AddUnknownMovieNames() {
	for _, m := range mh.UnknownMovies {
		m.AddActorNames()
	}
}

func (mh *MovieHandler) RenameUnknownMovies() {
	logger.Debug("MovieHandler::RenameUnknownMovies renaming:", len(mh.UnknownMovies))
	mh.RenameMovies(mh.UnknownMovies)
}
