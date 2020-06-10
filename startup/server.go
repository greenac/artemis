package startup

import (
	"fmt"
	"github.com/greenac/artemis/api"
	"github.com/greenac/artemis/config"
	"github.com/greenac/artemis/db"
	"github.com/greenac/artemis/logger"
	"log"
	"net/http"
)

const (
	actor                       string = "/api/actor"
	allActors                   string = "/api/actor/all"
	allActorsWithMovies         string = "/api/actor/all-with-movies"
	newActor                    string = "/api/actor/new"
	actorsMatch                 string = "/api/actor/match"
	actorsMatchWithMovies       string = "/api/actor/match-with-movies"
	actorsSimpleMatchWithMovies string = "/api/actor/simple-match-with-movies"
	actorsMovies                string = "/api/actor/movies"
	actorsRecent                string = "/api/actor/recent"
	addActorsToMovie            string = "/api/movie/add-actors"
	openMovie                   string = "/api/movie/open"
	unknownMovies               string = "/api/movie/unknown"
	moviesForIds                string = "/api/movie/ids"
	deleteMovie                 string = "/api/movie/delete"
	searchMovieByDate           string = "/api/movie/search-date"
	actorsForMovie              string = "/api/movie/actors"
)

func StartServer(ac *config.ArtemisConfig) {
	db.SetupMongo(&ac.Mongo)

	url := fmt.Sprintf("%s:%d", ac.Url, ac.Port)

	logger.Log("Starting artemis server on", url)

	// Actor routes
	http.HandleFunc(actor, api.GetActor)
	http.HandleFunc(allActors, api.AllActors)
	http.HandleFunc(allActorsWithMovies, api.AllActorsWithMovies)
	http.HandleFunc(actorsMatch, api.ActorsMatchingInput)
	http.HandleFunc(actorsMatchWithMovies, api.ActorsMatchingInputWithMovies)
	http.HandleFunc(actorsSimpleMatchWithMovies, api.ActorsSimpleMatchingInputWithMovies)
	http.HandleFunc(newActor, api.CreateActorWithName)
	http.HandleFunc(actorsMovies, api.GetMoviesForActor)
	http.HandleFunc(actorsRecent, api.ActorsByDate)

	// movie routes
	http.HandleFunc(unknownMovies, api.UnknownMovies)
	http.HandleFunc(addActorsToMovie, api.AddActorsToMovie)
	http.HandleFunc(openMovie, api.OpenMovie)
	http.HandleFunc(moviesForIds, api.MoviesForIds)
	http.HandleFunc(deleteMovie, api.RemoveMovie)
	http.HandleFunc(searchMovieByDate, api.SearchMovieByDate)
	http.HandleFunc(actorsForMovie, api.GetActorsForMovie)

	log.Fatal(http.ListenAndServe(url, nil))
}
