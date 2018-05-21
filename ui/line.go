package ui

import (
	"github.com/nsf/termbox-go"
)

type Line struct {
	X       int
	Y       int
	Text    []rune
	FgColor termbox.Attribute
	BgColor termbox.Attribute
}

func (l *Line) addText(t rune, x int) {
	if x == len(l.Text) {
		l.Text = append(l.Text, t)
	} else if x < len(l.Text) {
		txt := make([]rune, len(l.Text)+1)
		for i := 0; i < x; i += 1 {
			txt[i] = l.Text[i]
		}

		txt[x] = t

		for i := x + 1; i < len(txt); i += 1 {
			txt[i] = l.Text[i-1]
		}

		l.Text = txt
	} else {
		debug("Cannot add text to line. x:", x, "is greater than the line length:", len(l.Text))
	}
}

func (l *Line) removeText(x int) {
	if x > len(l.Text) {
		debug("Cannot remove text to line. x:", x, "is greater than the line length:", len(l.Text))
		return
	}

	txt := make([]rune, len(l.Text)-1)
	for i := 0; i < x-1; i += 1 {
		txt[i] = l.Text[i]
	}

	for i := x; i < len(l.Text); i += 1 {
		txt[i-1] = l.Text[i]
	}

	l.Text = txt
}
