package tools

type FilePath struct {
	Path *string
}

func (fp *FilePath)PathDefined() bool {
	return fp.Path != nil && len(*(fp.Path)) > 0
}

func (fp *FilePath)PathAsBytes() *[]byte {
	if !fp.PathDefined() {
		panic ("File Path Not Set")
	}

	p := []byte(*fp.Path)
	return &p
}

func (fp *FilePath)PathAsString() string {
	if !fp.PathDefined() {
		panic ("File Path Not Set")
	}

	return *fp.Path
}
