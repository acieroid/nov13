package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"time"
	"reflect"
	"flag"
	"fmt"
)

var SCROLLSTEP int = 5
var MapName = flag.String("map", "foo", "Map to play")
var Width = flag.Int("width", 640, "Width of the window")
var Height = flag.Int("height", 480, "Height of the window")
var Fullscreen = flag.Bool("fullscreen", false, "Fullscreen")

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func LoadImage(name string) *sdl.Surface {
	image := sdl.Load(name)

	if image == nil {
		panic(sdl.GetError())
	}

	return image
}

func main() {
	flag.Parse()
	if sdl.Init(sdl.INIT_VIDEO) != 0 {
		panic(sdl.GetError())
	}

	defer sdl.Quit()

	var videoMode uint32 = 0
	if *Fullscreen {
		videoMode = sdl.FULLSCREEN
	}
	var screen = sdl.SetVideoMode(*Width, *Height, 32, videoMode)

	if screen == nil {
		panic(sdl.GetError())
	}

	sdl.WM_SetCaption("Novendiales 13", "")
	sdl.EnableKeyRepeat(10, 10)

	scrollX := 0
	scrollY := 0
	m, units := LoadMap(*MapName)

	for true {
		select {
		case ev := <-sdl.Events:
			/* TODO: something more clean for the cases ? */
			switch reflect.TypeOf(ev) {
			case reflect.TypeOf(sdl.QuitEvent{}):
				return
			case reflect.TypeOf(sdl.KeyboardEvent{}):
				k := ev.(sdl.KeyboardEvent)
				switch k.Keysym.Sym {
				case sdl.K_LEFT:
					scrollX = Max(scrollX-SCROLLSTEP, 0)
				case sdl.K_RIGHT:
					scrollX = Min(scrollX+SCROLLSTEP, *Width)
				case sdl.K_UP:
					scrollY = Max(scrollY-SCROLLSTEP, 0)
				case sdl.K_DOWN:
					scrollY = Min(scrollY+SCROLLSTEP, *Height)
				case sdl.K_ESCAPE:
					return
				}
			case reflect.TypeOf(sdl.MouseButtonEvent{}):
				m := ev.(sdl.MouseButtonEvent)
				if m.Type == sdl.MOUSEBUTTONDOWN && m.Button == 1 {
					fmt.Println("Mouse clicked at ", m.X, m.Y)
				}
			}
		default:
		}

		screen.FillRect(nil, 0x000000)

		m.Draw(scrollX, scrollY, screen)
		for i := 0; i < len(units); i++ {
			units[i].Draw(scrollX, scrollY, screen)
		}

		screen.Flip()
		time.Sleep(250)
	}

}
