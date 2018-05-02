package main

import (
	"github.com/greenac/artemis/tools"
)

func main() {
	//dirPaths := [2]string{
	//	"/Users/andre/Downloads/p/03-27",
	//	"/Users/andre/Downloads/p/04-13",
	//}

	//dPaths := make([]tools.FilePath, len(dirPaths))
	//for i, p := range dirPaths {
	//	var fp tools.FilePath
	//	fp = tools.FilePath{Path:p}
	//	dPaths[i] = fp
	//}

	p := tools.FilePath{Path: "/Users/andre/Downloads/p/03-27"}
	fh := tools.FileHandler{BasePath: p}

	//fPaths := make([]tools.FilePath, len(filePaths))
	//for i, p := range filePaths {
	//	p := tools.FilePath{Path:p}
	//	fPaths[i] = p
	//}
	//
	//ah := handlers.ActorHandler{DirPaths: &dPaths, FilePaths: &fPaths}
	//err := ah.FillActors()
	//if err == nil {
	//	logger.Log("Got actors successfully")
	//	ah.PrintActors()
	//} else {
	//	logger.Error("Failed to get actors")
	//}
	//
	//pDirPaths := [1]string{
	//	"/Users/andre/Downloads/p/01-07",
	//}
	//
	//pPaths := make([]tools.FilePath, len(pDirPaths))
	//for i, p := range pDirPaths {
	//	var fp tools.FilePath
	//	fp = tools.FilePath{Path:p}
	//	pPaths[i] = fp
	//}
	//
	//mh := handlers.MovieHandler{DirPaths: &pPaths}
	//mh.GetMovies()
}
