package handlers

import (
	"errors"
	"github.com/greenac/artemis/pkg/artemiserror"
	"github.com/greenac/artemis/pkg/dbinteractors"
	"github.com/greenac/artemis/pkg/logger"
	"github.com/greenac/artemis/pkg/models"
	"sort"
	"strings"
)

type actorMap map[string]models.Actor

type ActorHandler struct {
	Actors    actorMap
	NamesPath *models.FilePath
}

func (ah *ActorHandler) SaveActorsToDb(acts *[]models.Actor) {
	for _, a := range *acts {
		err := a.Save()
		if err != nil {
			logger.Warn("ActorHandler::SaveActorsToDb could not save actor:", a.GetIdentifier(), err)
		}
	}
}

func (ah *ActorHandler) FillActors() error {
	ah.Actors = make(actorMap)

	err := ah.FillActorsFromDb()
	if err != nil {
		return err
	}

	//if err := ah.FillActorsFromDirs(); err != nil { return err }
	//
	//if err := ah.FillActorsFromFile(); err != nil { return err }
	//
	//if err := ah.fillActorsFromCachedFile(); err != nil { return err }

	return nil
}

func (ah *ActorHandler) FillActorsFromDb() error {
	acts, err := dbinteractors.AllActors()
	if err != nil {
		return err
	}

	for _, a := range *acts {
		ah.Actors[a.FullName()] = a
	}

	return nil
}

func (ah *ActorHandler) FillActorsFromFile() error {
	if ah.NamesPath == nil {
		logger.Error("ActorHandler::FillActorsFromFile Cannot fill actors from file. FilePath not initialized")
		return artemiserror.New(artemiserror.ArgsNotInitialized)
	}

	fh := FileHandler{BasePath: *ah.NamesPath}
	names, err := fh.ReadNameFile(ah.NamesPath)
	if err != nil {
		logger.Error("Cannot read name file at path:", ah.NamesPath.PathAsString(), "error:", err)
		return err
	}

	for _, n := range *names {
		a, err := ah.CreateActor(string(n))
		if err != nil {
			logger.Warn("`ActorHandler::FillActorsFromFiles` could not create actor from name:", string(n), err)
			continue
		}

		ah.Actors[a.FullName()] = a
	}

	return nil
}

func (ah *ActorHandler) SortedActors() *[]models.Actor {
	acts := make([]models.Actor, len(ah.Actors))
	i := 0
	for _, a := range ah.Actors {
		acts[i] = a
		i += 1
	}

	sort.Slice(acts, func(i int, j int) bool {
		return acts[i].FullName() < acts[j].FullName()
	})

	return &acts
}

func (ah *ActorHandler) CreateActor(name string) (models.Actor, error) {
	return CreateNewActor(name)
}

func (ah *ActorHandler) Matches(name string) []*models.Actor {
	actors := make([]*models.Actor, 0)
	for _, a := range ah.Actors {
		if a.IsMatch(name) {
			actors = append(actors, &a)
		}
	}

	sort.Slice(actors, func(i int, j int) bool {
		return actors[i].FullName() < actors[j].FullName()
	})

	return actors
}

func (ah *ActorHandler) NameMatches(name string) (actors []*models.Actor, common string) {
	acts := make([]*models.Actor, 0)
	n := strings.ToLower(strings.Replace(name, " ", "_", -1))
	for _, a := range ah.Actors {
		if a.MatchWhole(n) {
			acts = append(acts, &a)
		}
	}

	sort.Slice(acts, func(i int, j int) bool {
		return acts[i].FullName() < acts[j].FullName()
	})

	comp := make([]rune, 0)
	if len(acts) == 1 {
		comp = []rune(acts[0].FullName())
	} else if len(acts) > 1 {
		actor := acts[0]
		actName := actor.FullName()
		for i, c := range actName {
			add := true
			for j := 1; j < len(acts); j += 1 {
				a := acts[j]
				aName := a.FullName()
				if i >= len(aName) || byte(c) != aName[i] {
					add = false
					break
				}
			}

			if !add {
				break
			}

			comp = append(comp, c)
		}
	}

	return acts, string(comp)
}

func (ah *ActorHandler) PrintActors() {
	i := 1
	for _, actor := range ah.Actors {
		logger.Log(i, actor.FullName())
		i += 1
	}
}

func (ah *ActorHandler) ActorForName(name string) (actor *models.Actor, error error) {
	a, has := ah.Actors[name]
	if has {
		return &a, nil
	}

	return nil, artemiserror.New(artemiserror.InvalidName)
}

func (ah *ActorHandler) AddNewActor(name string) (*models.Actor, error) {
	a, has := ah.Actors[name]
	if has {
		return &a, nil
	}

	act, err := ah.CreateActor(name)
	if err != nil {
		return nil, err
	}

	ah.Actors[name] = act

	_, _, err = models.CreateIfDoesNotExist(&act)
	if err != nil {
		return nil, err
	}

	return &act, nil
}

func (ah *ActorHandler) AddActorsToMovieWithInput(input string, movie *models.SysMovie) {
	input = strings.ToLower(input)
	nms := strings.Split(input, ",")
	for _, n := range nms {
		nn := strings.ReplaceAll(strings.TrimSpace(n), " ", "_")
		a, err := ah.ActorForName(nn)
		if err != nil {
			a, err = ah.AddNewActor(nn)
		}

		if err == nil {
			movie.AddActor(*a)
			movie.UpdateNewName(a)
		}
	}
}

func CreateNewActor(name string) (models.Actor, error) {
	if len(name) == 0 {
		logger.Error("CreateNewActor invalid name:", name)
		return models.Actor{}, artemiserror.New(artemiserror.ArgsNotInitialized)
	}

	n := strings.TrimSpace(name)
	parts := strings.Split(n, " ")
	if len(parts) == 1 {
		parts = strings.Split(n, "_")
	}

	var a models.Actor
	switch len(parts) {
	case 1:
		a = dbinteractors.NewActor(parts[0], "", "")
	case 2:
		a = dbinteractors.NewActor(parts[0], "", parts[1])
	case 3:
		a = dbinteractors.NewActor(parts[0], parts[1], parts[2])
	default:
		logger.Error("Cannot parse actor name:", name)
		return models.Actor{}, errors.New("ActorNameInvalid")
	}

	_, err := a.Create()
	if err != nil {
		logger.Error("ActorHandler::createActor cannot create actor:", name, err)
		return a, err
	}

	return a, nil
}
