package tools

import (
	"github.com/greenac/artemis/logger"
	"io"
	"io/ioutil"
	"os"
)

type FileHandler struct {
	Files    *[]File
	BasePath FilePath
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

	files := make([]File, len(fi))
	for i, f := range fi {
		files[i] = File{Info: f}
	}

	fh.Files = &files
	return nil
}

func (fh *FileHandler) FileNames() *[][]byte {
	names := make([][]byte, len(*fh.Files))
	for i, f := range *fh.Files {
		names[i] = []byte(*f.Name())
	}

	return &names
}

func (fh *FileHandler) DirFiles() *[]File {
	dFiles := make([]File, 0)
	for _, f := range *fh.Files {
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
		names[i] = []byte(*f.Name())
	}

	return &names
}

func (fh *FileHandler) ReadNameFile(p *FilePath) (*[][]byte, error) {
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

func (fh *FileHandler) Rename(oldName string, newName string) error {
	_, err := os.Stat(newName)
	if err != nil {
		if os.IsNotExist(err) {
			return os.Rename(oldName, newName)
		}

		return err
	}

	return nil
}
