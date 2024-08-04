package utils

import (
	"fmt"
	"github.com/greenac/artemis/pkg/logger"
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

func MoveFile(from string, to string, delOrig bool) error {
	const BufferSize = 1e6

	ff, err := os.Open(from)
	if err != nil {
		return err
	}

	st, err := ff.Stat()
	if err != nil {
		logger.Error("MoveFile::Could not get stat for file:", from, err)
		return err
	}

	size := st.Size()

	logger.Log("MoveFile::Size:", size, "for file:", from)

	tf, err := os.OpenFile(to, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("MoveFile::Failed to open file:", to, "error:", err)
		return err
	}

	defer ff.Close()
	defer tf.Close()

	var buf []byte

	var offset int64 = 0

	run := true

	for run {
		var bufSize int64
		if offset+BufferSize > size {
			bufSize = size - BufferSize
			run = false
		} else {
			bufSize = BufferSize
		}

		buf = make([]byte, bufSize)

		br, err := ff.ReadAt(buf, offset)
		if err != nil {
			if err.Error() == "negative offset" {
				run = false
			}
		}

		_, err = tf.Write(buf)
		if err != nil {
			return err
		}

		offset += int64(br)
	}

	if delOrig {
		err := os.Remove(from)
		if err != nil {
			return err
		}
	}

	return nil
}
