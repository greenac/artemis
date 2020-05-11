package dbinteractors

import (
	"github.com/greenac/artemis/db"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

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
		//var r bson.M
		err := cur.Decode(&a)
		if err != nil {
			logger.Warn("AllActors failed to decode actor with error:", err)
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
