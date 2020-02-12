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
	Movies        *[]models.Movie
	NewToPath     *models.FilePath
	UnknownMovies []*models.Movie
	unkIndex int
}

func (mh *MovieHandler) SetMovies() error {
	if mh.DirPaths == nil {
		logger.Error("Cannot fill movies from dirs. DirPaths not initialized")
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

		for _, f := range *fh.Files {
			if utils.IsMovie(&f) {
				m := models.Movie{File: f}
				m.Path = path.Join(p.Path, *m.Name())
				mvs = append(mvs, m)
			}
		}
	}

	mh.Movies = &mvs
	return nil
}

func (mh *MovieHandler) RenameMovies(mvs []*models.Movie) {
	for _, m := range mvs {
		mh.RenameMovie(m)
	}
}

func (mh *MovieHandler) RenameMovie(m *models.Movie) error {
	if m.Path == "" {
		logger.Warn("`MovieHandler::RenameMovie` movie:", m.Name(), "does not have path set")
		return artemiserror.New(artemiserror.PathNotSet)
	}

	fh := FileHandler{}
	err := fh.Rename(m.Path, m.RenamePath())
	if err != nil {
		logger.Warn("`MovieHandler::RenameMovie` movie:", m.Name(), "failed to be renamed with error:", err)
		return err
	}

	return nil
}

func (mh *MovieHandler) AddUnknownMovie(m models.Movie) {
	mh.UnknownMovies = append(mh.UnknownMovies, &m)
	//logger.Debug("start", len(mh.UnknownMovies), m)

	//var unknowns []*models.Movie
	//
	//switch len(mh.UnknownMovies) {
	//case 0:
	//	//logger.Debug("0", m.NewName)
	//	unknowns = []*models.Movie{&m}
	//case 1:
	//	if m.NewName < mh.UnknownMovies[0].NewName {
	//		//logger.Debug("case 1 m.NewName < mh.UnknownMovies[0].NewName", m.NewName,  mh.UnknownMovies[0].NewName)
	//		unknowns = append([]*models.Movie{&m}, mh.UnknownMovies...)
	//	} else {
	//		//logger.Debug("case 1  m.NewName > mh.UnknownMovies[0].NewName",  m.NewName,  mh.UnknownMovies[0].NewName)
	//		unknowns = append(mh.UnknownMovies, &m)
	//	}
	//default:
	//	//logger.Debug("case default before", len(mh.UnknownMovies), mh.UnknownMovies)
	//
	//	t := 0
	//	for i := 0; i < len(mh.UnknownMovies) - 1; i += 1 {
	//		//logger.Debug(mh.UnknownMovies[i].NewName, m.NewName, mh.UnknownMovies[i+1].NewName)
	//		if m.NewName > mh.UnknownMovies[i].NewName && m.NewName < mh.UnknownMovies[i + 1].NewName {
	//			t = i
	//			break
	//		}
	//	}
	//
	//	//logger.Debug("Target index is:", t, "unknown length:", len(mh.UnknownMovies))
	//	unknowns = make([]*models.Movie, len(mh.UnknownMovies))
	//	copy(unknowns, mh.UnknownMovies)
	//
	//	unknowns = append(unknowns[:t], append([]*models.Movie{&m}, unknowns[t:]...)...)


		//logger.Debug("Unknowns has length after;", len(unknowns), unknowns)
	//}

	//mh.UnknownMovies = unknowns
}

func (mh *MovieHandler) UpdateUnknownMovies(unMvs *[]*models.Movie) {
	mh.UnknownMovies = *unMvs
}

func (mh *MovieHandler) AddUnknownMovieNames() {
	for _, m := range mh.UnknownMovies {
		m.AddActorNames()
	}
}

func (mh *MovieHandler) RenameUnknownMovies() {
	mvs := make([]*models.Movie, 0)
	for _, m := range(mh.UnknownMovies) {
		if m.NewName != m.Info.Name() {
			mvs = append(mvs, m)
		}
	}

	logger.Debug("MovieHandler::RenameUnknownMovies renaming:", len(mvs))

	mh.RenameMovies(mvs)
}

func (mh *MovieHandler) IncrementUnknownIndex() {
	mh.unkIndex += 1
}

func (mh *MovieHandler) CurrentUnknownMovie() *models.Movie {
	return mh.UnknownMovies[mh.unkIndex]
}

func (mh *MovieHandler) MoreUnknowns() bool {
	return mh.unkIndex >=  len(mh.UnknownMovies)
}
