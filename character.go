package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"fmt"
)

const (
	CHARACTERSIZE = TILESIZE-8
)

/* TODO: bonus/malus */
type Character struct {
	team int
	moveSpeed, attackSpeed int
	forestBonus, roadBonus int
	life, maxLife int
	damage, damageSize, damageRange int
	x, y int
	image *sdl.Surface
	Type int
	nextAction Action
}

var WarriorImage *sdl.Surface
var ArcherImage *sdl.Surface
var BoatImage *sdl.Surface
var RedBorder *sdl.Surface
var GreenBorder *sdl.Surface

const (
	WARRIOR = iota
	ARCHER
	BOAT
)

const (
	WMOVESPEED = 5
	WATTACKSPEED = 3
	WFORESTBONUS = -2
	WROADBONUS = 2
	WDAMAGESIZE = 30
	WDAMAGERANGE = 50
	WDAMAGE = 3
	WLIFE = 20
	AMOVESPEED = 5
	AATTACKSPEED = 4
	AFORESTBONUS = 3
	AROADBONUS = 0
	ADAMAGESIZE = 25
	ADAMAGERANGE = 150
	ADAMAGE = 2
	ALIFE = 20
	BMOVESPEED = 8
	BATTACKSPEED = 1
	BDAMAGESIZE = 45
	BDAMAGERANGE = 250
	BDAMAGE = 6
	BLIFE = 30
)

func NewCharacter(Type, team, ms, as, fb, rb, life, damage, ds, dr, x, y int, img *sdl.Surface) *Character {
	if RedBorder == nil {
		RedBorder = sdl.CreateRGBSurface(sdl.HWSURFACE,
			TILESIZE, TILESIZE, 32, 0, 0, 0, 0)
		RedBorder.FillRect(&sdl.Rect{0, 0, TILESIZE, TILESIZE},
			0x00FF0000)
		RedBorder.SetAlpha(sdl.SRCALPHA, 100)
		GreenBorder = sdl.CreateRGBSurface(sdl.HWSURFACE,
			TILESIZE, TILESIZE, 32, 0, 0, 0, 0)
		GreenBorder.FillRect(&sdl.Rect{0, 0, TILESIZE, TILESIZE},
			0x0000FF00)
		GreenBorder.SetAlpha(sdl.SRCALPHA, 100)
	}
	return &Character{team, ms, as, fb, rb, life, life, damage, ds, dr, x, y, img, Type, nil}
}

func NewWarrior(team, x, y int) *Character {
	if WarriorImage == nil {
		WarriorImage = LoadImage("img/warrior.png")
	}
	return NewCharacter(WARRIOR, team,
		WMOVESPEED, WATTACKSPEED, WFORESTBONUS, WROADBONUS,
		WLIFE, WDAMAGE, WDAMAGESIZE, WDAMAGERANGE,
		x, y,
		WarriorImage)
}

func NewArcher(team, x, y int) *Character {
	if ArcherImage == nil {
		ArcherImage = LoadImage("img/archer.png")
	}
	return NewCharacter(ARCHER, team,
		AMOVESPEED, AATTACKSPEED, AFORESTBONUS, AROADBONUS,
		ALIFE, ADAMAGE, ADAMAGESIZE, ADAMAGERANGE,
		x, y,
		ArcherImage)
}

func NewBoat(team, x, y int) *Character {
	if BoatImage == nil {
		BoatImage = LoadImage("img/boat.png")
	}
	return NewCharacter(BOAT, team,
		BMOVESPEED, BATTACKSPEED, 0, 0,
		BLIFE, BDAMAGE, BDAMAGESIZE, BDAMAGERANGE,
		x, y,
		BoatImage)
}


func (c *Character) Draw(scrollX, scrollY int, surf *sdl.Surface) {
	if !c.Alive() {
		return
	}
	border := RedBorder
	if c.team == 1 {
		border = GreenBorder
	}
	DrawImage(c.x - TILESIZE/2 - scrollX,
		c.y - TILESIZE/2 - scrollY,
		border, surf)
	DrawImage(c.x - CHARACTERSIZE/2 - scrollX,
		c.y - CHARACTERSIZE/2 - scrollY,
		c.image, surf)
	DrawText(fmt.Sprintf("%d/%d", c.life, c.maxLife),
		c.x - TILESIZE/2 - scrollX,
		c.y - TILESIZE/2 - scrollY,
		false, surf)
	if c.nextAction != nil {
		DrawText(c.nextAction.Name(),
			c.x - TILESIZE/2 - scrollX,
			c.y - TILESIZE/2 - scrollY + 14,
			false, surf)
	}
}

func (c *Character) Contains(x, y int) bool {
	return (x >= c.x - CHARACTERSIZE/2 && x <= c.x + CHARACTERSIZE/2 &&
		y >= c.y - CHARACTERSIZE/2 && y <= c.y + CHARACTERSIZE/2)
}

func (c *Character) Alive() bool {
	return c.life > 0
}

func (c *Character) Name() string {
	switch c.Type {
	case WARRIOR:
		return "Guerrier"
	case ARCHER:
		return "Archer"
	case BOAT:
		return "Navire"
	}
	return "Inconnu ?!!"
}