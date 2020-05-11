package handlers

import (
	"github.com/greenac/artemis/dbinteractors"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
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

func (ah *ArtemisHandler) Sort() {
	for _, m := range ah.MovieHandler.Movies {
		m.GetNewName()

		mm := dbinteractors.NewMovie(m.Name(), m.Path())
		ex, err := dbinteractors.DoesMovieExist(mm.GetIdentifier())
		if err != nil {
			ex = false
		}

		if !ex {
			_, err = mm.Create()
			if err != nil {
				continue
			}
		}

		save := false
		for _, a := range ah.ActorHandler.Actors {
			if a.IsIn(m.Name()) {
				logger.Log("Movie:", m.Name(), "has actor:", a.FullName())
				save = true
				mm.AddActor(a.Id)
				a.AddMovie(mm.Id)
				_ = a.Save()
			}
		}

		if save {
			logger.Debug("Would save movie", mm.Name, mm.GetIdentifier())
			_ = mm.Save()
		}
	}
}
