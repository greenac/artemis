package models

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/utils"
	"regexp"
	"strings"
)



type Movie struct {
	File
	Actors []*Actor
}

func (m *Movie) AddActor(a Actor) {
	m.Actors = append(m.Actors, &a)
}

func (m *Movie) AddActorNames() {
	m.GetNewName()

	for _, a := range m.Actors {
		m.UpdateNewName(a)
	}
}

func (m *Movie) FormattedName() (formattedName string, error error) {
	name := m.removeRepeats()
	parts := strings.Split(name, ".")
	ext := parts[len(parts)-1]
	name = strings.Join(parts[:len(parts)-1], ".")

	re, err := regexp.Compile(`[-\s\t!@#$%^&*()[\]<>,.?~]`)
	if err != nil {
		logger.Error("Movie::FormattedName failed to compile regex with error:", err)
		return "", err
	}

	rs := re.ReplaceAll([]byte(name), []byte{'_'})

	return string(*(m.cleanUnderscores(&rs))) + "." + ext, nil
}

func (m *Movie) cleanUnderscores(name *[]byte) *[]byte {
	cln := make([]byte, 0)
	fndUn := false
	for _, c := range *name {
		if c == '_' {
			if !fndUn {
				cln = append(cln, c)
				fndUn = true
			}
		} else {
			fndUn = false
			cln = append(cln, c)
		}
	}

	if cln[len(cln)-1] == '_' {
		var cut int
		for i := len(cln) - 1; i >= 0; i -= 1 {
			if cln[i] == '_' {
				cut = i
			} else {
				break
			}
		}

		cln = cln[:cut]
	}

	if cln[0] == '_' {
		var cut int
		for i := 0; i < len(cln); i += 1 {
			if cln[i] == '_' {
				cut = i
			} else {
				break
			}
		}

		cln = cln[cut+1:]
	}

	return &cln
}

func (m *Movie) AddName(a *Actor) string {
	newName := m.NewName
	newName = utils.AddPrecedingUnderscore(*a.FirstName, newName)
	newName = utils.AddFollowingUnderscore(*a.FirstName, newName)

	if a.MiddleName != nil {
		newName = utils.AddMiddleUnderscore(*a.FirstName, *a.MiddleName, newName)
		newName = utils.AddFollowingUnderscore(*a.MiddleName, newName)
		if a.LastName != nil {
			newName = utils.AddMiddleUnderscore(*a.MiddleName, *a.LastName, newName)
			newName = utils.AddFollowingUnderscore(*a.LastName, newName)
		}
	} else if a.LastName != nil {
		newName = utils.AddMiddleUnderscore(*a.FirstName, *a.LastName, newName)
		newName = utils.AddFollowingUnderscore(*a.LastName, newName)
	}

	return newName
}

func (m *Movie) UpdateNewName(a *Actor) {
	m.NewName = m.AddName(a)
}

func (m *Movie) GetNewName() string {
	if m.NewName == "" {
		nn, err := m.FormattedName()
		if err != nil {
			return ""
		}

		m.NewName = nn
	}

	return m.NewName
}

func (m *Movie) removeRepeats() string {
	nn := make([]byte, len(*m.Name()))
	copy(nn, *m.Name())
	name := strings.ToLower(string(nn))
	if strings.Contains(name, "scene_") {
		re, err := regexp.Compile(`\\(.+?\\)`)
		if err != nil {
			logger.Warn("Movie::RemoveRepeats failed to compile regex with error:", err)
		}

		name = re.ReplaceAllString(name, "")
	}

	return strings.ReplaceAll(name, " copy", "")
}
