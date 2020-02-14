package utils

import (
	"fmt"
	"github.com/greenac/artemis/logger"
	"os"
)

func AppendTxtToFile(filePath string, txt string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("AppendTxtToFile failed to open file:", filePath, err)
		return err
	}

	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("\n%s", txt)); err != nil {
		logger.Error("AppendTxtToFile failed to write line:", txt, "to file:", filePath, err)
	}

	return err
}

func CreateDir(dirPath string) error {
	fi, err := os.Stat(dirPath)

	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0775)
		if err != nil {
			logger.Error("CreateDir` could not make directory:", dirPath)
		}
	} else if err != nil {
		logger.Error("CreateDir error checking file:", err)
	} else if !fi.IsDir() {
		logger.Error("CreateDir File at path:", dirPath, "is not a directory")
	}

	return err
}

func RenameFile(oldPath string, newPath string) error {
	return os.Rename(oldPath, newPath)
}
