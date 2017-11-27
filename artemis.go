package main

import (
	"github.com/greenac/artemis/tools"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/handlers"
)

func main() {
	paths := [...]string{
		"/Users/andre/Downloads/ptest1",
		"/Users/andre/Downloads/ptest2",
	}

	fPaths := make([]tools.FilePath, len(paths))
	for i, p := range paths {
		p := tools.FilePath{Path:&p}
		fPaths[i] = p
	}

	ah := handlers.ActorHandler{Paths: &fPaths}
	err := ah.FillActors()
	if err == nil {
		logger.Log("Got actors successfully")
	} else {
		logger.Error("Failed to get actors")
	}
}
