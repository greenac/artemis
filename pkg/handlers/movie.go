package handlers

import (
	"errors"
	"fmt"
	"github.com/greenac/artemis/pkg/artemiserror"
	"github.com/greenac/artemis/pkg/db"
	"github.com/greenac/artemis/pkg/dbinteractors"
	"github.com/greenac/artemis/pkg/logger"
	"github.com/greenac/artemis/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"path"
	"strings"
	"time"
)

type MovieHandler struct {
	DirPaths      *[]models.FilePath
	Movies        []models.SysMovie
	KnownMovies   []*models.Movie
	UnknownMovies []*models.Movie
	unkIndex      int
}

func (mh *MovieHandler) SetMovies() error {
	if mh.DirPaths == nil {
		logger.Error("MovieHandler::SetMovies Cannot fill movies from dirs. DirPaths not initialized")
		return artemiserror.New(artemiserror.ArgsNotInitialized)
	}

	mvs := make([]models.SysMovie, 0)
	for _, p := range *mh.DirPaths {
		fh := FileHandler{BasePath: p}
		err := fh.SetFiles()
		if err != nil {
			logger.Warn("MovieHandler::SetMovies Could not fill movies from path:", p.PathAsString())
			continue
		}

		for _, f := range fh.Files {
			if f.IsMovie() {
				m := models.SysMovie{File: f}
				mvs = append(mvs, m)
			}
		}
	}

	mh.Movies = mvs

	return nil
}

func (mh *MovieHandler) RenameMovies(mvs []*models.SysMovie) {
	for _, m := range mvs {
		_ = mh.RenameMovie(m)
	}
}

func (mh *MovieHandler) RenameMovie(m *models.SysMovie) error {
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

func (mh *MovieHandler) AddKnownMovie(m models.SysMovie) {
	dbm, err := dbinteractors.GetMovieByIdentifier(m.Path())
	if err != nil {
		logger.Warn("MovieHandler::AddKnownMovie could not add movie:", m.Name(), "Failed with error:", err)
		return
	}

	if dbm == nil {
		nm := dbinteractors.NewMovie(m.Name(), m.Path())
		_ = nm.Save()
		dbm = &nm
	}

	mh.KnownMovies = append(mh.KnownMovies, dbm)
}

func AddActorsToMovie(movieId string, actorIds []string) error {
	logger.Debug("movie id:", movieId, "actor ids:", actorIds)

	movId, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		logger.Error("AddActorsToMovie::failed to create ObjectId from:", movieId, "error:", err)
		return err
	}

	m, err := dbinteractors.GetMovieById(movId)
	if err != nil {
		return err
	}

	save := false
	for _, aId := range actorIds {
		actId, err := primitive.ObjectIDFromHex(aId)
		if err != nil {
			logger.Warn("AddActorsToMovie::failed to create ObjectId from actorId:", aId, "error:", err)
			continue
		}

		a, err := dbinteractors.GetActorById(actId)
		if err != nil {
			logger.Warn("AddActorsToMovie::Could not get actor with id:", actId, err)
			continue
		}

		a.AddMovie(movId)
		a.Updated = time.Now()
		err = a.Save()
		if err != nil {
			logger.Warn("AddActorsToMovie::Could not add movie:", m.Name, "to actor:", a.FullName(), "error:", err)
			continue
		}

		m.AddActor(actId)
		save = true
	}

	if save {
		m.Updated = time.Now()
		_ = m.Save()
	}

	return nil
}

func GetMovieWithIds(ids []string) (*[]models.Movie, error) {
	objIds := make([]primitive.ObjectID, len(ids))
	for i, id := range ids {
		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			logger.Warn("GetMovieIds could not make object id from:", id)
			continue
		}

		objIds[i] = objId
	}

	return dbinteractors.MoviesForIds(objIds)
}

func GetMoviesForActor(actorId string) (*[]models.Movie, error) {
	actorObjId, err := primitive.ObjectIDFromHex(actorId)
	if err != nil {
		logger.Warn("GetMoviesForActor could not make object id from:", actorId)
		return nil, err
	}

	act, err := dbinteractors.GetActorById(actorObjId)
	if err != nil {
		return nil, err
	}

	return dbinteractors.MoviesForIds(act.MovieIds)
}

