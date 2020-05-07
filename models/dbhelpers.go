package models

import (
	"github.com/greenac/artemis/db"
	"github.com/greenac/artemis/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


func Create(m ModelDef) (interface{}, error) {
	cAndT, err := db.GetCollectionAndContext(m.GetCollectionType())
	if err != nil {
		return nil, err
	}

	m.GetIdentifier()

	res, err := cAndT.Col.InsertOne(cAndT.Ctx, m)

	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil
}

func CreateIfDoesNotExist(m ModelDef) (bool, interface{}, error) {
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
) (bson.M, error) {
	cAndT, err := db.GetCollectionAndContext(ct)
	if err != nil {
		return  nil, err
	}

	var m bson.M
	err = cAndT.Col.FindOne(cAndT.Ctx, bson.D{{"identifier", identifier}}).Decode(&m)
	if err != nil {
		logger.Error("FindByIdentifier::could not retrieve identifier:", identifier, "from collection:", ct)
	}

	return m,  err
}

func FindById(
	id string,
	ct db.CollectionType,
) (bson.M, error) {
	cAndT, err := db.GetCollectionAndContext(ct)
	if err != nil {
		return  nil, err
	}

	var m bson.M
	err = cAndT.Col.FindOne(cAndT.Ctx, bson.D{{"_id", id}}).Decode(&m)
	if err != nil {
		logger.Error("FindById::could not retrieve _id:", id, "from collection:", ct)
	}

	return m,  err
}

func Update(m ModelDef,) (*mongo.UpdateResult, error)  {
	cAndT, err := db.GetCollectionAndContext(m.GetCollectionType())
	if err != nil {
		return  nil, err
	}

	res, err := cAndT.Col.UpdateOne(cAndT.Ctx, m.IdentifierFilter(), m)
	if err != nil {
		logger.Error("Update::failed to update", m, "with error:", err)
	}

	return res, err
}

func Upsert(m ModelDef) (interface{}, error) {
	id := m.GetId()
	if id.String() == "" {
		mm, err := FindByIdentifier(m.GetIdentifier(), m.GetCollectionType())
		if err != nil {
			return nil, err
		}

		if mm == nil {
			return Create(m)
		}

		return Update(m)
	}

	mb, err := FindById(id.String(), db.MovieCollection)
	if err != nil {
		return nil, err
	}

	if mb == nil {
		return Create(m)
	}

	res, err := Update(m)
	if err != nil {
		return nil, err
	}

	return res.UpsertedID, nil
}
