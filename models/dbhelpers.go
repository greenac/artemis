package models

import (
	"github.com/greenac/artemis/db"
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

	id := m.GetIdentifier()
	f := map[string]string{"filter": id}
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

func NewActor(firstName string, middleName string, lastName string) Actor {
	a := Actor{
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
	}

	a.ColType = db.ActorCollection
	a.Identifier = a.GetIdentifier()

	return a
}
