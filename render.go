package main

import (
	"gofishies/ansi"
	"os"
	"slices"
	"strings"

	"golang.org/x/term"
)

type Renderable interface {
	Render() (string, string)
	Tick(*Renderer)
	GetPos() ansi.Pos
	DefaultColor() *int
	// setPos(Pos)
}

// func ()

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Renderer struct {
	paused bool

	entities []Renderable
	// Trails or bubbles that fish can leave behind
	fleeting []Renderable
}

func (r *Renderer) Tick() {
	for _, e := range r.entities {
		e.Tick(r)
	}
}

func (r *Renderer) Draw() {
	for _, e := range r.entities {
		art, colors := e.Render()

		rendered, printPos := cutAndColorize(art, colors, e.DefaultColor(), e.GetPos())
		if rendered == nil {
			continue
		}

		lines := strings.Split(*rendered, "\n")
		ansi.PrintLines(printPos, lines)
	}
}

// returns nil if the rendered string is completely out of view
func cutAndColorize(art string, colors string, base *int, pos ansi.Pos) (*string, ansi.Pos) {
	cutArt, printPos := cutVisible(art, pos)
	cutColors, _ := cutVisible(colors, pos)
	if cutArt == nil || cutColors == nil {
		return nil, printPos
	}
	colored := ansi.ColorizeArt(*cutArt, *cutColors, *base)
	return &colored, printPos
}

func cutVisible(content string, pos ansi.Pos) (*string, ansi.Pos) {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic("unable to get terminal size")
	}

	lines := strings.Split(content, "\n")
	lineLen := len(lines[0]) - 1
	lastIdx := len(lines) - 1
	printPos := pos

	if pos.Y < 0 {
		skip := pos.Y * -1
		printPos.Y = 0
		lines = slices.Delete(lines, 0, skip)
	}

	if (pos.Y + lastIdx) > height {
		skip := (pos.Y + lastIdx) - height
		if skip > lastIdx {
			return nil, printPos
		}
		start := lastIdx - skip
		lines = slices.Delete(lines, start, lastIdx)
	}

	if pos.X < 0 {
		start := pos.X * -1
		if start > lineLen {
			return nil, printPos
		}
		printPos.X = 0
		for idx, line := range lines {
			lines[idx] = line[start:]
		}
	} else if (pos.X + lineLen) > width {
		offscreenCells := lineLen - (width-(pos.X+lineLen))*-1
    if offscreenCells < 0 {
      return nil, printPos
    }
		for idx, line := range lines {
			lines[idx] = line[:offscreenCells]
		}
	}

	res := strings.Join(lines, "\n")
	return &res, printPos
}
