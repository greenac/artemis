package handlers

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"strings"
)

type ArtemisHandler struct {
	MovieHandler  *MovieHandler
	ActorHandler  *ActorHandler
	UnknownMovies []models.Movie
}

func (ah *ArtemisHandler) Setup(
	movieDirPaths *[]models.FilePath,
	actorDirPaths *[]models.FilePath,
	actorFilePath *models.FilePath,
	cachedNamePath *models.FilePath,
	toPath *models.FilePath,
) {
	if ah.ActorHandler == nil {
		actHand := ActorHandler{
			DirPaths:   actorDirPaths,
			NamesPath:  actorFilePath,
			CachedPath: cachedNamePath,
			ToPath:     toPath,
		}
		err := actHand.FillActors()
		if err != nil {
			logger.Error("`ArtemisHandler::Setup` getting actors", err)
		}

		ah.ActorHandler = &actHand
	}

	if ah.MovieHandler == nil {
		mh := MovieHandler{DirPaths: movieDirPaths, NewToPath: toPath}
		err := mh.SetMovies()
		if err != nil {
			logger.Error("`ArtemisHandler::Setup` could not set movies.", err)
			panic(err)
		}

		ah.MovieHandler = &mh
	}

	ah.UnknownMovies = make([]models.Movie, 0)
}

func (ah *ArtemisHandler) Sort() {
	for _, m := range *ah.MovieHandler.Movies {
		found := false
		for _, a := range ah.ActorHandler.Actors {
			if a.IsIn(&m) {
				err := a.AddMovie(&m)
				if err != nil {
					logger.Warn("`ArtemisHandler::Sort` could not add movie:", m, "for actor:", a.FullName())
					continue
				}
				found = true
			}
		}

		if !found {
			ah.UnknownMovies = append(ah.UnknownMovies, m)
		}
	}
}

func (ah *ArtemisHandler) Actors() *map[string]*models.Actor {
	return &ah.ActorHandler.Actors
}

func (ah *ArtemisHandler) AddMovie(names string, movie *models.Movie) {
	nms := strings.Split(names, ",")
	for _, n := range nms {
		ah.ActorHandler.AddMovie(n, movie)
	}
}

func (ah *ArtemisHandler) RenameMovies() {
	for _, a := range ah.ActorHandler.Actors {
		if len(a.Movies) == 0 {
			continue
		}

		mvs := make([]*models.Movie, 0)
		for _, m := range a.Movies {
			mvs = append(mvs, m)
		}

		ah.MovieHandler.RenameMovies(mvs)
	}

	ah.MovieHandler.AddUnknownMovieNames()
	ah.MovieHandler.RenameUnknownMovies()
}
