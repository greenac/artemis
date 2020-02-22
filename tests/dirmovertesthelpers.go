package artemis_tests

import (
	"github.com/greenac/artemis/mocks"
	"github.com/greenac/artemis/models"
)

func CreateMovieAndNumber(movieName string, num int) *models.MovieAndNumber {
	fi := mocks.MockFileInfo{MockName: movieName}

	m := models.Movie{
		File: models.File{
			BasePath:    "/path/to/movie",
			Info:        fi,
			NewName:     "",
			NewBasePath: "",
		},
	}

	m.GetNewName()

	return &models.MovieAndNumber{
		Movie:  &m,
		Number: num,
	}
}
