package api

import (
	"github.com/greenac/artemis/pkg/models"
)

const PaginatedSize = 25

type PaginatedResponse struct {
	Page   int   `json:"page"`
	Length int   `json:"length"`
	Size   int   `json:"size"`
	Total  int64 `json:"total"`
}

type PaginatedMovieResponse struct {
	Movies *[]models.Movie `json:"movies"`
	PaginatedResponse
}

type PaginatedActorResponse struct {
	Actors []models.Actor `json:"actors"`
	PaginatedResponse
}
