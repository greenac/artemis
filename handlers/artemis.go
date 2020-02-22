package handlers

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
)

type ArtemisHandler struct {
	MovieHandler *MovieHandler
	ActorHandler *ActorHandler
	ToPath       *models.FilePath
}

func (ah *ArtemisHandler) Setup(
	movieDirPaths *[]models.FilePath,
	actorDirPaths *[]models.FilePath,
	actorFilePath *models.FilePath,
	cachedNamePath *models.FilePath,
	toPath *models.FilePath,
) error {
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
		mh := MovieHandler{DirPaths: movieDirPaths, NewToPath: ah.ToPath}
		err := mh.CleanseInitialNames()
		if err != nil {
			return err
		}

		err = mh.SetMovies()
		if err != nil {
			return err
		}

		ah.MovieHandler = &mh
	}

	ah.ToPath = toPath

	return nil
}

func (ah *ArtemisHandler) Sort() {
	for _, m := range ah.MovieHandler.Movies {
		m.GetNewName()

		for _, a := range ah.ActorHandler.Actors {
			if a.IsIn(&m) {
				m.AddActor(*a)
				m.UpdateNewName(a)
			}
		}

		if m.IsKnown() {
			ah.MovieHandler.AddKnownMovie(m)
		} else {
			ah.MovieHandler.AddUnknownMovie(m)
		}
	}
}

func (ah *ArtemisHandler) RenameMovies() {
	ah.MovieHandler.AddKnownMovieNames()
	ah.MovieHandler.AddUnknownMovieNames()
}

func (ah *ArtemisHandler) MoveMovies() {
	ah.MovieHandler.MoveMovies(ah.ToPath.PathAsString())
}
