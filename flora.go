package main

import (
	"github.com/gdamore/tcell/v2"
)

type Seaweed struct {
	Pos    Pos
	length int
	cycle  int
}

func (f *Seaweed) DefaultColor() tcell.Color {
  return tcell.ColorGreen
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

  compareArtStrings(join(art), join(colors))

	return join(art), join(colors)
}

func (f *Seaweed) GetPos() Pos {
	return f.Pos
}

func (f *Seaweed) Tick(r *Renderer) {
	if f.cycle == 20 {
		f.cycle = 0
	} else {
		f.cycle++
	}
}
