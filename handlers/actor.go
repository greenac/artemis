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
	DirPaths *[]tools.FilePath
	FilePaths *[]tools.FilePath
	Actors *map[string]movie.Actor
}

func (ah *ActorHandler)FillActors() error {
	if err := ah.FillActorsFromDirs(); err != nil {
		return err
	}

	if err := ah.FillActorsFromFiles(); err != nil {
		return err
	}

	return nil
}

func (ah *ActorHandler)FillActorsFromFiles() error {
	if ah.FilePaths == nil {
		logger.Error("Cannot fill actors from file. FilePaths not initialized")
		return artemiserror.GetArtemisError(artemiserror.ArgsNotInitialized, nil)
	}

	fNames := make([][]byte, 0)
	for _, p := range *ah.FilePaths {
		fh := tools.FileHandler{BasePath: p}
		names, err := fh.ReadNameFile(&p)
		if err != nil {
			continue
		}

		fNames = append(fNames, *names...)
	}

	actors := make(map[string]movie.Actor, len(fNames))
	for _, n := range fNames {
		a, err := ah.createActor(&n)
		if err != nil {
			continue
		}

		actors[a.FullName()] = *a
	}

	if ah.Actors == nil {
		ah.Actors = &actors
	} else {
		for name, actor := range actors {
			a, hasName := (*ah.Actors)[name]
			if hasName {
				a.AddFiles(actor.Files)
			} else {
				(*ah.Actors)[name] = actor
			}
		}
	}

	return nil
}

func (ah *ActorHandler)FillActorsFromDirs() error {
	if ah.DirPaths == nil {
		logger.Error("Cannot fill actors from dirs. DirPaths not initialized")
		return artemiserror.GetArtemisError(artemiserror.ArgsNotInitialized, nil)
	}

	fNames := make([][]byte, 0)
	for _, p := range *ah.DirPaths {
		logger.Log("Files for base path:", p.PathAsString())
		fh := tools.FileHandler{BasePath: p}
		err := fh.SetFiles()
		if err != nil {
			logger.Warn("Could not fill actors from path:", p.PathAsString())
			continue
		}

		names := fh.DirFileNames()
		fNames = append(fNames, *names...)
	}

	actors := make(map[string]movie.Actor, len(fNames))
	for _, n := range fNames {
		a, err := ah.createActor(&n)
		if err != nil {
			continue
		}

		actors[a.FullName()] = *a
	}

	if ah.Actors == nil {
		ah.Actors = &actors
	} else {
		for name, actor := range actors {
			(*ah.Actors)[name] = actor
		}
	}

	return nil
}

func (ah *ActorHandler)createActor(name *[]byte) (*movie.Actor, error) {
	if name == nil || len(*name) == 0 {
		logger.Error("Cannot create actor from name:", name)
		return nil, artemiserror.GetArtemisError(artemiserror.ArgsNotInitialized, nil)
	}

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

func (ah *ActorHandler)PrintActors() {
	i := 1
	for _, actor := range *ah.Actors {
		logger.Log(i, actor.FullName())
		i += 1
	}
}