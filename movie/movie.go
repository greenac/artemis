package movie

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/tools"
	"regexp"
	"strings"
)

type Movie struct {
	tools.File
	Actors []*Actor
}

func (m *Movie) FormattedName() (formattedName string, error error) {
	nn := make([]byte, len(*m.Name()))
	copy(nn, *m.Name())
	name := strings.ToLower(string(nn))
	parts := strings.Split(name, ".")
	ext := parts[len(parts)-1]
	name = strings.Join(parts[:len(parts)-1], ".")

	re, err := regexp.Compile(`[-\s\t!@#$%^&*()[\]<>,.?~]`)
	if err != nil {
		logger.Error("Movie::FormattedName failed to compile regex with error:", err)
		return "", err
	}

	rs := re.ReplaceAll([]byte(name), []byte{'_'})

	return string(*(m.CleanUnderscores(&rs))) + "." + ext, nil
}

func (m *Movie) CleanUnderscores(name *[]byte) *[]byte {
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

	return nn
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
		m.NewName = m.AddName(a)
	}
}
