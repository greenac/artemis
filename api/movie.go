package api

import (
	"encoding/json"
	"github.com/greenac/artemis/dbinteractors"
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/utils"
	"net/http"
	"strconv"
)

func UnknownMovies(w http.ResponseWriter, r *http.Request) {
	res := utils.Response{Code: http.StatusOK}

	qry := r.URL.Query()

	if len(qry) != 1 {
		logger.Error("UnknownMovies::query string has incorrect params:", qry)
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	logger.Log("getting from query:", qry.Get("page"))

	page, err := strconv.Atoi(qry.Get("page"))
	if err != nil {
		logger.Error("UnknownMovies::failed to parse page:", qry.Get("page"), err)
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	mvs, total, err := dbinteractors.UnknownMovies(page, PaginatedSize)
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Respond(w)
		return
	}

	pr := PaginatedMovieResponse{
		Movies: mvs,
		PaginatedResponse: PaginatedResponse{
			Page: page,
			Length: len(*mvs),
			Size: PaginatedSize,
			Total: total,
		},
	}

	res.SetPayloadNoKey(&pr)
	res.Respond(w)
}

func AddActorsToMovie(w http.ResponseWriter, r *http.Request) {
	res := utils.Response{Code: http.StatusOK}

	var body struct {
		MovieId  string   `json:"movieId"`
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

func MoviesForIds(w http.ResponseWriter, r *http.Request) {
	res := utils.Response{Code: http.StatusOK}

	var body struct {
		MovieIds []string `json:"movieIds"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	mvs, err := handlers.GetMovieWithIds(body.MovieIds)
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	res.SetPayload("movies", mvs)
	res.Respond(w)
}

func RemoveMovie(w http.ResponseWriter, r *http.Request) {
	logger.Log("RemoveMovie...")

	res := utils.Response{Code: http.StatusOK}
	qry := r.URL.Query()

	if len(qry) != 1 {
		logger.Error("RemoveMovie::query string has incorrect params:", qry)
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	logger.Log("deleting movie:", qry.Get("movieId"))

	err := handlers.DeleteMovie(qry.Get("movieId"))
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	res.SetPayload("success", true)
	res.Respond(w)
}

func SearchMovieByDate(w http.ResponseWriter, r *http.Request) {
	res := utils.Response{Code: http.StatusOK}
	qry := r.URL.Query()

	if len(qry) != 1 {
		logger.Error("SearchMovieByDate::query string has incorrect params:", qry)
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	mvs, err := handlers.SearchMoviesByDate(qry.Get("name"))
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	res.SetPayload("data", map[string]interface{}{"movies": mvs, "count": len(*mvs), "page": 1})
	res.Respond(w)
}

func GetActorsForMovie(w http.ResponseWriter, r *http.Request) {
	res := utils.Response{Code: http.StatusOK}
	qry := r.URL.Query()

	if len(qry) != 1 {
		logger.Error("GetActorsForMovie::query string has incorrect params:", qry)
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	acts, err := handlers.ActorsInMovie(qry.Get("movieId"))
	if err != nil {
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	res.SetPayload("data", map[string]interface{}{"actors": acts, "count": len(*acts), "page": 1})
	res.Respond(w)
}
