package api

import (
	"encoding/json"
	"github.com/greenac/artemis/dbinteractors"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"net/http"
	"strings"
)

func AllActors(w http.ResponseWriter, r *http.Request) {
	logger.Log("Getting all actors...")

	acts, err := dbinteractors.AllActors()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(acts)
	if err != nil {
		logger.Error("allActors::Failed to encode actor json", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ActorsMatchingInput(w http.ResponseWriter, r *http.Request) {
	logger.Log("ActorsMatchingInput::", r.URL.Query())

	qry := r.URL.Query()

	if len(qry) != 1 {
		logger.Error("ActorsMatchingInput::query string has incorrect params:", qry)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Log("getting from query:", qry.Get("q"))

	acts, err := dbinteractors.GetActorsForInput(
		strings.Trim(qry.Get("q"), " "),
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]*[]models.Actor{"actors": acts})
	if err != nil {
		logger.Error("ActorsForInput::Failed to encode actor json", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}
