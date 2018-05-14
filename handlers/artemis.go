package handlers

import "github.com/greenac/artemis/tools"
import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/movie"
)

type ArtemisHandler struct {
	movieHandler  *MovieHandler
	actorHandler  *ActorHandler
	UnknownMovies []movie.Movie
}

func (ah *ArtemisHandler) Setup(movieDirPaths *[]tools.FilePath, actorDirPaths *[]tools.FilePath, actorFilePaths *[]tools.FilePath) {
	if ah.actorHandler == nil {
		actHand := ActorHandler{DirPaths: actorDirPaths, FilePaths: actorFilePaths}
		err := actHand.FillActors()
		if err != nil {
			logger.Error("`ArtemisHandler::Setup` getting actors", err)
		}

		ah.actorHandler = &actHand
	}

	if ah.movieHandler == nil {
		mh := MovieHandler{DirPaths: movieDirPaths}
		mh.SetMovies()

		ah.movieHandler = &mh
	}

	ah.UnknownMovies = make([]movie.Movie, 0)
}

func (ah *ArtemisHandler) Sort() {
	logger.Log("Sorting:", len(*ah.movieHandler.Movies), "movies")
	for _, m := range *ah.movieHandler.Movies {
		found := false
		for _, a := range ah.actorHandler.Actors {
			if a.IsIn(&m) {
				a.AddMovie(&m)
				found = true
			}
		}

		if !found {
			ah.UnknownMovies = append(ah.UnknownMovies, m)
		}
	}
}

func (ah *ArtemisHandler) Actors() *map[string]*movie.Actor {
	return &ah.actorHandler.Actors
}
