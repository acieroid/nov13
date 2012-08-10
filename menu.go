package main

import (
	"github.com/acieroid/go-sfml"
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

func (m *CharacterMenu) Draw(scrollX, scrollY int, win sfml.RenderWindow) {
	DrawText("Attaquer",
		m.c.x + TILESIZE/2 - scrollX,
		m.c.y - TILESIZE/2 - scrollY,
		false, win)
	DrawText("Déplacer",
		m.c.x + TILESIZE/2 - scrollX,
		m.c.y - TILESIZE/2 + 14 - scrollY,
		false, win)
	DrawText("Attendre",
		m.c.x + TILESIZE/2 - scrollX,
		m.c.y - TILESIZE/2 + 28 - scrollY,
		false, win)
	DrawText("Retour",
		m.c.x + TILESIZE/2 - scrollX,
		m.c.y - TILESIZE/2 + 42 - scrollY,
		false, win)
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
	green, red sfml.RectangleShape
}

func NewAttackMenu(c *Character) *AttackMenu {
	green := sfml.NewRectangleShape()
	green.SetSize(float32(c.damageSize), float32(c.damageSize))
	green.SetFillColor(sfml.FromRGBA(0, 255, 0, 150))
	red := green.Copy()
	red.SetFillColor(sfml.FromRGBA(255, 0, 0, 150))
	return &AttackMenu{c, green, red}
}

func (m *AttackMenu) Draw(scrollX, scrollY int, win sfml.RenderWindow) {
	x, y := sfml.MousePositionAbsolute()
	s := m.green
	if Square(x + scrollX - m.c.x) + Square(y + scrollY- m.c.y) > Square(m.c.damageRange) {
		/* TODO: red if outside map */
		s = m.red
	}
	s.SetPosition(float32(x - m.c.damageSize/2), float32(y - m.c.damageSize/2))
	win.DrawRectangleShapeDefault(s)
}

func (m *AttackMenu) Contains(x, y int) bool {
	/* We always want to receive the clicks */
	return true
}

func (m *AttackMenu) Clicked(x, y int) Menu {
	if Square(x - m.c.x) + Square(y-m.c.y) < Square(m.c.damageRange) {
		m.c.nextAction = NewAttackAction(x, y, m.c.attackSpeed)
		return nil
	}
	return m
}

type MoveMenu struct {
	c *Character
}

var Left, Right, Up, Down, BottomLeft, BottomRight, TopLeft, TopRight sfml.Sprite

func NewMoveMenu(c *Character) (m *MoveMenu) {
	if Left.Cref == nil {
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


func (m *MoveMenu) Draw(scrollX, scrollY int, win sfml.RenderWindow) {
	DrawImage(m.c.x - 3*TILESIZE/2 - scrollX, m.c.y - TILESIZE/2 - scrollY, Left, win)
	DrawImage(m.c.x - 3*TILESIZE/2 - scrollX, m.c.y - 3*TILESIZE/2 - scrollY, TopLeft, win)
	DrawImage(m.c.x - 3*TILESIZE/2 - scrollX, m.c.y + TILESIZE/2 - scrollY, BottomLeft, win)
	DrawImage(m.c.x - TILESIZE/2 - scrollX, m.c.y - 3*TILESIZE/2 - scrollY, Up, win)
	DrawImage(m.c.x - TILESIZE/2 - scrollX, m.c.y + TILESIZE/2 - scrollY, Down, win)
	DrawImage(m.c.x + TILESIZE/2 - scrollX, m.c.y - 3*TILESIZE/2 - scrollY, TopRight, win)
	DrawImage(m.c.x + TILESIZE/2 - scrollX, m.c.y - TILESIZE/2 - scrollY, Right, win)
	DrawImage(m.c.x + TILESIZE/2 - scrollX, m.c.y + TILESIZE/2 - scrollY, BottomRight, win)
}

func (m *MoveMenu) Contains(x, y int) bool {
	return (x > m.c.x - 3*TILESIZE/2 && x < m.c.x + 3*TILESIZE/2 &&
		y > m.c.y - 3*TILESIZE/2 && y < m.c.y + 3*TILESIZE/2)
}

func (m *MoveMenu) Clicked(x, y int) Menu {
	dx := 0
	dy := 0
	if m.c.Contains(x, y) {
		return m
	}
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
