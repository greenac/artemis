package tools

type FilePath struct {
	Path string
}

func (fp *FilePath)PathDefined() bool {
	return fp.Path != ""
}

func (fp *FilePath)PathAsBytes() *[]byte {
	p := []byte(fp.Path)
	return &p
}

func (fp *FilePath)PathAsString() string {
	return fp.Path
}
