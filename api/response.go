package api

import (
	"github.com/greenac/artemis/models"
)

type PaginatedResponse struct {
	Movies *[]models.Movie `json:"movies"`
	Page int `json:"page"`
	Length int `json:"length"`
	Size int `json:"size"`
	Total int `json:"total"`
}

const PaginatedSize = 50
