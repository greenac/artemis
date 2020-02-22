package handlers

import (
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"github.com/greenac/artemis/utils"
	"path"
)

type MovieHandler struct {
	DirPaths      *[]models.FilePath
	Movies        []models.Movie
	NewToPath     *models.FilePath
	KnownMovies   []*models.Movie
	UnknownMovies []*models.Movie
	unkIndex      int
}

func (mh *MovieHandler) CleanseInitialNames() error {
	if mh.DirPaths == nil {
		logger.Error("MovieHandler::CleanseInitialNames Cannot fill movies from dirs. DirPaths not initialized")
		return artemiserror.New(artemiserror.ArgsNotInitialized)
	}

	for _, p := range *mh.DirPaths {
		err := OrganizeRepeatNamesInDir(p.Path)
		if err != nil {
			logger.Warn("MovieHandler::SetMovies Could organize names in:", p.PathAsString())
		}
	}

	return nil
}

func (mh *MovieHandler) SetMovies() error {
	if mh.DirPaths == nil {
		logger.Error("MovieHandler::SetMovies Cannot fill movies from dirs. DirPaths not initialized")
		return artemiserror.New(artemiserror.ArgsNotInitialized)
	}

	mvs := make([]models.Movie, 0)
	for _, p := range *mh.DirPaths {
		fh := FileHandler{BasePath: p}
		err := fh.SetFiles()
		if err != nil {
			logger.Warn("MovieHandler::SetMovies Could not fill movies from path:", p.PathAsString())
			continue
		}

		for _, f := range fh.Files {
			if f.IsMovie() {
				m := models.Movie{File: f}
				mvs = append(mvs, m)
			}
		}
	}

	mh.Movies = mvs

	return nil
}

func (mh *MovieHandler) MoveMovies(toPath string) {
	mvs := make([]*models.Movie, 0)

	for _, m := range mh.KnownMovies {
		if m.IsKnown() {
			mvs = append(mvs, m)
		}
	}

	for _, m := range mh.UnknownMovies {
		if m.NewName != "" && m.Name() != m.NewName && len(m.Actors) > 0 {
			mvs = append(mvs, m)
		}
	}

	for _, m := range mvs {
		a := m.Actors[0]
		ap := path.Join(toPath, a.FullName())

		err := utils.CreateDir(ap)
		if err != nil {
			logger.Warn("`MovieHandler::MoveMovies` create directory:", ap)
			continue
		}

		m.GetNewName()
		m.NewBasePath = ap

		err = MoveMovie(m, External)
		if err != nil {
			logger.Warn("`MovieHandler::MoveMovies` could not rename:", m.Path(), "to:", m.NewPath(), err)
		}
	}
}

func (mh *MovieHandler) RenameMovies(mvs []*models.Movie) {
	for _, m := range mvs {
		mh.RenameMovie(m)
	}
}

func (mh *MovieHandler) RenameMovie(m *models.Movie) error {
	if m.BasePath == "" {
		logger.Warn("`MovieHandler::RenameMovie` movie:", m.Name(), "does not have path set")
		return artemiserror.New(artemiserror.PathNotSet)
	}

	m.NewBasePath = m.BasePath

	err := MoveMovie(m, Internal)
	if err != nil {
		logger.Warn("`MovieHandler::RenameMovie` movie:", m.Path, "failed to be renamed with error:", err)
		return err
	}

	return nil
}

func (mh *MovieHandler) AddKnownMovie(m models.Movie) {
	mh.KnownMovies = append(mh.KnownMovies, &m)
}

func (mh *MovieHandler) AddUnknownMovie(m models.Movie) {
	mh.UnknownMovies = append(mh.UnknownMovies, &m)
}

func (mh *MovieHandler) UpdateUnknownMovies(unMvs *[]*models.Movie) {
	mh.UnknownMovies = *unMvs
}

func (mh *MovieHandler) AddKnownMovieNames() {
	for _, m := range mh.KnownMovies {
		m.AddActorNames()
	}
}

func (mh *MovieHandler) AddUnknownMovieNames() {
	for _, m := range mh.UnknownMovies {
		m.AddActorNames()
	}
}

func (mh *MovieHandler) RenameAllMovies() {
	mvs := make([]*models.Movie, 0)
	for _, m := range mh.KnownMovies {
		if m.NewName != m.Info.Name() {
			mvs = append(mvs, m)
		}
	}

	for _, m := range mh.UnknownMovies {
		if m.NewName != m.Info.Name() {
			mvs = append(mvs, m)
		}
	}

	mh.RenameMovies(mvs)
}

func (mh *MovieHandler) IncrementUnknownIndex() {
	mh.unkIndex += 1
}

func (mh *MovieHandler) CurrentUnknownMovie() *models.Movie {
	return mh.UnknownMovies[mh.unkIndex]
}

func (mh *MovieHandler) MoreUnknowns() bool {
	return mh.unkIndex >= len(mh.UnknownMovies)
}
