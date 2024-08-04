package dbinteractors

import (
	"github.com/greenac/artemis/pkg/db"
	"github.com/greenac/artemis/pkg/logger"
	"github.com/greenac/artemis/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"strings"
	"time"
)

const MaxMoviesToReturn int64 = 50

func NewMovie(name string, path string) models.Movie {
	m := models.Movie{
		Name: name,
		Path: path,
	}

	m.SetIdentifier()
	m.ActorIds = make([]primitive.ObjectID, 0)
	m.Updated = time.Now()

	return m
}

func FindOrCreate(name string, path string) (*models.Movie, error) {
	mv, err := GetMovieByIdentifier(models.MovieIdentifier(path))
	if err != nil && err.Error() != "mongo: no documents in result" {
		return nil, err
	}

	if mv != nil {
		return mv, nil
	}

	m := NewMovie(name, path)
	_, err = m.Create()

	return &m, nil
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
		logger.Error("GetMovieByIdentifier::Failed to decode movie with identifier:", id, "error:", err)
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

func UnknownMovies(page int, size int) (*[]models.Movie, int64, error) {
	cAndT, err := db.GetCollectionAndContext(db.MovieCollection)
	if err != nil {
		return nil, 0, err
	}

	q := map[string]interface{}{"$size": 0}

	c, err := cAndT.Col.Find(cAndT.Ctx, bson.M{"actorIds": q})
	if err != nil {
		logger.Error("UnknownMovies::Failed to find unknown movies:", err)
		return nil, 0, err
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

	start := page * size
	end := start + size
	if len(mvs)-int(start) < int(size) {
		size = len(mvs) - start
		end = start + size
	}

	pMvs := make([]models.Movie, size)

	for i := start; i < end; i += 1 {
		pMvs[i-start] = mvs[i]
	}

	return &pMvs, int64(len(mvs)), nil
}

func MoviesForIds(ids []primitive.ObjectID) (*[]models.Movie, error) {
	mvs := make([]models.Movie, 0)

	if len(ids) == 0 {
		return &mvs, nil
	}

	cAndT, err := db.GetCollectionAndContext(db.MovieCollection)
	if err != nil {
		return nil, err
	}

	params := primitive.A{}

	for _, id := range ids {
		v := bson.D{
			{
				Key:   "_id",
				Value: id,
			},
		}

		params = append(params, v)
	}

	q := bson.D{
		{
			Key:   "$or",
			Value: params,
		},
	}

	c, err := cAndT.Col.Find(cAndT.Ctx, q)
	if err != nil {
		logger.Error("MoviesForIds::Failed to find unknown movies:", err)
		return nil, err
	}

	defer c.Close(cAndT.Ctx)

	for c.Next(cAndT.Ctx) {
		var m models.Movie
		err := c.Decode(&m)
		if err != nil {
			logger.Warn("MoviesForIds::Failed to decode movie with error:", err)
			continue
		}

		mvs = append(mvs, m)
	}

	sort.SliceStable(mvs, func(i, j int) bool {
		return strings.ToLower(mvs[i].Name) < strings.ToLower(mvs[j].Name)
	})

	return &mvs, nil
}

func DeleteMovie(id primitive.ObjectID) error {
	cAndT, err := db.GetCollectionAndContext(db.MovieCollection)
	if err != nil {
		return err
	}

	par := bson.D{{Key: "_id", Value: id}}
	_, err = cAndT.Col.DeleteOne(cAndT.Ctx, par)
	return err
}

func FetchAllMovies() (*[]models.Movie, error) {
	cAndT, err := db.GetCollectionAndContext(db.MovieCollection)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.SetBatchSize(50)

	c, err := cAndT.Col.Find(cAndT.Ctx, bson.D{}, opts)
	defer c.Close(cAndT.Ctx)

	mvs := make([]models.Movie, 0)

	for c.Next(cAndT.Ctx) {
		var m models.Movie
		err := c.Decode(&m)
		if err != nil {
			logger.Warn("FetchAllMovies::Failed to decode movie with error:", err)
			continue
		}

		mvs = append(mvs, m)
	}

	logger.Log("FetchAllMovies::got:", len(mvs), "movies")

	return &mvs, nil
}

func MoviesWith21Path() (*[]models.Movie, error) {
	cAndT, err := db.GetCollectionAndContext(db.MovieCollection)
	if err != nil {
		return nil, err
	}

	q := bson.D{
		{
			Key: "path",
			Value: bson.D{
				{
					"$regex",
					primitive.Regex{
						Pattern: "organized21",
						Options: "i",
					},
				},
			},
		},
	}

	c, err := cAndT.Col.Find(cAndT.Ctx, q)

	defer c.Close(cAndT.Ctx)

	mvs := make([]models.Movie, 0)

	for c.Next(cAndT.Ctx) {
		var m models.Movie
		err := c.Decode(&m)
		if err != nil {
			logger.Warn("FetchAllMovies::Failed to decode movie with error:", err)
			continue
		}

		mvs = append(mvs, m)
	}

	return &mvs, nil
}

func MoviesWithGolfPath() (*[]models.Movie, error) {
	cAndT, err := db.GetCollectionAndContext(db.MovieCollection)
	if err != nil {
		return nil, err
	}

	q := bson.D{
		{
			Key: "path",
			Value: bson.D{
				{
					"$regex",
					primitive.Regex{
						Pattern: `\/Golf\/`,
						Options: "i",
					},
				},
			},
		},
	}

	c, err := cAndT.Col.Find(cAndT.Ctx, q)

	defer c.Close(cAndT.Ctx)

	mvs := make([]models.Movie, 0)

	for c.Next(cAndT.Ctx) {
		var m models.Movie
		err := c.Decode(&m)
		if err != nil {
			logger.Warn("FetchAllMovies::Failed to decode movie with error:", err)
			continue
		}

		mvs = append(mvs, m)
	}

	return &mvs, nil
}

func GetMoviesForInput(input string, page int) (*[]models.Movie, error) {
	logger.Log("GetMoviesForInput::input is:", input, "page:", page)

	movies := []models.Movie{}
	if input == "" {
		return &movies, nil
	}

	filter := bson.D{
		{
			Key: "$or",
			Value: bson.A{
				bson.D{
					{
						Key: "name",
						Value: bson.D{
							{
								"$regex",
								primitive.Regex{
									Pattern: input,
									Options: "i",
								},
							},
						},
					},
				},
			},
		},
	}

	cAndT, err := db.GetCollectionAndContext(db.MovieCollection)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.SetSort(bson.D{
		{"name", 1},
	})
	opts.SetSkip(int64(page) * MaxMoviesToReturn)
	opts.SetLimit(MaxMoviesToReturn)

	cur, err := cAndT.Col.Find(cAndT.Ctx, filter, opts)
	if err != nil {
		logger.Error("GetMoviesForInput::Failed with error:", err)
		return nil, err
	}

	defer cur.Close(cAndT.Ctx)

	for cur.Next(cAndT.Ctx) {
		var m models.Movie

		err := cur.Decode(&m)
		if err != nil {
			logger.Warn("GetMoviesForInput::Failed to decode movie with error:", err)
			continue
		}

		movies = append(movies, m)
	}

	sort.Slice(movies, func(i int, j int) bool {
		return movies[i].Name < movies[j].Name
	})

	return &movies, nil
}

func GetCountOfMoviesForInput(input string) (int64, error) {
	if input == "" {
		return 0, nil
	}

	filter := bson.D{
		{
			Key: "$or",
			Value: bson.A{
				bson.D{
					{
						Key: "name",
						Value: bson.D{
							{
								"$regex",
								primitive.Regex{
									Pattern: input,
									Options: "i",
								},
							},
						},
					},
				},
			},
		},
	}

	cAndT, err := db.GetCollectionAndContext(db.MovieCollection)
	if err != nil {
		return 0, err
	}

	total, err := cAndT.Col.CountDocuments(cAndT.Ctx, filter, nil)
	if err != nil {
		logger.Error("GetMoviesForInput::Failed with error:", err)
		return 0, err
	}

	logger.Log("GetCountOfMoviesForInput->got total:", total, "from input:", input)

	return total, nil
}
