package movie

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/tools"
	"strings"
)

type Movie struct {
	tools.File
}

func (m *Movie) AddName(name string) {
	n := ""
	if m.NewName == "" {
		n = strings.ToLower(*m.Name())
		pts := strings.Split(n, ".")
		if len(pts) != 2 {
			logger.Warn("`Movie::Rename`", n, "not in proper format")
			return
		}
	} else {

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
