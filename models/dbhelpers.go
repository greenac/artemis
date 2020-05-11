package models

import (
	"github.com/greenac/artemis/db"
	"github.com/greenac/artemis/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create(m ModelDef) (*primitive.ObjectID, error) {
	cAndT, err := db.GetCollectionAndContext(m.GetCollectionType())
	if err != nil {
		logger.Debug("Failed to get collection and context with error:", err)
		return nil, err
	}

	m.SetIdentifier()

	res, err := cAndT.Col.InsertOne(cAndT.Ctx, m)
	if err != nil {
		logger.Error("Create::Failed to insert model:", m, "with error", err)
		return nil, err
	}

	id := res.InsertedID.(primitive.ObjectID)

	return &id, nil
}

func CreateIfDoesNotExist(m ModelDef) (bool, *primitive.ObjectID, error) {
	cAndT, err := db.GetCollectionAndContext(m.GetCollectionType())
	if err != nil {
		return false, nil, err
	}

	f := m.IdentifierFilter()
	a := cAndT.Col.FindOne(cAndT.Ctx, f)
	if a == nil {
		id, err := Create(m)
		if err != nil {
			return false, nil, err
		}

		return true, id, nil
	}

	return false, nil, nil
}

func FindByIdentifier(
	identifier string,
	ct db.CollectionType,
) (*mongo.SingleResult, error) {
	cAndT, err := db.GetCollectionAndContext(ct)
	if err != nil {
		return nil, err
	}

	return cAndT.Col.FindOne(
		cAndT.Ctx,
		bson.D{{"identifier", identifier}},
	), nil
}

func FindById(
	id primitive.ObjectID,
	ct db.CollectionType,
) (*mongo.SingleResult, error) {
	cAndT, err := db.GetCollectionAndContext(ct)
	if err != nil {
		return nil, err
	}

	return cAndT.Col.FindOne(cAndT.Ctx, bson.D{{"_id", id}}), nil
}

func Update(m ModelDef) (*primitive.ObjectID, error) {
	cAndT, err := db.GetCollectionAndContext(m.GetCollectionType())
	if err != nil {
		return nil, err
	}

	res, err := cAndT.Col.UpdateOne(cAndT.Ctx, m.IdentifierFilter(), m)
	if err != nil {
		logger.Error("Update::failed to update", m, "with error:", err)
	}

	id := res.UpsertedID.(primitive.ObjectID)

	return &id, err
}

func Upsert(m ModelDef) (*primitive.ObjectID, error) {
	id := m.GetId()
	if id.String() == "" {
		res, err := FindByIdentifier(m.GetIdentifier(), m.GetCollectionType())
		if err != nil {
			return nil, err
		}

		var r interface{}

		err = res.Decode(r)
		if err != nil {
			return Create(m)
		}

		return Update(m)
	}

	res, err := FindById(id, db.MovieCollection)
	if err != nil {
		return nil, err
	}

	var r interface{}

	err = res.Decode(r)
	if err != nil {
		return Create(m)
	}

	return Upsert(m)
}

func Save(m ModelDef) error {
	cAndT, err := db.GetCollectionAndContext(m.GetCollectionType())
	if err != nil {
		return err
	}

	opts := options.FindOneAndUpdate().SetUpsert(true)
	res := cAndT.Col.FindOneAndUpdate(cAndT.Ctx, m.IdFilter(), bson.M{"$set": m}, opts)
	if res.Err() != nil {
		logger.Error("Save::failed to save model def:", m, "error", res.Err())
		return res.Err()
	}

	return nil
}
