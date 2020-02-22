package artemis_tests

import (
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDirMover_MovieNumber(t *testing.T) {
	mn := "scene_480p_34_brad_pitt.mp4"

	np := handlers.NameUpdater{DirPath: "some/dir/path"}

	m := CreateMovieAndNumber(mn, 0)
	n, err := np.GetMovieNumber(m)
	if err != nil {
		logger.Error("TestDirMover_MovieNumber got err:", err)
		panic(err)
	}

	assert.Equal(t, 34, n, "Movie numbers should match")
}

func TestDirMover_MovieNumberSingleDigit(t *testing.T) {
	mn := "scene_480p_4_brad_pitt.mp4"

	np := handlers.NameUpdater{DirPath: "some/dir/path"}

	m := CreateMovieAndNumber(mn, 0)
	n, err := np.GetMovieNumber(m)
	if err != nil {
		logger.Error("TestDirMover_MovieNumber got err:", err)
		panic(err)
	}

	assert.Equal(t, 4, n, "Movie numbers should match")
}

func TestDirMover_MovieWithoutNumber(t *testing.T) {
	mn := "scene_480p_brad_pitt.mp4"

	np := handlers.NameUpdater{DirPath: "some/dir/path"}

	m := CreateMovieAndNumber(mn, -1)
	n, err := np.GetMovieNumber(m)
	if err != nil {
		logger.Error("TestDirMover_MovieNumber got err:", err)
		panic(err)
	}

	assert.Equal(t, -1, n, "Movie numbers should match")
}

func TestDirMover_UpdateMovieNumber(t *testing.T) {
	np := handlers.NameUpdater{DirPath: "some/dir/path"}

	m := CreateMovieAndNumber("scene_480p_34_brad_pitt.mp4", 34)
	nn := np.UpdateMovieNameWithNumber(m, 99)

	assert.Equal(t, "scene_480p_99_brad_pitt.mp4", nn, "Movie number should update correctly")
}

func TestDirMover_UpdateMovieNumberSingleDigit(t *testing.T) {
	np := handlers.NameUpdater{DirPath: "some/dir/path"}

	m := CreateMovieAndNumber("scene_480p_1_brad_pitt.mp4", 1)
	nn := np.UpdateMovieNameWithNumber(m, 99)

	assert.Equal(t, "scene_480p_99_brad_pitt.mp4", nn, "Movie number should update correctly")
}

func TestDirMover_UpdateMovieNumberNoNumber(t *testing.T) {
	np := handlers.NameUpdater{DirPath: "some/dir/path"}

	m := CreateMovieAndNumber("scene_480p_brad_pitt.mp4", -1)
	nn := np.UpdateMovieNameWithNumber(m, 1)

	assert.Equal(t, "scene_480p_brad_pitt.mp4", nn, "Movie number should update correctly")
}
