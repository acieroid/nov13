package main

import "github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"

type Drawable interface {
	Draw(scrollX, scrollY int, surf *sdl.Surface)
}