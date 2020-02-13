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

	if _, err := f.WriteString(fmt.Sprintf("%s\n", txt))
	err != nil {
		logger.Error("AppendTxtToFile failed to write line:", txt, "to file:", filePath, err)
	}

	return err
}
