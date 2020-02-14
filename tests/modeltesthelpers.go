package artemis_tests

import (
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/mocks"
	"github.com/greenac/artemis/models"
)

func createActorAndMovie(fileName string, middleName string) (actor *models.Actor, movie *models.Movie) {
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
	a, m := createActorAndMovie(mn, middleName)

	m.GetNewName()

	return m.AddName(a)
}

func CreateMovieWithActors(fileName string, names *[]string) *models.Movie {
	ah := handlers.ActorHandler{
		DirPaths:   nil,
		NamesPath:  nil,
		CachedPath: nil,
		Actors:     nil,
		ToPath:     nil,
	}

	acts := make([]*models.Actor, len(*names))
	for i, n := range *names {
		bn := []byte(n)
		a, err := ah.CreateActor(&bn)
		if err != nil {
			panic(err)
		}

		acts[i] = &a
	}

	fi := mocks.MockFileInfo{MockName: fileName}

	m := models.Movie{
		File: models.File{
			Path:    "/path/to/movie",
			Info:    fi,
			NewName: "",
			NewPath: "",
		},
		Actors: acts,
	}

	return &m
}
