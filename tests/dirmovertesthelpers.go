package artemis_tests

import (
	"github.com/greenac/artemis/mocks"
	"github.com/greenac/artemis/models"
)

func CreateMovie(movieName string, num int) *models.MovieAndNumber {
	fi := mocks.MockFileInfo{MockName: movieName}

	m := models.Movie{
		File: models.File{
			BasePath:    "/path/to/movie",
			Info:        fi,
			NewName:     "",
			NewBasePath: "",
		},
	}

	return &models.MovieAndNumber{
		Movie:  &m,
		Number: num,
	}
}
