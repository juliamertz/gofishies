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

type Entity struct {
	Id           string
	Pos          Pos
	Facing       Direction
	CurrentFrame int
	Tick         int

	defaultColor tcell.Color
	frames       []Frame

	width  int
	height int

	update func(*Entity, *Renderer)
}

func (e *Entity) Move(dir Direction) {
	switch dir {
	case Left:
		e.Pos.X--
	case Right:
		e.Pos.X++
	}
}

func (e *Entity) NextFrame() {
	e.CurrentFrame++
	if e.CurrentFrame >= len(e.frames) {
		e.CurrentFrame = 0
	}
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

func createEntity(
	id string,
	art []string,
	colorMap []string,
	defaultColor tcell.Color,
	pos Pos,
	facing Direction,
	update func(*Entity, *Renderer),
) Entity {
	var frames []Frame
	// generate frames for entity
	for i, frame := range art {
		var cMap string
		if i < len(colorMap) {
			cMap = colorMap[i]
		} else {
			// fallback to first entry if frame doesn't have corresponding color map
			cMap = colorMap[0]
		}

		// assume all art is left facing
		if facing == Right {
			frame = mirrorAsciiArt(frame)
			cMap = reverseArt(cMap)
		}

		frames = append(frames, generateFrame(frame, cMap, defaultColor))
	}

	return Entity{
		Id:           id,
		Pos:          pos,
		Facing:       facing,
		defaultColor: defaultColor,
		frames:       frames,
		update:       update,

		// assume all frames are of same height/width
		height: len(frames[0].cells),
		width:  len(frames[0].cells[0]),
	}
}

type Renderer struct {
	screen tcell.Screen
	frame  Frame

	debug       bool
	maxEntities uint32
	paused      bool
	seaLevel    int
	tickRate    int

	entities []Entity
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

func (r *Renderer) IsOffscreen(e Entity) bool {
	cols, lines := r.screen.Size()

	if e.Pos.X >= cols || e.Pos.X+e.width <= 0 {
		return true
	}

	if e.Pos.Y > lines || e.Pos.Y+e.height <= 0 {
		return true
	}

	return false
}

func (r *Renderer) KillEntity(v Entity) {
	for i, e := range r.entities {
		// TODO: find better way to determine uniqueness
		if e.Id == v.Id && e.Pos == v.Pos {
			r.entities = slices.Delete(r.entities, i, i+1)
			break
		}
	}
}

func (r *Renderer) Tick() {
	for i, e := range slices.Backward(r.entities) {
		// TODO: figure out why it doesn't update if `e` is passed instead if indexing
		r.entities[i].Tick++
		if r.IsOffscreen(e) {
			r.KillEntity(e)
			continue
		}

		e.update(&r.entities[i], r)
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
		// TODO: figure out why this is needed
		if e.CurrentFrame > len(e.frames)-1 {
			continue
		}
		r.frame.MergeAt(e.frames[e.CurrentFrame], e.Pos)
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
