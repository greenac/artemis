package handlers

import (
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"github.com/greenac/artemis/utils"
	"path"
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

	nu.RenameMovies(false)

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

	err = nu.AddMovie(m)
	if err != nil {
		return err
	}

	nu.RenameMovies(false)

	return nil
}

func MoveMovies(fromDir string, toDir string) error {
	fh := FileHandler{BasePath: models.FilePath{Path: fromDir}}
	err := fh.SetFiles()
	if err != nil {
		return err
	}

	for _, f := range fh.Files {
		if !f.IsDir() {
			continue
		}

		logger.Log("\n\nMoving movies for:", f.Name())

		fh2 := FileHandler{BasePath: models.FilePath{Path: f.Path()}}
		err = fh2.SetFiles()
		if err != nil {
			logger.Error("MoveMovies file handler could not set movies for:", f.Path(), err)
			continue
		}

		np := path.Join(toDir, f.Name())

		err = utils.CreateDir(np)
		if err != nil {
			continue
		}

		for i, f2 := range fh2.Files {
			if !f2.IsMovie() {
				continue
			}

			f2.NewBasePath = np

			m := models.Movie{File:   f2}
			m.GetNewName()

			err = MoveMovie(&m, External)
			if err != nil {
				logger.Warn("Failed to move movie:", i, ":", m.NewPath())
				continue
			}

			logger.Log("Moved movie:", i, ":", m.NewPath())
		}

		err = fh2.RemoveDir(f.Path())
		if err != nil {
			logger.Warn("Could not remove directory:", f.Path(), "there are still movies present")
		}
	}

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
		nn, err := nu.UpdateMovieNameWithNumber(&m, i+1)
		if err == nil {
			m.NewName = nn
		}
	}

	return nil
}

func (nu *NameUpdater) sortMovies() {
	sort.SliceStable(nu.moviesAndNumbers, func(i, j int) bool {
		return nu.moviesAndNumbers[i].Number < nu.moviesAndNumbers[j].Number
	})
}

func (nu *NameUpdater) GetMovieNumber(m *models.MovieAndNumber) (int, error) {
	if !m.IsRepeat() {
		return -1, nil
	}

	parts := strings.Split(m.NewName, ".")
	if len(parts) != 2 {
		logger.Error("NameUpdater::GetMovieNumber movie is improper format:", m.NewName)
		return -1, artemiserror.New(artemiserror.InvalidName)
	}

	n := parts[0]

	re, err := regexp.Compile(`_[0-9]*_`)
	if err != nil {
		logger.Error("NameUpdater::GetMovieNumber Could not compile regex", err)
		return -1, err
	}

	matches := re.FindAllString(n, -1)
	if len(matches) == 0 {
		return -1, nil
	}

	match := matches[len(matches)-1]
	match = strings.ReplaceAll(match, "_", "")
	mi, err := strconv.Atoi(match)
	if err != nil {
		return -1, nil
	}

	return mi, nil
}

func (nu *NameUpdater) UpdateMovieNameWithNumber(m *models.MovieAndNumber, newNum int) (string, error) {
	parts := strings.Split(m.NewNameOrName(), ".")
	if len(parts) != 2 {
		logger.Error("NameUpdater::UpdateMovieNameWithNumber movie is improper format:", m.NewNameOrName())
		return "", artemiserror.New(artemiserror.InvalidName)
	}

	name := parts[0]
	on := strconv.Itoa(m.Number)
	i := strings.LastIndex(name, on)
	if i == -1 {
		return m.NewName, nil
	}

	rn := []rune(name)

	return string(append(rn[:i], append([]rune(strconv.Itoa(newNum)), rn[i+len(on):]...)...)) + "." + parts[1], nil
}

func (nu *NameUpdater) RenameMovies(replace bool) {
	for _, m := range nu.moviesAndNumbers {
		op := m.Path()
		np := m.NewPath()

		if op == np {
			continue
		}

		err := nu.fh.Rename(op, np, replace)
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

		nName, err := nu.UpdateMovieNameWithNumber(&mn, nn)
		if err == nil {
			mn.Movie.NewName = nName
		}
	}

	nu.moviesAndNumbers = append(nu.moviesAndNumbers, mn)

	return nil
}
