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
	OrganizeStagingDir ArtemisRunType = "ORGANIZE_STAGING_DIR"
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

func OrganizeStagingDirectory(ac *models.ArtemisConfig) {
	err := handlers.OrganizeAllRepeatNamesInDir(ac.StagingDir)
	if err != nil {
		logger.Error("OrganizeStagingDirectory failed with error:", err)
		panic(err)
	}
}
