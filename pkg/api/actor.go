package api

import (
	"encoding/json"
	"fmt"
	"github.com/greenac/artemis/pkg/dbinteractors"
	"github.com/greenac/artemis/pkg/handlers"
	"github.com/greenac/artemis/pkg/logger"
	"github.com/greenac/artemis/pkg/middleware"
	"github.com/greenac/artemis/pkg/utils"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

func GetActor(w http.ResponseWriter, r *http.Request) {
	logger.Log("GetActor::", r.URL.Query())

	res := utils.Response{Code: http.StatusOK}
	qry := r.URL.Query()

	if len(qry) != 1 {
		logger.Error("GetActor::query string has incorrect params:", qry)
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	act, err := dbinteractors.GetActorByIdString(
		strings.Trim(qry.Get("actorId"), " "),
	)

	if err != nil {
		res.Code = http.StatusInternalServerError
	} else {
		res.SetPayload("actor", act)
	}

	res.Respond(w)
}

func AllActors(w http.ResponseWriter, r *http.Request) {
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

func PaginatedActors(w http.ResponseWriter, r *http.Request) {
	res := utils.Response{Code: http.StatusOK}
	qry := r.URL.Query()

	if len(qry) < 1 || len(qry) > 3 {
		logger.Error("PaginatedActors->query string has incorrect params:", qry)
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	logger.Log("got all querry params", qry)
	logger.Log("got page search param:", strings.Trim(qry.Get("p"), " "))
	page := 0
	p := strings.Trim(qry.Get("p"), " ")
	if p != "" {
		pg, err := strconv.Atoi(p)
		if err != nil {
			logger.Error("PaginatedActors->failed to convert query to int:", err)
			res.Code = http.StatusBadRequest
			res.Respond(w)
			return
		}
		logger.Log("setting page to:", pg)
		page = pg
	}

	input := strings.Trim(qry.Get("q"), " ")
	t := strings.Trim(qry.Get("t"), " ")

	logger.Log("PaginatedActors->querying for input:", input, "and page:", page, "and type:", t)

	var result dbinteractors.PaginatedQueryResult
	var err error
	if input == "" && t == "" {
		result, err = dbinteractors.ActorsAtPage(page)
	} else if t == "firstName" {
		result, err = dbinteractors.GetActorsForInputSimple(input, page, false)
	} else {
		result, err = dbinteractors.GetActorsForInput(input, false, page)
	}

	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Respond(w)
		return
	}

	logger.Log("PaginatedActors->got # actors::", len(result.Actors), "total:", result.Total)
	logger.Log("firstActor:", result.Actors[0].FullName(), "lastActor:", result.Actors[len(result.Actors)-1].FullName())

	ar := PaginatedActorResponse{
		Actors: result.Actors,
		PaginatedResponse: PaginatedResponse{
			Page:   page,
			Length: len(result.Actors),
			Size:   PaginatedSize,
			Total:  result.Total,
		},
	}

	res.SetPayloadNoKey(ar)
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

	actors := *acts
	sort.SliceStable(actors, func(i, j int) bool {
		return strings.ToLower(actors[i].FullName()) < strings.ToLower(actors[j].FullName())
	})

	res.SetPayload("actors", acts)
	res.Respond(w)
}

func ActorsMatchingInput(w http.ResponseWriter, r *http.Request) {
	logger.Log("ActorsMatchingInput::", r.URL.Query())

	res := utils.Response{Code: http.StatusOK}
	qry := r.URL.Query()

	page := 0
	if qry.Get("page") != "" {
		p, err := strconv.Atoi(strings.Trim(qry.Get("page"), " "))
		if err != nil {
			logger.Error("ActorsMatchingInputWithMovies::failed to convert query to int:", err)
			res.Code = http.StatusBadRequest
			res.Respond(w)
			return
		}
		page = p
	}

	result, err := dbinteractors.GetActorsForInput(
		strings.Trim(qry.Get("q"), " "),
		false,
		page,
	)

	if err != nil {
		res.Code = http.StatusInternalServerError
	} else {
		ar := PaginatedActorResponse{
			Actors: result.Actors,
			PaginatedResponse: PaginatedResponse{
				Page:   page,
				Length: len(result.Actors),
				Size:   PaginatedSize,
				Total:  result.Total,
			},
		}
		res.SetPayloadNoKey(ar)
	}

	res.Respond(w)
}

func ActorsMatchingInputWithMovies(w http.ResponseWriter, r *http.Request) {
	logger.Log("ActorsMatchingInputWithMovies::", r.URL.Query())

	res := utils.Response{Code: http.StatusOK}
	qry := r.URL.Query()

	page := 0
	if qry.Get("page") != "" {
		p, err := strconv.Atoi(strings.Trim(qry.Get("page"), " "))
		if err != nil {
			logger.Error("ActorsMatchingInputWithMovies::failed to convert query to int:", err)
			res.Code = http.StatusBadRequest
			res.Respond(w)
			return
		}
		page = p
	}

	result, err := dbinteractors.GetActorsForInput(
		strings.Trim(qry.Get("q"), " "),
		true,
		page,
	)

	if err != nil {
		res.Code = http.StatusInternalServerError
	} else {
		ar := PaginatedActorResponse{
			Actors: result.Actors,
			PaginatedResponse: PaginatedResponse{
				Page:   page,
				Length: len(result.Actors),
				Size:   PaginatedSize,
				Total:  result.Total,
			},
		}
		res.SetPayloadNoKey(ar)
	}

	res.Respond(w)
}

func ActorsSimpleMatchingInputWithMovies(w http.ResponseWriter, r *http.Request) {
	logger.Log("ActorsSimpleMatchingInputWithMovies::", r.URL.Query())

	res := utils.Response{Code: http.StatusOK}
	qry := r.URL.Query()

	if len(qry) != 1 {
		logger.Error("ActorsSimpleMatchingInputWithMovies::query string has incorrect params:", qry)
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	logger.Log("ActorsSimpleMatchingInputWithMovies::getting from query:", qry.Get("q"))

	acts, err := dbinteractors.GetActorsForInputSimple(strings.Trim(qry.Get("q"), " "), 0, false)
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

func GetMoviesForActor(w http.ResponseWriter, r *http.Request) {
	logger.Log("GetMoviesForActor::", r.URL.Query())

	res := utils.Response{Code: http.StatusOK}
	qry := r.URL.Query()

	if len(qry) != 1 {
		logger.Error("GetMoviesForActor::query string has incorrect params:", qry)
		res.Code = http.StatusBadRequest
		res.Respond(w)
		return
	}

	mvs, err := handlers.GetMoviesForActor(strings.Trim(qry.Get("actorId"), " "))
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Respond(w)
		return
	}

	res.SetPayload("movies", mvs)
	res.Respond(w)
}

func ActorsByDate(w http.ResponseWriter, r *http.Request) {
	res := utils.Response{Code: http.StatusOK}

	acts, err := dbinteractors.ActorsByDate()
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Respond(w)
		return
	}

	res.SetPayload("data", map[string]interface{}{"actors": acts, "count": len(*acts), "page": 1})
	res.Respond(w)
}

func GetActorProfilePicture(w http.ResponseWriter, r *http.Request) {
	qry := r.URL.Query()
	if len(qry) != 1 {
		logger.Error("GetActorProfilePicture->Query has incorrect length:", len(qry))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	basePath, ok := ctx.Value(middleware.ProfilePickMiddleWarePathKey).(string)
	if !ok {
		logger.Error("GetActorProfilePicture->Missing profile pic path in context", basePath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	actorId := strings.Trim(qry.Get("actorId"), " ")

	imagePaths := []string{
		path.Join(basePath, "profile-pics/"),
		path.Join(basePath, "profile-pics-manual/"),
	}

	var imageData []byte
	for _, p := range imagePaths {
		imagePath := fmt.Sprintf("%s.jpg", path.Join(p, actorId))
		data, err := os.ReadFile(imagePath)
		if err != nil {
			continue
		}
		imageData = data
		break
	}

	if len(imageData) == 0 {
		//logger.Error("GetActorProfilePicture->Failed to get image for actor:", actorId)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(imageData)
	if err != nil {
		logger.Error("GetActorProfilePicture->Failed to write response with error:", err)
	}

	logger.Log("GetActorProfilePicture->Successfully retrieved profile pic for actor:", actorId)
}
