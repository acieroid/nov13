package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
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
		/* TODO: move */
		return m
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