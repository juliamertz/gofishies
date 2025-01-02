package main

import (
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

type Cell struct {
	Content byte
	Fg      tcell.Color
	Bg      int
}

type Canvas struct {
	cells [][]Cell
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
			c.cells[y+i][x+j] = cell
		}
	}
}

func CanvasFromArt(art string, colors string, defaultColor tcell.Color) Canvas {
	lines := strings.Split(art, "\n")
	colorLines := strings.Split(colors, "\n")
	buff := NewCanvas(len(lines[0]), len(lines))

	for y, line := range lines {
		for x, ch := range line {
			color := ColorFromRune(rune(colorLines[y][x]))
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

type Renderable interface {
	Render(r *Renderer) (string, string)
	Tick(*Renderer)
	GetPos() Pos
	DefaultColor() tcell.Color
}

type Pos struct {
	X int
	Y int
}

type Renderer struct {
	screen tcell.Screen
	canvas Canvas

  maxEntities uint32
	paused   bool
	tickRate int

	entities  []Renderable
	particles []Renderable
}

// Creates new 2 dimensional cell slice and discards the previous one
func (r *Renderer) InitCells() {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic("unable to get terminal size")
	}
	r.canvas = NewCanvas(width, height)
}

func (r *Renderer) Tick() {
	for _, e := range r.entities {
		e.Tick(r)
	}
	for _, e := range r.particles {
		e.Tick(r)
	}
}

func (r *Renderer) Draw(screen tcell.Screen) {
  // render entities
	for _, e := range r.entities {
		if r == nil {
			panic("draw was called but r is nil")
		}
		art, colors := e.Render(r)
		c := CanvasFromArt(art, colors, e.DefaultColor())
		r.canvas.MergeAt(c, e.GetPos())
	}
  // render particles
	for _, e := range r.particles {
		if r == nil {
			panic("draw was called but r is nil")
		}
		art, colors := e.Render(r)
		c := CanvasFromArt(art, colors, e.DefaultColor())
		r.canvas.MergeAt(c, e.GetPos())
	}

	// print each cell of final canvas
	// ansi.Print(ansi.MoveHome)
	for y, line := range r.canvas.cells {
		for x, cell := range line {
			style := tcell.StyleDefault.Foreground(cell.Fg)
			screen.SetContent(x, y, rune(cell.Content), nil, style)
			// screen.SetContent(x, y, rune(cell.Content), nil, style)
		}
	}

	// screen.Show()
	screen.Sync()
}

// Check if art and color map are of the exact same lengths and are normalized
func compareArtStrings(art string, colors string) {
	artLines := strings.Split(art, "\n")
	colorLines := strings.Split(colors, "\n")
	firstLineLength := len(artLines[0])

	for i, line := range artLines {
		if len(line) != firstLineLength || len(colorLines[i]) != firstLineLength {
			panic("invalid art strings: " + art + colors)
		}
	}
}
