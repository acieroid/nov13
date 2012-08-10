package main

import (
	"github.com/acieroid/go-sfml"
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

func (g *Game) Run(win sfml.RenderWindow) int {
	e, b := win.PollEvent()
	for b {
		switch e.(type) {
		case sfml.KeyEvent:
			ev := e.(sfml.KeyEvent)
			if g.userAction != nil &&
				int(time.Since(g.lastKey)/1e6) > 250 {
				quit := g.userAction.KeyPress(ev.Code())
				g.userAction = nil
				if quit {
					return MENU
				}
				g.lastKey = time.Now()
			}
			switch ev.Code() {
			case sfml.KeyLeft:
				g.scrollX = Max(g.scrollX-SCROLLSTEP, 0)
			case sfml.KeyRight:
				g.scrollX = Min(g.scrollX+SCROLLSTEP, *Width)
			case sfml.KeyUp:
				g.scrollY = Max(g.scrollY-SCROLLSTEP, 0)
			case sfml.KeyDown:
				g.scrollY = Min(g.scrollY+SCROLLSTEP, *Height)
			case sfml.KeyN:
				g.ScrollToNext()
			case sfml.KeyEscape:
				if g.userAction == nil && int(time.Since(g.lastKey)/1e6) > 250 {
					if g.menu != nil {
						g.menu = nil
					} else {
						g.userAction = NewQuitUserAction()
					}
					g.lastKey = time.Now()
				}
			case sfml.KeyReturn:
				if g.userAction == nil && g.menu == nil && int(time.Since(g.lastKey)/1e6) > 250 {
					if g.mode == GAME {
						if g.AllUnitsGood() {
							g.StartWatch()
						} else {
							g.userAction = NewWatchUserAction(g)
						}
					} else if g.mode == WATCH {
						g.watchButton.Finish()
					}
					g.lastKey = time.Now()
				}
			}
		case sfml.MouseButtonEvent:
			ev := e.(sfml.MouseButtonEvent)
			if ev.Type == sfml.EvtMouseButtonPressed && ev.Button() == sfml.MouseLeft {
				x := ev.X() + g.scrollX
				y := ev.Y() + g.scrollY

				if g.menu != nil && g.menu.Contains(x, y) {
					g.menu = g.menu.Clicked(x, y)
				} else if g.watchButton.Contains(ev.X(), ev.Y()) {
					if g.mode == GAME {
						if g.AllUnitsGood() {
							g.StartWatch()
						} else {
							g.userAction = NewWatchUserAction(g)
						}
					} else {
						AddMessage("Fin du tour")
						g.watchButton.Finish()
					}
				} else if g.mode == GAME && x < g.m.width * TILESIZE && y < g.m.height * TILESIZE {
					for _, unit := range g.units {
						if unit.Alive() && unit.Contains(x, y) && unit.team == 1 {
							g.menu = NewCharacterMenu(unit)
							break
						}
					}
				}
			}
		}
		e, b = win.PollEvent()
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
				if unit.Alive() && unit.nextAction != nil {
					unit.nextAction.Apply(unit,
						g.units,
						g.m,
						int(time.Since(g.lastWatchUpdate)/1e7))
				}
			}
			g.lastWatchUpdate = time.Now()
		}
	}

	g.m.Draw(g.scrollX, g.scrollY, win)
	for _, unit := range g.units {
		unit.Draw(g.scrollX, g.scrollY, win)
	}
	if g.menu != nil {
		g.menu.Draw(g.scrollX, g.scrollY, win)
	}
	if g.userAction != nil {
		DrawUserAction(g.userAction, win)
	}
	g.watchButton.Draw(win)
	return GAME
}

func (g *Game) StartWatch() {
	RunAI(g.units, g.m)
	AddMessage("Début du tour")
	g.mode = WATCH
	g.menu = nil
	g.watchButton.Enabled()
	g.lastWatchUpdate = time.Now()
}

func (g *Game) AllUnitsGood() bool {
	for _, unit := range g.units {
		if unit.Alive() && unit.team == 1 && unit.nextAction == nil {
			return false
		}
	}
	return true
}

func (g *Game) ScrollToNext() {
	var (
		x int
		y int
		found bool
	)
	for _, unit := range g.units {
		if unit.team == 1 && unit.Alive() && unit.nextAction == nil {
			x = unit.x
			y = unit.y
			found = true
		}
	}

	if !found {
		return
	}

	g.scrollX = Max(0, x - *Width/2)
	g.scrollY = Max(0, y - *Height/2)
}