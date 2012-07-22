package main

func RunAI(units []*Character, m *Map) {
	for _, unit := range units {
		if unit.team != 2 || !unit.Alive() {
			continue
		}

		ennemy, dist := NearestEnnemy(unit, units)
		if ennemy != nil && dist <= Square(unit.damageRange) {
			Attack(unit, ennemy)
			continue
		}

		if ennemy != nil {
			dx := ennemy.x - unit.x
			if dx != 0 {
				dx = dx/Abs(dx)
			}
			dy := ennemy.y - unit.y
			if dy != 0 {
				dy = dy/Abs(dy)
			}
			Move(unit, dx, dy)
			continue
		}
	}
}

func Attack(unit *Character, ennemy *Character) {
	unit.nextAction = NewAttackAction(ennemy.x, ennemy.y, unit.attackSpeed)
}

func Move(unit *Character, dx, dy int) {
	unit.nextAction = NewMoveAction(dx, dy)
}

func NearestEnnemy(unit *Character, units []*Character) (ennemy *Character, dist int) {
	var d int
	for _, u := range units {
		if u.team != unit.team && unit.Alive() {
			d = Distance(unit, u)
			if (d < dist || dist == 0) && d > 0 {
				dist = d
				ennemy = u
			}
		}
	}
	return
}

func Distance(a *Character, b *Character) int {
	return Square(a.x - b.x) + Square(a.y - b.y)
}