package db

import (
	"context"
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/config"
	"github.com/greenac/artemis/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CollectionType string

const (
	ActorCollection CollectionType = "actor"
	MovieCollection CollectionType = "movie"
)

type CollectionAndContext struct {
	Col *mongo.Collection
	Ctx context.Context
}

const contextTime = 10 * time.Second

type database struct {
	Config *config.MongoConfig
	client *mongo.Client
	ctx    *context.Context
}

func (db *database) getClient() (*mongo.Client, error) {
	if db.client == nil {
		ctx := db.getContext()
		cl, err := mongo.Connect(ctx, options.Client().ApplyURI(db.Config.Url))
		if err != nil {
			logger.Error("Database::setClient failed to connect to mongo on:", db.Config.Url, err)
			return nil, err
		}

		db.client = cl
	}

	return db.client, nil
}

func (db *database) getContext() context.Context {
	c, _ := context.WithTimeout(context.Background(), contextTime)
	return c
}

func (db *database) getCollection(ct CollectionType) *mongo.Collection {
	var c string
	switch ct {
	case ActorCollection:
		c = db.Config.Collections.Actors
	case MovieCollection:
		c = db.Config.Collections.Movies
	}

	clt, _ := db.getClient()
	return clt.Database(db.Config.Database).Collection(c)
}

var db *database = nil

func SetupMongo(config *config.MongoConfig) {
	if db == nil {
		d := database{Config: config}
		db = &d
	}
}

func GetCollection(ct CollectionType) (*mongo.Collection, error) {
	if db == nil {
		return nil, artemiserror.New(artemiserror.MongoNotSetUp)
	}

	c := db.getCollection(ct)

	return c, nil
}

func GetCollectionAndContext(ct CollectionType) (*CollectionAndContext, error) {
	col, err := GetCollection(ct)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	return &CollectionAndContext{
		Col: col,
		Ctx: ctx,
	}, nil
}
