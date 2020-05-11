package handlers

import (
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/dbinteractors"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
)

type MovieHandler struct {
	DirPaths      *[]models.FilePath
	Movies        []models.SysMovie
	KnownMovies   []*models.Movie
	UnknownMovies []*models.Movie
	unkIndex      int
}

func (mh *MovieHandler) SetMovies() error {
	if mh.DirPaths == nil {
		logger.Error("MovieHandler::SetMovies Cannot fill movies from dirs. DirPaths not initialized")
		return artemiserror.New(artemiserror.ArgsNotInitialized)
	}

	mvs := make([]models.SysMovie, 0)
	for _, p := range *mh.DirPaths {
		fh := FileHandler{BasePath: p}
		err := fh.SetFiles()
		if err != nil {
			logger.Warn("MovieHandler::SetMovies Could not fill movies from path:", p.PathAsString())
			continue
		}

		for _, f := range fh.Files {
			if f.IsMovie() {
				m := models.SysMovie{File: f}
				mvs = append(mvs, m)
			}
		}
	}

	mh.Movies = mvs

	return nil
}

func (mh *MovieHandler) RenameMovies(mvs []*models.SysMovie) {
	for _, m := range mvs {
		_ = mh.RenameMovie(m)
	}
}

func (mh *MovieHandler) RenameMovie(m *models.SysMovie) error {
	if m.BasePath == "" {
		logger.Warn("`MovieHandler::RenameMovie` movie:", m.Name(), "does not have path set")
		return artemiserror.New(artemiserror.PathNotSet)
	}

	m.NewBasePath = m.BasePath

	err := MoveMovie(m, Internal)
	if err != nil {
		logger.Warn("`MovieHandler::RenameMovie` movie:", m.Path, "failed to be renamed with error:", err)
		return err
	}

	return nil
}

func (mh *MovieHandler) AddKnownMovie(m models.SysMovie) {
	dbm, err := dbinteractors.GetMovieByIdentifier(m.Path())
	if err != nil {
		logger.Warn("MovieHandler::AddKnownMovie could not add movie:", m.Name(), "Failed with error:", err)
		return
	}

	if dbm == nil {
		nm := dbinteractors.NewMovie(m.Name(), m.Path())
		_ = nm.Save()
		dbm = &nm
	}

	mh.KnownMovies = append(mh.KnownMovies, dbm)
}
