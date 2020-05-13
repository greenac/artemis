package api

import (
	"encoding/json"
	"github.com/greenac/artemis/dbinteractors"
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/utils"
	"net/http"
	"sort"
	"strings"
)

func AllActors(w http.ResponseWriter, r *http.Request) {
	logger.Log("Getting all actors...")

	res := utils.Response{Code: http.StatusOK}

	acts, err := dbinteractors.AllActors()
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	res.SetPayload("actors", acts)
	res.Respond(w)
}

func AllActorsWithMovies(w http.ResponseWriter, r *http.Request) {
	logger.Log("Getting all actors...")

	res := utils.Response{Code: http.StatusOK}

	acts, err := dbinteractors.AllActorsWithMovies()
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	sort.SliceStable(acts, func(i, j int) bool {
		return strings.ToLower((*acts)[i].FullName()) < strings.ToLower((*acts)[j].FullName())
	})

	res.SetPayload("actors", acts)
	res.Respond(w)
}

func ActorsMatchingInput(w http.ResponseWriter, r *http.Request) {
	logger.Log("ActorsMatchingInput::", r.URL.Query())

	res := utils.Response{Code: http.StatusOK}
	qry := r.URL.Query()

	if len(qry) != 1 {
		logger.Error("ActorsMatchingInput::query string has incorrect params:", qry)
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	logger.Log("getting from query:", qry.Get("q"))

	acts, err := dbinteractors.GetActorsForInput(
		strings.Trim(qry.Get("q"), " "),
	)

	if err != nil {
		res.Code = http.StatusInternalServerError
	} else {
		res.SetPayload("actors", acts)
	}

	res.Respond(w)
}

func CreateActorWithName(w http.ResponseWriter, r *http.Request) {
	logger.Log("CreateActorWithName::Creating actor")

	res := utils.Response{Code: http.StatusOK}

	var body struct {
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	a, err := handlers.CreateNewActor(body.Name)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Respond(w)
		return
	}

	res.SetPayload("actor", a)
	res.Respond(w)
}
