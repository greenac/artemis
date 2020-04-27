package models

type MongoConfig struct {
	Url         string `json:"targetDirs"`
	Database    string `json:"database"`
	Collections struct {
		Actors string `json:"actors"`
		Movies string `json:"movies"`
	} `json:"collections"`
}

type ArtemisConfig struct {
	TargetDirs      []string    `json:"targetDirs"`
	ActorDirs       []string    `json:"actorDirs"`
	ActorNamesFile  string      `json:"actorNamesFile"`
	CachedNamesFile string      `json:"cachedNamesFile"`
	StagingDir      string      `json:"stagingDir"`
	ToDir           string      `json:"toDir"`
	Url             string      `json:"url"`
	Port            int         `json:"port"`
	Mongo           MongoConfig `json:"mongo"`
}
