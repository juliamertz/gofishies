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
	Pos       Pos
	tick      int
	direction Direction
}

func (w *Whale) DefaultColor() tcell.Color {
	return tcell.ColorGreen
}

func (w *Whale) Render(r *Renderer) (string, string) {
	art := join([]string{
		"   ______/~/~/~/__           /((",
		" // __            ====__    /_((",
		"//  @))       ))))      ===/__((",
		"))           )))))))        __((",
		"\\\\     \\)     ))))    __===\\ _((",
		" \\\\_______________====      \\_((",
		"                             \\((",
	})

	colors := join([]string{
		"   ______/~/~/~/__           /((",
		" //d__dddddddddddd====__    /_(w",
		"//ddwggddddddd))))dddddd===/__(w",
		"))ddddddddddd)))))))dddddddd__(w",
		"wwdddddwwdddddwwwwssss__===\\ _(w",
		" wwwwwwwwwwwwwwwwwwwww      \\_(w",
		"                             \\(( ",
	})

	if w.direction == Right {
		art = flipAsciiArt(art)
		colors = reverseArt(colors)
	}

	return art, colors
}

func (w *Whale) GetPos() Pos {
	return w.Pos
}

func (w *Whale) Spawn(r *Renderer) {
	cols, lines := r.screen.Size()
	height := rand.IntN(lines - r.seaLevel - 5)
	dir := Direction(rand.IntN(2))
	w.direction = dir
	if w.direction == Left {
		w.Pos.X = cols - 5
	}
	w.Pos.Y = r.seaLevel + height
	r.entities = append(r.entities, w)
}

func (w *Whale) Tick(r *Renderer) {
	if w.tick == 10 {
		w.tick = 0
		switch w.direction {
		case Left:
			w.Pos.X--
		case Right:
			w.Pos.X++
		}

		if RandOneIn(20) {
			r.entities = append(r.entities, &Bubble{Pos: w.Pos})
		}
	} else {
		w.tick++
	}
}

func (w *Whale) Clone() Spawnable {
	return &Whale{
		direction: w.direction,
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
      /
  ,../...
 /       '\  /
< '  )rwx  =<
 \ \      /  \
  ''\'"'"'`
		colors = `
      /
  ,../...
 /       '\  /
< '  )rwx  =<
 \ \      /  \
  ''\'"'"'`

	}

	if f.direction == Right {
		art = flipAsciiArt(art)
		colors = reverseArt(colors)
	}

	return art, colors
}

func (f *Fish) GetPos() Pos {
	return f.Pos
}

func (f *Fish) Tick(r *Renderer) {
	if f.cycle == 10 {
		if RandOneIn(20) {
			r.entities = append(r.entities, &Bubble{Pos: f.Pos})
		}
		f.cycle = 0
		switch f.direction {
		case Left:
			f.Pos.X--
		case Right:
			f.Pos.X++
		}
	} else {
		f.cycle++
	}
}

func (f *Fish) Spawn(r *Renderer) {
	_, lines := r.screen.Size()
	height := rand.IntN(lines - r.seaLevel)
	f.variation = rand.IntN(3)
	f.direction = Direction(rand.IntN(2))
	f.Pos = Pos{Y: r.seaLevel + height}
	r.entities = append(r.entities, f)
}

func (f *Fish) Clone() Spawnable {
	return &Fish{
		variation: f.variation,
		direction: f.direction,
	}
}
