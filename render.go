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

func (f *Frame) width() int  { return len(f.cells[0]) }
func (f *Frame) height() int { return len(f.cells) }

type Pos struct {
	X int
	Y int
}

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

// TODO: rename to engine or something
type Renderer struct {
	screen tcell.Screen
	frame  Frame

	debug      bool
	entityCaps *EntityCap
	paused     bool
	seaLevel   int
	tickRate   int

	entities []Entity
}

func (r *Renderer) KillEntity(idx int) {
	r.entities = slices.Delete(r.entities, idx, idx+1)
}

func (r *Renderer) Tick() {
	for i, e := range slices.Backward(r.entities) {
		// TODO: figure out why it doesn't update if `e` is passed instead if indexing into r.entities
		r.entities[i].Tick++
		if e.IsOffscreen(r.screen.Size()) || e.shouldKill {
			r.KillEntity(i)
			continue
		}

		e.update(&r.entities[i], r)
	}
}

func (r *Renderer) SpawnEntity(e Entity) {
	e.Id = fmt.Sprintf("%s-%d", e.Id, RNG.IntN(1000000))
	r.entities = append(r.entities, e)
}

func (r *Renderer) DrawText(content string, pos Pos) {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		for j, ch := range line {
			style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
			r.screen.SetContent(pos.X+j, pos.Y+i, ch, nil, style)
		}
	}
}

func (r *Renderer) Draw() error {
	if r == nil {
		return fmt.Errorf("Draw was called but the renderer has no screen set")
	}

	// render entities
	for _, e := range r.entities {
		// TODO: figure out why this is needed
		if e.currentFrame > len(*e.frames)-1 {
			continue
		}
		r.frame.MergeAt((*e.frames)[e.currentFrame], e.pos)
	}

	// print each cell of final canvas
	for y, line := range r.frame.cells {
		for x, cell := range line {
			if isWhitespace(cell.Content) {
				continue
			}
			style := tcell.StyleDefault.Foreground(cell.Fg)
			r.screen.SetContent(x, y, rune(cell.Content), nil, style)
		}
	}

	return nil
}

func (r *Renderer) spawnRandomEntity() {
	// assume small fish for now
	facing := Direction(RNG.IntN(2))
	var x int
	switch facing {
	case Left:
		x = r.frame.width() - 5
	case Right:
		x = -5
	}

	f := Fish(RNG.IntN(6), facing, RandColor(), Pos{
		Y: r.seaLevel + RNG.IntN(r.frame.height()-r.seaLevel),
		X: x,
	})

	f.Id = fmt.Sprintf("%s_%d", f.Id, len(r.entities))
	r.SpawnEntity(f)

	// r.entities = append(r.entities, f)
}

func (c *Frame) toString() string {
	buff := ""
	for _, line := range c.cells {
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
func (c *Frame) MergeAt(art Frame, pos Pos) {
	y := pos.Y
	x := pos.X

	for i, line := range art.cells {
		// skip line if it falls outside of base canvas
		if len(c.cells) <= y+i || y+i < 0 {
			continue
		}
		for j, cell := range line {
			// skip column if it falls outside of base canvas
			if len(c.cells[i]) <= x+j || x+j < 0 {
				continue
			}
			// if content is set to 0 we can assume this cell wasn't initialized so we can ignore it
			if cell.Content == 0 {
				continue
			}

			c.cells[y+i][x+j] = cell
		}
	}
}
