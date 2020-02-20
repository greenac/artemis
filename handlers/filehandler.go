package handlers

import (
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"io"
	"io/ioutil"
	"os"
)

type FileHandler struct {
	Files    []models.File
	BasePath models.FilePath
}

func (fh *FileHandler) SetFiles() error {
	if !fh.BasePath.PathDefined() {
		panic("File Handler Base Path Not Set")
	}

	fi, err := ioutil.ReadDir(fh.BasePath.PathAsString())
	if err != nil {
		logger.Error("Failed to set file handler file names at path:", fh.BasePath.PathAsString(), "error:", err)
		return err
	}

	files := make([]models.File, len(fi))
	for i, f := range fi {
		files[i] = models.File{Info: f, BasePath: fh.BasePath.PathAsString()}
	}

	fh.Files = files

	return nil
}

func (fh *FileHandler) FileNames() *[][]byte {
	names := make([][]byte, len(fh.Files))
	for i, f := range fh.Files {
		names[i] = []byte(f.Name())
	}

	return &names
}

func (fh *FileHandler) DirFiles() *[]models.File {
	dFiles := make([]models.File, 0)
	for _, f := range fh.Files {
		if f.IsDir() {
			dFiles = append(dFiles, f)
		}
	}

	return &dFiles
}

func (fh *FileHandler) DirFileNames() *[][]byte {
	dFiles := fh.DirFiles()
	names := make([][]byte, len(*dFiles))
	for i, f := range *dFiles {
		names[i] = []byte(f.Name())
	}

	return &names
}

func (fh *FileHandler) ReadNameFile(p *models.FilePath) (*[][]byte, error) {
	f, err := os.Open(p.PathAsString())
	if err != nil {
		logger.Error("Failed to open name file at path:", p.PathAsString(), err)
		return nil, err
	}

	const buffer = 1000
	var offset int64 = 0
	cont := true
	data := make([]byte, 0)
	for cont {
		d := make([]byte, buffer)
		n, err := f.ReadAt(d, offset)
		if err != nil {
			if err == io.EOF {
				cont = false
			} else {
				logger.Error("Error reading name file at path:", p.PathAsString(), err)
				return nil, err
			}
		}

		offset += int64(n)
		data = append(data, d[:n]...)
	}

	names := make([][]byte, 0)
	i := 0
	for _, b := range data {
		if b == '\n' {
			i += 1
		} else {
			if i == len(names) {
				names = append(names, []byte{b})
			} else {
				names[i] = append(names[i], b)
			}
		}
	}

	return &names, nil
}

func (fh *FileHandler) Rename(oldPath string, newPath string, replaceExisting bool) error {
	exists, err := fh.DoesFileExistAtPath(newPath)
	if err != nil {
		logger.Error("FileHandler::Rename could not rename file at:", oldPath, "to:", newPath)
		return err
	}

	if exists && !replaceExisting {
		logger.Warn("FileHandler::Rename could not rename:", oldPath, "to:", newPath, "file exists already")
		return nil
	}

	return os.Rename(oldPath, newPath)
}

func (fh *FileHandler) DoesFileExistAtPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true, nil
		} else if os.IsNotExist(err) {
			return false, nil
		}

		logger.Debug("`FileHandler::DoesFileExistAtPath` Error checking if file exists at path:", path)
		return false, err
	}

	return false, nil
}
