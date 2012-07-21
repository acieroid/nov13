package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"container/ring"
	"reflect"
	"fmt"
	"os"
	"strings"
	"time"
)

type MainMenu struct {
	maps *ring.Ring
	lastModif time.Time
	w, h int
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
	for _, fi := range(fis) {
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

func (m *MainMenu) Run(screen *sdl.Surface) (string, int) {
	select {
	case ev := <-sdl.Events:
		switch reflect.TypeOf(ev) {
		case reflect.TypeOf(sdl.QuitEvent{}):
			return "", QUIT
		case reflect.TypeOf(sdl.KeyboardEvent{}):
			e := ev.(sdl.KeyboardEvent)
			switch e.Keysym.Sym {
			case sdl.K_LEFT:
				if int(time.Since(m.lastModif)) > 250e6 {
					m.maps = m.maps.Move(1)
					m.lastModif = time.Now()
				}
			case sdl.K_RIGHT:
				if int(time.Since(m.lastModif)) > 250e6 {
					m.maps = m.maps.Move(-1)
					m.lastModif = time.Now()
				}
			case sdl.K_RETURN:
				return m.maps.Value.(string), GAME
			case sdl.K_ESCAPE:
				return "", QUIT
			}
		}
	}

	DrawTextBig(fmt.Sprintf("< Map: %s >", m.maps.Value.(string)),
		/* TODO: use font metrics, again. */
		m.w/2 - 100, m.h/2 - 10,
		screen)

	return "", MENU
}
