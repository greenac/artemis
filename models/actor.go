package models

import (
	"github.com/fatih/structs"
	"strings"
)

type NamePart string

const (
	FirstName  NamePart = "firstName"
	MiddleName NamePart = "middleName"
	LastName   NamePart = "lastName"
)

type Actor struct {
	FirstName  *string
	LastName   *string
	MiddleName *string
}

func (a *Actor) GetFirstName() string {
	if a.FirstName == nil {
		return ""
	}

	return *a.FirstName
}

func (a *Actor) GetMiddleName() string {
	if a.MiddleName == nil {
		return ""
	}

	return *a.MiddleName
}

func (a *Actor) GetLastName() string {
	if a.LastName == nil {
		return ""
	}

	return *a.LastName
}

func (a *Actor) FullName() string {
	name := ""
	if a.FirstName != nil {
		name += *a.FirstName
	}

	if a.MiddleName != nil {
		if name == "" {
			name = *a.MiddleName
		} else {
			name += "_" + *a.MiddleName
		}
	}

	if a.LastName != nil {
		if name == "" {
			name = *a.LastName
		} else {
			name += "_" + *a.LastName
		}
	}

	return strings.ToLower(name)
}

func (a *Actor) IsIn(m *Movie) bool {
	n := strings.ToLower(*m.Name())
	isIn := false
	if a.FirstName != nil {
		isIn = strings.Contains(n, *a.FirstName)
	}

	if isIn && a.MiddleName != nil {
		isIn = strings.Contains(n, *a.MiddleName)
	}

	if isIn && a.LastName != nil {
		isIn = strings.Contains(n, *a.LastName)
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
		match = len(pName) <= len(*a.FirstName)
		name = strings.ToLower(*a.FirstName)
	case MiddleName:
		match = len(pName) <= len(*a.MiddleName)
		name = strings.ToLower(*a.MiddleName)
	case LastName:
		match = len(pName) <= len(*a.LastName)
		name = strings.ToLower(*a.LastName)
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
	return a.FirstName != nil && a.MiddleName != nil && a.LastName != nil
}

func (a *Actor) HasFirstLastName() bool {
	return a.FirstName != nil && a.LastName != nil
}

func (a *Actor) HasFirstName() bool {
	return a.FirstName != nil
}

func (a *Actor) FullNameNoUnderscores() string {
	return strings.ReplaceAll(a.FullName(), "_", "")
}
