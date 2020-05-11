package api

import (
	"encoding/json"
	"github.com/greenac/artemis/dbinteractors"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"net/http"
)

func UnknownMovies(w http.ResponseWriter, r *http.Request) {
	logger.Log("Getting all actors...")

	mvs, err := dbinteractors.UnknownMovies()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]*[]models.Movie{"movies": mvs})
	if err != nil {
		logger.Error("UnknownMovies::Failed to encode movie json", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
