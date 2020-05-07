package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShowType string

const (
	ShowTypeMovie ShowType = "movie"
	ShowTypeSeries ShowType = "series"
)

type Movie struct {
	Id         primitive.ObjectID  `json:"id,omitempty" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Path  string `json:"path" bson:"path"`
	Studio string  `json:"studio" bson:"studio"`
	Series string  `json:"series" bson:"series"`
	Type ShowType `json:"type" bson:"type"`
	MetaData string `json:"meta" bson:"meta"`
	RepeatNum int `json:"repeatNum" bson:"repeatNum"`
	ActorIds []string `json:"actorIds" bson:"actorIds"`
	Actors *[]*Actor `json:"actors" bson:"-"`
	Model
}
