package main

import "github.com/acieroid/go-sfml"

type Drawable interface {
	Draw(scrollX, scrollY int, win sfml.RenderWindow)
}