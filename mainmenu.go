package main

import (
	"container/ring"
	"fmt"
	"github.com/acieroid/go-sfml"
	"os"
	"strings"
	"time"
)

type MainMenu struct {
	maps      *ring.Ring
	lastModif time.Time
	w, h      int
}

func NewMainMenu(mapDir string, w, h int) (m *MainMenu) {
	file, err := os.Open(mapDir)
	defer file.Close()
	if err != nil {
		panic(err)
	}

	fis, err := file.Readdir(0)
	if err != nil {
		panic(err)
	}

	var maps *ring.Ring
	for _, fi := range fis {
		if strings.HasSuffix(fi.Name(), ".txt") {
			name := fi.Name()
			mapName := name[:len(name)-4]
			r := ring.New(1)
			r.Value = mapName
			if maps != nil {
				maps.Link(r)
			} else {
				maps = r
			}
		}
	}
	return &MainMenu{maps, time.Now(), w, h}
}

func (m *MainMenu) Run(win sfml.RenderWindow) (string, int) {
	e, b := win.PollEvent()
	if b {
		switch e.(type) {
		case sfml.KeyEvent:
			ev := e.(sfml.KeyEvent)
			switch ev.Code() {
			case sfml.KeyLeft:
				if int(time.Since(m.lastModif)/1e6) > 250 {
					m.maps = m.maps.Move(1)
					m.lastModif = time.Now()
				}
			case sfml.KeyRight:
				if int(time.Since(m.lastModif)/1e6) > 250 {
					m.maps = m.maps.Move(-1)
					m.lastModif = time.Now()
				}
			case sfml.KeyReturn:
				return m.maps.Value.(string), GAME
			case sfml.KeyEscape:
				if int(time.Since(m.lastModif)/1e6) > 250 {
					return "", QUIT
				}
			}
		}
	}

	DrawTextBig(fmt.Sprintf("< Map: %s >", m.maps.Value.(string)),
		m.w/2, m.h/2,
		true, win)

	return "", MENU
}
