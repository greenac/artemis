package mocks


import (
"os"
"time"
)


type MockFileInfo struct {
	MockName string
}

func (mfi MockFileInfo) Name() string {
	return mfi.MockName
}

func (mfi MockFileInfo) Size() int64 {
	return 0
}

func (mfi MockFileInfo) Mode() os.FileMode {
	return 0
}

func (mfi MockFileInfo) ModTime() time.Time {
	return time.Now()
}

func (mfi MockFileInfo) IsDir() bool {
	return false
}

func (mfi MockFileInfo) Sys() interface{}  {
	return nil
}
