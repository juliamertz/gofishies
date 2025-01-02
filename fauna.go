package main

import (
	"github.com/gdamore/tcell/v2"
)

// Whale

type Whale struct {
	Pos   Pos
	cycle int
  tick int
}

func (f *Whale) DefaultColor() tcell.Color {
	return tcell.ColorGreen
}

func (f *Whale) Render(r *Renderer) (string, string) {
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

func (f *Whale) GetPos() Pos {
	return f.Pos
}

func (f *Whale) Tick(r *Renderer) {
	if f.tick == 10 {
		f.tick = 0
		f.Pos.X--
    if RandOneIn(20) {
      r.particles = append(r.particles, &Bubble{Pos: f.Pos})
    }
	} else {
		f.tick++
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
	if f.cycle == 20  {
		f.cycle = 0
		f.Pos.Y--
	} else {
		f.cycle++
	}
}

// Goldfish

type Goldfish struct {
	Pos   Pos
	cycle int
}

func (f *Goldfish) DefaultColor() tcell.Color {
	return tcell.ColorOrange
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
