package main

import (
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"github.com/greenac/artemis/ui"
)

type ArtemisRunType string

const (
	Rename             ArtemisRunType = "RENAME"
	MoveMovies         ArtemisRunType = "MOVE_MOVIES"
	OrganizeStagingDir ArtemisRunType = "ORGANIZE_STAGING_DIR"
	WriteNames         ArtemisRunType = "WRITE_NAMES_TO_FILE"
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
