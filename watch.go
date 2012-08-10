package main

import (
	"github.com/acieroid/go-sfml"
	"time"
)

type WatchButton struct {
	x, y           int
	enabled        bool
	bg             sfml.RectangleShape
	text, grayText sfml.Text
	lastStart      time.Time
	finished       bool
}

func NewWatchButton(x, y int) *WatchButton {
	bg := sfml.NewRectangleShape()
	bg.SetSize(80, 17)
	bg.SetFillColor(sfml.FromRGB(012, 123, 234))
	bg.SetPosition(float32(x), float32(y))

	text, err := sfml.NewText()
	if err != nil {
		panic(err)
	}
	text.SetFont(Font)
	text.SetString("Regarder !")
	text.SetColor(sfml.FromRGB(255, 255, 255))
	text.SetPosition(float32(x+2), float32(y+2))
	text.SetCharacterSize(12)

	grayText, err := sfml.NewText()
	if err != nil {
		panic(err)
	}
	grayText.SetFont(Font)
	grayText.SetString("Jouer !")
	grayText.SetColor(sfml.FromRGB(128, 20, 20))
	grayText.SetPosition(float32(x+2), float32(y+2))
	grayText.SetCharacterSize(12)
	return &WatchButton{x, y, false, bg, text, grayText, time.Now(), false}
}

func (w *WatchButton) Draw(win sfml.RenderWindow) {
	win.DrawRectangleShapeDefault(w.bg)
	if w.enabled {
		win.DrawTextDefault(w.grayText)
	} else {
		win.DrawTextDefault(w.text)
	}
}

func (w *WatchButton) Contains(x, y int) bool {
	return x > w.x && x < w.x+80 && y > w.y && w.y < w.y+17
}

func (w *WatchButton) Enabled() {
	w.enabled = true
	w.finished = false
	w.lastStart = time.Now()
}

func (w *WatchButton) Disabled() {
	w.enabled = false
}

func (w *WatchButton) Finish() {
	w.finished = true
}

func (w *WatchButton) WatchFinished() bool {
	return w.finished || int64(time.Since(w.lastStart)) > WATCHTIME*1e9
}