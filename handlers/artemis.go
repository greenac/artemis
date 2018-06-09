package handlers

import "github.com/greenac/artemis/tools"
import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/movie"
)

type ArtemisHandler struct {
	MovieHandler  *MovieHandler
	ActorHandler  *ActorHandler
	UnknownMovies []movie.Movie
}

func (ah *ArtemisHandler) Setup(
	movieDirPaths *[]tools.FilePath,
	actorDirPaths *[]tools.FilePath,
	actorFilePath *tools.FilePath,
	cachedNamePath *tools.FilePath,
	toPath *tools.FilePath,
) {
	if ah.ActorHandler == nil {
		actHand := ActorHandler{
			DirPaths: actorDirPaths,
			NamesPath: actorFilePath,
			CachedPath: cachedNamePath,
			ToPath: toPath,
		}
		err := actHand.FillActors()
		if err != nil {
			logger.Error("`ArtemisHandler::Setup` getting actors", err)
		}

		ah.ActorHandler = &actHand
	}

	if ah.MovieHandler == nil {
		mh := MovieHandler{DirPaths: movieDirPaths}
		mh.SetMovies()

		ah.MovieHandler = &mh
	}

	ah.UnknownMovies = make([]movie.Movie, 0)
}

func (ah *ArtemisHandler) Sort() {
	logger.Log("Sorting:", len(*ah.MovieHandler.Movies), "movies")
	for _, m := range *ah.MovieHandler.Movies {
		found := false
		for _, a := range ah.ActorHandler.Actors {
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
	return &ah.ActorHandler.Actors
}

func (ah *ArtemisHandler) AddMovie(names string, movie *movie.Movie) {
	//parts := strings.Split(names, ",")
	//
	//ah.actorHandler.AddMovie(name, movie)
}

func (ah *ArtemisHandler) RenameMovies() {
	for _, a := range ah.ActorHandler.Actors {
		if len(a.Movies) == 0 {
			continue
		}

		mvs := make([]*movie.Movie, 0)
		for _, m := range a.Movies {
			mvs = append(mvs, m)
		}

		ah.MovieHandler.RenameMovies(mvs)
	}
}