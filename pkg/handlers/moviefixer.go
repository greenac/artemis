package handlers

import (
	"github.com/greenac/artemis/pkg/logger"
	"github.com/greenac/artemis/pkg/models"
	"strings"
)

type MoviePathRenamer struct {
	Movie      models.Movie
	Seg        string
	ReplaceSeg string
}

func (mpr *MoviePathRenamer) Fix() error {
	if !strings.Contains(mpr.Movie.Path, mpr.Seg) {
		return nil
	}

	oldId := mpr.Movie.Identifier

	mpr.Movie.Path = strings.Replace(mpr.Movie.Path, mpr.Seg, mpr.ReplaceSeg, 1)
	mpr.Movie.SetIdentifier()

	logger.Warn("MoviePathRenamer::Fix::new path:", mpr.Movie.Path, "old id:", oldId, "new Id:", mpr.Movie.Identifier)

	return nil
}
