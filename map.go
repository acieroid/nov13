package main

import (
	"github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
	"os"
	"log"
	"bufio"
	"strconv"
	"fmt"
)

const (
	GRASS = 1
	ROAD = 2
	FOREST = 3
	WATER = 4
	TILESIZE = 32
)
	
type Map struct {
	width, height int
	contents [][]int
	images []*sdl.Surface
	surf *sdl.Surface
}

func LoadMap(name string) (m *Map, units []*Character) {
	m = &Map{}
	file, err := os.Open(fmt.Sprintf("maps/%s.txt", name))
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	widthLine, _ := reader.ReadString(' ')
	heightLine, _ := reader.ReadString('\n')
	m.width, err = strconv.Atoi(widthLine[:len(widthLine)-1])
	if err != nil { log.Fatal(err) }
	m.height, err = strconv.Atoi(heightLine[:len(heightLine)-1])
	if err != nil { log.Fatal(err) }

	var c byte
	m.contents = make([][]int, m.width)
	for x := 0; x < m.width; x++ {
		m.contents[x] = make([]int, m.height)
	}

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			c, err = reader.ReadByte()
			if err != nil {
				log.Fatal(err)
			}
			m.contents[x][y] = int(c - '0')
		}
		c, err = reader.ReadByte()
		if err != nil {
			log.Fatal(err)
		}
		if c != '\n' {
			log.Fatal("Bad map format for map")
		}
	}

	unitsLine, _ := reader.ReadString('\n')
	nUnits, err := strconv.Atoi(unitsLine[:len(unitsLine)-1])
	if err != nil { log.Fatal(err) }
	units = make([]*Character, nUnits)

	for i := 0; i < nUnits; i++ {
		unitType, _ := reader.ReadString(' ')
		teamLine, _ := reader.ReadString(' ')
		XLine, _ := reader.ReadString(' ')
		YLine, _ := reader.ReadString('\n')
		team, err := strconv.Atoi(teamLine[:len(teamLine)-1])
		if err != nil { log.Fatal(err) }
		x, err := strconv.Atoi(XLine[:len(XLine)-1])
		if err != nil { log.Fatal(err) }
		y, err := strconv.Atoi(YLine[:len(YLine)-1])
		if err != nil { log.Fatal(err) }
		x = x*TILESIZE + TILESIZE/2
		y = y*TILESIZE + TILESIZE/2
		switch unitType[0] {
		case 'W':
			units[i] = NewWarrior(team, x, y)
		case 'A':
			units[i] = NewArcher(team, x, y)
		case 'B':
			units[i] = NewBoat(team, x, y)
		}
	}

	m.images = []*sdl.Surface{
		LoadImage("img/grass.png"),
		LoadImage("img/road.png"),
		LoadImage("img/forest.png"),
		LoadImage("img/water.png")}

	m.surf = sdl.CreateRGBSurface(sdl.HWSURFACE,
		m.width*TILESIZE, m.height*TILESIZE,
		32, 0, 0, 0, 0)
	for x := 0; x < m.width; x++ {
		for y := 0; y < m.height; y++ {
			m.surf.Blit(&sdl.Rect{
				int16(x * TILESIZE),
				int16(y * TILESIZE),
				0, 0},
				m.images[m.contents[x][y]-1], nil)
		}
	}

	return
}

func (m *Map) Draw(scrollX, scrollY int, surf *sdl.Surface) {
	surf.Blit(&sdl.Rect{
		int16(-scrollX),
		int16(-scrollY),
		0, 0},
		m.surf, nil)
}

func (m *Map) TileAt(x, y int) int {
	return m.contents[x/TILESIZE][y/TILESIZE]
}

func (m *Map) CanMove(c *Character, dx, dy int) bool {
	x := c.x + dx
	y := c.y + dy
	left := x - CHARACTERSIZE/2
	right := x + CHARACTERSIZE/2
	top := y - CHARACTERSIZE/2
	bottom := y + CHARACTERSIZE/2
	/* still in the map ? */
	if left < 0 || right >= m.width*TILESIZE ||
		top < 0 || bottom >= m.height*TILESIZE {
		return false
	}

	/* can go on this case ? */
	if c.Type == BOAT {
		return m.TileAt(left, top) == WATER &&
			m.TileAt(right, top) == WATER &&
			m.TileAt(left, bottom) == WATER &&
			m.TileAt(right, bottom) == WATER
	}
	return m.TileAt(left, top) != WATER &&
		m.TileAt(right, top) != WATER &&
		m.TileAt(left, bottom) != WATER &&
		m.TileAt(right, bottom) != WATER
}