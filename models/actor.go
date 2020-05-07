package models

import (
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type NamePart string

const (
	FirstName  NamePart = "firstName"
	MiddleName NamePart = "middleName"
	LastName   NamePart = "lastName"
)

type Actor struct {
	Id         primitive.ObjectID  `json:"id,omitempty" bson:"_id"`
	FirstName  string `json:"firstName" bson:"firstName"`
	MiddleName string `json:"middleName" bson:"middleName"`
	LastName   string `json:"lastName" bson:"lastName"`
	Model
}

func (a *Actor) GetId() primitive.ObjectID {
	return a.Id
}

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

func (a *Actor) GetIdentifier() string {
	a.Identifier = a.FullName()
	return a.Identifier
}

func (a *Actor) IsIn(m *SysMovie) bool {
	n := strings.ToLower(m.Name())
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
