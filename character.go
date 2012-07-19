package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"fmt"
)

/* TODO: bonus/malus */
type Character struct {
	team int
	moveSpeed, attackSpeed int
	life, maxLife int
	damage, damageSize int
	x, y int
	image *sdl.Surface
	Type int
	nextAction Action
}

var WarriorImage *sdl.Surface
var ArcherImage *sdl.Surface
var BoatImage *sdl.Surface

const (
	WARRIOR = iota
	ARCHER
	BOAT
)

const (
	WMOVESPEED = 5
	WATTACKSPEED = 7
	WDAMAGESIZE = 25
	WDAMAGE = 7
	WLIFE = 20
	AMOVESPEED = 7
	AATTACKSPEED = 5
	ADAMAGESIZE = 15
	ADAMAGE = 6
	ALIFE = 20
	BMOVESPEED = 8
	BATTACKSPEED = 3
	BDAMAGESIZE = 40
	BDAMAGE = 15
	BLIFE = 30
)

func NewWarrior(team, x, y int) *Character {
	if WarriorImage == nil {
		WarriorImage = LoadImage("img/warrior.png")
	}
	return &Character{team,
		WMOVESPEED, WATTACKSPEED,
		WLIFE, WLIFE, WDAMAGE, WDAMAGESIZE,
		x, y,
		WarriorImage, WARRIOR, nil}
}

func NewArcher(team, x, y int) *Character {
	if ArcherImage == nil {
		ArcherImage = LoadImage("img/archer.png")
	}
	return &Character{team,
		AMOVESPEED, AATTACKSPEED,
		ALIFE, ALIFE, ADAMAGE, ADAMAGESIZE,
		x, y,
		ArcherImage, ARCHER, nil}
}

func NewBoat(team, x, y int) *Character {
	if BoatImage == nil {
		BoatImage = LoadImage("img/boat.png")
	}
	return &Character{team,
		BMOVESPEED, BATTACKSPEED,
		BLIFE, BLIFE, BDAMAGE, BDAMAGESIZE,
		x, y,
		BoatImage, BOAT, nil}
}


func (c *Character) Draw(scrollX, scrollY int, surf *sdl.Surface) {
	surf.Blit(&sdl.Rect{
		int16(c.x - TILESIZE/2 - scrollX),
		int16(c.y - TILESIZE/2 - scrollY),
		0, 0},
		c.image, nil)
	DrawText(fmt.Sprintf("%d/%d", c.life, c.maxLife),
		c.x - TILESIZE/2 - scrollX,
		c.y - TILESIZE/2 - scrollY,
		surf)
	if c.nextAction != nil {
		DrawText(c.nextAction.Name(),
			c.x - TILESIZE/2 - scrollX,
			c.y - TILESIZE/2 - scrollY + 14,
			surf)
	}
}

func (c *Character) Contains(x, y int) bool {
	return (x > c.x - TILESIZE/2 && x < c.x + TILESIZE/2 &&
		y > c.y - TILESIZE/2 && y < c.y + TILESIZE/2)
}