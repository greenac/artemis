package ui

import (
	"fmt"
	"github.com/greenac/artemis/handlers"
	"github.com/greenac/artemis/tools"
	"strings"
)

type AddNamesHandler struct {
	artemisHandler handlers.ArtemisHandler
	uiHandler      Handler
	unkIndex       int
	numUpdated     int
	inputStart     int
}

func (anh *AddNamesHandler) Setup(movDirPaths *[]tools.FilePath, actDirPaths *[]tools.FilePath, actFilePaths *[]tools.FilePath) {
	anh.artemisHandler.Setup(movDirPaths, actDirPaths, actFilePaths)
	anh.uiHandler.Setup()
	anh.uiHandler.Kickoff = anh.ShowUnknown
	anh.uiHandler.KeyPress = anh.onKeyPress
	anh.uiHandler.AfterReturn = anh.readInput
	anh.uiHandler.Tab = anh.handleTab
}

func (anh *AddNamesHandler) Run() {
	anh.artemisHandler.Sort()
	anh.uiHandler.Run()
	anh.ShowUnknown()
}

func (anh *AddNamesHandler) ShowUnknown() {
	if anh.unkIndex == len(anh.artemisHandler.UnknownMovies) {
		anh.showComplete()
		return
	}

	m := anh.artemisHandler.UnknownMovies[anh.unkIndex]
	txts := []string{"Add name(s) to:", *m.Name()}
	anh.uiHandler.SetHeader(txts, true)
	anh.inputStart = anh.uiHandler.NumOfLines(Body) - 1
}

func (anh *AddNamesHandler) showComplete() {
	anh.uiHandler.ClearAll()
	anh.uiHandler.AddToFooter(fmt.Sprintf("Updated %d movies", anh.numUpdated))
	anh.uiHandler.Flush()
}

func (anh *AddNamesHandler) onKeyPress() {
  txt := make([]rune, 0)
  lns := anh.uiHandler.GetLines(Body)
  for _, l := range *lns {
    txt = append(txt, l.Text...)
  }

  matches := anh.artemisHandler.ActorHandler.Matches(string(txt))
  acts := ""
  for i, a := range matches {
    acts += a.FullName()
    if i < len(matches) - 1 {
      acts += ", "
    }
  }

  anh.uiHandler.Draw(Footer)
}

func (anh *AddNamesHandler) readInput() {
	txt := make([]rune, 0)
	lns := anh.uiHandler.GetLines(Body)
	for _, l := range *lns {
		txt = append(txt, l.Text...)
	}

	anh.uiHandler.Clear(Body)
	anh.uiHandler.Clear(Footer)
	anh.uiHandler.Section = Input
	//anh.uiHandler.AddToInput(fmt.Sprint("Add value:", " some value"))
	anh.uiHandler.AddToFooter(fmt.Sprint("got text: ", string(txt)))
	anh.uiHandler.DrawAll()
}

func (anh *AddNamesHandler) handleTab() {
	txt := make([]rune, 0)
	lns := anh.uiHandler.GetLines(Body)
	for _, l := range *lns {
		txt = append(txt, l.Text...)
	}

	pts := strings.Split(string(txt), ",")
	name := strings.Trim(pts[len(pts)-1], " ")
	matches := anh.artemisHandler.ActorHandler.Matches(name)
	names := ""
	for i, actor := range matches {
		names += actor.FullName()
		if i < len(matches)-1 {
			names += ", "
		}
	}

	anh.uiHandler.Clear(Footer)
	anh.uiHandler.AddToFooter(names)
	anh.uiHandler.ClearUI()
	anh.uiHandler.DrawAll()
	anh.uiHandler.Flush()
}
