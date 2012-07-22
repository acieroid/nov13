package main

func RunAI(units []*Character, m *Map) {
	for _, unit := range units {
		if unit.team != 2 {
			continue
		}

		/* TODO: try to maximize the damages */
		ennemy, dist := NearestEnnemy(unit, units)
		if ennemy != nil && dist <= Square(unit.damageRange) {
			Attack(unit, ennemy)
		}
	}
}

func Attack(unit *Character, ennemy *Character) {
	unit.nextAction = NewAttackAction(ennemy.x, ennemy.y, unit.attackSpeed)
}

func NearestEnnemy(unit *Character, units []*Character) (ennemy *Character, dist int) {
	var d int
	for _, u := range units {
		if u.team != unit.team {
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