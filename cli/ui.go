package ui

import (
	"fmt"
	"github.com/greenac/artemis/logger"
	"github.com/nsf/termbox-go"
)

const (
	SpacePerTab = 2
)

type Container string

const (
	Header Container = "header"
	Body   Container = "body"
	Input  Container = "input"
	Footer Container = "footer"
)

type Handler struct {
	lines       map[Container][]*Line
	positions   map[Container]int
	debugLine   *Line
	CursorPosX  int
	CursorPosY  int
	ContIndex   int
	run         bool
	Section     Container
	Kickoff     func()
	AfterReturn func()
	Tab         func()
	KeyPress    func()
	Escape      func()
}

func (uih *Handler) Setup() {
	lines := make(map[Container][]*Line, 3)
	lines[Header] = make([]*Line, 0)
	lines[Body] = make([]*Line, 0)
	lines[Input] = make([]*Line, 0)
	lines[Footer] = make([]*Line, 0)

	uih.lines = lines
	uih.CursorPosY = 0
	uih.CursorPosX = 0
	uih.Section = Body
	uih.run = true
}

func (uih *Handler) GetLines(c Container) *[]*Line {
	ls := uih.lines[c]
	return &ls
}

func (uih *Handler) NumOfLines(c Container) int {
	return len(*uih.GetLines(c))
}

func (uih *Handler) Run() error {
	err := termbox.Init()
	if err != nil {
		logger.Error("`Handler::setup` Failed to initialize termbox:", err)
		return err
	}

	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	if uih.Kickoff != nil {
		uih.Kickoff()
	}

mainloop:
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				if uih.Escape != nil {
					uih.Escape()
				}

				break mainloop
			case termbox.KeyCtrlC:
				uih.run = false
				break mainloop
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				uih.arrowRight()
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				uih.arrowLeft()
			case termbox.KeyArrowUp:
				uih.arrowUp()
			case termbox.KeyArrowDown:
				uih.arrowDown()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				uih.backspace()
			case termbox.KeyEnter:
				uih.handleReturn()
			case termbox.KeyTab:
				uih.tab()
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
			}
		case termbox.EventError:
			logger.Error("Failed to run UI Handler with error:", ev.Err)
			panic(ev.Err)
		}

		if !uih.run {
			//break mainloop
		}
	}

	return nil
}

func (uih *Handler) newLine(c Container) *Line {
	var clr termbox.Attribute
	switch c {
	case Header:
		clr = termbox.ColorCyan
	case Body:
		clr = termbox.ColorWhite
	case Input:
		clr = termbox.ColorYellow
	case Footer:
		clr = termbox.ColorGreen
	}

	l := Line{}
	l.FgColor = clr
	l.BgColor = termbox.ColorDefault

	return &l
}

func (uih *Handler) currentLine() *Line {
	ls := uih.GetLines(uih.Section)
	var l *Line
	if len(*ls) == 0 {
		l = uih.newLine(uih.Section)
		l.Y = uih.CursorPosY
		uih.addLine(l, uih.Section)
	} else {
		l = (*ls)[uih.ContIndex]
	}

	return l
}

func (uih *Handler) addLine(l *Line, c Container) {
	ls := uih.GetLines(c)
	*ls = append(*ls, l)
	uih.lines[c] = *ls
}

func (uih *Handler) AddToHeader(txt string) *Line {
	l := uih.newLine(Header)
	l.Text = []rune(txt)
	uih.addLine(l, Header)

	return l
}

func (uih *Handler) AddToBody(txt string) *Line {
	l := uih.newLine(Body)
	l.Y = uih.CursorPosY
	l.Text = []rune(txt)
	uih.addLine(l, Body)

	return l
}

func (uih *Handler) AddToInput(txt string) *Line {
	l := uih.newLine(Input)
	l.Y = uih.CursorPosY
	l.Text = []rune(txt)
	uih.addLine(l, Input)

	return l
}

func (uih *Handler) AddToFooter(txt string) *Line {
	l := uih.newLine(Footer)
	l.Text = []rune(txt)
	uih.addLine(l, Footer)

	return l
}

func (uih *Handler) arrowLeft() {
	if uih.CursorPosX == 0 {
		return
	}

	uih.CursorPosX -= 1
	uih.SetCursorPosition()
	uih.Flush()
}

func (uih *Handler) arrowRight() {
	l := uih.currentLine()
	w, _ := termbox.Size()
	if uih.CursorPosX >= w || uih.CursorPosX >= len(l.Text) {
		return
	}

	uih.CursorPosX += 1
	uih.SetCursorPosition()
	uih.Flush()
}

func (uih *Handler) arrowUp() {
	ls := uih.GetLines(Header)
	if uih.CursorPosY <= len(*ls) {
		return
	}

	uih.CursorPosY -= 1
	uih.ContIndex -= 1
	l := uih.currentLine()
	if uih.CursorPosX > len(l.Text) {
		uih.CursorPosX = len(l.Text)
	}

	uih.SetCursorPosition()
	uih.Flush()
}

