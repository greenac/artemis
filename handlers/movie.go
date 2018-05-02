package handlers

import (
	"github.com/greenac/artemis/tools"
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/logger"
	"fmt"
)

type MovieHandler struct {
	DirPaths *[]tools.FilePath
	Movies *[]tools.FilePath
}

func (mh *MovieHandler)GetMovies() error {
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
		for _, f := range *names {
			fmt.Println(string(f))
		}

		fNames = append(fNames, *names...)
	}

	return nil
}
