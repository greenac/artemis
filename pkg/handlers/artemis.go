package handlers

import (
	"github.com/greenac/artemis/pkg/dbinteractors"
	"github.com/greenac/artemis/pkg/logger"
	"github.com/greenac/artemis/pkg/models"
)

type ArtemisHandler struct {
	MovieHandler *MovieHandler
	ActorHandler *ActorHandler
}

func (ah *ArtemisHandler) Setup(movieDirPaths *[]models.FilePath) error {
	if ah.ActorHandler == nil {
		actHand := ActorHandler{}
		err := actHand.FillActors()
		if err != nil {
			logger.Error("`ArtemisHandler::Setup` getting actors", err)
		}

		ah.ActorHandler = &actHand
	}

	if ah.MovieHandler == nil {
		mh := MovieHandler{DirPaths: movieDirPaths}

		err := mh.SetMovies()
		if err != nil {
			return err
		}

		ah.MovieHandler = &mh
	}

	return nil
}

func (ah *ArtemisHandler) Save(saveIfExts bool) {
	for _, m := range ah.MovieHandler.Movies {
		m.GetNewName()

		mm := dbinteractors.NewMovie(m.Name(), m.Path())
		exists, err := dbinteractors.DoesMovieExist(mm.GetIdentifier())
		if err != nil {
			exists = false
		}

		if !exists {
			_, err = mm.Create()
			if err != nil {
				continue
			}
		} else if !saveIfExts {
			continue
		}

		movieName := m.Name()
		for _, a := range ah.ActorHandler.Actors {
			if a.IsIn(movieName) {
				logger.Log("Movie:", movieName, "has actor:", a.FullName())

				mm.AddActor(a.Id)
				a.AddMovie(mm.Id)
				err = a.Save()
				if err == nil {
					logger.Log("ArtemisHandler::Save::Saved movie", mm.Name, "for actor:", a.FullName())
				} else {
					logger.Error("ArtemisHandler::Save::Failed to save movie", mm.Name, "for actor:", a.FullName())
					continue
				}

				logger.Log("ArtemisHandler::Save::Saving Movie", mm.Name)
				err = mm.Save()
				if err != nil {
					logger.Error("ArtemisHandler::Save::Failed to save movie", mm.Name, "with actors:", mm.ActorIds)
				}
			}
		}
	}
}
