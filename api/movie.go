package api

import (
	"encoding/json"
	"github.com/greenac/artemis/dbinteractors"
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/utils"
	"net/http"
)

func UnknownMovies(w http.ResponseWriter, r *http.Request) {
	logger.Log("Getting unknown movies...")

	res := utils.Response{Code: http.StatusOK}

	mvs, err := dbinteractors.UnknownMovies()
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Respond(w)
		return
	}

	res.SetPayload("movies", mvs)
	res.Respond(w)
}

func AddActorsToMovie(w http.ResponseWriter, r *http.Request) {
	logger.Log("Adding actors to movie")

	res := utils.Response{Code: http.StatusOK}

	var body struct {
		MovieId string `json:"movieId"`
		ActorIds []string `json:"actorIds"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	logger.Debug("Got body:", body, body.MovieId, body.ActorIds)

	err = handlers.AddActorsToMovie(body.MovieId, body.ActorIds)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Respond(w)
		return
	}

	res.SetPayload("success", true)
	res.Respond(w)
}

func OpenMovie(w http.ResponseWriter, r *http.Request) {
	logger.Log("Getting unknown movies...")

	res := utils.Response{Code: http.StatusOK}
	qry := r.URL.Query()

	if len(qry) != 1 {
		logger.Error("OpenMovie::query string has incorrect params:", qry)
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	logger.Log("getting from query:", qry.Get("movieId"))

	m, err := dbinteractors.GetMovieByIdString(qry.Get("movieId"))
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	err = handlers.OpenMovie(m)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Respond(w)
		return
	}

	res.SetPayload("success", true)
	res.Respond(w)
}

