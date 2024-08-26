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
		RestClient: &rc,
		BaseUrl:    "https://www.data18.com/name/",
		HtmlTarget: "https://cdn.dt18.com/images/names",
		Separator:  "src=",
		DestPath:   "/Users/andre/Documents/artemis/profile-pics/",
		Cookies:    &[]http.Cookie{{Name: "data_user_captcha", Value: "1"}},
	}

	err = startup.SaveImages(&ac, input)
	if err != nil {
		logger.Error("Failed to save images with error:", err)
		os.Exit(1)
	}

	logger.Log("all done!!!")
}
