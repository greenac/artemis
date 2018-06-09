package main

import (
	"github.com/greenac/artemis/tools"
	"github.com/greenac/artemis/ui"
	"github.com/joho/godotenv"
	"github.com/greenac/artemis/logger"
	"os"
	"io/ioutil"
	"encoding/json"
)

func main() {
	err := godotenv.Load(); if err != nil {
		logger.Error("Error loading .env file:", err)
	}

	cp := os.Getenv("CONFIG_PATH")
	if cp == "" {
		logger.Error("No config path set")
		panic("NO_CONFIG_PATH")
	}

	data, err := ioutil.ReadFile(cp); if err != nil {
		logger.Error("Failed to config file")
		panic(err)
	}

	config := make(map[string]interface{}, 0)
	err = json.Unmarshal(data, &config); if err != nil {
		logger.Error("failed to unmarshal config file json")
		panic(err)
	}


	targs, has := config["targetDirs"].([]string); if !has {
		logger.Error("No target directories in config")
		panic("INVALID_CONFIG")
	}

	acts, has := config["actorDirs"].([]string); if !has {
		logger.Error("No actor directories in config")
		panic("INVALID_CONFIG")
	}

	anfp, has := config["actorNamesFile"].(string); if !has {
		logger.Error("No actor names file path in config")
		panic("INVALID_CONFIG")
	}

	cnfp, has := config["cachedNamesFile"].(string); if !has {
		logger.Error("No cached names file path in config")
		panic("INVALID_CONFIG")
	}

	sdfp, has := config["stagingDir"].(string); if !has {
		logger.Error("No cached names file path in config")
		panic("INVALID_CONFIG")
	}

	actNameFile := tools.FilePath{Path: anfp}
	cachedPath := tools.FilePath{Path: cnfp}
	stagingPath := tools.FilePath{Path: sdfp}

	targetPaths := make([]tools.FilePath, len(targs))
	actorPaths := make([]tools.FilePath, len(acts))

	for i, p := range targs {
		targetPaths[i] = tools.FilePath{Path: p}
	}

	for i, p := range acts {
		actorPaths[i] = tools.FilePath{Path: p}
	}

	anh := ui.AddNamesHandler{}
	anh.Setup(&targetPaths, &actorPaths, &actNameFile, &cachedPath, &stagingPath)
	anh.Run()
}
