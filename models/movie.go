package models

import (
	"crypto/md5"
	"fmt"
	"github.com/greenac/artemis/db"
	"github.com/greenac/artemis/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShowType string

const (
	ShowTypeMovie  ShowType = "movie"
	ShowTypeSeries ShowType = "series"
)

func MovieIdentifier(path string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(path)))
}

type Movie struct {
	Id         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Identifier string               `json:"identifier" bson:"identifier"`
	Name       string               `json:"name" bson:"name"`
	Path       string               `json:"path" bson:"path"`
	Studio     string               `json:"studio" bson:"studio"`
	Series     string               `json:"series" bson:"series"`
	Type       ShowType             `json:"type" bson:"type"`
	MetaData   string               `json:"meta" bson:"meta"`
	RepeatNum  int                  `json:"repeatNum" bson:"repeatNum"`
	ActorIds   []primitive.ObjectID `json:"actorIds" bson:"actorIds"`
	Actors     *[]*Actor            `json:"actors" bson:"-"`
}

func (m *Movie) GetId() primitive.ObjectID {
	return m.Id
}

func (m *Movie) SetIdentifier() string {
	m.Identifier = MovieIdentifier(m.Path)

	return m.Identifier
}

func (m *Movie) GetIdentifier() string {
	if m.Identifier == "" {
		m.SetIdentifier()
	}

	return m.Identifier
}

func (m *Movie) IdFilter() bson.M {
	return bson.M{"_id": m.Id}
}

func (m *Movie) IdentifierFilter() bson.M {
	return bson.M{"identifier": m.Identifier}
}

func (m *Movie) GetCollectionType() db.CollectionType {
	return db.MovieCollection
}

func (m *Movie) Create() (*primitive.ObjectID, error) {
	id, err := Create(m)
	if err != nil {
		return nil, err
	}

	m.Id = *id

	return id, nil
}

func (m *Movie) Upsert() (*primitive.ObjectID, error) {
	return Upsert(m)
}

func (m *Movie) Save() error {
	return Save(m)
}

func (m *Movie) AddActor(actorId primitive.ObjectID) bool {
	for _, aID := range m.ActorIds {
		if actorId == aID {
			return false
		}
	}

	logger.Debug("Movie:", m.Name, "Adding actor:", actorId)
	m.ActorIds = append(m.ActorIds, actorId)

	return true
}
