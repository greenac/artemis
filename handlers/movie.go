package handlers

import (
	"fmt"
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/db"
	"github.com/greenac/artemis/dbinteractors"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"path"
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
	if err != nil { return err }

	err = os.Remove(m.Path)
	if err != nil {
		logger.Error("DeleteMovie::Failed to delete movie at path:", m.Path, "error:", err)
		return err
	}

	return dbinteractors.DeleteMovie(m.Id)
}

func AddOrganizedMovies(basePath string) error {
	fh := FileHandler{BasePath: models.FilePath{Path: basePath}}
	err := fh.SetFiles()
	if err != nil { return err }

	acts, err := dbinteractors.AllActors()
	if err != nil { return err }

	actors := *acts
	for _, df := range fh.Files {
		if !df.IsDir() {
			continue
		}

		mfh := FileHandler{BasePath: models.FilePath{Path: path.Join(df.Path())}}
		err := mfh.SetFiles()
		if err != nil { return err }

		for _, f := range mfh.Files {
			if !f.IsMovie() {
				continue
			}

			m, err := dbinteractors.FindOrCreate(f.Name(), f.Path())
			if err != nil { continue }

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
