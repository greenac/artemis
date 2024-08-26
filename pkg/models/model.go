package models

import (
	"github.com/greenac/artemis/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModelDef interface {
	GetId() primitive.ObjectID
	GetIdentifier() string
	SetIdentifier() string
	IdFilter() bson.M
	IdentifierFilter() bson.M
	GetCollectionType() db.CollectionType
	Save() error
	Upsert() (*primitive.ObjectID, error)
	Create() (*primitive.ObjectID, error)
}
