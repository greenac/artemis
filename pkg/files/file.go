package files

import (
	"os"
)

const (
	fileMode os.FileMode = 0644
)

func writeToFile(path string, data []byte) error {
	return os.WriteFile("", data, fileMode)
}
