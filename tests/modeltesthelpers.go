package artemis_tests

import (
	"github.com/greenac/artemis/mocks"
	"github.com/greenac/artemis/models"
)

func createActorAndModel(fileName string, middleName string) (actor *models.Actor, movie *models.Movie) {
	fn := "brad"
	ln := "pitt"

	fi := mocks.MockFileInfo{MockName: fileName}

	var mn *string
	if middleName != "" {
		mn = &middleName
	}

	a := models.Actor{
		FirstName:  &fn,
		LastName:   &ln,
		MiddleName: mn,
	}

	m := models.Movie{
		File: models.File{
			Path:    "/path/to/movie",
			Info:    fi,
			NewName: "",
			NewPath: "",
		},
		Actors: []*models.Actor{&a},
	}

	return &a, &m
}

func AddName(mn string, middleName string) string {
	a, m := createActorAndModel(mn, middleName)

	m.GetNewName()

	return m.AddName(a)
}
