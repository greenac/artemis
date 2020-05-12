package handlers

import (
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"os"
	"os/exec"
)

const MoviePlayerEnvName = "ARTEMIS_MOVIE_PLAYER"

func OpenMovie(m *models.Movie) error {
	pp := os.Getenv(MoviePlayerEnvName)
	if pp == "" {
		logger.Error("OpenMovie::No designated movie player")
		return artemiserror.New(artemiserror.ArgsNotInitialized)
	}

	cmd := exec.Command(pp, m.Path)
	err := cmd.Start()
	if err != nil {
		logger.Error("OpenMovieWithId::Failed to open movie:", m.Path, "command:", pp, "error:", err)
		return err
	}

	return nil
}
