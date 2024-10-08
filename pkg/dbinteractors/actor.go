package dbinteractors

import (
	"fmt"
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

const MaxActorsToReturn = 25

type PaginatedQueryResult struct {
	Actors []models.Actor
	Total  int64
}

func AllActors() (*[]models.Actor, error) {
	cAndT, err := db.GetCollectionAndContext(db.ActorCollection)
	if err != nil {
		return nil, err
	}

	cur, err := cAndT.Col.Find(cAndT.Ctx, bson.M{})
	if err != nil {
		logger.Error("AllActors::failed with error:", err)
		return nil, err
	}

	acts := make([]models.Actor, 0)

	defer cur.Close(cAndT.Ctx)

	for cur.Next(cAndT.Ctx) {
		var a models.Actor
		err := cur.Decode(&a)
		if err != nil {
			logger.Warn("AllActors::failed to decode actor with error:", err)
			continue
		}

		acts = append(acts, a)
	}

	return &acts, nil
}

func ActorsAtPage(page int) (PaginatedQueryResult, error) {
	cAndT, err := db.GetCollectionAndContext(db.ActorCollection)
	if err != nil {
		return PaginatedQueryResult{}, err
	}

	opts := options.Find()
	opts.SetSkip(int64(page * MaxActorsToReturn))
	opts.SetLimit(MaxActorsToReturn)

	cur, err := cAndT.Col.Find(cAndT.Ctx, bson.M{}, opts)
	if err != nil {
		logger.Error("ActorsAtPage::failed with error:", err)
		return PaginatedQueryResult{}, err
	}

	acts := make([]models.Actor, 0)

	defer cur.Close(cAndT.Ctx)

	for cur.Next(cAndT.Ctx) {
		var a models.Actor
		err := cur.Decode(&a)
		if err != nil {
			logger.Error("ActorsAtPage::failed to decode actor with error:", err)
			continue
		}

		acts = append(acts, a)
	}

	// FIXME: cache count
	total, err := cAndT.Col.CountDocuments(cAndT.Ctx, bson.M{})
	if err != nil {
		logger.Error("ActorsAtPage::failed to count actors:", err)
		return PaginatedQueryResult{Actors: acts}, err
	}

	return PaginatedQueryResult{Actors: acts, Total: total}, nil
}

func AllActorsWithMovies() (*[]models.Actor, error) {
	cAndT, err := db.GetCollectionAndContext(db.ActorCollection)
	if err != nil {
		return nil, err
	}

	qry := bson.D{
		{
			Key:   "$where",
			Value: "this.movieIds.length>0",
		},
	}

	cur, err := cAndT.Col.Find(cAndT.Ctx, qry)
	if err != nil {
		logger.Error("AllActorsWithMovies::failed with error:", err)
		return nil, err
	}

	acts := make([]models.Actor, 0)

	defer cur.Close(cAndT.Ctx)

	for cur.Next(cAndT.Ctx) {
		var a models.Actor

		err := cur.Decode(&a)
		if err != nil {
			logger.Warn("AllActorsWithMovies failed to decode actor with error:", err)
			continue
		}

		acts = append(acts, a)
	}

	return &acts, nil
}

func NewActor(firstName string, middleName string, lastName string) models.Actor {
	a := models.Actor{
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
	}

	a.Identifier = a.GetIdentifier()
	a.MovieIds = make([]primitive.ObjectID, 0)
	a.Updated = time.Now()

	return a
}

func GetActorById(id primitive.ObjectID) (*models.Actor, error) {
	var a models.Actor

	res, err := models.FindById(id, db.ActorCollection)
	if err != nil {
		logger.Error("GetActorById::Failed to fetch model with id:", id, "error:", err)
		return nil, err
	}

	err = res.Decode(&a)
	if err != nil {
		logger.Error("GetActorById::Failed to decode model with id:", id, "error:", err)
		return nil, err
	}

	return &a, nil
}

func GetActorByIdentifier(id string) (*models.Actor, error) {
	var a models.Actor
	res, err := models.FindByIdentifier(id, db.MovieCollection)
	if err != nil {
		logger.Error("GetActorByIdentifier::Failed to fetch model with identifier:", id, "error:", err)
		return nil, err
	}

	err = res.Decode(&a)
	if err != nil {
		logger.Error("GetActorByIdentifier::Failed to decode model with identifier:", id, "error:", err)
		return nil, err
	}

	return &a, nil
}

func GetActorByIdString(actorId string) (*models.Actor, error) {
	id, err := primitive.ObjectIDFromHex(actorId)
	if err != nil {
		logger.Warn("GetActorByIdString could not make object id from:", actorId)
		return nil, err
	}

	return GetActorById(id)
}

func GetActorsForIds(ids []primitive.ObjectID) (*[]models.Actor, error) {
	cAndT, err := db.GetCollectionAndContext(db.ActorCollection)
	if err != nil {
		return nil, err
	}

	filter := bson.D{
		{
			Key: "_id",
			Value: bson.D{
				{
					Key:   "$in",
					Value: ids,
				},
			},
		},
	}

	cur, err := cAndT.Col.Find(cAndT.Ctx, filter)
	if err != nil {
		logger.Error("GetActorsForIds::Failed to make query for actor ids:", ids, err)
		return nil, err
	}

	defer cur.Close(cAndT.Ctx)

	acts := make([]models.Actor, 0)

	for cur.Next(cAndT.Ctx) {
		var a models.Actor

		err := cur.Decode(&a)
		if err != nil {
			logger.Warn("GetActorsForIds::Failed to decode actor with error:", err)
			continue
		}

		acts = append(acts, a)
	}

	return &acts, nil
}

func GetActorsForInput(input string, withMovs bool, page int) (PaginatedQueryResult, error) {
	if input == "" {
		return PaginatedQueryResult{}, nil
	}

	var filter bson.D

	nms := strings.Split(input, " ")

	switch len(nms) {
	case 1:
		filter = bson.D{
			{
				Key: "$or",
				Value: bson.A{
					bson.D{
						{
							Key: "firstName",
							Value: bson.D{
								{
									"$regex",
									primitive.Regex{
										Pattern: fmt.Sprintf("^%s", nms[0]),
										Options: "i",
									},
								},
							},
						},
					},
					bson.D{
						{
							Key: "middleName",
							Value: bson.D{
								{
									"$regex",
									primitive.Regex{
										Pattern: fmt.Sprintf("^%s", nms[0]),
										Options: "i",
									},
								},
							},
						},
					},
					bson.D{
						{
							Key: "lastName",
							Value: bson.D{
								{
									"$regex",
									primitive.Regex{
										Pattern: fmt.Sprintf("^%s", nms[0]),
										Options: "i",
									},
								},
							},
						},
					},
				},
			},
		}
	case 2:
		filter = bson.D{
			{
				Key: "firstName",
				Value: bson.D{
					{
						"$regex",
						primitive.Regex{
							Pattern: fmt.Sprintf("^%s", nms[0]),
							Options: "i",
						},
					},
				},
			},
			{
				Key: "$or",
				Value: bson.A{
					bson.D{
						{
							Key: "middleName",
							Value: bson.D{
								{
									"$regex",
									primitive.Regex{
										Pattern: fmt.Sprintf("^%s", nms[1]),
										Options: "i",
									},
								},
							},
						},
					},
					bson.D{
						{
							Key: "lastName",
							Value: bson.D{
								{
									"$regex",
									primitive.Regex{
										Pattern: fmt.Sprintf("^%s", nms[1]),
										Options: "i",
									},
								},
							},
						},
					},
				},
			},
		}
	default:
		filter = bson.D{
			{
				Key: "firstName",
				Value: bson.D{
					{
						"$regex",
						primitive.Regex{
							Pattern: fmt.Sprintf("^%s", nms[0]),
							Options: "i",
						},
					},
				},
			},
			{
				Key: "middleName",
				Value: bson.D{
					{
						"$regex",
						primitive.Regex{
							Pattern: fmt.Sprintf("^%s", nms[1]),
							Options: "i",
						},
					},
				},
			},
			{
				Key: "lastName",
				Value: bson.D{
					{
						"$regex",
						primitive.Regex{
							Pattern: fmt.Sprintf("^%s", nms[2]),
							Options: "i",
						},
					},
				},
			},
		}
	}

	if withMovs {
		filter = append(filter, bson.E{
			Key:   "$where",
			Value: "this.movieIds.length>0",
		})
	}

	cAndT, err := db.GetCollectionAndContext(db.ActorCollection)
	if err != nil {
		return PaginatedQueryResult{}, err
	}

	opts := options.Find()
	opts.SetSort(bson.D{
		{"fistName", 1},
	})
	opts.SetSkip(int64(page * MaxActorsToReturn))
	opts.SetLimit(MaxActorsToReturn)

	cur, err := cAndT.Col.Find(cAndT.Ctx, filter, opts)
	if err != nil {
		logger.Error("GetActorsForInput::Failed with error:", err)
		return PaginatedQueryResult{}, err
	}
	defer cur.Close(cAndT.Ctx)

	acts := make([]models.Actor, 0)
	for cur.Next(cAndT.Ctx) {
		var a models.Actor
		err := cur.Decode(&a)
		if err != nil {
			logger.Warn("GetActorsForInput::Failed to decode actor with error:", err)
			continue
		}
		acts = append(acts, a)
	}

	sort.Slice(acts, func(i int, j int) bool {
		return acts[i].FullName() < acts[j].FullName()
	})

	total, err := cAndT.Col.CountDocuments(cAndT.Ctx, filter)
	if err != nil {
		logger.Error("ActorsAtPage::failed to count actors:", err)
		return PaginatedQueryResult{Actors: acts}, err
	}

	return PaginatedQueryResult{Actors: acts, Total: total}, nil
}

func GetActorsForInputSimple(input string, page int, withMovs bool) (PaginatedQueryResult, error) {
	acts := make([]models.Actor, 0)
	if input == "" {
		return PaginatedQueryResult{}, nil
	}

	var filter bson.D

	nms := strings.Split(input, " ")

	switch len(nms) {
	case 1:
		filter = bson.D{
			{
				Key: "firstName",
				Value: bson.D{
					{
						"$regex",
						primitive.Regex{
							Pattern: fmt.Sprintf("^%s", nms[0]),
							Options: "i",
						},
					},
				},
			},
		}
	case 2:
		filter = bson.D{
			{
				Key: "firstName",
				Value: bson.D{
					{
						"$regex",
						primitive.Regex{
							Pattern: fmt.Sprintf("^%s", nms[0]),
							Options: "i",
						},
					},
				},
			},
			{
				Key: "$or",
				Value: bson.A{
					bson.D{
						{
							Key: "middleName",
							Value: bson.D{
								{
									"$regex",
									primitive.Regex{
										Pattern: fmt.Sprintf("^%s", nms[1]),
										Options: "i",
									},
								},
							},
						},
					},
					bson.D{
						{
							Key: "lastName",
							Value: bson.D{
								{
									"$regex",
									primitive.Regex{
										Pattern: fmt.Sprintf("^%s", nms[1]),
										Options: "i",
									},
								},
							},
						},
					},
				},
			},
		}
	default:
		filter = bson.D{
			{
				Key: "firstName",
				Value: bson.D{
					{
						"$regex",
						primitive.Regex{
							Pattern: fmt.Sprintf("^%s", nms[0]),
							Options: "i",
						},
					},
				},
			},
			{
				Key: "middleName",
				Value: bson.D{
					{
						"$regex",
						primitive.Regex{
							Pattern: fmt.Sprintf("^%s", nms[1]),
							Options: "i",
						},
					},
				},
			},
			{
				Key: "lastName",
				Value: bson.D{
					{
						"$regex",
						primitive.Regex{
							Pattern: fmt.Sprintf("^%s", nms[2]),
							Options: "i",
						},
					},
				},
			},
		}
	}

	if withMovs {
		filter = append(filter, bson.E{
			Key:   "$where",
			Value: "this.movieIds.length>0",
		})
	}

	opts := options.Find()
	opts.SetSkip(int64(page * MaxActorsToReturn))
	opts.SetLimit(MaxActorsToReturn)
	opts.SetSort(bson.D{
		{"firstName", 1},
		{"middleName", 1},
		{"lastName", 1},
	})

	cAndT, err := db.GetCollectionAndContext(db.ActorCollection)
	if err != nil {
		return PaginatedQueryResult{}, err
	}

	logger.Log("GetActorsForInputSimple->starting search at:", page*MaxActorsToReturn)

	cur, err := cAndT.Col.Find(cAndT.Ctx, filter, opts)
	if err != nil {
		logger.Error("GetActorsForInputSimple->Failed with error:", err)
		return PaginatedQueryResult{}, err
	}
	defer cur.Close(cAndT.Ctx)

	for cur.Next(cAndT.Ctx) {
		var a models.Actor
		err := cur.Decode(&a)
		if err != nil {
			logger.Warn("GetActorsForInputSimple->Failed to decode actor with error:", err)
			continue
		}

		acts = append(acts, a)
	}

	sort.Slice(acts, func(i int, j int) bool {
		return acts[i].FullName() < acts[j].FullName()
	})

	total, err := cAndT.Col.CountDocuments(cAndT.Ctx, filter)
	if err != nil {
		logger.Error("ActorsAtPageGetActorsForInputSimple->failed to count actors:", err)
		return PaginatedQueryResult{Actors: acts}, err
	}

	return PaginatedQueryResult{Actors: acts, Total: total}, nil
}

func ActorsByDate() (*[]models.Actor, error) {
	acts := make([]models.Actor, 0)

	var filter = bson.D{}

	cAndT, err := db.GetCollectionAndContext(db.ActorCollection)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.SetSort(bson.D{
		{"updated", -1},
	})
	opts.SetLimit(300)

	cur, err := cAndT.Col.Find(cAndT.Ctx, filter, opts)
	if err != nil {
		logger.Error("ActorsByDate::Failed with error:", err)
		return nil, err
	}

	defer cur.Close(cAndT.Ctx)

	for cur.Next(cAndT.Ctx) {
		var a models.Actor

		err := cur.Decode(&a)
		if err != nil {
			logger.Warn("ActorsByDate::Failed to decode movie with error:", err)
			continue
		}

		if len(a.MovieIds) > 0 {
			mvs, err := MoviesForIds(a.MovieIds)
			if err == nil {
				a.Movies = mvs
			}
		}

		acts = append(acts, a)
	}

	return &acts, nil
}
