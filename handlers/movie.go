package handlers

import (
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/movie"
	"github.com/greenac/artemis/tools"
	"path"
)

type MovieHandler struct {
	DirPaths *[]tools.FilePath
	Movies   *[]movie.Movie
}

func (mh *MovieHandler) SetMovies() error {
	if mh.DirPaths == nil {
		logger.Error("Cannot fill movies from dirs. DirPaths not initialized")
		return artemiserror.GetArtemisError(artemiserror.ArgsNotInitialized, nil)
	}

	mvs := make([]movie.Movie, 0)
	for _, p := range *mh.DirPaths {
		logger.Log("Movies for base path:", p.PathAsString())
		fh := tools.FileHandler{BasePath: p}
		err := fh.SetFiles()
		if err != nil {
			logger.Warn("Could not fill movies from path:", p.PathAsString())
			continue
		}

		for _, f := range *fh.Files {
			if movie.IsMovie(&f) {
				m := movie.Movie{File: f}
				m.Path = path.Join(p.Path, *m.Name())
				mvs = append(mvs, m)
			}
		}
	}

	mh.Movies = &mvs
	return nil
}

func (mh *MovieHandler) RenameMovies(mvs []*movie.Movie) {
	for _, m := range mvs {
		err := mh.RenameMovie(m)
		if err != nil {
			continue
		}
	}
}

func (mh *MovieHandler) RenameMovie(m *movie.Movie) error {
	if m.Path == "" {
		logger.Warn("`MovieHandler::RenameMovie` movie:", m.Name(), "does not have path set")
		return artemiserror.New(artemiserror.PathNotSet)
	}

	fh := tools.FileHandler{}
	err := fh.Rename(m.Path, m.GetNewName())
	if err != nil {
		logger.Warn("`MovieHandler::RenameMovie` movie:", m.Name(), "failed to be renamed with error:", err)
	}

	return nil
}
