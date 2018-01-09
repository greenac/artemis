package tools

type FilePath struct {
	Path string
}

func (fp *FilePath)PathDefined() bool {
	return len(fp.Path) > 0
}

func (fp *FilePath)PathAsBytes() *[]byte {
	p := []byte(fp.Path)
	return &p
}

func (fp *FilePath)PathAsString() string {
	return fp.Path
}
