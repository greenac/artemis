package movie

import (
	"github.com/fatih/structs"
	"strings"
)

type Actor struct {
	FirstName  *string
	LastName   *string
	MiddleName *string
	Movies     map[string]*Movie
}

func (a *Actor) setup() {
	if a.Movies == nil {
		a.Movies = make(map[string]*Movie, 0)
	}
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

	return name
}

func (a *Actor) AddMovie(m *Movie) {
	a.setup()
	a.Movies[*m.Name()] = m
}

func (a *Actor) AddFiles(mvs []*Movie) {
	for _, m := range mvs {
		a.AddMovie(m)
	}
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
