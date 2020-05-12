package dbinteractors

import (
	"github.com/greenac/artemis/db"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sort"
	"strings"
)

func NewMovie(name string, path string) models.Movie {
	m := models.Movie{
		Name: name,
		Path: path,
	}

	m.SetIdentifier()
	m.ActorIds = make([]primitive.ObjectID, 0)

	return m
}

func GetMovieById(id primitive.ObjectID) (*models.Movie, error) {
	var m models.Movie
	res, err := models.FindById(id, db.MovieCollection)
	if err != nil {
		logger.Error("GetMovieById::Failed to fetch model with id:", id, "error:", err)
		return nil, err
	}

	err = res.Decode(&m)
	if err != nil {
		logger.Error("GetMovieById::Failed to decode model with id:", id, "error:", err)
		return nil, err
	}

	return &m, nil
}

func GetMovieByIdString(id string) (*models.Movie, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("GetMovieByIdString::failed to create ObjectId from:", id, "error:", err)
		return nil, err
	}

	return GetMovieById(objId)
}

func GetMovieByIdentifier(id string) (*models.Movie, error) {
	var m models.Movie
	res, err := models.FindByIdentifier(id, db.MovieCollection)
	if err != nil {
		logger.Error("GetMovieByIdentifier::Failed to fetch model with identifier:", id, "error:", err)
		return nil, err
	}

	err = res.Decode(&m)
	if err != nil {
		logger.Error("GetMovieByIdentifier::Failed to decode model with identifier:", id, "error:", err)
		return nil, err
	}

	return &m, nil
}

func DoesMovieExist(identifier string) (bool, error) {
	m, err := GetMovieByIdentifier(identifier)
	if err != nil {
		return false, err
	}

	return m != nil, nil
}

func UnknownMovies() (*[]models.Movie, error) {
	cAndT, err := db.GetCollectionAndContext(db.MovieCollection)
	if err != nil {
		return nil, err
	}

	q := map[string]interface{}{"$size": 0}

	c, err := cAndT.Col.Find(cAndT.Ctx, bson.M{"actorIds": q})
	if err != nil {
		logger.Error("UnknownMovies::Failed to find unknown movies:", err)
		return nil, err
	}

	mvs := make([]models.Movie, 0)

	defer c.Close(cAndT.Ctx)

	for c.Next(cAndT.Ctx) {
		var m models.Movie
		err := c.Decode(&m)
		if err != nil {
			logger.Warn("UnknownMovies::Failed to decode movie with error:", err)
			continue
		}

		mvs = append(mvs, m)
	}

	sort.SliceStable(mvs, func(i, j int) bool {
		return strings.ToLower(mvs[i].Name) < strings.ToLower(mvs[j].Name)
	})

	return &mvs, nil
}
