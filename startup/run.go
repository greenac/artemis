package startup

import (
	"github.com/greenac/artemis/config"
	"github.com/greenac/artemis/db"
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
)

type ArtemisRunType string

const (
	SaveActors         ArtemisRunType = "save-actors"
	SaveMovies         ArtemisRunType = "save-movies"
	MoveMovies         ArtemisRunType = "move-movies"
	OrganizeStagingDir ArtemisRunType = "organize-staging-dir"
	WriteNames         ArtemisRunType = "write-names=-to-file"
	Server             ArtemisRunType = "server"
	Test               ArtemisRunType = "test"
	MoveDir            ArtemisRunType = "move-dir"
)

func SaveActorsFromFile(ac *config.ArtemisConfig) {
	db.SetupMongo(&ac.Mongo)

	p := models.FilePath{Path: ac.ActorNamesFile}

	ah := handlers.ActorHandler{NamesPath: &p}
	err := ah.FillActors()
	if err != nil {
		logger.Error("SaveActorsFromFile::Failed to fill up actor handler with error:", err)
		panic(err)
	}

	acts := ah.SortedActors()

	ah.SaveActorsToDb(acts)
}

func SaveMoviesInDirs(ac *config.ArtemisConfig) {
	db.SetupMongo(&ac.Mongo)

	paths := make([]models.FilePath, len(ac.TargetDirs))

	for i, p := range ac.TargetDirs {
		paths[i] = models.FilePath{Path: p}
	}

	ah := handlers.ArtemisHandler{}
	err := ah.Setup(&paths)
	if err != nil {
		logger.Error("SaveUnknown::Failed to set up artemis handler with error:", err)
		panic(err)
	}

	ah.Save()
}

func OrganizeStagingDirectory(ac *config.ArtemisConfig) {
	err := handlers.OrganizeAllRepeatNamesInDir(ac.StagingDir)
	if err != nil {
		logger.Error("OrganizeStagingDirectory failed with error:", err)
		panic(err)
	}
}

func WriteNamesToFile(ac *config.ArtemisConfig) {
	targetPaths := make([]models.FilePath, len(ac.TargetDirs))
	actorPaths := make([]models.FilePath, len(ac.ActorDirs))

	for i, p := range ac.TargetDirs {
		targetPaths[i] = models.FilePath{Path: p}
	}

	for i, p := range ac.ActorDirs {
		actorPaths[i] = models.FilePath{Path: p}
	}

	ah := handlers.ActorHandler{}

	err := ah.FillActors()
	if err != nil {
		panic(err)
	}
}

func MoveMoviesFromStagingToMaster(ac *config.ArtemisConfig) {
	err := handlers.MoveMovies(ac.StagingDir, ac.ToDir)
	if err != nil {
		logger.Error("MoveMoviesFromStagingToMaster could not move movies. Failed with error:", err)
	}
}

func RunServer(ac *config.ArtemisConfig) {
	StartServer(ac)
}

func TestRun(ac *config.ArtemisConfig) {

}

func MoveMovieDirs(ac *config.ArtemisConfig) {
	err := handlers.MoveDirAndUpdateMovies(ac.FromDir, ac.ToDir, &ac.Mongo)
	if err != nil {
		logger.Error("Failed to move dir from", ac.FromDir, "to:", ac.ToDir)
		return
	}

	logger.Log("Successfully moved directories from:", ac.FromDir, "to:", ac.ToDir)
}
