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
	allActors        string = "/api/actor/all"
	newActor         string = "/api/actor/new"
	actorsMatch      string = "/api/actor/match"
	addActorsToMovie string = "/api/movie/add-actors"
	openMovie string = "/api/movie/open"
	unknownMovies    string = "/api/movie/unknown"
)

func StartServer(ac *config.ArtemisConfig) {
	db.SetupMongo(&ac.Mongo)

	url := fmt.Sprintf("%s:%d", ac.Url, ac.Port)

	logger.Log("Starting artemis server on", url)

	// Actor routes
	http.HandleFunc(allActors, api.AllActors)
	http.HandleFunc(actorsMatch, api.ActorsMatchingInput)
	http.HandleFunc(newActor, api.CreateActorWithName)

	// movie routes
	http.HandleFunc(unknownMovies, api.UnknownMovies)
	http.HandleFunc(addActorsToMovie, api.AddActorsToMovie)
	http.HandleFunc(openMovie, api.OpenMovie)

	log.Fatal(http.ListenAndServe(url, nil))
}
