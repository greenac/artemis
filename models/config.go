package models

type ArtemisConfig struct {
	TargetDirs      []string `json:"targetDirs"`
	ActorDirs       []string `json:"actorDirs"`
	ActorNamesFile  string   `json:"actorNamesFile"`
	CachedNamesFile string   `json:"cachedNamesFile"`
	StagingDir      string   `json:"stagingDir"`
}
