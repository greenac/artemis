package main

import (
	"encoding/json"
	"fmt"
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/config"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/startup"
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

	ac := config.ArtemisConfig{}
	err = json.Unmarshal(data, &ac)
	if err != nil {
		logger.Error("failed to unmarshal config file json")
		panic(err)
	}

	rt := startup.ArtemisRunType(os.Getenv("ARTEMIS_RUN_TYPE"))

	logger.Log("Running in mode:", rt)

	switch rt {
	case startup.SaveActors:
		startup.SaveActorsFromFile(&ac)
	case startup.WriteNames:
		startup.WriteNamesToFile(&ac)
	case startup.Server:
		startup.RunServer(&ac)
	case startup.SaveMovies:
		startup.SaveMoviesInDirs(&ac)
	case startup.MoveDir:
		startup.MoveMovieDirs(&ac)
	case startup.Test:
		startup.TestRun(&ac)
	case startup.ConvertOrganized:
		startup.UpdateOrganizedMovies(&ac)
	case startup.RemoveLinks:
		startup.RemoveSymLinks(&ac)
	default:
		logger.Error("Unknown run type:", rt)
		panic(artemiserror.New(artemiserror.InvalidParameter))
	}

	logger.Log("Finished running:", rt)
}
