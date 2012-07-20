package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"fmt"
)

type Menu interface {
	Drawable
	Contains(x, y int) bool
	Clicked(x, y int) Menu
}

type CharacterMenu struct {
	c *Character
}

func NewCharacterMenu(c *Character) *CharacterMenu {
	return &CharacterMenu{c}
}

func (m *CharacterMenu) Draw(scrollX, scrollY int, surf *sdl.Surface) {
	DrawText("Attaquer",
		m.c.x + TILESIZE/2 - scrollX,
		m.c.y - TILESIZE/2 - scrollY,
		surf)
	DrawText("Déplacer",
		m.c.x + TILESIZE/2 - scrollX,
		m.c.y - TILESIZE/2 + 14 - scrollY,
		surf)
	DrawText("Attendre",
		m.c.x + TILESIZE/2 - scrollX,
		m.c.y - TILESIZE/2 + 28 - scrollY,
		surf)
	DrawText("Retour",
		m.c.x + TILESIZE/2 - scrollX,
		m.c.y - TILESIZE/2 + 42 - scrollY,
		surf)
}

func (m *CharacterMenu) Contains(x, y int) bool {
	return (x > m.c.x + TILESIZE/2 && x < m.c.x + TILESIZE/2 + 50 && /* TODO: use font metrics ? */
		y > m.c.y - TILESIZE/2 && y < m.c.y - TILESIZE/2 + 60)
}

func (m *CharacterMenu) Clicked(x, y int) Menu {
	element := (y - m.c.y + TILESIZE/2)/14
	if element == 0 {
		return NewAttackMenu(m.c)
	} else if element == 1 {
		return NewMoveMenu(m.c)
	} else if element == 2 {
		m.c.nextAction = NewWaitAction()
		return nil
	} else {
		return nil
	}
	return nil
}

type AttackMenu struct {
	c *Character
	surface *sdl.Surface
}

func NewAttackMenu(c *Character) *AttackMenu {
	surf :=	sdl.CreateRGBSurface(sdl.HWSURFACE,
		c.damageSize, c.damageSize,
		32, 0, 0, 0, 0)
	surf.FillRect(&sdl.Rect{0, 0,
		uint16(c.damageSize), uint16(c.damageSize)},
		0x00FF0000)
  surf.SetAlpha(sdl.SRCALPHA, 200)
	return &AttackMenu{c, surf}
}

func (m *AttackMenu) Draw(scrollX, scrollY int, surf *sdl.Surface) {
	var x, y int
	/* this sucks in SDL binding, we should use a multiple value return */
	sdl.GetMouseState(&x, &y)
	surf.Blit(&sdl.Rect{
		int16(x - m.c.damageSize/2),
		int16(y - m.c.damageSize/2),
		0, 0},
		m.surface, nil)
}

func (m *AttackMenu) Contains(x, y int) bool {
	/* We always want to receive the clicks */
	return true
}

func (m *AttackMenu) Clicked(x, y int) Menu {
	m.c.nextAction = NewAttackAction(x, y, m.c.attackSpeed)
	return nil
}

type MoveMenu struct {
	c *Character
}

var Left, Right, Up, Down, BottomLeft, BottomRight, TopLeft, TopRight *sdl.Surface

func NewMoveMenu(c *Character) (m *MoveMenu) {
	if Left == nil {
		Left = LoadImage("img/left.png")
		Right = LoadImage("img/right.png")
		Up = LoadImage("img/up.png")
		Down = LoadImage("img/down.png")
		BottomLeft = LoadImage("img/bottomleft.png")
		BottomRight = LoadImage("img/bottomright.png")
		TopLeft = LoadImage("img/topleft.png")
		TopRight = LoadImage("img/topright.png")
	}
	return &MoveMenu{c}
}


func (m *MoveMenu) Draw(scrollX, scrollY int, surf *sdl.Surface) {
	DrawImage(m.c.x - 3*TILESIZE/2, m.c.y - TILESIZE/2, Left, surf)
	DrawImage(m.c.x - 3*TILESIZE/2, m.c.y - 3*TILESIZE/2, TopLeft, surf)
	DrawImage(m.c.x - 3*TILESIZE/2, m.c.y + TILESIZE/2, BottomLeft, surf)
	DrawImage(m.c.x - TILESIZE/2, m.c.y - 3*TILESIZE/2, Up, surf)
	DrawImage(m.c.x - TILESIZE/2, m.c.y + TILESIZE/2, Down, surf)
	DrawImage(m.c.x + TILESIZE/2, m.c.y - 3*TILESIZE/2, TopRight, surf)
	DrawImage(m.c.x + TILESIZE/2, m.c.y - TILESIZE/2, Right, surf)
	DrawImage(m.c.x + TILESIZE/2, m.c.y + TILESIZE/2, BottomRight, surf)
}

func (m *MoveMenu) Contains(x, y int) bool {
	return (x > m.c.x - 3*TILESIZE/2 && x < m.c.x + 3*TILESIZE/2 &&
		y > m.c.y - 3*TILESIZE/2 && y < m.c.y + 3*TILESIZE/2)
}

func (m *MoveMenu) Clicked(x, y int) Menu {
	dx := 0
	dy := 0
	if x > m.c.x - 3*TILESIZE/2 && x < m.c.x - TILESIZE/2 {
		dx = -1
	}
	if x > m.c.x + TILESIZE/2 && x < m.c.x + 3*TILESIZE/2 {
		dx = 1
	}
	if y > m.c.y - 3*TILESIZE/2 && y < m.c.y - TILESIZE/2 {
		dy = -1
	}
	if y > m.c.y + TILESIZE/2 && y < m.c.y + 3*TILESIZE/2 {
		dy = 1
	}
	m.c.nextAction = NewMoveAction(dx, dy)
	return nil
}
