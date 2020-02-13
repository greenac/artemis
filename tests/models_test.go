package models_test

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/mocks"
	"github.com/greenac/artemis/models"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestMovie_AddName(t *testing.T) {
	input := "a_river_runs_through_it_bradpitt.mp4"
	output := "a_river_runs_through_it_brad_pitt.mp4"
	fn := "Brad"
	ln := "Pitt"
	fi := mocks.MockFileInfo{MockName: input}

	a := models.Actor{
		FirstName: &fn,
		LastName: &ln,
		MiddleName: nil,
		Movies: nil,
	}

	m := models.Movie{
		File: models.File{
			Path: "/path/to/movie",
			Info: fi,
			NewName: "",
			NewPath: "",
		},
		Actors:nil,
	}

	nn := m.GetNewName()
	logger.Log("New name is:", nn)

	n := m.AddName(&a)
	logger.Log("name is:", n)

	assert.Equal(t, output, n, "Movie names should match")
}
