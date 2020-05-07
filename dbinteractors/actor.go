package dbinteractors

import (
	"github.com/greenac/artemis/db"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"gopkg.in/mgo.v2/bson"
)

func AllActors() (*[]models.Actor, error) {
	cAndT, err := db.GetCollectionAndContext(db.ActorCollection)
	if err != nil { return nil, err }

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

	a.ColType = db.ActorCollection
	a.Identifier = a.GetIdentifier()

	return a
}

//func ActorForId(id string) (error) {
//	a, err := models.FindByIdentifier(id, db.ActorCollection)
//	if err != nil { return nil, err }
//
//
//	return a
//
//}
