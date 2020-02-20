package ui

import (
	"fmt"
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/logger"
	"github.com/greenac/artemis/models"
	"strings"
)

type AddNamesHandler struct {
	artemisHandler handlers.ArtemisHandler
	uiHandler      Handler
	numUpdated     int
	inputStart     int
	addNames       string
}

func (anh *AddNamesHandler) Setup(
	movDirPaths *[]models.FilePath,
	actDirPaths *[]models.FilePath,
	actFilePath *models.FilePath,
	cachedNamePath *models.FilePath,
	toPath *models.FilePath,
) {
	anh.artemisHandler.Setup(movDirPaths, actDirPaths, actFilePath, cachedNamePath, toPath)
	anh.uiHandler.Setup()
	anh.uiHandler.Kickoff = anh.ShowUnknown
	anh.uiHandler.KeyPress = anh.onKeyPress
	anh.uiHandler.AfterReturn = anh.readInput
	anh.uiHandler.Tab = anh.handleTab
	anh.uiHandler.Escape = anh.AddNamesToMovies
}

func (anh *AddNamesHandler) Run() {
	anh.artemisHandler.Sort()
	err := anh.uiHandler.Run()
	if err != nil {
		logger.Error("`AddNamesHandler::Run` failed with error:", err)
		panic(err)
	}

	anh.ShowUnknown()
}

func (anh *AddNamesHandler) ShowUnknown() {
	anh.uiHandler.ClearAll()
	anh.uiHandler.ClearUI()

	if anh.artemisHandler.MovieHandler.MoreUnknowns() {
		logger.Log("AddNamesHandler::ShowUnknown no unknown movies")
		anh.showComplete()
		return
	}

	m := anh.artemisHandler.MovieHandler.CurrentUnknownMovie()
	txts := []string{"Add name(s) to:", m.Name()}
	anh.uiHandler.SetHeader(txts, true)
	anh.inputStart = anh.uiHandler.NumOfLines(Body) - 1
}

func (anh *AddNamesHandler) showComplete() {
	anh.uiHandler.ClearAll()
	anh.uiHandler.AddToFooter(fmt.Sprintf("Updated %d movies", anh.numUpdated))
	anh.uiHandler.Flush()
}

func (anh *AddNamesHandler) onKeyPress() {
	if anh.uiHandler.Section == Input {
		txt := make([]rune, 0)
		lns := anh.uiHandler.GetLines(Input)
		for _, l := range *lns {
			txt = append(txt, l.Text...)
		}

		txtStr := strings.ToLower(string(txt))
		if txtStr == "y" || txtStr == "yes" {
			m := anh.artemisHandler.MovieHandler.CurrentUnknownMovie()
			anh.artemisHandler.ActorHandler.AddActorsToMovieWithInput(anh.addNames, m)
		}

		anh.artemisHandler.MovieHandler.IncrementUnknownIndex()
		anh.ShowUnknown()
	}
}

func (anh *AddNamesHandler) readInput() {
	if anh.uiHandler.Section == Body {
		txt := make([]rune, 0)
		lns := anh.uiHandler.GetLines(Body)
		for _, l := range *lns {
			txt = append(txt, l.Text...)
		}

		namesTxt := string(txt)
		anh.addNames = namesTxt

		anh.uiHandler.Clear(Body)
		anh.uiHandler.Clear(Input)
		anh.uiHandler.Clear(Footer)
		anh.uiHandler.ClearUI()
		m := anh.artemisHandler.MovieHandler.CurrentUnknownMovie()
		ftTxt := fmt.Sprint("Add name(s) ", string(txt), " to: ", m.Name(), " (Y/N)?")
		anh.uiHandler.AddToBody(ftTxt)
		anh.uiHandler.CursorPosX = 0
		anh.uiHandler.CursorPosY += 1
		anh.uiHandler.SetCursorPosition()
		anh.uiHandler.Section = Input
		anh.uiHandler.ContIndex = 0
		anh.uiHandler.AddToInput("")
		anh.uiHandler.DrawAll()
	} else if anh.uiHandler.Section == Input {
		anh.uiHandler.Clear(Input)
		anh.uiHandler.SiftLines()
		anh.uiHandler.CursorPosX = 0
		anh.uiHandler.AddToFooter(fmt.Sprint("going to write: ", anh.addNames))
		anh.uiHandler.DrawAll()
	}
}

func (anh *AddNamesHandler) handleTab() {
	txt := make([]rune, 0)
	lns := anh.uiHandler.GetLines(Body)
	for _, l := range *lns {
		txt = append(txt, l.Text...)
	}

	pts := strings.Split(string(txt), ",")
	name := strings.ToLower(strings.Trim(pts[len(pts)-1], " "))
	matches, common := anh.artemisHandler.ActorHandler.NameMatches(name)
	names := ""
	for i, actor := range matches {
		n := actor.FullName()
		names += n
		if i < len(matches)-1 {
			names += ", "
		}
	}

	if common != "" && common != name && len(common) > len(name) {
		for i, pt := range pts {
			pts[i] = strings.Trim(pt, " ")
		}

		pts[len(pts)-1] = common
		l := strings.Join(pts, ", ")
		anh.uiHandler.ReplaceLastLine(l, Body)
		anh.uiHandler.CursorPosX = len(l)
		anh.uiHandler.SetCursorPosition()
	}

	anh.uiHandler.Clear(Footer)
	anh.uiHandler.AddToFooter(names)
	anh.uiHandler.DrawAll()
	anh.uiHandler.Flush()
}

func (anh *AddNamesHandler) AddNamesToMovies() {
	anh.artemisHandler.RenameMovies()
	anh.artemisHandler.MoveMovies()
}
