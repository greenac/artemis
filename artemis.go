package main

import (
	"encoding/json"
	"fmt"
	"github.com/greenac/artemis/artemiserror"
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"github.com/greenac/artemis/ui"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
)


type ArtemisRunType string

const (
	Rename ArtemisRunType = "RENAME"
	OrganizeSingleDir  ArtemisRunType = "ORGANIZE_SINGLE_DIR"
)


func RenameMovies(ac *models.ArtemisConfig) {
	actNameFile := models.FilePath{Path: ac.ActorNamesFile}
	cachedPath := models.FilePath{Path: ac.CachedNamesFile}
	stagingPath := models.FilePath{Path: ac.StagingDir}

	targetPaths := make([]models.FilePath, len(ac.TargetDirs))
	actorPaths := make([]models.FilePath, len(ac.ActorDirs))

	for i, p := range ac.TargetDirs {
		targetPaths[i] = models.FilePath{Path: p}
	}

	for i, p := range ac.ActorDirs {
		actorPaths[i] = models.FilePath{Path: p}
	}

	anh := ui.AddNamesHandler{}
	anh.Setup(&targetPaths, &actorPaths, &actNameFile, &cachedPath, &stagingPath)
	anh.Run()
}

func OrganizeSingleDirectory(ac *models.ArtemisConfig) {
	err := handlers.OrganizeRepeatNamesInDir("/Users/andre/Downloads/p/organized/abella_danger_copy")
	if err != nil {
		logger.Error("Failed to run organize single directory with error:", err)
		panic(err)
	}
}


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
	case OrganizeSingleDir:
		OrganizeSingleDirectory(&ac)
	default:
		logger.Error("Unknown run type:", rt)
		panic(artemiserror.New(artemiserror.InvalidParameter))
	}
}
