package handlers

import (
	"fmt"
	"github.com/greenac/artemis/pkg/clients/rest"
	"github.com/greenac/artemis/pkg/dbinteractors"
	"github.com/greenac/artemis/pkg/errs"
	"github.com/greenac/artemis/pkg/logger"
	"github.com/greenac/artemis/pkg/models"
	"github.com/greenac/artemis/pkg/worker"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type SaveImageInput struct {
	RestClient  rest.IClient
	BaseUrl     string
	BaseUrl2    string
	SubBaseUrl2 string
	HtmlTarget  string
	HtmlTarget2 string
	Separator   string
	DestPath    string
	Cookies     *[]http.Cookie
}

type saveImageResult struct {
	Actor models.Actor
	Error error
}

func SaveImages(input SaveImageInput) error {
	acts, err := dbinteractors.AllActors()
	if err != nil {
		return err
	}

	actors := *acts
	ch := make(chan worker.Message[saveImageResult], len(actors))
	wrkr := worker.NewWorker(5, ch)
	wrkr.Work()

	go func() {
		for _, a := range actors {
			wrkr.AddTask(saveImageTask(input, a))
		}
	}()

	var erroredResults []saveImageResult
	count := 0
	for count < len(actors) {
		m := <-ch
		if m.Result.Error != nil {
			logger.Error(fmt.Sprintf("SaveImages->failed for actor: %s with error: %s", m.Result.Actor.FullName(), m.Result.Error.Error()))
			erroredResults = append(erroredResults, m.Result)
		} else {
			logger.Log(fmt.Sprintf("SaveImages->succeeded for actor: %s", m.Result.Actor.FullName()))
		}

		count += 1
		logger.Log("jobs completed:", count, "/", len(actors))
	}

	for _, r := range erroredResults {
		logger.Log("failed actor:", r.Actor.FullName(), "with id:", r.Actor.Id.Hex(), "has error:", r.Error.Error())
	}

	return nil
}

func saveImageTask(input SaveImageInput, actor models.Actor) func() saveImageResult {
	return func() saveImageResult {
		r := scrapeTaskSite2(input, actor)
		if r.Error == nil {
			return r
		}

		logger.Warn("saveImageTask->failed to make request to site1 for actor:", actor.FullName(), "with error:", r.Error.Error())

		return scrapeTaskSite1(input, actor)
	}
}

func scrapeTaskSite1(input SaveImageInput, actor models.Actor) saveImageResult {
	htmlUrl, err := url.JoinPath(input.BaseUrl, actor.FullNameWithHyphens())
	if err != nil {
		logger.Error("saveImageTask->failed join base url and name with error:", err)
		return saveImageResult{Actor: actor, Error: err}
	}

	res, err := input.RestClient.Get(htmlUrl, nil, nil, input.Cookies)
	if err != nil {
		logger.Error("saveImageTask->failed to get html with error:", err)
		return saveImageResult{Actor: actor, Error: err}
	}

	if res.StatusCode != http.StatusOK {
		logger.Error("saveImageTask->bad response code getting html:", res.StatusCode, "for url:", htmlUrl)
		return saveImageResult{Actor: actor, Error: errs.NewGenError("bad response code getting html")}
	}

	lines := strings.Split(string(res.Body), "\n")

	var targetUrl string
	for _, l := range lines {
		if strings.Contains(l, input.HtmlTarget) {
			parts := strings.Split(l, " ")
			for _, p := range parts {
				if strings.Contains(p, input.Separator) {
					targetUrl = strings.Replace(strings.ReplaceAll(p, `"`, ""), input.Separator, "", 1)
				}
			}
		}
	}

	if targetUrl == "" {
		logger.Error("saveImageTask->failed to find target for actor:", actor.FullName())
		return saveImageResult{Actor: actor, Error: errs.NewGenError("no target found")}
	}

	res, err = input.RestClient.Get(targetUrl, nil, nil, input.Cookies)
	if err != nil {
		logger.Error("saveImageTask->failed to get image for actor", actor.FullName(), "with error:", err)
		return saveImageResult{Actor: actor, Error: err}
	}

	if res.StatusCode != http.StatusOK {
		logger.Error("saveImageTask->bad response code getting image:", res.StatusCode)
		return saveImageResult{Actor: actor, Error: errs.NewGenError("bad response code getting html")}
	}

	filePath, e := url.JoinPath(input.DestPath, actor.Id.Hex()+".jpg")
	if e != nil {
		logger.Error("saveImageTask->failed to join image path with error:", err)
		return saveImageResult{Actor: actor, Error: err}
	}

	e = os.WriteFile(filePath, res.Body, 0644)
	if e != nil {
		logger.Error("saveImageTask->failed to write image to file with error:", e)
		return saveImageResult{Actor: actor, Error: errs.NewGenError(e.Error())}
	}

	return saveImageResult{Actor: actor}
}

func scrapeTaskSite2(input SaveImageInput, actor models.Actor) saveImageResult {
	actorName := strings.ReplaceAll(actor.FullName(), "_", "")
	htmlUrl, err := url.JoinPath(input.BaseUrl2, actorName+".html")
	if err != nil {
		logger.Error("scrapeTaskSite2->failed join base url and name with error:", err)
		return saveImageResult{Actor: actor, Error: err}
	}

	res, err := input.RestClient.Get(htmlUrl, nil, nil, nil)
	if err != nil {
		logger.Error("scrapeTaskSite2->failed to get html with error:", err)
		return saveImageResult{Actor: actor, Error: err}
	}

	if res.StatusCode != http.StatusOK {
		logger.Error("scrapeTaskSite2->bad response code getting html:", res.StatusCode, "for url:", htmlUrl)
		return saveImageResult{Actor: actor, Error: errs.NewGenError("bad response code getting html")}
	}

	lines := strings.Split(string(res.Body), "\n")

	var targetUrl string
	for _, l := range lines {
		if strings.Contains(l, input.HtmlTarget2) {
			parts := strings.Split(l, " ")
			for _, p := range parts {
				if strings.Contains(p, input.Separator) {
					targetUrl = strings.Replace(strings.ReplaceAll(p, `"`, ""), input.Separator, "", 1)
				}
			}
		}
	}

	if targetUrl == "" {
		logger.Error("scrapeTaskSite2->failed to find target for actor:", actor.FullName())
		return saveImageResult{Actor: actor, Error: errs.NewGenError("no target found")}
	}

	targetUrl, e := url.JoinPath(input.SubBaseUrl2, targetUrl)
	if e != nil {
		logger.Error("scrapeTaskSite2->failed to make target url for actor:", actor.FullName())
		return saveImageResult{Actor: actor, Error: errs.NewGenError("making target url failed")}
	}

	logger.Log("scrapeTaskSite2->making request to:", targetUrl)

	res, err = input.RestClient.Get(targetUrl, nil, nil, input.Cookies)
	if err != nil {
		logger.Error("scrapeTaskSite2->failed to get image for actor", actor.FullName(), "with error:", err)
		return saveImageResult{Actor: actor, Error: err}
	}

	if res.StatusCode != http.StatusOK {
		logger.Error("scrapeTaskSite2->bad response code getting image:", res.StatusCode)
		return saveImageResult{Actor: actor, Error: errs.NewGenError("bad response code getting html")}
	}

	filePath, e := url.JoinPath(input.DestPath, actor.Id.Hex()+".jpg")
	if e != nil {
		logger.Error("scrapeTaskSite2->failed to join image path with error:", err)
		return saveImageResult{Actor: actor, Error: err}
	}

	e = os.WriteFile(filePath, res.Body, 0644)
	if e != nil {
		logger.Error("scrapeTaskSite2->failed to write image to file with error:", e)
		return saveImageResult{Actor: actor, Error: errs.NewGenError(e.Error())}
	}

	return saveImageResult{Actor: actor}
}
