package handlers

import (
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/movie"
	"github.com/greenac/artemis/tools"
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
				mvs = append(mvs, movie.Movie{File: f})
			}
		}
	}

	mh.Movies = &mvs
	return nil
}
