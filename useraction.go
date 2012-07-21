package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
)

type UserAction interface {
	KeyPress(sym uint32) bool
	Text() string
}

var bgSurf *sdl.Surface

func DrawUserAction(a UserAction, surf *sdl.Surface) {
	if bgSurf == nil {
		bgSurf = sdl.CreateRGBSurface(sdl.HWSURFACE,
			int(surf.W), 24, 32, 0, 0, 0, 0)
		bgSurf.FillRect(&sdl.Rect{0, 0, uint16(surf.W), 24},
			0x00FF0000)
		bgSurf.SetAlpha(sdl.SRCALPHA, 150)
	}
	surf.Blit(&sdl.Rect{0, int16(surf.H/2 - 12 + 50), 0, 0},
		bgSurf, nil)
	DrawTextBig(a.Text(), int(surf.W/2 - 100),
		int(surf.H/2 - 10 + 50), surf)
}

type QuitUserAction struct {
}

func NewQuitUserAction() *QuitUserAction {
	return &QuitUserAction{}
}

func (a *QuitUserAction) KeyPress(sym uint32) bool {
	return sym == sdl.K_ESCAPE
}

func (a *QuitUserAction) Text() string {
	return "Appuyez (encore) sur ESC pour quitter"
}

type EOGUserAction struct {
	won bool
}

func NewEOGUserAction(won bool) *EOGUserAction {
	return &EOGUserAction{won}
}

func (a *EOGUserAction) KeyPress(sym uint32) bool {
	return sym == sdl.K_ESCAPE
}

func (a *EOGUserAction) Text() string {
	if a.won {
		return "Vous avez gagné, appuyez sur ESC pour revenir au menu"
	}
	return "Vous avez perdu, appuyez sur ESC pour revenir au menu"
}
