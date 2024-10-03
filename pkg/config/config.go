package config

type MongoConfig struct {
	Url         string `json:"url"`
	Database    string `json:"database"`
	Collections struct {
		Actors string `json:"actors"`
		Movies string `json:"movies"`
	} `json:"collections"`
}

type ProfileImageConfig struct {
	ImageSiteBaseUrl     string `json:"imageSiteBaseUrl"`
	ImageSiteBaseUrl2    string `json:"imageSiteBaseUrl2"`
	ImageSiteSubBaseUrl2 string `json:"imageSiteSubBaseUrl2"`
	HtmlTarget           string `json:"htmlTarget"`
	HtmlTarget2          string `json:"htmlTarget2"`
	DestPath             string `json:"destPath"`
}

type ArtemisConfig struct {
	TargetDirs         []string           `json:"targetDirs"`
	ActorDirs          []string           `json:"actorDirs"`
	ActorNamesFile     string             `json:"actorNamesFile"`
	CachedNamesFile    string             `json:"cachedNamesFile"`
	OrganizedDir       string             `json:"organizedDir"`
	ToDir              string             `json:"toDir"`
	FromDir            string             `json:"fromDir"`
	Url                string             `json:"url"`
	Port               int                `json:"port"`
	VlcPath            string             `json:"vlcPath"`
	Mongo              MongoConfig        `json:"mongo"`
	ProfilePicPath     string             `json:"profilePicPath"`
	ImageSiteBaseUrl   string             `json:"imageSiteBaseUrl"`
	ProfileImageConfig ProfileImageConfig `json:"profileImageConfig"`
}
