package handlers

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//func MoveDir(dir models.File) error {
//	if !dir.IsDir() {
//		return nil
//	}
//
//	fh := FileHandler{BasePath: models.FilePath{Path: dir.Path()}}
//
//	err := fh.SetFiles()
//	if err != nil {
//		logger.Error("MoveDir failed to read files in dir:", dir.Path, err)
//		return err
//	}
//
//	exists, err := fh.DoesFileExistAtPath(dir.NewPath())
//	if err != nil {
//		logger.Warn("MoveDir failed to move directory to:", dir.GetNewTotalPath(), err)
//		return err
//	}
//
//	if exists {
//		for _, f := range fh.Files {
//
//			f.NewBasePath = dir.NewPath()
//			if f.IsMovie() {
//				err := fh.Rename(f.Path(), f.GetNewTotalPath())
//				if err != nil {
//					continue
//				}
//			}
//		}
//	} else {
//		err = fh.Rename(dir.Path(), dir.GetNewTotalPath())
//		return err
//	}
//
//	return nil
//}

//func AddMovieToDir(m *models.Movie, dir models.File) error {
//	nu := NameUpdater{DirPath: dir.Path()}
//	nu.SetUp()
//
//
//}


type NameUpdater struct {
	DirPath string
	fh FileHandler
	moviesAndNumbers []models.MovieAndNumber
	isSorted bool
}

func (nu *NameUpdater) SetUp() {
	nu.fh = FileHandler{BasePath: models.FilePath{Path: nu.DirPath}}
}

func (nu *NameUpdater) FillMovies(dirPath string) error {
	err := nu.fh.SetFiles()
	if err != nil {
		logger.Error("UpdateRepeatNames failed to read files in dir:", dirPath, err)
		return err
	}

	mvs := make([]models.MovieAndNumber, 0)

	for _, f := range nu.fh.Files {
		if f.IsMovie() && strings.Contains(f.Name(), "scene_") {
			m := models.Movie{
				File:   f,
				Actors: nil,
			}

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
	nu.SortMovies()
	nu.isSorted = true

	for i, m := range nu.moviesAndNumbers {
		m.NewName = nu.UpdateMovieNameWithNumber(&m, i + 1)
	}


	return nil
}

func (nu *NameUpdater) SortMovies() {
	sort.SliceStable(nu.moviesAndNumbers, func(i, j int) bool {
		return nu.moviesAndNumbers[i].Number < nu.moviesAndNumbers[j].Number
	})
}

func (nu *NameUpdater) GetMovieNumber(m *models.MovieAndNumber) (int, error) {
	parts := strings.Split(m.Name(), ".")
	n := parts[0]

	if strings.Contains(n, "scene_", ) {
		re, err := regexp.Compile(`_.[0-9]+_`)
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
	name := m.Name()
	on := strconv.Itoa(m.Number)
	i := strings.LastIndex(name, on)
	logger.Log("index for name:", name, "is:", i, "for number:", m.Number)
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

		nu.fh.Rename(op, np, true)
	}
}

func (nu *NameUpdater) AddMovie(m *models.Movie) error {
	mn := models.MovieAndNumber{Movie:  m, Number: 0,}

	on, err := nu.GetMovieNumber(&mn)
	if err != nil {
		return err
	}

	if on > -1 {
		nn := 1
		if len(nu.moviesAndNumbers) > 0 {
			if !nu.isSorted {
				nu.SortMovies()
			}

			nn = nu.moviesAndNumbers[len(nu.moviesAndNumbers)-1].Number + 1
		}

		mn.Movie.NewName = nu.UpdateMovieNameWithNumber(&mn, nn)
	}


	nu.moviesAndNumbers = append(nu.moviesAndNumbers, mn)

	return nil
}
