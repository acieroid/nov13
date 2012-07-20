package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
	"time"
	"reflect"
	"flag"
)

const (
	GAME = iota
	WATCH
	SCROLLSTEP = 5
)

var Font *ttf.Font
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

func LoadFont(name string) *ttf.Font {
	font := ttf.OpenFont(name, 12)
	if font == nil {
		panic(sdl.GetError())
	}
	return font
}

func main() {
	flag.Parse()
	if sdl.Init(sdl.INIT_VIDEO) != 0 || ttf.Init() != 0 {
		panic(sdl.GetError())
	}

	defer sdl.Quit()
	defer ttf.Quit()

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

	Font = LoadFont("font.ttf")
	defer Font.Close()

	scrollX := 0
	scrollY := 0
	m, units := LoadMap(*MapName)
	var menu Menu
	mode := GAME
	watchButton := NewWatchButton(*Width-100, *Height - 30)

	for true {
		select {
		case ev := <-sdl.Events:
			/* TODO: something more clean for the cases ? */
			switch reflect.TypeOf(ev) {
			case reflect.TypeOf(sdl.QuitEvent{}):
				return
			case reflect.TypeOf(sdl.KeyboardEvent{}):
				e := ev.(sdl.KeyboardEvent)
				switch e.Keysym.Sym {
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
				e := ev.(sdl.MouseButtonEvent)
				if e.Type == sdl.MOUSEBUTTONDOWN && e.Button == 1 {
					x := int(e.X) + scrollX
					y := int(e.Y) + scrollY
					if menu != nil && menu.Contains(x, y) {
						menu = menu.Clicked(x, y)
					} else if mode == GAME && watchButton.Contains(x, y) {
						mode = WATCH
						menu = nil
						watchButton.Enabled()
					} else if mode == GAME && x < m.width * TILESIZE && y < m.height * TILESIZE {
						for _, unit := range(units) {
							if unit.Contains(x, y) && unit.team == 1 {
								menu = NewCharacterMenu(unit)
								break
							}
						}
					}
				}
			}
		default:
		}

		screen.FillRect(nil, 0x000000)

		m.Draw(scrollX, scrollY, screen)
		for _, unit := range(units) {
			unit.Draw(scrollX, scrollY, screen)
		}
		if menu != nil {
			menu.Draw(scrollX, scrollY, screen)
		}
		watchButton.Draw(screen)

		screen.Flip()
		time.Sleep(250)
	}

}

func DrawText(text string, x, y int, surf *sdl.Surface) {
	surf.Blit(&sdl.Rect{int16(x), int16(y), 0, 0},
		ttf.RenderUTF8_Solid(Font, text, sdl.Color{0, 0, 0, 0}),
		nil)
}

func DrawImage(x, y int, img, surf *sdl.Surface) {
	surf.Blit(&sdl.Rect{int16(x), int16(y), 0, 0}, img, nil)
}

func Square(x int) int {
	return x*x
}