func DeleteMovie(movieId string) error {
	m, err := dbinteractors.GetMovieByIdString(movieId)
	if err != nil {
		return err
	}

	err = os.Remove(m.Path)
	if err != nil {
		logger.Error("DeleteMovie::Failed to delete movie at path:", m.Path, "error:", err)
		return err
	}

	return dbinteractors.DeleteMovie(m.Id)
}

func DeleteMovieById(id primitive.ObjectID) error {
	m, err := dbinteractors.GetMovieById(id)
	if err != nil {
		return err
	}

	err = os.Remove(m.Path)
	if err != nil {
		logger.Error("DeleteMovieById::Failed to delete movie at path:", m.Path, "error:", err)
		return err
	}

	return dbinteractors.DeleteMovie(m.Id)
}

func AddOrganizedMovies(basePath string) error {
	fh := FileHandler{BasePath: models.FilePath{Path: basePath}}
	err := fh.SetFiles()
	if err != nil {
		return err
	}

	acts, err := dbinteractors.AllActors()
	if err != nil {
		return err
	}

	actors := *acts
	for _, df := range fh.Files {
		if !df.IsDir() {
			continue
		}

		mfh := FileHandler{BasePath: models.FilePath{Path: path.Join(df.Path())}}
		err := mfh.SetFiles()
		if err != nil {
			return err
		}

		for _, f := range mfh.Files {
			if !f.IsMovie() {
				continue
			}

			m, err := dbinteractors.FindOrCreate(f.Name(), f.Path())
			if err != nil {
				continue
			}

			save := false
			for i, a := range actors {
				if a.IsIn(m.Name) {
					logger.Log("Adding actor:", a.FullName(), "to movie:", m.Name, "id:", m.Id)
					m.AddActor(a.Id)
					a.AddMovie(m.Id)
					a.Updated = time.Now()
					logger.Debug("actor movies before save:", a.MovieIds)
					_ = a.Save()
					save = true
					actors[i] = a
				}
			}

			if save {
				logger.Log("Saving movie:", m.Name, "with actors:", m.ActorIds)
				m.Updated = time.Now()
				_ = m.Save()
			} else {
				logger.Warn("Could not add any actors for movie:", m.Name)
			}
		}
	}

	return nil
}

func SearchMoviesByDate(name string) (*[]models.Movie, error) {
	logger.Debug("Searching movies with name:", name)
	movs := make([]models.Movie, 0)

	var filter = bson.D{}

	if name != "" {
		filter = bson.D{
			{
				Key: "name",
				Value: bson.D{
					{
						"$regex",
						primitive.Regex{
							Pattern: fmt.Sprintf("^%s", name),
							Options: "i",
						},
					},
				},
			},
		}
	}

	logger.Debug("filter is:", filter)

	cAndT, err := db.GetCollectionAndContext(db.MovieCollection)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.SetSort(bson.D{
		{"updated", -1},
	})
	opts.SetLimit(200)

	cur, err := cAndT.Col.Find(cAndT.Ctx, filter, opts)
	if err != nil {
		logger.Error("SearchMoviesByDate::Failed with error:", err)
		return nil, err
	}

	defer cur.Close(cAndT.Ctx)

	for cur.Next(cAndT.Ctx) {
		var m models.Movie

		err := cur.Decode(&m)
		if err != nil {
			logger.Warn("SearchMoviesByDate::Failed to decode movie with error:", err)
			continue
		}

		if len(m.ActorIds) > 0 {
			acts, err := dbinteractors.GetActorsForIds(m.ActorIds)
			if err == nil {
				m.Actors = acts
			}
		}

		movs = append(movs, m)
	}

	return &movs, nil
}

func ActorsInMovie(movieId string) (*[]models.Actor, error) {
	m, err := dbinteractors.GetMovieByIdString(movieId)
	if err != nil {
		return nil, err
	}

	acts, err := dbinteractors.GetActorsForIds(m.ActorIds)
	if err != nil {
		return nil, err
	}

	return acts, nil
}

func RemoveActorFromMovie(movieId string, actorId string) (*models.Movie, error) {
	// FIXME: Add a db transaction
	a, err := dbinteractors.GetActorByIdString(actorId)
	if err != nil {
		return nil, err
	}

	m, err := dbinteractors.GetMovieByIdString(movieId)
	if err != nil {
		return nil, err
	}

	a.RemoveMovie(m.Id)
	err = a.Save()
	if err != nil {
		return nil, err
	}

	m.RemoveActor(a.Id)
	err = m.Save()

	return m, err
}

