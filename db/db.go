package db

import (
	"context"
	"errors"
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type CollectionType string

const (
	ActorCollection CollectionType = "actor"
	MovieCollection CollectionType = "movie"
)

const contextTime = 10 * time.Second

type database struct {
	Config *models.MongoConfig
	client *mongo.Client
	ctx    *context.Context
}

func (db *database) getClient() (*mongo.Client, error) {
	if db.client == nil {
		ctx := db.getContext()
		cl, err := mongo.Connect(*ctx, options.Client().ApplyURI(db.Config.Url))
		if err != nil {
			logger.Error("Database::setClient failed to connect to mongo on:", db.Config.Url, err)
			return nil, err
		}

		db.client = cl
	}

	return db.client, nil
}

func (db *database) getContext() *context.Context {
	if db.ctx == nil {
		ctx, _ := context.WithTimeout(context.Background(), contextTime)
		db.ctx = &ctx
	}

	return db.ctx
}

func (db *database) getCollection(col CollectionType) *mongo.Collection {
	var c string
	switch col {
	case ActorCollection:
		c = db.Config.Collections.Actors
	case MovieCollection:
		c = db.Config.Collections.Movies
	}

	clt, _ := db.getClient()
	return clt.Database(db.Config.Database).Collection(c)
}

var db *database = nil

func SetupMongo(config *models.MongoConfig) {
	if db == nil {
		d := database{Config: config}
		db = &d
	}
}

func GetCollection(col CollectionType) (*mongo.Collection, error) {
	if db == nil {
		return nil, artemiserror.New(artemiserror.MongoNotSetUp)
	}

	return db.getCollection(col), nil
}
