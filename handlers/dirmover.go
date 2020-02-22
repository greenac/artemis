package handlers

import (
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type MoveMovieType int

const (
	Internal MoveMovieType = 0
	External MoveMovieType = 1
)

func OrganizeAllRepeatNamesInDir(dirPath string) error {
	fh := FileHandler{BasePath: models.FilePath{Path: dirPath}}
	err := fh.SetFiles()
	if err != nil {
		return err
	}

	for _, f := range fh.Files {
		if !f.IsDir() {
			continue
		}

		err = OrganizeRepeatNamesInDir(f.Path())
		if err != nil {
			logger.Warn("OrganizeAllRepeatNamesInDir failed to organize dir:", f.Path(), "with error:", err)
		}
	}

	return nil
}

func OrganizeRepeatNamesInDir(dirPath string) error {
	nu := NameUpdater{DirPath: dirPath, isSorted: false}

	err := nu.FillMovies()
	if err != nil {
		return err
	}

	for _, m := range nu.moviesAndNumbers {
		m.Movie.BasePath = dirPath
		m.Movie.NewBasePath = dirPath
	}

	nu.RenameMovies()

	return nil
}

func MoveMovie(m *models.Movie, ty MoveMovieType) error {
	// FIXME: Think of a way to make this more efficient. This should not have to
	// pull all the files in dirPath every time a movie is moved

	var p string
	switch ty {
	case Internal:
		p = m.BasePath
	case External:
		p = m.NewBasePath
	}

	nu := NameUpdater{DirPath: p, isSorted: false}

	err := nu.FillMovies()
	if err != nil {
		return err
	}

	//nu.RenameMovies()

	err = nu.AddMovie(m)
	if err != nil {
		return err
	}

	nu.RenameMovies()

	return nil
}

type NameUpdater struct {
	DirPath          string
	fh               FileHandler
	moviesAndNumbers []models.MovieAndNumber
	isSorted         bool
}

func (nu *NameUpdater) FillMovies() error {
	nu.fh = FileHandler{BasePath: models.FilePath{Path: nu.DirPath}}
	err := nu.fh.SetFiles()
	if err != nil {
		logger.Error("UpdateRepeatNames failed to read files in dir:", nu.DirPath, err)
		return err
	}

	mvs := make([]models.MovieAndNumber, 0)

	for _, f := range nu.fh.Files {
		if f.IsMovie() {
			m := models.Movie{
				File:   f,
				Actors: nil,
			}

			m.GetNewName()

			if !m.IsRepeat() {
				continue
			}

			m.NewBasePath = nu.DirPath

			mn := models.MovieAndNumber{Movie: &m, Number: 0}
			on, err := nu.GetMovieNumber(&mn)
			if err != nil {
				continue
			}

			mn.Number = on

			mvs = append(mvs, mn)
		}
	}

	nu.moviesAndNumbers = mvs
	nu.sortMovies()
	nu.isSorted = true

	for i, m := range nu.moviesAndNumbers {
		m.NewName = nu.UpdateMovieNameWithNumber(&m, i+1)
	}

	logger.Log("Filled:", len(nu.moviesAndNumbers), "movies")
	return nil
}

func (nu *NameUpdater) sortMovies() {
	sort.SliceStable(nu.moviesAndNumbers, func(i, j int) bool {
		return nu.moviesAndNumbers[i].Number < nu.moviesAndNumbers[j].Number
	})
}

func (nu *NameUpdater) GetMovieNumber(m *models.MovieAndNumber) (int, error) {
	parts := strings.Split(m.NewName, ".")
	if len(parts) != 2 {
		logger.Error("NameUpdater::GetMovieNumber movie is improper format:", m.NewName)
		return -1, artemiserror.New(artemiserror.InvalidName)
	}

	n := parts[0]

	if strings.Contains(n, "scene_") {
		re, err := regexp.Compile(`_[0-9]*_`)
		if err != nil {
			logger.Error("GetMovieNumber Could not compile regex", err)
			return -1, err
		}

		matches := re.FindAllString(n, -1)
		if len(matches) == 0 {
			return -1, nil
		}

		m := matches[len(matches)-1]
		m = strings.ReplaceAll(m, "_", "")
		mi, err := strconv.Atoi(m)
		if err != nil {
			return -1, nil
		}

		return mi, nil
	}

	return -1, nil
}

func (nu *NameUpdater) UpdateMovieNameWithNumber(m *models.MovieAndNumber, newNum int) string {
	name := m.NewNameOrName()
	on := strconv.Itoa(m.Number)
	i := strings.LastIndex(name, on)
	if i == -1 {
		return name
	}

	rn := []rune(name)

	return string(append(rn[:i], append([]rune(strconv.Itoa(newNum)), rn[i+len(on):]...)...))
}

func (nu *NameUpdater) RenameMovies() {
	for _, m := range nu.moviesAndNumbers {
		op := m.Path()
		np := m.NewPath()

		if op == np {
			continue
		}

		err := nu.fh.Rename(op, np, true)
		if err != nil {
			logger.Warn("NameUpdater::RenameMovies could not rename movie at:", op, "to:", np, err)
		}
	}
}

func (nu *NameUpdater) AddMovie(m *models.Movie) error {
	mn := models.MovieAndNumber{Movie: m, Number: 0}

	on, err := nu.GetMovieNumber(&mn)
	if err != nil {
		return err
	}

	if on > -1 {
		mn.Number = on
		nn := 1
		if len(nu.moviesAndNumbers) > 0 {
			if !nu.isSorted {
				nu.sortMovies()
			}

			nn = nu.moviesAndNumbers[len(nu.moviesAndNumbers)-1].Number + 1
		}

		mn.Movie.NewName = nu.UpdateMovieNameWithNumber(&mn, nn)
	}

	nu.moviesAndNumbers = append(nu.moviesAndNumbers, mn)

	return nil
}
