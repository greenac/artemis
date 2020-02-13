package handlers

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"os"
	"path"
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
		mh := MovieHandler{DirPaths: movieDirPaths, NewToPath: ah.ToPath}
		err := mh.SetMovies()
		if err != nil {
			logger.Error("`ArtemisHandler::Setup` could not set movies.", err)
			panic(err)
		}

		ah.MovieHandler = &mh
	}

	ah.ToPath = toPath
}

func (ah *ArtemisHandler) Sort() {
	for _, m := range ah.MovieHandler.Movies {
		m.GetNewName()
		isKnown := false

		for _, a := range ah.ActorHandler.Actors {
			if a.IsIn(&m) {
				m.AddActor(*a)
				m.UpdateNewName(a)
				isKnown = true
			}
		}

		if isKnown {
			ah.MovieHandler.AddKnownMovie(m)
		} else {
			ah.MovieHandler.AddUnknownMovie(m)
		}
	}
}

func (ah *ArtemisHandler) Actors() *map[string]*models.Actor {
	return &ah.ActorHandler.Actors
}

func (ah *ArtemisHandler) RenameMovies() {
	ah.MovieHandler.AddKnownMovieNames()
	ah.MovieHandler.AddUnknownMovieNames()
	ah.MovieHandler.RenameAllMovies()
}

func (ah *ArtemisHandler) MoveMovies() {
	mvs := make([]*models.Movie, 0)

	for _, m := range ah.MovieHandler.KnownMovies {
		if len(m.Actors) > 0 {
			mvs = append(mvs, m)
		}
	}

	for _, m := range ah.MovieHandler.UnknownMovies {
		if m.NewName != "" && *m.Name() != m.NewName && len(m.Actors) > 0{
			mvs = append(mvs, m)
		}
	}

	for _, m := range mvs {
		a := m.Actors[0]
		ap := path.Join(ah.ToPath.PathAsString(), a.FullName())

		fi, err := os.Stat(ap)
		if err != nil && os.IsNotExist(err) {
			err = os.Mkdir(ap, 0775)
			if err != nil {
				logger.Error("`ArtemisHandler::MoveMovies` could not make directory:", ap)
				panic(err)
			}
		} else if err != nil {
			logger.Error("`ArtemisHandler::MoveMovies` error checking file:", err)
			continue
		} else if !fi.IsDir() {
			logger.Error("`ArtemisHandler::MoveMovies` File at path:", ap, "is not a directory")
			continue
		}

		m.NewPath = path.Join(ap, m.GetNewName())
		err = os.Rename(m.Path, m.NewPath)

		if err != nil {
			logger.Error("`ArtemisHandler::MoveMovies` could not rename:", m.Path, "to:", m.NewPath, err)
			panic(err)
		}
	}
}
