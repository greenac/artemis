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
	allActors     string = "/api/actor/all"
	actorsMatch   string = "/api/actor/match"
	unknownMovies string = "/api/movie/unknown"
)

func StartServer(ac *config.ArtemisConfig) {
	db.SetupMongo(&ac.Mongo)

	url := fmt.Sprintf("%s:%d", ac.Url, ac.Port)

	logger.Log("Starting artemis server on", url)

	http.HandleFunc(allActors, api.AllActors)
	http.HandleFunc(unknownMovies, api.UnknownMovies)
	http.HandleFunc(actorsMatch, api.ActorsMatchingInput)

	log.Fatal(http.ListenAndServe(url, nil))
}
