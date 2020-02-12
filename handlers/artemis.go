package handlers

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"os"
	"path"
	"strings"
)

type ArtemisHandler struct {
	MovieHandler  *MovieHandler
	ActorHandler  *ActorHandler
	ToPath *models.FilePath
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
	for _, m := range *ah.MovieHandler.Movies {
		m.GetNewName()
		found := false

		for _, a := range ah.ActorHandler.Actors {
			if a.IsIn(&m) {
				err := a.AddMovie(m)
				if err != nil {
					logger.Warn("`ArtemisHandler::Sort` could not add movie:", m, "for actor:", a.FullName())
					continue
				}

				found = true
			}
		}

		if !found {
			ah.MovieHandler.AddUnknownMovie(m)
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

func (ah *ArtemisHandler) MoveMovies() {
	for _, m := range ah.MovieHandler.UnknownMovies {
		if len(m.Actors) == 0 {
			continue
		}

		a := m.Actors[0]
		m.Actors = nil
		err := a.AddMovie(*m)
		if err != nil {
			logger.Warn("ArtemisHandler::MoveMovies could not add unknown movie with error:", err)
		}
	}

	for _, a := range ah.ActorHandler.Actors {
		if len(a.Movies) == 0 {
			continue
		}

		ap := path.Join(ah.ToPath.PathAsString(), a.FullName())
		logger.Debug("moving to actor path:", ap)
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

		for _, m := range a.Movies {
			m.NewPath = path.Join(ap, m.GetNewName())
			logger.Log("actor path:", ap, m.GetNewName(), m.NewPath)
			err = os.Rename(m.Path, m.NewPath)
			if err != nil {
				logger.Error("`ArtemisHandler::MoveMovies` could not rename:", m.Path, "to:", m.NewPath, err)
				panic(err)
			}
		}
	}
}
