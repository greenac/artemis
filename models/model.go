package models

import (
	"github.com/greenac/artemis/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IdentifierFilter map[string]string

type ModelDef interface {
	GetId() primitive.ObjectID
	GetIdentifier() string
	IdentifierFilter() IdentifierFilter
	GetCollectionType() db.CollectionType
}

type Model struct {
	Identifier string            `json:"identifier" bson:"identifier"`
	ColType    db.CollectionType `json:"collectionType" bson:"collectionType"`
}

func (m *Model) GetId() primitive.ObjectID {
	// This method should return the struct that has nested `Model`
	// when it implements `GetId`
	id, _ := primitive.ObjectIDFromHex("0")
	return id
}

func (m *Model) GetIdentifier() string {
	return m.Identifier
}

func (m *Model) IdentifierFilter() IdentifierFilter {
	f := IdentifierFilter{"identifier": m.GetIdentifier()}
	return f
}

func (m *Model) GetCollectionType() db.CollectionType {
	return m.ColType
}

func (m *Model) Save() error {
	cAndT, err := db.GetCollectionAndContext(m.ColType)
	if err != nil {
		return err
	}

	cAndT.Col.FindOneAndUpdate(cAndT.Ctx, m.GetIdentifier(), m)

	return nil
}

func (m *Model) Upsert() (interface{}, error) {
	return Upsert(m)
}