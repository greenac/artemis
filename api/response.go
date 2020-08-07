package api

import (
	"github.com/greenac/artemis/models"
)

const PaginatedSize = 50

type PaginatedResponse struct {
	Page int `json:"page"`
	Length int `json:"length"`
	Size int `json:"size"`
	Total int `json:"total"`
}

type PaginatedMovieResponse struct {
	Movies *[]models.Movie `json:"movies"`
	PaginatedResponse
}
