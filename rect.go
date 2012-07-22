package main

type Rect struct {
	top, left, bottom, right int
}

func (r Rect) ContainsPoint(x, y int) bool {
	return x >= r.left && x <= r.right &&
		y >= r.top && y <= r.bottom
}

func (r Rect) Contains(c *Character) bool {
	left := c.x - CHARACTERSIZE/2
	right := c.x + CHARACTERSIZE/2
	top := c.y - CHARACTERSIZE/2
	bottom := c.y + CHARACTERSIZE/2
	return r.ContainsPoint(left, top) ||
		r.ContainsPoint(left, bottom) ||
		r.ContainsPoint(right, top) ||
		r.ContainsPoint(right, bottom)
}