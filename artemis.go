package main

import (
	"github.com/greenac/artemis/tools"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/handlers"
)

func main() {
	dirPaths := [2]string{
		"/Users/andre/Downloads/ptest",
		"/Users/andre/Downloads/ptest2",
	}

	filePaths := [1]string{
		"/Users/andre/Downloads/names.txt",
	}

	dPaths := make([]tools.FilePath, len(dirPaths))
	for i, p := range dirPaths {
		var fp tools.FilePath
		fp = tools.FilePath{Path:p}
		dPaths[i] = fp
	}

	fPaths := make([]tools.FilePath, len(filePaths))
	for i, p := range filePaths {
		p := tools.FilePath{Path:p}
		fPaths[i] = p
	}

	ah := handlers.ActorHandler{DirPaths: &dPaths, FilePaths: &fPaths}
	err := ah.FillActors()
	if err == nil {
		logger.Log("Got actors successfully")
		ah.PrintActors()
	} else {
		logger.Error("Failed to get actors")
	}

	pDirPath := [1]string{
		"/Users/andre/Downloads/p/01-07",
	}

	pPaths := make([]tools.FilePath, len())
	mh := handlers.MovieFormatter{}
	mh.

}