func (uih *Handler) arrowDown() {
	ls := uih.GetLines(Body)
	if uih.CursorPosY >= len(*ls)-1 {
		return
	}

	uih.CursorPosY += 1
	uih.ContIndex += 1
	l := uih.currentLine()
	if uih.CursorPosX > len(l.Text) {
		uih.CursorPosX = len(l.Text)
	}

	uih.SetCursorPosition()
	uih.Flush()
}

func (uih *Handler) SetCursorPosition() {
	termbox.SetCursor(uih.CursorPosX, uih.CursorPosY)
}

func (uih *Handler) handleReturn() {
	// TODO: handle case return is hit and the cursor is not at the end of lines
	if uih.AfterReturn != nil {
		uih.AfterReturn()
	} else {
		uih.AddBlankLine(Body)
	}
}

func (uih *Handler) AddBlankLine(c Container) {
	uih.CursorPosX = 0
	uih.CursorPosY += 1
	uih.ContIndex += 1

	l := uih.newLine(c)
	l.Y = uih.CursorPosY
	uih.addLine(l, c)
	uih.SetCursorPosition()
	uih.Flush()
}

func (uih *Handler) handleSpace() {
	l := uih.currentLine()
	l.addText(' ', uih.CursorPosX)
	uih.Print(l)
	uih.CursorPosX += 1
	uih.SetCursorPosition()
	uih.Flush()
}

func (uih *Handler) tab() {
	if uih.Tab != nil {
		uih.Tab()
		uih.Flush()
		return
	}

	for i := 0; i < SpacePerTab; i += 1 {
		uih.handleSpace()
	}
}

func (uih *Handler) backspace() {
	if uih.CursorPosX == 0 {
		return
	}

	uih.Flush()
	l := uih.currentLine()

	for i := uih.CursorPosX - 1; i < len(l.Text); i += 1 {
		termbox.SetCell(i, l.Y, 0, termbox.ColorDefault, termbox.ColorDefault)
	}

	l.removeText(uih.CursorPosX)
	uih.Print(l)
	uih.CursorPosX -= 1
	uih.SetCursorPosition()
	uih.Flush()
}

func (uih *Handler) keyPress(ch rune) {
	l := uih.currentLine()
	l.addText(ch, uih.CursorPosX)
	uih.Print(l)
	uih.CursorPosX += 1
	uih.SetCursorPosition()
	uih.SiftLines()

	if uih.KeyPress != nil {
		uih.KeyPress()
	}

	uih.Flush()
}

func (uih *Handler) Print(l *Line) {
	x := l.X
	for _, c := range l.Text {
		termbox.SetCell(x, l.Y, c, l.FgColor, l.BgColor)
		x += 1
	}
}

func (uih *Handler) DrawAll() {
	uih.ClearUI()
	uih.SiftLines()
	for _, cnt := range uih.orderedContainers() {
		for _, l := range *uih.GetLines(cnt) {
			uih.Print(l)
		}
	}

	uih.Flush()
}

func (uih *Handler) Draw(c Container) {
	uih.SiftLines()
	for _, l := range *uih.GetLines(c) {
		uih.Print(l)
	}

	uih.Flush()
}

func (uih *Handler) orderedContainers() []Container {
	return []Container{Header, Body, Input, Footer}
}

func (uih *Handler) SiftLines() {
	y := 0
	for _, c := range uih.orderedContainers() {
		for _, l := range *(uih.GetLines(c)) {
			l.Y = y
			y += 1
		}
	}
}

func (uih *Handler) Flush() {
	if !uih.run {
		return
	}

	termbox.Flush()
}

func (uih *Handler) Clear(c Container) {
	ls := uih.GetLines(c)
	for _, l := range *ls {
		l.Text = make([]rune, 0)
		uih.Print(l)
	}

	if c == uih.Section {
		uih.ContIndex = 0
	}

	uih.lines[c] = make([]*Line, 0)
}

func (uih *Handler) ClearAll() {
	uih.Setup()
	uih.lines = make(map[Container][]*Line)

	uih.ClearUI()
}

func (uih *Handler) ClearUI() {
	err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if err != nil {
		logger.Warn("`Handler::ClearUI` failed with error:", err)
	}
}

func (uih *Handler) SetHeader(txts []string, updateCursor bool) {
	ls := uih.GetLines(Header)
	for i, t := range txts {
		var l *Line
		if i < len(*ls) {
			l = (*ls)[i]
		} else {
			l = uih.newLine(Header)
			uih.addLine(l, Header)
		}

		l.Y = i
		l.Text = []rune(t)
		uih.Print(l)
	}

	if updateCursor {
		uih.CursorPosY = len(txts)
		uih.AddBlankLine(Header)
		uih.AddBlankLine(Header)
		uih.ContIndex = 0
	}

	uih.SiftLines()
	uih.Flush()
}

func (uih *Handler) ReplaceLastLine(txt string, c Container) {
	ls := uih.GetLines(c)
	if len(*ls) == 0 {
		return
	}

	l := (*ls)[len(*ls)-1]
	l.Text = []rune(txt)
}

func (uih *Handler) Debug(a ...interface{}) {
	debug(a)
}

var dStart = 15

func debug(a ...interface{}) {
	msg := fmt.Sprint(a)
	x := 0
	for _, c := range msg {
		termbox.SetCell(x, dStart, c, termbox.ColorYellow, termbox.ColorDefault)
		x++
	}

	dStart += 1
}
