package movie

type MovieExt string
var movieExts = [12]MovieExt {
"mp4",
"wmv",
"avi",
"mpg",
"mpeg",
"mov",
"asf",
"mkv",
"flv",
"m4v",
"rmvb",
"si",
}

func MovieExts() *[]MovieExt {
	me := make([]MovieExt, len(movieExts))
	for i := 0; i < len(movieExts); i++ {
		me[i] = movieExts[i]
	}

	return &me
}