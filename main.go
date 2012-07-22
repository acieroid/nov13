package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/ttf"
	"github.com/0xe2-0x9a-0x9b/Go-SDL/gfx"
	"flag"
	"time"
	"os"
	"runtime/pprof"
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
var Width = flag.Int("width", 800, "Width of the window")
var Height = flag.Int("height", 600, "Height of the window")
var Fullscreen = flag.Bool("fullscreen", false, "Fullscreen")
var CPUProfile = flag.String("cpuprofile", "", "Write CPU Profile to file")

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

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

	if *CPUProfile != "" {
		f, err := os.Create(*CPUProfile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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
	BigFont = LoadFont("font.ttf", 16)
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

func DrawText(text string, x, y int, center bool, surf *sdl.Surface) {
	var w int
	var h int
	if center {
		w, h, _ = Font.SizeText(text)
	}
	surf.Blit(&sdl.Rect{int16(x - w/2), int16(y - h/2), 0, 0},
		ttf.RenderUTF8_Solid(Font, text, sdl.Color{0, 0, 0, 0}),
		nil)
}

func DrawTextBig(text string, x, y int, center bool, surf *sdl.Surface) {
	var w int
	var h int
	if center {
		w, h, _ = Font.SizeText(text)
	}
	surf.Blit(&sdl.Rect{int16(x - w/2),
		int16(y - h/2), 0, 0},
		ttf.RenderUTF8_Solid(BigFont, text, sdl.Color{255, 255, 255, 0}),
		nil)
}


func DrawImage(x, y int, img, surf *sdl.Surface) {
	surf.Blit(&sdl.Rect{int16(x), int16(y), 0, 0}, img, nil)
}

func Square(x int) int {
	return x*x
}