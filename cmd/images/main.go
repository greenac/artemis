package main

import (
	"github.com/greenac/artemis/pkg/clients/rest"
	"github.com/greenac/artemis/pkg/handlers"
	"github.com/greenac/artemis/pkg/logger"
	"github.com/greenac/artemis/pkg/startup"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ac, err := startup.GetConfig()
	if err != nil {
		os.Exit(1)
	}

	rc := rest.Client{
		HttpClient: &http.Client{Timeout: 30 * time.Second},
		BodyReader: io.ReadAll,
		GetRequest: http.NewRequest,
	}

	input := handlers.SaveImageInput{
		RestClient:  &rc,
		BaseUrl:     ac.ProfileImageConfig.ImageSiteBaseUrl,
		BaseUrl2:    ac.ProfileImageConfig.ImageSiteBaseUrl2,
		SubBaseUrl2: ac.ProfileImageConfig.ImageSiteBaseUrl2,
		HtmlTarget:  ac.ProfileImageConfig.HtmlTarget,
		HtmlTarget2: ac.ProfileImageConfig.HtmlTarget2,
		Separator:   "src=",
		DestPath:    ac.ProfileImageConfig.DestPath,
		Cookies:     &[]http.Cookie{{Name: "data_user_captcha", Value: "1"}},
	}

	err = startup.SaveImages(&ac, input)
	if err != nil {
		logger.Error("Failed to save images with error:", err)
		os.Exit(1)
	}

	logger.Log("all done!!!")
}
