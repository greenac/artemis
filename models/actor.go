package models

import (
	"crypto/md5"
	"fmt"
	"github.com/fatih/structs"
	"github.com/greenac/artemis/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"time"
)

type NamePart string

const (
	FirstName  NamePart = "firstName"
	MiddleName NamePart = "middleName"
	LastName   NamePart = "lastName"
)

type Actor struct {
	Id         primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Identifier string               `json:"identifier" bson:"identifier"`
	FirstName  string               `json:"firstName" bson:"firstName"`
	MiddleName string               `json:"middleName" bson:"middleName"`
	LastName   string               `json:"lastName" bson:"lastName"`
	MovieIds   []primitive.ObjectID `json:"movieIds" bson:"movieIds"`
	Updated    time.Time            `json:"updated" bson:"updated"`
}

// Model Interface methods
func (a *Actor) GetId() primitive.ObjectID {
	return a.Id
}

func (a *Actor) IdFilter() bson.M {
	return bson.M{"_id": a.Id}
}

func (a *Actor) IdentifierFilter() bson.M {
	return bson.M{"identifier": a.Identifier}
}

func (a *Actor) GetIdentifier() string {
	if a.Identifier == "" {
		a.SetIdentifier()
	}

	return a.Identifier
}

func (a *Actor) SetIdentifier() string {
	a.Identifier = fmt.Sprintf("%x", md5.Sum([]byte(a.FullName())))

	return a.Identifier
}

func (a *Actor) GetCollectionType() db.CollectionType {
	return db.ActorCollection
}

func (a *Actor) Save() error {
	a.Updated = time.Now()
	return Save(a)
}

func (a *Actor) Upsert() (*primitive.ObjectID, error) {
	return Upsert(a)
}

func (a *Actor) Create() (*primitive.ObjectID, error) {
	id, err := Create(a)
	if err != nil {
		return nil, err
	}

	a.Id = *id
	a.Updated = time.Now()

	return id, nil
}

// Class methods
func (a *Actor) FullName() string {
	name := ""
	if a.FirstName != "" {
		name += a.FirstName
	}

	if a.MiddleName != "" {
		if name == "" {
			name = a.MiddleName
		} else {
			name += "_" + a.MiddleName
		}
	}

	if a.LastName != "" {
		if name == "" {
			name = a.LastName
		} else {
			name += "_" + a.LastName
		}
	}

	return strings.ToLower(name)
}

func (a *Actor) IsIn(name string) bool {
	n := strings.ToLower(name)
	isIn := false
	if a.FirstName != "" {
		isIn = strings.Contains(n, a.FirstName)
	}

	if isIn && a.MiddleName != "" {
		isIn = strings.Contains(n, a.MiddleName)
	}

	if isIn && a.LastName != "" {
		isIn = strings.Contains(n, a.LastName)
	}

	return isIn
}

func (a *Actor) AsMap() map[string]interface{} {
	return structs.Map(a)
}

func (a *Actor) FormatName(name string) string {
	return strings.ToLower(strings.Replace(name, " ", "_", -1))
}

func (a *Actor) IsMatch(name string) bool {
	fmtName := a.FormatName(name)
	return strings.Contains(strings.ToLower(a.FullName()), fmtName)
}

func (a *Actor) MatchPartial(pName string, np NamePart) bool {
	match := true
	name := ""

	switch np {
	case FirstName:
		match = len(pName) <= len(a.FirstName)
		name = strings.ToLower(a.FirstName)
	case MiddleName:
		match = len(pName) <= len(a.MiddleName)
		name = strings.ToLower(a.MiddleName)
	case LastName:
		match = len(pName) <= len(a.LastName)
		name = strings.ToLower(a.LastName)
	}

	if !match {
		return match
	}

	for i, c := range name {
		if byte(c) != name[i] {
			match = false
			break
		}
	}

	return match
}

func (a *Actor) MatchWhole(frag string) bool {
	n := a.FullName()
	if len(frag) > len(n) {
		return false
	}

	match := true
	for i, c := range frag {
		if byte(c) != n[i] {
			match = false
			break
		}
	}

	return match
}

func (a *Actor) HasFirstMiddleLastName() bool {
	return a.FirstName != "" && a.MiddleName != "" && a.LastName != ""
}

func (a *Actor) HasFirstLastName() bool {
	return a.FirstName != "" && a.LastName != ""
}

func (a *Actor) HasFirstName() bool {
	return a.FirstName != ""
}

func (a *Actor) FullNameNoUnderscores() string {
	return strings.ReplaceAll(a.FullName(), "_", "")
}

func (a *Actor) AddMovie(id primitive.ObjectID) bool {
	for _, mid := range a.MovieIds {
		if mid == id {
			return false
		}
	}

	a.MovieIds = append(a.MovieIds, id)

	return true
}

func (a *Actor) RemoveMovie(id primitive.ObjectID) bool {
	t := -1
	for i, mid := range a.MovieIds {
		if mid == id {
			t = i
			break
		}
	}

	if t == -1 {
		return false
	}

	a.MovieIds = append(a.MovieIds[:t], a.MovieIds[t+1:]...)

	return true
}
