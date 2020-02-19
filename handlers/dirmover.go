package handlers

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"regexp"
	"strconv"
	"strings"
)


func MoveDir(dir models.File) error {
	if !dir.IsDir() {
		return nil
	}

	fh := FileHandler{BasePath: models.FilePath{Path: dir.Path()}}

	err := fh.SetFiles()
	if err != nil {
		logger.Error("MoveDir failed to read files in dir:", dir.Path, err)
		return err
	}

	exists, err := fh.DoesFileExistAtPath(dir.NewPath())
	if err != nil {
		logger.Warn("MoveDir failed to move directory to:", dir.GetNewTotalPath(), err)
		return err
	}

	if exists {
		for _, f := range fh.Files {
			f.NewBasePath = dir.NewPath()
			if f.IsMovie() {
				err := fh.Rename(f.Path(), f.GetNewTotalPath())
				if err != nil {
					continue
				}
			}
		}
	} else {
		err = fh.Rename(dir.Path(), dir.GetNewTotalPath())
		return err
	}

	return nil
}

func UpdateRepeatNames(dirPath string) error {
	fh := FileHandler{BasePath: models.FilePath{Path: dirPath}}

	err := fh.SetFiles()
	if err != nil {
		logger.Error("UpdateRepeatNames failed to read files in dir:", dirPath, err)
		return err
	}

	//name := strings.ToLower(string(nn))
	//fls := make([]models.File, 0)
	//
	//for _, f := range fh.Files {
	//	if f.IsMovie() && strings.Contains(f.Name(), "scene_") {
	//
	//		fls = append(fls, f)
	//	}
	//}
	//
	//
	//if strings.Contains(name, "scene_") {
	//	re, err := regexp.Compile(`\\(.+?\\)`)
	//	if err != nil {
	//		logger.Warn("UpdateRepeatNames failed to compile regex with error:", err)
	//	}
	//
	//	name = re.ReplaceAllString(name, "")
	//}

	return nil
}

func GetMovieNumber(name string) (int, error) {
	parts := strings.Split(name, ".")
	n := parts[0]

	if strings.Contains(n, "scene_", ) {
		re, err := regexp.Compile(`_.[0-9]+_`)
		if err != nil {
			logger.Error("GetMovieNumber Could not compile regex", err)
			return -1, err
		}

		matches := re.FindAllString(n, -1)
		if len(matches) == 0 {
			return -1, nil
		}

		m := matches[len(matches)-1]
		m = strings.ReplaceAll(m, "_", "")
		mi, err := strconv.Atoi(m)
		if err != nil {
			return -1, nil
		}

		return mi, nil
	}

	return -1, nil
}

func UpdateMovieNumber(name string, oldNum int, newNum int) string {
	on := strconv.Itoa(oldNum)
	i := strings.LastIndex(name, on)
	if i == -1 {
		return name
	}

	rn := []rune(name)
	return string(append(rn[:i], append([]rune(strconv.Itoa(newNum)), rn[i+len(on):]...)...))
}
