package handlers

import (
	"encoding/json"
	"errors"
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/movie"
	"github.com/greenac/artemis/tools"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
)

type ActorHandler struct {
	DirPaths   *[]tools.FilePath
	NamesPath  *tools.FilePath
	CachedPath *tools.FilePath
	Actors     map[string]*movie.Actor
	ToPath     *tools.FilePath
}

func (ah *ActorHandler) FillActors() error {
	ah.Actors = make(map[string]*movie.Actor)
	if err := ah.FillActorsFromDirs(); err != nil {
		return err
	}

	if err := ah.FillActorsFromFile(); err != nil {
		return err
	}

	if err := ah.fillActorsFromCachedFile(); err != nil {
		return err
	}

	return nil
}

func (ah *ActorHandler) FillActorsFromFile() error {
	if ah.NamesPath == nil {
		logger.Error("Cannot fill actors from file. FilePath not initialized")
		return artemiserror.New(artemiserror.ArgsNotInitialized)
	}

	fh := tools.FileHandler{BasePath: *ah.NamesPath}
	names, err := fh.ReadNameFile(ah.NamesPath)
	if err != nil {
		logger.Error("Cannot read name file at path:", ah.NamesPath.PathAsString(), "error:", err)
		return err
	}

	for _, n := range *names {
		a, err := ah.createActor(&n)
		if err != nil {
			logger.Error("`ActorHandler::FillActorsFromFiles` could not create actor from name:", string(n), err)
			continue
		}

		ah.Actors[a.FullName()] = &a
	}

	return nil
}

func (ah *ActorHandler) FillActorsFromDirs() error {
	if ah.DirPaths == nil {
		logger.Error("Cannot fill actors from dirs. DirPaths not initialized")
		return artemiserror.New(artemiserror.ArgsNotInitialized)
	}

	fNames := make([][]byte, 0)
	for _, p := range *ah.DirPaths {
		fh := tools.FileHandler{BasePath: p}
		err := fh.SetFiles()
		if err != nil {
			logger.Warn("Could not fill actors from path:", p.PathAsString())
			continue
		}

		names := fh.DirFileNames()
		fNames = append(fNames, *names...)
	}

	for _, n := range fNames {
		a, err := ah.createActor(&n)
		if err != nil {
			continue
		}

		ah.Actors[a.FullName()] = &a
	}

	return nil
}

func (ah *ActorHandler) fillActorsFromCachedFile() error {
	if ah.CachedPath == nil || !ah.CachedPath.PathDefined() {
		return nil
	}

	data, err := ioutil.ReadFile(ah.CachedPath.Path)
	if err != nil {
		logger.Error("Failed to read temp names file with error:", err)
		return err
	}

	names := make([]string, 0)
	err = json.Unmarshal(data, &names)
	if err != nil {
		logger.Error("`ActorHandler::fillActorsFromCachedFile` failed to unmarshal json file with error:", err)
		return err
	}

	for _, n := range names {
		nb := []byte(n)
		a, err := ah.createActor(&nb)
		if err != nil {
			continue
		}

		ah.Actors[a.FullName()] = &a
	}

	return nil
}

func (ah *ActorHandler) createActor(name *[]byte) (movie.Actor, error) {
	if name == nil || len(*name) == 0 {
		logger.Error("Cannot create actor from name:", name)
		return movie.Actor{}, artemiserror.New(artemiserror.ArgsNotInitialized)
	}

	n := strings.TrimSpace(string(*name))
	parts := strings.Split(n, " ")
	if len(parts) == 1 {
		parts = strings.Split(n, "_")
	}

	var a movie.Actor
	switch len(parts) {
	case 1:
		a = movie.Actor{FirstName: &parts[0]}
	case 2:
		a = movie.Actor{FirstName: &parts[0], LastName: &parts[1]}
	case 3:
		a = movie.Actor{FirstName: &parts[0], MiddleName: &parts[1], LastName: &parts[2]}
	default:
		logger.Error("Cannot parse actor name:", name)
		return movie.Actor{}, errors.New("ActorNameInvalid")
	}

	return a, nil
}

func (ah *ActorHandler) AddMovie(name string, m *movie.Movie) error {
	n := strings.Replace(strings.Trim(strings.ToLower(name), " "), " ", "_", -1)
	a, has := ah.Actors[n]
	if !has {
		logger.Warn("Cannot add movie:", *m.Name(), "to actor:", n, "no actor with that name")
		return errors.New("ActorNameInvalid")
	}

	return a.AddMovie(m)
}

func (ah *ActorHandler) Matches(name string) []*movie.Actor {
	actors := make([]*movie.Actor, 0)
	for _, a := range ah.Actors {
		if a.IsMatch(name) {
			actors = append(actors, a)
		}
	}

	sort.Slice(actors, func(i int, j int) bool {
		return actors[i].FullName() < actors[j].FullName()
	})

	return actors
}

func (ah *ActorHandler) NameMatches(name string) (actors []*movie.Actor, common string) {
	acts := make([]*movie.Actor, 0)
	n := strings.ToLower(strings.Replace(name, " ", "_", -1))
	for _, a := range ah.Actors {
		if a.MatchWhole(n) {
			acts = append(acts, a)
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

func (ah *ActorHandler) AddNameToMovies() {
	logger.Debug("`ActorHandler::AddNameToMovies`")
	for _, a := range ah.Actors {
		n := a.FullName()
		logger.Debug("`ActorHandler::AddNameToMovies` adding name:", n, "to actor:", a.FullName())
		for _, m := range a.Movies {
			m.AddName(n)
		}
	}
}

func (ah *ActorHandler) MoveMovies() {
	for _, a := range ah.Actors {
		ap := path.Join(ah.ToPath.PathAsString(), a.FullName())
		fi, err := os.Stat(ap)
		if err != nil && os.IsNotExist(err) {
			err = os.Mkdir(ap, 0644)
			if err != nil {
				logger.Error("`ActorHandler::MoveMovies` could not make directory:", ap)
				panic(err)
			}
		} else if err != nil {
			logger.Error("`ActorHandler::MoveMovies` error checking file:", err)
			continue
		} else if !fi.IsDir() {
			logger.Error("`ActorHandler::MoveMovies` File at path:", ap, "is not a directory")
			continue
		}

		for _, m := range a.Movies {
			if m.NewPath == "" {
				m.NewPath = path.Join(ap, m.GetNewName())
				err = os.Rename(m.Path, m.NewPath)
				if err != nil {
					logger.Error("`ActorHandler::MoveMovies` could not rename:", m.Path, "to:", m.NewPath)
					panic(err)
				}
			} else {
				np := path.Join(ap, m.GetNewName())
				err = os.Symlink(m.NewPath, np)
				if err != nil {
					logger.Error("`ActorHandler::MoveMovies` could not symlink:", m.Path, "to:", np)
					panic(err)
				}
			}
		}
	}
}

func (ah *ActorHandler) PrintActors() {
	i := 1
	for _, actor := range ah.Actors {
		logger.Log(i, actor.FullName())
		i += 1
	}
}
