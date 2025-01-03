package main

import (
	"github.com/gdamore/tcell/v2"
	"math/rand/v2"
	"strings"
)

type Waves struct {
	Pos   Pos
	cycle int
	ticks int
}

func (f *Waves) DefaultColor() tcell.Color {
	return tcell.ColorGreen
}

func (f *Waves) Render(r *Renderer) (string, string) {
	width, _ := r.screen.Size()

	src := rand.NewPCG(uint64(width), uint64(f.cycle))
	rng := rand.New(src)

	art := strings.Repeat("~", width) + "\n"
	colors := strings.Repeat("b", width) + "\n"
	waveHeight := 3
	for i := 0; i < waveHeight; i++ {
		for j := 0; j < width; j++ {
			r := rng.IntN(4)
			if r == 1 {
				art += "^"
				colors += "b"
			} else {
				art += " "
				colors += "w"
			}
		}
		art += "\n"
		colors += "\n"
	}

	return art, colors
}

func (f *Waves) GetPos() Pos {
	return f.Pos
}

func (f *Waves) Tick(r *Renderer) {
	if f.ticks == 20 {
		f.cycle++
		f.ticks = 0
	} else {
		f.ticks++
	}
}

// Bubble

type Bubble struct {
	Pos   Pos
	cycle int
}

func (f *Bubble) DefaultColor() tcell.Color {
	return tcell.ColorWhite
}

func (f *Bubble) Render(r *Renderer) (string, string) {
	return "o", "w"
}

func (f *Bubble) GetPos() Pos {
	return f.Pos
}

func (f *Bubble) Tick(r *Renderer) {
	if f.cycle == 20 {
		f.cycle = 0
		f.Pos.Y--
	} else {
		f.cycle++
	}
}
