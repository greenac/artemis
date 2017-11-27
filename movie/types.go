package movie

type MovieType string
var movieTypes = [12]MovieType {
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

func MovieTypes() *[]MovieType {
mt := make([]MovieType, len(movieTypes))
for i := 0; i < len(movieTypes); i++ {
mt[i] = movieTypes[i]
}

return &mt
}