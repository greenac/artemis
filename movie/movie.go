package movie

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/tools"
	"strings"
	"regexp"
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
		logger.Error("Cannot format movie name. Failed to compile regex with error:", err)
		return "", err
	}

	rs := re.ReplaceAll(nn, []byte{'_'})
	return string(rs) + string(ext), nil
}

func (m *Movie) AddName(name string) {
	n := ""
	if m.NewName == "" {
		nn, err := m.FormattedName(); if err != nil {
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
		nn, err := m.FormattedName(); if err != nil {
			return ""
		}

		m.NewName = nn
	}

	return m.NewName
}
