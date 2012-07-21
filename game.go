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
	userAction UserAction
	lastKey time.Time
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
	g.userAction = nil
	g.lastKey = time.Now()
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
			if g.userAction != nil &&
				int(time.Since(g.lastKey)) > 250e6 {
				quit := g.userAction.KeyPress(e.Keysym.Sym)
				g.userAction = nil
				if quit {
					return MENU
				}
				g.lastKey = time.Now()
			}
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
				if g.userAction == nil {
					if g.menu != nil {
						g.menu = nil
					} else {
						g.userAction = NewQuitUserAction()
						g.lastKey = time.Now()
					}
				}
			}
		case reflect.TypeOf(sdl.MouseButtonEvent{}):
			if g.userAction != nil {
				break
			}

			e := ev.(sdl.MouseButtonEvent)
			if e.Type == sdl.MOUSEBUTTONDOWN && e.Button == 1 {
				x := int(e.X) + g.scrollX
				y := int(e.Y) + g.scrollY
				if g.menu != nil && g.menu.Contains(x, y) {
					g.menu = g.menu.Clicked(x, y)
				} else if g.watchButton.Contains(int(e.X), int(e.Y)) {
					if g.mode == GAME {
						AddMessage("Début du tour")
						g.mode = WATCH
						g.menu = nil
						g.watchButton.Enabled()
						g.lastWatchUpdate = time.Now()
					} else {
						AddMessage("Fin du tour")
						g.watchButton.Finish()
					}
				} else if g.mode == GAME && x < g.m.width * TILESIZE && y < g.m.height * TILESIZE {
					for _, unit := range g.units {
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
			myUnits := 0
			ennemyUnits := 0
			for _, unit := range g.units {
				if unit.team == 1 && unit.Alive() {
					myUnits += 1
				} else if unit.team == 2 && unit.Alive() {
					ennemyUnits += 1
				}
				unit.nextAction = nil
			}
			if myUnits == 0 {
				g.userAction = NewEOGUserAction(false)
			} else if ennemyUnits == 0 {
				g.userAction = NewEOGUserAction(true)
			}
		} else {
			for _, unit := range g.units {
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
	for _, unit := range g.units {
		unit.Draw(g.scrollX, g.scrollY, screen)
	}
	if g.menu != nil {
		g.menu.Draw(g.scrollX, g.scrollY, screen)
	}
	if g.userAction != nil {
		DrawUserAction(g.userAction, screen)
	}
	g.watchButton.Draw(screen)
	return GAME
}
