package api

import (
	"encoding/json"
	"github.com/greenac/artemis/dbinteractors"
	"github.com/greenac/artemis/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func UnknownMovies(w http.ResponseWriter, r *http.Request) {
	logger.Log("Getting unknown movies...")

	mvs, err := dbinteractors.UnknownMovies()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	mvIds := make(map[primitive.ObjectID]int, 0)
	for _, m := range *mvs {
		_, has := mvIds[m.Id]
		if has {
			logger.Warn("There are repeat movies with id:", m.Id)
		} else {
			mvIds[m.Id] = 0
		}
	}

	err = json.NewEncoder(w).Encode(mvs)
	if err != nil {
		logger.Error("UnknownMovies::Failed to encode movie json", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
