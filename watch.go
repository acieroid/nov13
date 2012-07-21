package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
	"time"
)

type WatchButton struct {
	x, y int
	enabled bool
	bg, text, grayText *sdl.Surface
	bgPos, textPos *sdl.Rect
	lastStart time.Time
	finished bool
}

func NewWatchButton(x, y int) *WatchButton {
	bg := sdl.CreateRGBSurface(sdl.HWSURFACE, 80, 17, 32, 0, 0, 0, 0)
	bg.FillRect(&sdl.Rect{0, 0, 80, 17}, 0x00123456)
	bgPos := &sdl.Rect{int16(x), int16(y), 0, 0}
	text := ttf.RenderUTF8_Solid(Font, "Regarder !", sdl.Color{255, 255, 255, 0})
	grayText := ttf.RenderUTF8_Solid(Font, "Jouer !", sdl.Color{128, 20, 20, 0})
	textPos := &sdl.Rect{int16(x+2), int16(y+2), 0, 0}
	return &WatchButton{x, y, false, bg, text, grayText, bgPos, textPos, time.Now(), false}
}

func (w *WatchButton) Draw(surf *sdl.Surface) {
	surf.Blit(w.bgPos, w.bg, nil)
	if w.enabled {
		surf.Blit(w.textPos, w.grayText, nil)
	} else {
		surf.Blit(w.textPos, w.text, nil)
	}
}

func (w *WatchButton) Contains(x, y int) bool {
	return x > w.x && x < w.x + 80 && y > w.y && w.y < w.y + 17
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