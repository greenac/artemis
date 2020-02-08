package movie

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/tools"
	"regexp"
	"strings"
)

type Movie struct {
	tools.File
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

		cln = cln[cut + 1:]
	}

	return &cln
}

func (m *Movie) AddName(name string) {
	logger.Debug("Adding name:", name, "to movie:", m.NewName, m.Name())
	n := ""
	if m.NewName == "" {
		nn, err := m.FormattedName()
		if err != nil {
			logger.Error("Could not add name:", name, "to movie name:", m.Name())
			return
		}

		m.NewName = nn
	}

	pts := strings.Split(m.NewName, ".")
	if len(pts) != 2 {
		logger.Warn("`Movie::Rename`", n, "not in proper format")
		return
	}

	nm := pts[0]
	nmb := []byte(nm)
	an := strings.ToLower(name)

	if strings.Contains(nm, an) {
		return
	}

	if nmb[len(nmb)-1] == '_' {
		nm += an
	} else {
		nm += "_" + an
	}

	m.NewName = nm
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
