package startup

import (
	"encoding/json"
	"fmt"
	"github.com/greenac/artemis/cli"
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"log"
	"net/http"
)

type ArtemisRunType string

const (
	Rename             ArtemisRunType = "RENAME"
	MoveMovies         ArtemisRunType = "MOVE_MOVIES"
	OrganizeStagingDir ArtemisRunType = "ORGANIZE_STAGING_DIR"
	WriteNames         ArtemisRunType = "WRITE_NAMES_TO_FILE"
	Server             ArtemisRunType = "SERVER"
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
	err := anh.Setup(&targetPaths, &actorPaths, &actNameFile, &cachedPath, &stagingPath)
	if err != nil {
		panic(err)
	}

	anh.Run()
}

func OrganizeStagingDirectory(ac *models.ArtemisConfig) {
	err := handlers.OrganizeAllRepeatNamesInDir(ac.StagingDir)
	if err != nil {
		logger.Error("OrganizeStagingDirectory failed with error:", err)
		panic(err)
	}
}

func WriteNamesToFile(ac *models.ArtemisConfig) {
	actNameFile := models.FilePath{Path: ac.ActorNamesFile}
	targetPaths := make([]models.FilePath, len(ac.TargetDirs))
	actorPaths := make([]models.FilePath, len(ac.ActorDirs))

	for i, p := range ac.TargetDirs {
		targetPaths[i] = models.FilePath{Path: p}
	}

	for i, p := range ac.ActorDirs {
		actorPaths[i] = models.FilePath{Path: p}
	}

	ah := handlers.ActorHandler{
		DirPaths:  &actorPaths,
		NamesPath: &actNameFile,
	}

	err := ah.FillActors()
	if err != nil {
		panic(err)
	}

	err = ah.WriteActorsToFile()
	if err != nil {
		panic(err)
	}
}

func MoveMoviesFromStagingToMaster(ac *models.ArtemisConfig) {
	err := handlers.MoveMovies(ac.StagingDir, ac.ToDir)
	if err != nil {
		logger.Error("MoveMoviesFromStagingToMaster could not move movies. Failed with error:", err)
	}
}

func RunServer(ac *models.ArtemisConfig) {
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

	ah := handlers.ArtemisHandler{}
	err := ah.Setup(&targetPaths, &actorPaths, &actNameFile, &cachedPath, &stagingPath)

	if err != nil {
		panic(err)
	}


	logger.Log("Starting artemis server on", fmt.Sprintf("%s:%d", ac.Url, ac.Port))
	http.HandleFunc("/api/all-actors", func (w http.ResponseWriter, r *http.Request) {
		logger.Log("Getting all actors...")

		err := json.NewEncoder(w).Encode(ah.ActorHandler.Actors)
		s := http.StatusOK
		if err != nil {
			logger.Error("Failed to encode actor json", err)
			s = http.StatusInternalServerError
		}

		w.WriteHeader(s)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", ac.Url, ac.Port), nil))
}
