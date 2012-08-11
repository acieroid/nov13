package main

import (
	"flag"
	"github.com/acieroid/go-sfml"
	"os"
	"runtime/pprof"
	"time"
)

const (
	GAME = iota
	WATCH
	MENU
	QUIT
	SCROLLSTEP = 5 /* scrolls 5px at a time */
	WATCHTIME  = 3 /* 3 seconds of watch */
)

var Font sfml.Font
var Text sfml.Text
var BigText sfml.Text
var MapDir = flag.String("maps", "maps", "Map directory")
var Width = flag.Int("width", 800, "Width of the window")
var Height = flag.Int("height", 600, "Height of the window")
var Fullscreen = flag.Bool("fullscreen", false, "Fullscreen")
var CPUProfile = flag.String("cpuprofile", "", "Write CPU Profile to file")

var RenderTexture sfml.RenderTexture

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

func LoadImage(name string) sfml.Sprite {
	texture, err := sfml.TextureFromFile(name, sfml.IntRect{})
	if err != nil {
		panic(err)
	}
	sprite, err := sfml.NewSprite()
	if err != nil {
		panic(err)
	}
	sprite.SetTexture(texture, false)
	texture.Cref = nil
	return sprite
}

func LoadFont(name string) sfml.Font {
	font, err := sfml.FontFromFile(name)
	if err != nil {
		panic(err)
	}
	return font
}

func LoadText(font sfml.Font, size int) sfml.Text {
	text, err := sfml.NewText()
	if err != nil {
		panic(err)
	}
	text.SetFont(font)
	text.SetCharacterSize(size)
	return text
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

	vm := sfml.NewVideoMode(uint(*Width), uint(*Height), 32)
	window := sfml.NewRenderWindow(vm, "Novendiales 13", sfml.StyleTitlebar|sfml.StyleClose, sfml.ContextSettings{nil})
	black := sfml.FromRGB(0, 0, 0)
	window.SetFramerateLimit(30)
	window.SetKeyRepeatEnabled(false)

	Font = LoadFont("font.ttf")
	defer Font.Destroy()
	Text = LoadText(Font, 12)
	defer Text.Destroy()
	BigText = LoadText(Font, 16)
	defer BigText.Destroy()

	var game *Game
	menu := NewMainMenu(*MapDir, *Width, *Height)
	InitMessages(*Width, *Height)
	lastUpdate := time.Now()
	delta := 0
	mode := MENU
	mapName := ""

	for window.IsOpen() {
		window.Clear(black)

		switch mode {
		case MENU:
			mapName, mode = menu.Run(window)
			if mode == GAME {
				game = NewGame(mapName, *Width, *Height)
			}
		case GAME:
			mode = game.Run(window)
			if mode == MENU {
				/* Recreate the menu to avoid quitting when esc is still pressed */
				menu = NewMainMenu(*MapDir, *Width, *Height)
			}
		case QUIT:
			return
		}

		delta = int(time.Since(lastUpdate) / 1e6)
		DrawMessages(delta, window)

		window.SetActive(false)

		window.Display()
		lastUpdate = time.Now()
	}
	RenderTexture.Destroy()
}

func DrawText(text string, x, y int, center bool, win sfml.RenderWindow) {
	Text.SetUnicodeString(text)
	if center {
		r := Text.LocalBounds()
		Text.SetPosition(float32(x)-r.Width()/2, float32(y)-r.Height())
	} else {
		Text.SetPosition(float32(x), float32(y))
	}
	win.DrawTextDefault(Text)
}

func DrawTextBig(text string, x, y int, center bool, win sfml.RenderWindow) {
	BigText.SetUnicodeString(text)
	if center {
		r := BigText.LocalBounds()
		BigText.SetPosition(float32(x)-r.Width()/2, float32(y)+5-r.Height())
	} else {
		BigText.SetPosition(float32(x), float32(y))
	}
	win.DrawTextDefault(BigText)
}

func DrawImage(x, y int, sprite sfml.Sprite, win sfml.RenderWindow) {
	sprite.SetPosition(float32(x), float32(y))
	win.DrawSpriteDefault(sprite)
}

func DrawShape(x, y int, shape sfml.RectangleShape, win sfml.RenderWindow) {
	shape.SetPosition(float32(x), float32(y))
	win.DrawRectangleShapeDefault(shape)
}


func Square(x int) int {
	return x*x
}
