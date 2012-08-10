package main

import (
	"github.com/acieroid/go-sfml"
)

type UserAction interface {
	KeyPress(code sfml.KeyCode) bool
	Text() string
}

var Bg sfml.RectangleShape

func DrawUserAction(a UserAction, win sfml.RenderWindow) {
	w, h := win.Size()
	if Bg.Cref == nil {
		Bg = sfml.NewRectangleShape()
		Bg.SetSize(float32(w), 24)
		Bg.SetFillColor(sfml.FromRGBA(255, 0, 0, 150))
		Bg.SetPosition(0, float32(h/2-12+50))
	}
	win.DrawRectangleShapeDefault(Bg)
	DrawTextBig(a.Text(), int(w)/2, int(h)/2+50, true, win)
}

type QuitUserAction struct {
}

func NewQuitUserAction() *QuitUserAction {
	return &QuitUserAction{}
}

func (a *QuitUserAction) KeyPress(code sfml.KeyCode) bool {
	return code == sfml.KeyEscape
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

func (a *EOGUserAction) KeyPress(code sfml.KeyCode) bool {
	return code == sfml.KeyEscape
}

func (a *EOGUserAction) Text() string {
	if a.won {
		return "Vous avez gagné, appuyez sur ESC pour revenir au menu"
	}
	return "Vous avez perdu, appuyez sur ESC pour revenir au menu"
}

type WatchUserAction struct {
	g *Game
}

func NewWatchUserAction(g *Game) *WatchUserAction {
	return &WatchUserAction{g}
}

func (a *WatchUserAction) KeyPress(code sfml.KeyCode) bool {
	if code == sfml.KeyReturn {
		a.g.StartWatch()
	}
	return false
}

func (a *WatchUserAction) Text() string {
	return "Certaines unités n'ont pas d'action assignée, appuyez sur ENTER pour continuer"
}
