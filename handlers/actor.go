package handlers

import (
	"github.com/greenac/artemis/tools"
	"github.com/greenac/artemis/movie"
	"strings"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/artemiserror"
	"errors"
)

type ActorHandler struct {
	Paths *[]tools.FilePath
	Actors *[]movie.Actor
}

func (ah *ActorHandler)FillActors() error {
	if ah.Paths == nil || len(*(ah.Paths)) == 0 {
		return artemiserror.GetArtemisError(artemiserror.ArgsNotInitialized, nil)
	}

	fNames := make([][]byte, 0)

	for _, p := range *ah.Paths {
		fh := tools.FileHandler{BasePath: p}
		err := fh.SetFiles()
		if err != nil {
			logger.Warn("Could not fill actors from path:", p.PathAsString())
			continue
		}

		ns := fh.DirFileNames()
		fNames = append(fNames, *ns...)
	}


	logger.Warn("got:", len(fNames), "names")
	actors := make([]movie.Actor, len(fNames))
	for i, n := range fNames {
		a, err := ah.createActor(&n)
		if err != nil {
			continue
		}

		actors[i] = *a

		logger.Log(a, i, a.FullName())
	}

	ah.Actors = &actors
	return nil
}

func (ah *ActorHandler)createActor(name *[]byte) (*movie.Actor, error) {
	n := string(*name)
	parts := strings.Split(n, " ")
	if len(parts) == 1 {
		parts = strings.Split(n, "_")
	}

	var a movie.Actor
	switch len(parts) {
	case 1:
		a = movie.Actor{FirstName: &parts[0]}
	case 2:
		a = movie.Actor{FirstName: &parts[0], LastName:&parts[1]}
	case 3:
		a = movie.Actor{FirstName: &parts[0], MiddleName:&parts[1], LastName:&parts[2]}
	default:
		logger.Error("Cannot parse actor name:", name)
		return nil, errors.New("ActorNameInvalid")
	}

	return &a, nil
}
