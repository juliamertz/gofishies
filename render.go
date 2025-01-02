package main

import (
	"gofishies/ansi"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

type Canvas struct {
	cells [][]ansi.Cell
}

func NewCanvas(width int, height int) Canvas {
	buff := make([][]ansi.Cell, height, height)
	for i := range buff {
		buff[i] = make([]ansi.Cell, width, width)
	}
	return Canvas{cells: buff}
}

// Merge canvas cells into parent canvas
func (c *Canvas) MergeAt(art Canvas, pos ansi.Pos) {
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

func CanvasFromArt(art string, colors string) Canvas {
	lines := strings.Split(art, "\n")
	colorLines := strings.Split(colors, "\n")
	buff := NewCanvas(len(lines[0]), len(lines))

	for y, line := range lines {
		for x, ch := range line {
			color := ansi.ColorFromByte(byte(colorLines[y][x]))
			buff.cells[y][x] = ansi.Cell{
				Fg:      ansi.ForegroundColor(*color),
				Content: byte(ch),
			}
		}
	}
	return buff
}

type Renderable interface {
	Render(r *Renderer) (string, string)
	Tick(*Renderer)
	GetPos() ansi.Pos
	DefaultColor() *int
}

type Renderer struct {
	paused bool
	canvas Canvas

	entities []Renderable
	// Trails or bubbles that fish can leave behind
	fleeting []Renderable
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
}

func (r *Renderer) Draw() {
	for _, e := range r.entities {
		art, colors := e.Render(r)
		c := CanvasFromArt(art, colors)
		r.canvas.MergeAt(c, e.GetPos())
	}
	screen, err := tcell.NewScreen()
  check(err)
	if err = screen.Init(); err != nil {
    panic(err)
	}
	defer screen.Fini()

	// print each cell of final canvas
	// ansi.Print(ansi.MoveHome)
	for y, line := range r.canvas.cells {
		for x, _ := range line {
      screen.SetContent(x,y, 'X', nil, tcell.StyleDefault)
			// if cell.Content == 0 {
			// 	buff += " "
			// }
			//
			// buff += ansi.Color(int(cell.Fg))
			// buff += string(cell.Content)

		}
	}
  screen.Show()

	// for _, e := range r.entities {
	// 	art, colors := e.Render()
	//
	// 	rendered, printPos := cutAndColorize(art, colors, e.DefaultColor(), e.GetPos())
	// 	if rendered == nil {
	// 		continue
	// 	}
	//
	// 	lines := strings.Split(*rendered, "\n")
	// 	ansi.PrintLines(printPos, lines)
	// }
}

// returns nil if the rendered string is completely out of view
// func cutAndColorize(art string, colors string, base *int, pos ansi.Pos) (*string, ansi.Pos) {
// 	cutArt, printPos := cutVisible(art, pos)
// 	cutColors, _ := cutVisible(colors, pos)
// 	if cutArt == nil || cutColors == nil {
// 		return nil, printPos
// 	}
// 	colored := ansi.ColorizeArt(*cutArt, *cutColors, *base)
// 	return &colored, printPos
// }
//
// func cutVisible(content string, pos ansi.Pos) (*string, ansi.Pos) {
// 	width, height, err := term.GetSize(int(os.Stdin.Fd()))
// 	if err != nil {
// 		panic("unable to get terminal size")
// 	}
//
// 	lines := strings.Split(content, "\n")
// 	lineLen := len(lines[0]) - 1
// 	lastIdx := len(lines) - 1
// 	printPos := pos
//
// 	if pos.Y < 0 {
// 		skip := pos.Y * -1
// 		printPos.Y = 0
// 		lines = slices.Delete(lines, 0, skip)
// 	}
//
// 	if (pos.Y + lastIdx) > height {
// 		skip := (pos.Y + lastIdx) - height
// 		if skip > lastIdx {
// 			return nil, printPos
// 		}
// 		start := lastIdx - skip
// 		lines = slices.Delete(lines, start, lastIdx)
// 	}
//
// 	if pos.X < 0 {
// 		start := pos.X * -1
// 		if start > lineLen {
// 			return nil, printPos
// 		}
// 		printPos.X = 0
// 		for idx, line := range lines {
// 			lines[idx] = line[start:]
// 		}
// 	} else if (pos.X + lineLen) > width {
// 		offscreenCells := lineLen - (width-(pos.X+lineLen))*-1
// 		if offscreenCells < 0 {
// 			return nil, printPos
// 		}
// 		for idx, line := range lines {
// 			lines[idx] = line[:offscreenCells]
// 		}
// 	}
//
// 	res := strings.Join(lines, "\n")
// 	return &res, printPos
// }