func UpdatedMoviePaths() error {
	mvs, err := dbinteractors.FetchAllMovies()
	if err != nil {
		return err
	}

	const targetPath = "/Volumes/Foxtrot"

	const orgPathSeg = "/Volumes/Foxtrot/.p/organized"
	const org1PathSeg = "/Volumes/Foxtrot/.p/organized1"
	const rawPathSeg = "/Volumes/Foxtrot/.raw"

	const orgRepSeg = "/Volumes/Golf/.p/organized2"
	const org1RepSeg = "/Volumes/Golf/.p/organized1"
	const rawRepSeg = "/Volumes/Golf/.p/raw"

	for _, m := range *mvs {
		if !strings.Contains(m.Path, targetPath) {
			//logger.Warn("UpdatedMoviePaths::target path not found in:", m.Path)
			continue
		}

		var originalSeg string
		var replSeg string

		if strings.Contains(m.Path, orgPathSeg) {
			originalSeg = orgPathSeg
			replSeg = orgRepSeg
		} else if strings.Contains(m.Path, org1PathSeg) {
			originalSeg = org1PathSeg
			replSeg = org1RepSeg
		} else if strings.Contains(m.Path, rawPathSeg) {
			originalSeg = rawPathSeg
			replSeg = rawRepSeg
		} else {
			logger.Warn("UpdatedMoviePaths::does not handle path:", m.Path)
			continue
		}

		m.Path = strings.Replace(m.Path, originalSeg, replSeg, 1)
		m.SetIdentifier()

		err = m.Save()
		if err != nil {
			logger.Error("UpdatedMoviePaths::Failed to save movie", m.Path)
			continue
		}

		//logger.Warn("UpdatedMoviePaths::new path:", m.Path, "old path:", oldPath, "old id:", oldId, "new Id:", m.Identifier )
	}

	return nil
}

func AddActorToMovies() error {
	movies, err := dbinteractors.FetchAllMovies()
	if err != nil {
		return err
	}

	actors, err := dbinteractors.AllActors()
	if err != nil {
		return err
	}

	for _, m := range *movies {
		for _, a := range *actors {
			if a.IsIn(m.Name) {
				if !a.HasMovie(m.Id) {
					a.AddMovie(m.Id)
					err := a.Save()
					if err != nil {
						logger.Error("AddActorToMovies::Save::Failed to save movie", m.Name, "for actor:", a.FullName())
						continue
					}

					logger.Log("AddActorToMovies->save movie", m.Name, "in actor:", a.FullName())
				}

				if !m.HasActor(a.Id) {
					logger.Log("Movie:", m.Name, "has actor:", a.FullName())
					m.AddActor(a.Id)
					err = m.Save()
					if err != nil {
						logger.Error("AddActorToMovies->failed to save movie", m.Name, "with actors:", m.ActorIds)
					}
				}
			}
		}
	}

	return nil
}

func FixMoviePath21() error {
	mvs, err := dbinteractors.MoviesWith21Path()
	if err != nil {
		logger.Error("FixMoviesWithPath21::failed to access movies in db with error:", err)
		return err
	}

	for _, m := range *mvs {
		logger.Log(m.Path)
		//	newPath := strings.ReplaceAll(m.Path, "organized21", "organized1")
		//	_, err := os.Stat(newPath)
		//	if err != nil {
		//		if os.IsNotExist(err) {
		//			logger.Error("Path does not exist for movie:", m.Name, "at path:", newPath)
		//		} else {
		//			logger.Error("Failed to get file info for movie:", m.Name, err)
		//		}
		//
		//		continue
		//	}
		//
		//	logger.Log("File", m.Name, "exists at path:", newPath)
		//	m.Path = newPath
		//
		//	err = m.Save()
		//	if err != nil {
		//		logger.Error("Failed to save movie:", m.Name, err)
		//	}
	}

	return nil
}

