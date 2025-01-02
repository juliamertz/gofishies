package main

import (
	"gofishies/ansi"
)

type Seaweed struct {
	Pos    ansi.Pos
	length int
	cycle  int
}

func (f *Seaweed) DefaultColor() *int {
	return ansi.ColorFromByte('c')
}

func (f *Seaweed) Render(r *Renderer) (string, string) {
	var art []string
	var colors []string

	for i := 0; i < f.length; i++ {
		if f.cycle < 10 {
			art = append(art, " )")
			art = append(art, "( ")
		} else {
			art = append(art, "( ")
			art = append(art, " )")
		}
		colors = append(colors, "gg")
		colors = append(colors, "gg")
	}

	return join(art), join(colors)
}

func (f *Seaweed) GetPos() ansi.Pos {
	return f.Pos
}

func (f *Seaweed) Tick(r *Renderer) {
	if f.cycle == 20 {
		f.cycle = 0
	} else {
		f.cycle++
	}
}
