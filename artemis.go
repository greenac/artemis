package main

import (
	"encoding/json"
	"github.com/greenac/artemis/config"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/tools"
	"github.com/greenac/artemis/ui"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file:", err)
	}

	lp := os.Getenv("LOG_PATH")
	if lp != "" {
		logger.Setup(lp)
	}

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

	logger.Log(ac)

	//targs, has := config["targetDirs"].([]string)
	//if !has {
	//	logger.Error("No target directories in config")
	//	panic("INVALID_CONFIG")
	//}
	//
	//acts, has := config["actorDirs"].([]string)
	//if !has {
	//	logger.Error("No actor directories in config")
	//	panic("INVALID_CONFIG")
	//}
	//
	//anfp, has := config["actorNamesFile"].(string)
	//if !has {
	//	logger.Error("No actor names file path in config")
	//	panic("INVALID_CONFIG")
	//}
	//
	//cnfp, has := config["cachedNamesFile"].(string)
	//if !has {
	//	logger.Error("No cached names file path in config")
	//	panic("INVALID_CONFIG")
	//}
	//
	//sdfp, has := config["stagingDir"].(string)
	//if !has {
	//	logger.Error("No cached names file path in config")
	//	panic("INVALID_CONFIG")
	//}

	logger.Debug("Actor name file:", ac.ActorNamesFile)

	actNameFile := tools.FilePath{Path: ac.ActorNamesFile}
	cachedPath := tools.FilePath{Path: ac.CachedNamesFile}
	stagingPath := tools.FilePath{Path: ac.StagingDir}

	targetPaths := make([]tools.FilePath, len(ac.TargetDirs))
	actorPaths := make([]tools.FilePath, len(ac.ActorDirs))

	for i, p := range ac.TargetDirs {
		targetPaths[i] = tools.FilePath{Path: p}
	}

	for i, p := range ac.ActorDirs {
		actorPaths[i] = tools.FilePath{Path: p}
	}

	anh := ui.AddNamesHandler{}
	anh.Setup(&targetPaths, &actorPaths, &actNameFile, &cachedPath, &stagingPath)
	anh.Run()
}