func DeleteNonExistentMovies() error {
	mvs, err := dbinteractors.MoviesWithGolfPath()
	if err != nil {
		logger.Error("DeleteNonExistentMovies::failed to fetch movies with error:", err)
		return err
	}

	for _, m := range *mvs {
		_, err := os.Stat(m.Path)

		if err == nil {
			continue
		}

		if !errors.Is(err, os.ErrNotExist) {
			logger.Error("DeleteNonExistentMovies::Failed with error:", err)
			continue
		}

		for _, id := range m.ActorIds {
			a, err := dbinteractors.GetActorById(id)
			if err != nil {
				logger.Error("DeleteNonExistentMovies::Could not get actor with id:", id, "for movie:", m.Path)
				continue
			}

			logger.Log("DeleteNonExistentMovies::will delete actor:", a.FullName())

			a.RemoveMovie(m.Id)
			_ = a.Save()
		}

		logger.Log("DeleteNonExistentMovies::will delete movie:", m.Path)

		err = dbinteractors.DeleteMovie(m.Id)
		if err != nil {
			logger.Error("DeleteNonExistentMovies::Could not delete movie:", m.Name, err)
		}
	}

	//return ioutil.WriteFile(
	//	"/Users/andre/Documents/missing-movies.txt",
	//	[]byte(out),
	//	0644,
	//)

	return nil
}

func AddMissingMoviesToActor() error {
	mvs, err := dbinteractors.FetchAllMovies()
	if err != nil {
		return err
	}

	count := 0
	for _, m := range *mvs {
		for _, id := range m.ActorIds {
			a, err := dbinteractors.GetActorById(id)
			if err != nil {
				logger.Error("AddMissingMoviesToActor::Could not get actor with id:", id, "for movie:", m.Path)
				continue
			}

			hasMovie := false
			for _, mId := range a.MovieIds {
				if mId == m.Id {
					hasMovie = true
					break
				}
			}

			if !hasMovie {
				logger.Log("AddMissingMoviesToActor::is missing movie:", m.Name, "actor:", a.FullName())
				a.AddMovie(m.Id)
				err = a.Save()
				if err != nil {
					logger.Error("AddMissingMoviesToActor::failed to save movie:", m.Name, "actor:", a.FullName())
				}
				count += 1
			}
		}
	}

	logger.Log("AddMissingMoviesToActor::number of missing movies", count)

	return nil
}

func AddSecondaryPaths() error {
	mvs, err := dbinteractors.FetchAllMovies()
	if err != nil {
		return err
	}

	for _, m := range *mvs {
		if m.Path == "" {
			logger.Error("AddSecondaryPaths->got movie with empty path:", m.Name, m.Identifier)
			return errors.New("movie missing path")
		}

		if m.SecondaryPath != "" {
			logger.Log("AddSecondaryPaths->secondary path is set for movie", m.Name, m.Identifier, m.SecondaryPath)
			continue
		}

		var secPath string
		if strings.Contains(m.Path, "/Volumes/Golf/.p/") {
			secPath = strings.Replace(m.Path, "/Volumes/Golf/.p/", "/Volumes/Hotel/.volumes/golf/", 1)
		} else if strings.Contains(m.Path, "/Volumes/Papa/.p/") {
			secPath = strings.Replace(m.Path, "/Volumes/Papa/.p/", "/Volumes/Hotel/.volumes/papa/", 1)
		} else {
			logger.Error("AddSecondaryPaths->got movie with invalid path:", m.Path, m.Identifier)
			return errors.New("invalid path")
		}

		m.SecondaryPath = secPath
		err = m.Save()
		if err != nil {
			logger.Error("AddSecondaryPaths->got movie with invalid path:", m.Name, m.Identifier, "with error:", err.Error())
			return err
		}

		logger.Log("AddSecondaryPaths->saved movie:", m.Name, "to:", m.SecondaryPath)
	}

	return nil
}

func SwapPaths() error {
	mvs, err := dbinteractors.FetchAllMovies()
	if err != nil {
		return err
	}

	for _, m := range *mvs {
		if m.Path == "" {
			logger.Error("SwapPaths->got movie with empty path:", m.Name, m.Id)
			continue
		}

		if m.SecondaryPath == "" {
			logger.Warn("SwapPaths->secondary path is not set for movie", m.Name, m.Id)
			continue
		}

		if strings.Contains(m.Path, "/Volumes/Hotel/.volumes/") {
			logger.Warn("SwapPaths->path contains target", m.Name, m.Id, m.Path)
			continue
		}

		p := m.SecondaryPath
		m.SecondaryPath = m.Path
		m.Path = p
		m.SetIdentifier()

		err = m.Save()
		if err != nil {
			logger.Error("SwapPaths->saving movie", m.Name, m.Id, "failed with error:", err.Error())
			continue
		}

		logger.Log("SwapPaths->saved movie:", m.String())
	}

	return nil
}
