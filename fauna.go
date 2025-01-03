package main

import (
	"math/rand/v2"

	"github.com/gdamore/tcell/v2"
)

// Whale

type Whale struct {
	Pos   Pos
	cycle int
	tick  int
}

func (w *Whale) DefaultColor() tcell.Color {
	return tcell.ColorGreen
}

func (w *Whale) Render(r *Renderer) (string, string) {
	art :=
		`                                    
       ______/~/~/~/__           /((
     // __            ====__    /_((
    //  @))       ))))      ===/__((
    ))           )))))))        __((
    \\     \)     ))))    __===\ _((
     \\_______________====      \_((
                                 \((`

	colors :=
		`                                    
       ______/~/~/~/__           /((
     // __            ====__    /_((
    //  wgg       ))))      ===/__((
    ))           )))))))        __((
    ww     ww     wwww    __===\ _((
     wwwwwwwwwwwwwwwwwwwww      \_((
                                 \((`

	compareArtStrings(art, colors)

	return art, colors
}

func (w *Whale) GetPos() Pos {
	return w.Pos
}

func (w *Whale) Spawn(r *Renderer) {
	_, lines := r.screen.Size()
	height := rand.IntN(lines - r.seaLevel)
	w.Pos = Pos{Y:  r.seaLevel + height}
	r.entities = append(r.entities, w)
}


func (w *Whale) Tick(r *Renderer) {
	if w.tick == 10 {
		w.tick = 0
		w.Pos.X--
		if RandOneIn(20) {
			r.entities = append(r.entities, &Bubble{Pos: w.Pos})
		}
	} else {
		w.tick++
	}
}

// Goldfish

type Goldfish struct {
	Pos   Pos
	cycle int
}

func (g *Goldfish) DefaultColor() tcell.Color {
	return tcell.ColorOrange
}

func (g *Goldfish) Spawn(r *Renderer) {
	_, lines := r.screen.Size()
	height := rand.IntN(lines - r.seaLevel)
	g.Pos = Pos{Y:  r.seaLevel + height}
	r.entities = append(r.entities, g)
}

func (g *Goldfish) Render(r *Renderer) (string, string) {
	var art []string
	if g.cycle < 5 {
		art = []string{
			"  _ ",
			"><_>",
		}
	} else {
		art = []string{
			"  _ ",
			"~<_>",
		}
	}

	colors := []string{
		"  c ",
		"yycc",
	}

	compareArtStrings(join(art), join(colors))

	return join(art), join(colors)
}

func (g *Goldfish) GetPos() Pos {
	return g.Pos
}

func (g *Goldfish) Tick(r *Renderer) {
	if g.cycle == 10 {
		g.cycle = 0
		g.Pos.X++
	} else {
		g.cycle++
	}
}
