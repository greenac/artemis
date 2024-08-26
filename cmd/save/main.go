package main

import (
	"github.com/greenac/artemis/pkg/startup"
	"os"
)

func main() {
	ac, err := startup.GetConfig()
	if err != nil {
		os.Exit(1)
	}

	startup.SaveMoviesInDirs(&ac)
}
