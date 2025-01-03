package main

import (
	"math/rand/v2"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Seaweed struct {
	Pos    Pos
	length int
	stage  bool
}

func (s *Seaweed) DefaultColor() tcell.Color {
	return tcell.ColorGreen
}

func (s *Seaweed) Render(r *Renderer) (string, string) {
	var art []string
	var colors []string
	stage := s.stage
	for i := 0; i < s.length; i++ {
		if stage {
			art = append(art, " )")
			art = append(art, "( ")
		} else {
			art = append(art, "( ")
			art = append(art, " )")
		}
		stage = !stage
		colors = append(colors, "gg")
		colors = append(colors, "gg")
	}

	compareArtStrings(join(art), join(colors))

	return join(art), join(colors)
}

func (s *Seaweed) GetPos() Pos {
	return s.Pos
}

func (s *Seaweed) Tick(r *Renderer) {
	src := rand.NewPCG(uint64(s.length), uint64(time.Now().UnixNano()))
	rng := rand.New(src)
	n := rng.IntN(60)

	if n == 1 {
		s.stage = !s.stage
	}
}
