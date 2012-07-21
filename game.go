package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"reflect"
	"time"
)

type Game struct {
	scrollX, scrollY int
	m *Map
	units []*Character
	mode int
	watchButton *WatchButton
	lastUpdate, lastWatchUpdate time.Time
	menu Menu
}

func NewGame(mapName string, w, h int) (g *Game) {
	g = &Game{}
	g.scrollX = 0
	g.scrollY = 0
	g.m, g.units = LoadMap(mapName)
	g.menu = nil
	g.mode = GAME
	g.watchButton = NewWatchButton(w-100, h-30)
	g.lastUpdate = time.Now()
	g.lastWatchUpdate = time.Now()
	return
}

func (g *Game) Run(screen *sdl.Surface) int {
	select {
	case ev := <-sdl.Events:
		/* TODO: something more clean for the cases ? */
		switch reflect.TypeOf(ev) {
		case reflect.TypeOf(sdl.QuitEvent{}):
			return QUIT
		case reflect.TypeOf(sdl.KeyboardEvent{}):
			e := ev.(sdl.KeyboardEvent)
			switch e.Keysym.Sym {
			case sdl.K_LEFT:
				g.scrollX = Max(g.scrollX-SCROLLSTEP, 0)
			case sdl.K_RIGHT:
				g.scrollX = Min(g.scrollX+SCROLLSTEP, *Width)
			case sdl.K_UP:
				g.scrollY = Max(g.scrollY-SCROLLSTEP, 0)
			case sdl.K_DOWN:
				g.scrollY = Min(g.scrollY+SCROLLSTEP, *Height)
			case sdl.K_ESCAPE:
				if g.menu != nil {
					g.menu = nil
				}
				/* TODO: else, propose to quit (w/ some msg like "Press ESC again to quit) */
			}
		case reflect.TypeOf(sdl.MouseButtonEvent{}):
			e := ev.(sdl.MouseButtonEvent)
			if e.Type == sdl.MOUSEBUTTONDOWN && e.Button == 1 {
				x := int(e.X) + g.scrollX
				y := int(e.Y) + g.scrollY
				if g.menu != nil && g.menu.Contains(x, y) {
					g.menu = g.menu.Clicked(x, y)
				} else if g.mode == GAME && g.watchButton.Contains(x, y) {
					AddMessage("Début du tour")
					g.mode = WATCH
					g.menu = nil
					g.watchButton.Enabled()
					g.lastWatchUpdate = time.Now()
				} else if g.mode == GAME && x < g.m.width * TILESIZE && y < g.m.height * TILESIZE {
					for _, unit := range(g.units) {
						if unit.Contains(x, y) && unit.team == 1 {
							g.menu = NewCharacterMenu(unit)
							break
						}
					}
				}
			}
		}
	default:
	}
	if g.mode == WATCH {
		if g.watchButton.WatchFinished() {
			AddMessage("Fin du tour")
			g.mode = GAME
			g.watchButton.Disabled()
			for _, unit := range(g.units) {
				unit.nextAction = nil
			}
		} else {
			for _, unit := range(g.units) {
				if unit.nextAction != nil {
					unit.nextAction.Apply(unit,
						g.units,
						g.m,
						int(time.Since(g.lastWatchUpdate)/1e7))
				}
			}
			g.lastWatchUpdate = time.Now()
		}
	}


	g.m.Draw(g.scrollX, g.scrollY, screen)
	for _, unit := range(g.units) {
		unit.Draw(g.scrollX, g.scrollY, screen)
	}
	if g.menu != nil {
		g.menu.Draw(g.scrollX, g.scrollY, screen)
	}
	g.watchButton.Draw(screen)
	return GAME
}
