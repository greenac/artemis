package handlers

import (
  "github.com/nsf/termbox-go"
  "github.com/greenac/artemis/logger"
  "fmt"
)

const (
  SpacePerTab = 2
)

type UILine struct {
  x int
  y int
  text []rune
  fgColor termbox.Attribute
  bgColor termbox.Attribute
}

func (l *UILine) addText(t rune, x int) {
  if x == len(l.text) {
    l.text = append(l.text, t)
  } else if x < len(l.text) {
    txt := make([]rune, len(l.text) + 1)
    for i := 0; i < x; i += 1 {
      txt[i] = l.text[i]
    }

    txt[x] = t

    for i := x + 1; i < len(txt); i += 1 {
      txt[i] = l.text[i - 1]
    }

    l.text = txt
  } else {
    debug("Cannot add text to line. x:", x, "is greater than the line length:", len(l.text))
  }
}

func (l *UILine) removeText(x int) {
  if x > len(l.text) {
    debug("Cannot remove text to line. x:", x, "is greater than the line length:", len(l.text))
    return
  }

  txt := make([]rune, len(l.text) - 1)
  for i := 0; i < x - 1; i += 1 {
    txt[i] = l.text[i]
  }

  for i := x; i < len(l.text); i += 1 {
    txt[i - 1] = l.text[i]
  }

  l.text = txt
}

type UIHandler struct {
  lines []*UILine
  line int
  debugLine *UILine
  cursorPosX int
  cursorPosY int
}

func (uih *UIHandler) Setup() {
  ls := make([]*UILine, 0)
  uih.lines = ls
  uih.line = 0
}

func (uih *UIHandler) Run() error {
  err := termbox.Init()
  if err != nil {
    logger.Error("`UIHandler::setup` Failed to initialize termbox:", err)
    return err
  }

  defer termbox.Close()
  termbox.SetInputMode(termbox.InputEsc)

mainloop:
  for {
    switch ev := termbox.PollEvent(); ev.Type {
    case termbox.EventKey:
      switch ev.Key {
      case termbox.KeyEsc, termbox.KeyCtrlC:
        break mainloop
      case termbox.KeyArrowRight, termbox.KeyCtrlF:
        uih.arrowRight()
      case termbox.KeyArrowLeft, termbox.KeyCtrlB:
        uih.arrowLeft()
      case termbox.KeyBackspace, termbox.KeyBackspace2:
        uih.backspace()
      case termbox.KeyEnter:
        uih.handleReturn()
      case termbox.KeyTab:
        for i := 0; i < SpacePerTab; i += 1 {
          uih.handleSpace()
        }
      case termbox.KeySpace:
        uih.handleSpace()
      case termbox.KeyCtrlK:
        //uih.arrowLeft()
      case termbox.KeyHome, termbox.KeyCtrlA:
        //uih.arrowLeft()
      case termbox.KeyEnd, termbox.KeyCtrlE:
        //uih.arrowLeft()
      default:
        uih.keyPress(ev.Ch)
        //if ev.Ch != 0 {
        //	uih.arrowLeft()
        //}
      }
    case termbox.EventError:
      panic(ev.Err)
    }
  }

  return nil
}

func (uih *UIHandler) newLine() *UILine {
  l := UILine{}
  l.fgColor = termbox.ColorCyan
  l.bgColor = termbox.ColorDefault
  l.y = uih.line

  return &l
}

func (uih *UIHandler) currentLine() *UILine {
  var l *UILine
  if len(uih.lines) == 0 {
    l = uih.newLine()
    uih.addLine(l)
  } else {
    l = uih.lines[uih.line]
  }

  return l
}

func (uih *UIHandler) addLine(l *UILine) {
  uih.lines = append(uih.lines, l)
}

func (uih *UIHandler) arrowLeft() {
  if uih.cursorPosX == 0 {
    return
  }

  uih.cursorPosX -= 1
  uih.setCursorPosition()
  termbox.Flush()
}

func (uih *UIHandler) arrowRight() {
  l := uih.currentLine()
  w, _ := termbox.Size()
  if uih.cursorPosX >= w || uih.cursorPosX >= len(l.text) {
    return
  }

  uih.cursorPosX += 1
  uih.setCursorPosition()
  termbox.Flush()
}

func (uih *UIHandler) setCursorPosition() {
  termbox.SetCursor(uih.cursorPosX, uih.cursorPosY)
}

func (uih *UIHandler) handleReturn() {
  // TODO: handle case where return is hit and the cursor is not at the end of lines
  uih.cursorPosX = 0
  uih.cursorPosY += 1
  uih.line += 1
  l := uih.newLine()
  uih.addLine(l)
  uih.setCursorPosition()
  termbox.Flush()
}

func (uih *UIHandler) handleSpace() {
  l := uih.currentLine()
  l.addText(' ', uih.cursorPosX)
  uih.print(l)
  uih.cursorPosX += 1
  uih.setCursorPosition()
  termbox.Flush()
}

func (uih *UIHandler) backspace() {
  if uih.cursorPosX == 0 {
    return
  }

  termbox.Flush()
  l := uih.currentLine()

  for i := uih.cursorPosX - 1; i < len(l.text); i += 1 {
    termbox.SetCell(i, l.y, 0, l.fgColor, l.bgColor)
  }

  l.removeText(uih.cursorPosX)
  uih.print(l)
  uih.cursorPosX -= 1
  uih.setCursorPosition()
  termbox.Flush()
}

func (uih *UIHandler) keyPress(ch rune) {
  l := uih.currentLine()
  l.addText(ch, uih.cursorPosX)
  uih.print(l)
  uih.cursorPosX += 1
  uih.setCursorPosition()
  termbox.Flush()
}

func (uih *UIHandler) print(l *UILine) {
  x := l.x
  for _, c := range l.text {
    termbox.SetCell(x, l.y, c, l.fgColor, l.bgColor)
    x += 1
  }
}

var dStart = 20
func debug(a ...interface{}) {
  msg := fmt.Sprint(a)
  x := 0
	for _, c := range msg {
		termbox.SetCell(x, dStart, c, termbox.ColorYellow, termbox.ColorDefault)
		x++
	}

  dStart += 1
}
