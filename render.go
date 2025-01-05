package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Cell struct {
	Content byte
	Fg      tcell.Color
	Bg      tcell.Color
}

type Canvas struct {
	cells [][]Cell
}

type Frames = []Canvas

type Pos struct {
	X int
	Y int
}

func (c *Canvas) toString() string {
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

func NewCanvas(width int, height int) Canvas {
	buff := make([][]Cell, height, height)
	for i := range buff {
		buff[i] = make([]Cell, width, width)
	}
	return Canvas{cells: buff}
}

// Merge canvas cells into parent canvas
func (c *Canvas) MergeAt(art Canvas, pos Pos) {
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

func CanvasFromArt(art string, colors string, defaultColor tcell.Color) Canvas {
	lines := strings.Split(art, "\n")
	colorLines := strings.Split(colors, "\n")
	buff := NewCanvas(findArtWidth(art), len(lines))

	for y, line := range lines {
		for x, ch := range line {
			var color *tcell.Color

			if y > len(colorLines)-1 || x > len(colorLines[y])-1 {
				color = &defaultColor
			} else {
				color = ColorFromRune(rune(colorLines[y][x]))
			}

			if color == nil {
				n := tcell.ColorNone
				color = &n
				if ch == ' ' {
					ch = 0
				}
			}
			if *color == tcell.ColorDefault {
				color = &defaultColor
			}

			buff.cells[y][x] = Cell{
				Fg:      *color,
				Content: byte(ch),
			}
		}
	}
	return buff
}

func ColorFromRune(r rune) *tcell.Color {
	var color tcell.Color
	switch r {
	case 'r':
		color = tcell.ColorRed
	case 'g':
		color = tcell.ColorGreen
	case 'y':
		color = tcell.ColorYellow
	case 'b':
		color = tcell.ColorBlue
	case 'p':
		color = tcell.ColorPurple
	case 'c':
		color = tcell.ColorLightCyan
	case 'w':
		color = tcell.ColorWhite
	case 'd':
		color = tcell.ColorDefault
	case ' ':
		return nil
	}
	return &color
}

type Entity interface {
	Render(r *Renderer) (string, string)
	Tick(*Renderer)
	GetPos() Pos
	DefaultColor() tcell.Color
}

type Renderer struct {
	screen tcell.Screen
	canvas Canvas

	debug       bool
	maxEntities uint32
	paused      bool
	seaLevel    int
	tickRate    int

	entities []Entity
}

func (r *Renderer) IsOffscreen(e Entity) bool {
	pos := e.GetPos()
	rendered, _ := e.Render(r)
	split := strings.Split(rendered, "\n")

	height := len(split)
	width := findArtWidth(rendered)

	cols, lines := r.screen.Size()

	if pos.X >= cols || pos.X+width <= 0 {
		return true
	}

	if pos.Y > lines || pos.Y+height <= 0 {
		return true
	}

	return false
}

func FilterNil(input []Entity) []Entity {
	result := make([]Entity, 0, len(input))
	for _, elem := range input {
		if elem != nil {
			result = append(result, elem)
		}
	}
	return result
}

func (r *Renderer) KillEntity(v Entity) {
	for i, e := range r.entities {
		if e == v {
			r.entities = FilterNil(slices.Delete(r.entities, i, i+1))
			break
		}
	}
}

func (r *Renderer) Tick() {
	for _, item := range r.entities {
		if item == nil {
			continue
		}
		item.Tick(r)
		if r.IsOffscreen(item) {
      r.KillEntity(item)
		}
	}
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
		if e == nil {
			continue
		}
		art, colors := e.Render(r)
		c := CanvasFromArt(art, colors, e.DefaultColor())
		r.canvas.MergeAt(c, e.GetPos())
	}

	// print each cell of final canvas
	for y, line := range r.canvas.cells {
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
