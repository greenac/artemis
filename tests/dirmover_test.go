package artemis_tests

import (
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirMover_MovieNumber(t *testing.T) {
	mn := "scene_480p_34_brad_pitt.mp4"

	n, err := handlers.GetMovieNumber(mn)
	if err != nil {
		logger.Error("TestDirMover_MovieNumber got err:", err)
		panic(err)
	}

	assert.Equal(t, 34, n, "Movie numbers should match")
}

func TestDirMover_MovieWithoutNumber(t *testing.T) {
	mn := "scene_480p_brad_pitt.mp4"

	n, err := handlers.GetMovieNumber(mn)
	if err != nil {
		logger.Error("TestDirMover_MovieNumber got err:", err)
		panic(err)
	}

	assert.Equal(t, -1, n, "Movie numbers should match")
}

func TestDirMover_UpdateMovieNumber(t *testing.T) {
	mn := "scene_480p_34_brad_pitt.mp4"
	exp := "scene_480p_1_brad_pitt.mp4"

	nn := handlers.UpdateMovieNumber(mn, 34, 1)

	assert.Equal(t, exp, nn, "Movie number should update correctly")
}
