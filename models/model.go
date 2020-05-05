package models

import (
	"github.com/greenac/artemis/db"
)

type IdentifierFilter map[string]string

type ModelDef interface {
	GetId() string
	GetIdentifier() string
	IdentifierFilter() IdentifierFilter
	GetCollectionType() db.CollectionType
}

type Model struct {
	Id         string            `json:"id" bson:"id"`
	Identifier string            `json:"identifier" bson:"identifier"`
	ColType    db.CollectionType `json:"collectionType" bson:"collectionType"`
}

func (m *Model) GetId() string {
	return m.Id
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