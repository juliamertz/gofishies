package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Direction = int

const (
	Left  Direction = 0
	Right Direction = 1
)

type Cell struct {
	Content byte
	Fg      tcell.Color
	Bg      tcell.Color
}

type Frame struct {
	cells [][]Cell
}

type Pos struct {
	X int
	Y int
}

type EntityCaps struct {
	SmallFish int
	LargeFish int
	WaterLine int
	Vehicle   int
}

type Engine struct {
	screen  tcell.Screen
	frame   Frame
	spawner Spawner

	debug    bool
	paused   bool
	seaLevel int
	tickRate int

	entities []Entity
}

func (f *Frame) width() int  { return len(f.cells[0]) }
func (f *Frame) height() int { return len(f.cells) }

func generateFrame(art string, colors string, defaultColor tcell.Color) Frame {
	lines := strings.Split(art, "\n")
	colorLines := strings.Split(colors, "\n")
	buff := makeFrame(findArtWidth(art), len(lines))

	for y, line := range lines {
		for x, ch := range line {
			var color *tcell.Color

			if y > len(colorLines)-1 || x > len(colorLines[y])-1 {
				color = &defaultColor
			} else {
				c := ColorFromRune(rune(colorLines[y][x]))
				color = &c
			}

			if *color == tcell.ColorNone {
				// if the cell doesn't have a color nor a charachter the char is set to 0
				// then the renderer knows it can be skipped over
				if ch == ' ' {
					ch = 0
				} else {
					// if it does hold a char which isn't whitespace we can apply the default color
					color = &defaultColor
				}
			}

			buff.cells[y][x] = Cell{
				Fg:      *color,
				Content: byte(ch),
			}
		}
	}
	return buff
}

func (e *Engine) Tick() {
	for idx, entity := range slices.Backward(e.entities) {
		// TODO: figure out why it doesn't update if `e` is passed instead if indexing into r.entities
		e.entities[idx].Tick++
		if entity.IsOffscreen(e.screen.Size()) || entity.shouldKill {
			e.KillEntity(idx)
			continue
		}

		entity.update(&e.entities[idx], e)
	}
}

func (e *Engine) Reset(entities []Entity) {
	e.entities = entities
	e.spawner.caps = EntityCaps{}
}

func (e *Engine) KillEntity(idx int) {
	e.spawner.caps.decrement(e.entities[idx].kind)
	e.entities = slices.Delete(e.entities, idx, idx+1)
}

func (e *Engine) SpawnEntity(entity Entity) {
	entity.Id = fmt.Sprintf("%s-%d", entity.Id, RNG.IntN(1000000))
	e.entities = append(e.entities, entity)
}

func (e *Engine) DrawText(content string, pos Pos) {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		for j, ch := range line {
			style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
			e.screen.SetContent(pos.X+j, pos.Y+i, ch, nil, style)
		}
	}
}

func (e *Engine) Draw() error {
	if e == nil {
		return fmt.Errorf("Draw was called but the renderer has no screen set")
	}

	// render entities
	for _, entity := range e.entities {
		// TODO: figure out why this is needed
		if entity.currentFrame > len(*entity.frames)-1 {
			continue
		}
		e.frame.MergeAt((*entity.frames)[entity.currentFrame], entity.pos)
	}

	// print each cell of final canvas
	for y, line := range e.frame.cells {
		for x, cell := range line {
			if isWhitespace(cell.Content) {
				continue
			}
			style := tcell.StyleDefault.Foreground(cell.Fg)
			e.screen.SetContent(x, y, rune(cell.Content), nil, style)
		}
	}

	return nil
}

func (f *Frame) toString() string {
	buff := ""
	for _, line := range f.cells {
		for _, cell := range line {
			buff += string(cell.Content)
		}
		buff += "\n"
	}
	return buff
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || byte(ch) == 0
}

// Create empty sized frame
func makeFrame(width int, height int) Frame {
	buff := make([][]Cell, height, height)
	for i := range buff {
		buff[i] = make([]Cell, width, width)
	}
	return Frame{cells: buff}
}

// Merge canvas cells into parent canvas
func (f *Frame) MergeAt(art Frame, pos Pos) {
	y := pos.Y
	x := pos.X

	for i, line := range art.cells {
		// skip line if it falls outside of base canvas
		if len(f.cells) <= y+i || y+i < 0 {
			continue
		}
		for j, cell := range line {
			// skip column if it falls outside of base canvas
			if len(f.cells[i]) <= x+j || x+j < 0 {
				continue
			}
			// if content is set to 0 we can assume this cell wasn't initialized so we can ignore it
			if cell.Content == 0 {
				continue
			}

			f.cells[y+i][x+j] = cell
		}
	}
}

func (c *EntityCaps) increment(kind EntityKind) {
	switch kind {
	case SmallFish:
		c.SmallFish++
	case LargeFish:
		c.LargeFish++
	case WaterLine:
		c.WaterLine++
	case Vehicle:
		c.Vehicle++
	}
}

func (c *EntityCaps) decrement(kind EntityKind) {
	switch kind {
	case SmallFish:
		c.SmallFish--
	case LargeFish:
		c.LargeFish--
	case WaterLine:
		c.WaterLine--
	case Vehicle:
		c.Vehicle--
	}
}
