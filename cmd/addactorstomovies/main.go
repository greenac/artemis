package main

import (
	"github.com/greenac/artemis/pkg/logger"
	"github.com/greenac/artemis/pkg/startup"
	"os"
)

func main() {
	ac, err := startup.GetConfig()
	if err != nil {
		os.Exit(1)
	}

	err = startup.AddActorsToMovies(&ac)
	if err != nil {
		logger.Error("failed to run main with error:", err)
		os.Exit(1)
	}

	logger.Log("run main successful")
}
