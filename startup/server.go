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
	allActors     string = "/api/all-actors"
	unknownMovies string = "/api/unknown-movies"
)

func StartServer(ac *config.ArtemisConfig) {
	db.SetupMongo(&ac.Mongo)

	url := fmt.Sprintf("%s:%d", ac.Url, ac.Port)
	logger.Log("Starting artemis server on", url)

	http.HandleFunc(allActors, api.AllActors)
	http.HandleFunc(unknownMovies, api.UnknownMovies)

	log.Fatal(http.ListenAndServe(url, nil))
}
