package handlers

import (
	"github.com/greenac/artemis/tools"
	"github.com/greenac/artemis/movie"
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/logger"
)

type MovieHandler struct {
	DirPaths *[]tools.FilePath
	FilePaths *[]tools.FilePath
}

func (mh *MovieHandler)getMovies() error {
	if mh.DirPaths == nil {
		logger.Error("Cannot fill movies from dirs. DirPaths not initialized")
		return artemiserror.GetArtemisError(artemiserror.ArgsNotInitialized, nil)
	}

	fNames := make([][]byte, 0)
	for _, p := range *mh.DirPaths {
		logger.Log("Files for base path:", p.PathAsString())
		fh := tools.FileHandler{BasePath: p}
		err := fh.SetFiles()
		if err != nil {
			logger.Warn("Could not fill actors from path:", p.PathAsString())
			continue
		}

		names := fh.DirFileNames()
		fNames = append(fNames, *names...)
	}

	actors := make(map[string]movie.Actor, len(fNames))
	for _, n := range fNames {
		a, err := mh.createActor(&n)
		if err != nil {
			continue
		}

		actors[a.FullName()] = *a
	}

	if mh.Actors == nil {
		mh.Actors = &actors
	} else {
		for name, actor := range actors {
			(*mh.Actors)[name] = actor
		}
	}

}
