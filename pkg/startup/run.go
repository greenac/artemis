package startup

import (
	"encoding/json"
	"fmt"
	"github.com/greenac/artemis/pkg/config"
	"github.com/greenac/artemis/pkg/db"
	"github.com/greenac/artemis/pkg/handlers"
	"github.com/greenac/artemis/pkg/logger"
	"github.com/greenac/artemis/pkg/models"
	"github.com/joho/godotenv"
	"os"
)

type ArtemisRunType string

const (
	SaveActors              ArtemisRunType = "save-actors"
	SaveMovies              ArtemisRunType = "save-movies"
	WriteNames              ArtemisRunType = "write-names=-to-file"
	Server                  ArtemisRunType = "server"
	Test                    ArtemisRunType = "test"
	MoveDir                 ArtemisRunType = "move-dir"
	ConvertOrganized        ArtemisRunType = "convert-organized"
	RenameMovies            ArtemisRunType = "rename-movies"
	Fix21Path               ArtemisRunType = "fix-21"
	CheckGolf               ArtemisRunType = "check-golf"
	AddMissingMoviesToActor ArtemisRunType = "add-missing-movies-to-actor"
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

	ah.Save(false)
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

func UpdateOrganizedMovies(ac *config.ArtemisConfig) {
	db.SetupMongo(&ac.Mongo)
	_ = handlers.AddOrganizedMovies(ac.OrganizedDir)
}

func RenameAllMovies(ac *config.ArtemisConfig) {
	db.SetupMongo(&ac.Mongo)
	_ = handlers.UpdatedMoviePaths()
}

func Fix21(ac *config.ArtemisConfig) {
	db.SetupMongo(&ac.Mongo)

	_ = handlers.FixMoviePath21()
}

func CheckGolfMovies(ac *config.ArtemisConfig) {
	db.SetupMongo(&ac.Mongo)

	_ = handlers.DeleteNonExistentMovies()
}

func FixMissingMoviesForActors(ac *config.ArtemisConfig) error {
	db.SetupMongo(&ac.Mongo)

	return handlers.AddMissingMoviesToActor()
}

func AddActorsToMovies(ac *config.ArtemisConfig) error {
	db.SetupMongo(&ac.Mongo)

	return handlers.AddActorToMovies()
}

func AddSecondaryPaths(ac *config.ArtemisConfig) error {
	db.SetupMongo(&ac.Mongo)

	return handlers.AddSecondaryPaths()
}

func SwapPaths(ac *config.ArtemisConfig) error {
	db.SetupMongo(&ac.Mongo)

	return handlers.SwapPaths()
}

func SaveImages(ac *config.ArtemisConfig, input handlers.SaveImageInput) error {
	db.SetupMongo(&ac.Mongo)

	return handlers.SaveImages(input)
}

func GetConfig() (config.ArtemisConfig, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return config.ArtemisConfig{}, err
	}

	lp := os.Getenv("LOG_PATH")
	if lp != "" {
		logger.Setup(lp)
	}

	logger.Log("Starting artemis...")

	cp := os.Getenv("CONFIG_PATH")
	if cp == "" {
		logger.Error("No config path set")
		return config.ArtemisConfig{}, err
	}

	data, err := os.ReadFile(cp)
	if err != nil {
		logger.Error("Failed to config file with err:", err)
		return config.ArtemisConfig{}, err
	}

	ac := config.ArtemisConfig{}
	err = json.Unmarshal(data, &ac)
	if err != nil {
		logger.Error("failed to unmarshal config file json withe err:", err)
		return config.ArtemisConfig{}, err
	}

	return ac, nil
}
