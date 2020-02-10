package models


type ArtemisConfig struct {
	TargetDirs      []string
	ActorDirs       []string
	ActorNamesFile  string
	CachedNamesFile string
	StagingDir      string
}
