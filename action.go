package main

type Action interface {
	Apply(c *Character, units []*Character, delta int)
	Name() string
}

type MoveAction struct {
	dirX, dirY int
}

func NewMoveAction(dx, dy int) *MoveAction {
	return &MoveAction{dx, dy}
}

func (a *MoveAction) Apply(c *Character, units []*Character, delta int) {
	dx := (a.dirX * delta * c.moveSpeed)/10
	dy := (a.dirY * delta * c.moveSpeed)/10
	c.x += dx
	c.y += dy
}

func (a *MoveAction) Name() string {
	return "Se déplace"
}

type AttackAction struct {
	x, y int
	nextAttack int
}

func NewAttackAction(x, y, speed int) *AttackAction {
	return &AttackAction{x, y, speed*10}
}

func (a *AttackAction) Apply(c *Character, units[]*Character, delta int) {
	a.nextAttack -= delta
	if a.nextAttack < 0 {
		a.nextAttack += c.attackSpeed*10
		for _, unit := range units {
			if unit.x > a.x - c.damageSize/2 &&
				unit.x < a.x + c.damageSize/2 &&
				unit.y > a.y - c.damageSize/2 &&
				unit.y < a.y + c.damageSize/2 {
				unit.life -= c.damage
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

func (a *WaitAction) Apply(c *Character, units[]*Character, delta int) {
	/* Do nothing */
}

func (a *WaitAction) Name() string {
	return "Se tourne les pouces"
}