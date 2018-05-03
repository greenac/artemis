package main

import (
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/tools"
	"sort"
)

func main() {
	adps := []tools.FilePath{
		{Path: "/Users/andre/Downloads/pnames"},
	}

	afps := []tools.FilePath{
		{Path: "/Users/andre/Downloads/names.txt"},
	}

	mps := []tools.FilePath{
		{Path: "/Users/andre/Downloads/p/03-27"},
		{Path: "/Users/andre/Downloads/p/04-13"},
		{Path: "/Users/andre/Downloads/p/05-03"},
	}

	ah := handlers.ArtemisHandler{}
	ah.Setup(&mps, &adps, &afps)
	ah.Sort()
	actors := ah.Actors()
	names := make([]string, len(*actors))
	i := 0
	for name := range *actors {
		names[i] = name
		i += 1
	}

	sort.Strings(names)

	for _, n := range names {
		a := (*actors)[n]
		logger.Log(a.FullName())
		for i, m := range a.Movies {
			logger.Log("\t", i, *m.Name())
		}
	}

	logger.Warn("Unknown movies")
	for i, m := range ah.UnknownMovies {
		logger.Log(i+1, *m.Name())
	}
}
