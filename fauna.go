package main

import (
	"math/rand/v2"

	"github.com/gdamore/tcell/v2"
)

type Direction = int

const (
	Left  Direction = 0
	Right Direction = 1
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
     // __            ====__    /_(w
    //  wgg       ))))      ===/__(w
    ))           )))))))        __(w
    ww     ww     wwww    __===\ _(w
     wwwwwwwwwwwwwwwwwwwww      \_(w
                                 \((`

	return art, colors
}

func (w *Whale) GetPos() Pos {
	return w.Pos
}

func (w *Whale) Spawn(r *Renderer) {
	_, lines := r.screen.Size()
	height := rand.IntN(lines - r.seaLevel)
	w.Pos = Pos{Y: r.seaLevel + height}
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

type Fish struct {
	Pos       Pos
	cycle     int
	variation int
	direction Direction
}

func (f *Fish) DefaultColor() tcell.Color {
	return tcell.ColorOrange
}

func (f *Fish) Spawn(r *Renderer) {
	_, lines := r.screen.Size()
	height := rand.IntN(lines - r.seaLevel)
	f.Pos = Pos{Y: r.seaLevel + height}
	r.entities = append(r.entities, f)
}

func (f *Fish) Render(r *Renderer) (string, string) {
	var art string
	var colors string

	switch f.variation {
	case 0:
		art = `
 _ 
<_><
`
		colors = `
 y 
y   
`
	case 1:
		art = `
      .:/
   ,,///;,   ,;/
 o:::::::;;///
>::::::::;;\\\\\
  ''\\\\\\\\\'' ';\
`
		colors = `
      .r/
   ,,///;,   ,;/
 orrrrrrr;;///
>rrrrrrrr;;\\\\\
  ''\\\\\\\\\'' ';\
`
	case 2:
		art = `
`
		colors = `
`

	}

	// if f.cycle < 5 {
	// 	art = []string{
	// 		"  _ ",
	// 		"><_>",
	// 	}
	// } else {
	// 	art = []string{
	// 		"  _ ",
	// 		"~<_>",
	// 	}
	// }

	// colors := []string{
	// 	"  c ",
	// 	"yycc",
	// }

  if f.direction == Right {
    art = flipAsciiArt(art)
    colors = reverseArt(colors)
  }

	return art, colors
}

func (g *Fish) GetPos() Pos {
	return g.Pos
}

func (g *Fish) Tick(r *Renderer) {
	if g.cycle == 10 {
		g.cycle = 0
		switch g.direction {
		case Left:
			g.Pos.X--
		case Right:
			g.Pos.X++
		}
	} else {
		g.cycle++
	}
}
