package models

import (
	"github.com/greenac/artemis/logger"
	"regexp"
	"strings"
)

type Movie struct {
	File
	Actors []*Actor
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

	return  string(*(m.cleanUnderscores(&rs))) + "." + ext, nil
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
	fn := a.FullName()
	if strings.Contains(m.NewName, fn) {
		return m.NewName
	}

	var nn string
	i := strings.Index(m.NewName, *a.FirstName)
	if i == -1 {
		pts := strings.Split(m.NewName, ".")
		if len(pts) != 2 {
			logger.Warn("`Movie::Rename`", m.NewName, "not in proper format")
			return m.NewName
		}

		nm := pts[0]
		nmb := []byte(nm)

		if !strings.Contains(nm, fn) {
			if nmb[len(nmb)-1] == '_' {
				nm += fn
			} else {
				nm += "_" + fn
			}
		}

		nn = nm + "." + pts[1]
	} else {
		nn = strings.ReplaceAll(m.NewName, *a.FirstName, fn)
	}

	logger.Debug("Before Adding name:", a.FullName(), "to movie:", nn)
	nn = m.addFullName(a, nn)
	logger.Debug("After Adding name:", a.FullName(), "to movie:", nn)

	return nn
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

func (m *Movie) AddActor(a *Actor) {
	m.Actors = append(m.Actors, a)
}

func (m *Movie) AddActorNames() {
	m.GetNewName()
	for _, a := range m.Actors {
		m.UpdateNewName(a)
	}
}

func (m *Movie) removeRepeats() string {
	nn := make([]byte, len(*m.Name()))
	copy(nn, *m.Name())
	name := strings.ToLower(string(nn))
	if strings.Contains(name, "scene_", ) {
		re, err := regexp.Compile(`\\(.+?\\)`)
		if err != nil {
			logger.Warn("Movie::RemoveRepeats failed to compile regex with error:", err)
		}

		name = re.ReplaceAllString(name, "")
	}

	return strings.ReplaceAll(name, " copy", "")
}

func (m *Movie) addFullName(a *Actor, newName string) string {
	fn := a.FullName()

	logger.Log("Running with new name:", newName, "for:", fn)

	if strings.Contains(newName, fn) {
		return newName
	}

	if a.MiddleName == nil {
		i := strings.Index(newName, *a.FirstName)
		if i != -1 {
			nrns := []rune(newName)
			if i > 0 {
				if nrns[i-1] != '_' {
					logger.Log("Should insert _ here at", i - 1)
					nrns = append(nrns[:i], append([]rune{'_'}, nrns[i:]...)...)
					newName = string(nrns)
					logger.Log("New name for nrns[i-1] != '_':", newName)
					i += 1
				}
			}

			if i < len(nrns) - 1 {
				t := i + len(*a.FirstName)
				if nrns[t] != '_' {
					if a.LastName == nil {
						nrns = append(nrns[:t], append([]rune{'_'}, nrns[t:]...)...)
						logger.Log("new name:", newName)
					} else {
						li := strings.Index(newName, *a.LastName)
						logger.Log("last name starts at:", li, "t:", t)
						if li == t {
							nrns = append(nrns[:li], append([]rune{'_'}, nrns[li:]...)...)
							li += 1
							logger.Log("new with _ in middle of name is:", string(nrns))
						}

						li += len(*a.LastName)
						logger.Log("target value:", string(nrns[li]))
						if nrns[li] != '.' && nrns[li] != '_' {
							nrns = append(nrns[:li], append([]rune{'_'}, nrns[li:]...)...)
							logger.Log("new at end of last name is:", string(nrns))
						}
					}

					newName = string(nrns)
				}
			}

			newName = strings.ReplaceAll(newName, *a.FirstName, fn)
		}
	} else if strings.Contains(newName, *a.FirstName + "_" + *a.LastName) {
		newName = strings.ReplaceAll(newName, *a.FirstName + "_" + *a.LastName, fn)
	}

	return newName
}
