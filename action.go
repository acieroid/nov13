package main

import (
	"fmt"
)

type Action interface {
	Apply(c *Character, units []*Character, m *Map, delta int)
	Name() string
}

type MoveAction struct {
	dirX, dirY int
}

func NewMoveAction(dx, dy int) *MoveAction {
	return &MoveAction{dx, dy}
}

func (a *MoveAction) Apply(c *Character, units []*Character, m *Map, delta int) {
	if !c.Alive() {
		return
	}
	moveSpeed := c.moveSpeed
	switch m.TileAt(c.x, c.y) {
	case FOREST:
		moveSpeed += c.forestBonus
	case ROAD:
		moveSpeed += c.roadBonus
	}
	dx := (a.dirX * Max(10, delta*moveSpeed)) / 10
	dy := (a.dirY * Max(10, delta*moveSpeed)) / 10
	left := c.x + dx - CHARACTERSIZE/2
	right := c.x + dx + CHARACTERSIZE/2
	top := c.y + dy - CHARACTERSIZE/2
	bottom := c.y + dy + CHARACTERSIZE/2
	/* will not walk on another character ? */
	for _, unit := range units {
		if unit.Alive() && unit != c &&
			(unit.Contains(left, top) || unit.Contains(right, top) ||
				unit.Contains(left, bottom) || unit.Contains(right, bottom)) {
			return
		}
	}
	if m.CanMove(c, dx, dy) {
		c.x += dx
		c.y += dy
	}
}

func (a *MoveAction) Name() string {
	return "Se déplace"
}

type AttackAction struct {
	x, y       int
	nextAttack int
}

func NewAttackAction(x, y, speed int) *AttackAction {
	return &AttackAction{x, y, (10 - speed) * 10}
}

func (a *AttackAction) Apply(c *Character, units []*Character, m *Map, delta int) {
	if !c.Alive() {
		return
	}
	a.nextAttack -= delta
	if a.nextAttack < 0 {
		a.nextAttack += (10 - c.attackSpeed) * 10
		for _, unit := range units {
			top := a.y - c.damageSize/2
			bottom := a.y + c.damageSize/2
			left := a.x - c.damageSize/2
			right := a.x + c.damageSize/2
			rect := Rect{top, left, bottom, right}
			if unit.Alive() && unit.team != c.team &&
				rect.Contains(unit) {
				damage := Min(unit.life, c.damage)
				unit.life -= damage
				if !unit.Alive() {
					switch unit.team {
					case 1:
						AddMessage(fmt.Sprintf("Une de vos unité (%s) est morte", unit.Name()))
					case 2:
						AddMessage(fmt.Sprintf("Une unité ennemie (%s) est morte", unit.Name()))
					}
				}
			}
		}
	}
}

func (a *AttackAction) Name() string {
	return "Attaque"
}

type WaitAction struct {
}

func NewWaitAction() *WaitAction {
	return &WaitAction{}
}

func (a *WaitAction) Apply(c *Character, units[]*Character, m *Map, delta int) {
	/* Do nothing */
}

func (a *WaitAction) Name() string {
	return "Se tourne les pouces"
}