package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/gfx"
	"flag"
	"time"
)

const (
	GAME = iota
	WATCH
	MENU
	QUIT
	SCROLLSTEP = 5 /* scrolls 5px at a time */
	WATCHTIME = 3 /* 3 seconds of watch */
)

var Font *ttf.Font
var BigFont *ttf.Font
var MapDir = flag.String("maps", "maps", "Map directory")
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

func LoadFont(name string, size int) *ttf.Font {
	font := ttf.OpenFont(name, size)
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

	Font = LoadFont("font.ttf", 12)
	defer Font.Close()
	BigFont = LoadFont("font.ttf", 20)
	defer BigFont.Close()

	var game *Game
	menu := NewMainMenu(*MapDir, *Width, *Height)
	fps := gfx.NewFramerate()
	fps.SetFramerate(30)
	InitMessages(*Width, *Height)
	lastUpdate := time.Now()
	delta := 0
	mode := MENU
	mapName := ""

	for true {
		screen.FillRect(nil, 0x000000)

		switch mode {
		case MENU:
			mapName, mode = menu.Run(screen)
			if mode == GAME {
				game = NewGame(mapName, *Width, *Height)
			}
		case GAME:
			mode = game.Run(screen)
			if mode == MENU {
				/* Recreate the menu to avoid quitting when esc is still pressed */
				menu = NewMainMenu(*MapDir, *Width, *Height)
			}
		case QUIT:
			return
		}

		delta = int(time.Since(lastUpdate)/1e6)
		DrawMessages(delta, screen)


		screen.Flip()
		fps.FramerateDelay()
		lastUpdate = time.Now()
	}

}

func DrawText(text string, x, y int, surf *sdl.Surface) {
	surf.Blit(&sdl.Rect{int16(x), int16(y), 0, 0},
		ttf.RenderUTF8_Solid(Font, text, sdl.Color{0, 0, 0, 0}),
		nil)
}

func DrawTextBig(text string, x, y int, surf *sdl.Surface) {
	surf.Blit(&sdl.Rect{int16(x) /* TODO: font metrics */,
		int16(y), 0, 0},
		ttf.RenderUTF8_Solid(BigFont, text, sdl.Color{255, 255, 255, 0}),
		nil)
}


func DrawImage(x, y int, img, surf *sdl.Surface) {
	surf.Blit(&sdl.Rect{int16(x), int16(y), 0, 0}, img, nil)
}

func Square(x int) int {
	return x*x
}