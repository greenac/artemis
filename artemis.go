package main

import (
	"encoding/json"
	"fmt"
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	lp := os.Getenv("LOG_PATH")
	if lp != "" {
		logger.Setup(lp)
	}

	logger.Log("Starting artemis...")

	cp := os.Getenv("CONFIG_PATH")
	if cp == "" {
		logger.Error("No config path set")
		panic("NO_CONFIG_PATH")
	}

	data, err := ioutil.ReadFile(cp)
	if err != nil {
		logger.Error("Failed to config file")
		panic(err)
	}

	ac := models.ArtemisConfig{}
	err = json.Unmarshal(data, &ac)
	if err != nil {
		logger.Error("failed to unmarshal config file json")
		panic(err)
	}

	rt := ArtemisRunType(os.Getenv("ARTEMIS_RUN_TYPE"))

	logger.Log("Running in mode:", rt)

	switch rt {
	case Rename:
		RenameMovies(&ac)
	case OrganizeStagingDir:
		OrganizeStagingDirectory(&ac)
	default:
		logger.Error("Unknown run type:", rt)
		panic(artemiserror.New(artemiserror.InvalidParameter))
	}

	logger.Log("Finished running:", rt)
}
