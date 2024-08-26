package models

import (
	"crypto/md5"
	"fmt"
	"github.com/greenac/artemis/pkg/db"
	"github.com/greenac/artemis/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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
	Id            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Identifier    string               `json:"identifier" bson:"identifier"`
	Name          string               `json:"name" bson:"name"`
	Path          string               `json:"path" bson:"path"`
	SecondaryPath string               `json:"secondaryPath" bson:"secondaryPath"`
	Studio        string               `json:"studio" bson:"studio"`
	Series        string               `json:"series" bson:"series"`
	Type          ShowType             `json:"type" bson:"type"`
	MetaData      string               `json:"meta" bson:"meta"`
	RepeatNum     int                  `json:"repeatNum" bson:"repeatNum"`
	ActorIds      []primitive.ObjectID `json:"actorIds" bson:"actorIds"`
	Actors        *[]Actor             `json:"actors" bson:"-"`
	Updated       time.Time            `json:"updated" bson:"updated"`
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
	m.Updated = time.Now()

	return id, nil
}

func (m *Movie) Upsert() (*primitive.ObjectID, error) {
	return Upsert(m)
}

func (m *Movie) Save() error {
	m.Updated = time.Now()
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

func (m *Movie) HasActor(actorId primitive.ObjectID) bool {
	for _, id := range m.ActorIds {
		if id == actorId {
			return true
		}
	}

	return false
}

func (m *Movie) RemoveActor(actorId primitive.ObjectID) {
	if len(m.ActorIds) == 0 {
		return
	}

	if len(m.ActorIds) == 1 && actorId == m.ActorIds[0] {
		m.ActorIds = []primitive.ObjectID{}
		return
	}

	for i, aID := range m.ActorIds {
		if actorId == aID {
			logger.Debug("Move:RemoveActor::removing actor id:", actorId, "from movie:", m.Name)

			switch i {
			case 0:
				m.ActorIds = m.ActorIds[1:]
			case len(m.ActorIds):
				m.ActorIds = m.ActorIds[:len(m.ActorIds)-1]
			default:
				m.ActorIds = append(m.ActorIds[:i], m.ActorIds[i+1:]...)
			}

			break
		}
	}
}

func (m *Movie) String() string {
	numActs := 0
	if m.Actors != nil {
		numActs = len(*m.Actors)
	}

	return fmt.Sprintf("name: %s, path: %s, secondary path: %s, id: %s, num actors: %d", m.Name, m.Path, m.SecondaryPath, m.Id, numActs)
}
