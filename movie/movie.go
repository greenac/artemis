package movie

import (
	"github.com/greenac/artemis/tools"
	"strings"
	"github.com/greenac/artemis/logger"
)

type Movie struct {
	tools.File
}

func (m *Movie) Rename(a *Actor) {
	n := *m.Name()
	pts := strings.Split(n, ".")
	if len(pts) != 2 {
		logger.Warn("`Movie::Rename", n, "not in proper format")
		return
	}

	nm := pts[0]
	nmb := []byte(nm)
	if nmb[len(nmb) - 1] == '_' {
		nm += a.FullName()
	} else {
		nm += "_" + a.FullName()
	}

	m.NewName = nm
}
