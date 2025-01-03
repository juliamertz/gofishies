package main

import (
	"encoding/json"
	"fmt"
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

type Pos struct {
	X int
	Y int
}

// check if position in canvas is whitespace at an edge
// meaning it should be cut while rendering
func (p *Pos) IsEdge(c *Canvas) bool {
	if p.Y < len(c.cells) && p.X < len(c.cells[p.Y]) && c.cells[p.Y][p.X].Content != ' ' {
    return false 
	}

	hasTop := false
	hasBottom := false
	hasRight := false
	hasLeft := false

	for y := p.Y - 1; y >= 0; y-- {
		if c.cells[y][p.X].Content != ' ' {
			hasTop = true
			break
		}
	}

	for y := p.Y + 1; y < len(c.cells); y++ {
		if c.cells[y][p.X].Content != ' ' {
			hasBottom = true
			break
		}
	}

	for x := p.X - 1; x >= 0; x-- {
		if c.cells[p.Y][x].Content != ' ' {
			hasLeft = true
			break
		}
	}

	for x := p.X + 1; x < len(c.cells[p.Y]); x++ {
		if c.cells[p.Y][x].Content != ' ' {
			hasRight = true
			break
		}
	}

	return hasTop && hasBottom && hasRight && hasLeft
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
			// p := Pos{Y: y + i, X: x + j}
			if cell.Content == 0  {
				continue
			}
			// check if cell is surrounded by characthers
			// if so we fill it with the default bg color, otherwise its ignored

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

	maxEntities uint32
	paused      bool
	seaLevel    int
	tickRate    int

	entities []Entity

	stdin []byte
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

func removeIdx(slice []Entity, idx int) []Entity {
	var output []Entity
	for i, element := range slice {
		if i != idx {
			output = append(output, element)
		}
	}
	return output
}

func (r *Renderer) Tick() {
	for i, item := range r.entities {
		if item == nil {
			continue
		}
		item.Tick(r)
		if r.IsOffscreen(item) {
			r.entities = removeIdx(r.entities, i)
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
	width, height := r.screen.Size()
	r.canvas = NewCanvas(width, height)

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
			style := tcell.StyleDefault.Foreground(cell.Fg)
			r.screen.SetContent(x, y, rune(cell.Content), nil, style)
		}
	}

	r.DrawText(fmt.Sprintf("entities: %d", len(r.entities)), Pos{})
	ser, err := json.Marshal(r.entities)
	check(err)
	r.DrawText(fmt.Sprintf("entities: %s", string(ser)), Pos{Y: 1})
	r.DrawText(fmt.Sprintf("tickRate: %d", r.tickRate), Pos{Y: 2})
	if len(r.stdin) != 0 {
		r.DrawText(fmt.Sprintf("lastKey: %d ", r.stdin[len(r.stdin)-1]), Pos{Y: 3})
	}

	r.screen.Show()
	return nil
}
