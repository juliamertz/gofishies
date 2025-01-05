package main

import (
	"math/rand/v2"
	"strings"

	"github.com/gdamore/tcell/v2"
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
	ticks int
	stage int
}

func (f *Bubble) DefaultColor() tcell.Color {
	return tcell.ColorWhite
}

func (f *Bubble) Render(r *Renderer) (string, string) {
	var art string
	switch f.stage {
	case 0:
		art = "."
	case 1:
		art = "o"
	}

	return art, "w"
}

func (f *Bubble) GetPos() Pos {
	return f.Pos
}

func (b *Bubble) Tick(r *Renderer) {
	if b.Pos.Y < r.seaLevel+3 {
		r.KillEntity(b)
	}
	if b.ticks%20 == 0 {
		b.Pos.Y--
	}
	if b.ticks > 200 {
		b.stage = 1
	}
	b.ticks++
}

// Castle

type Castle struct {
	Pos Pos
}

func (f *Castle) DefaultColor() tcell.Color {
	return tcell.ColorWhite
}

func (f *Castle) Render(r *Renderer) (string, string) {
	art := `
               T~~
               |
              /^\
             /   \
 _   _   _  /     \  _   _   _
[ ]_[ ]_[ ]/ _   _ \[ ]_[ ]_[ ]
|_=__-_ =_|_[ ]_[ ]_|_=-___-__|
 | _- =  | =_ = _    |= _=   |
 |= -[]  |- = _ =    |_-=_[] |
 | =_    |= - ___    | =_ =  |
 |=  []- |-  /| |\   |=_ =[] |
 |- =_   | =| | | |  |- = -  |
 |_______|__|_|_|_|__|_______|
  `

	colors := `
                rr           
                             
              yyy            
             ydddy           
  ddd ddd ddydddddyd ddd ddd 
 d   d   d yd ddd dyd   d   d
       d     d   d           
  d  d dd d  d d dddd  d  ddd
   d   dd  d d d dddd       d
  d  dddd  d dyyydddd d  d dd
   dd   d  ddyydyydd   d   d 
   d  ddd d ydydydydd  d d dd
            yyyyyyy           `

	return art, colors
}

func (f *Castle) GetPos() Pos {
	return f.Pos
}

func (f *Castle) Tick(r *Renderer) {
}

// Boat

type Boat struct {
	Pos       Pos
	ticks     int
	variation int
}

func (f *Boat) DefaultColor() tcell.Color {
	return tcell.ColorGray
}

func (b *Boat) Render(r *Renderer) (string, string) {
	var art string
	var colors string

	switch b.variation {
	case 0:
		art = `
              __/___            
        _____/______|           
_______/_____\_______\_____     
\              < < <       |`
		colors = `
              wwwwww            
        wwwwwwwwwwwww           
_______w_____w_______w_____     
\dddddddddddddd<d<d<ddddddd|`

	case 1:
		art = `
   |    |    |                 
   )_)  )_)  )_)              
  )___))___))___)\            
 )____)____)_____)\\
_____|____|____|____\\\__
\                   /`
		colors = `
   |    |    |                 
   )_)  )_)  )_)              
  )___))___))___)\            
 )____)____)_____)\\
_____|____|____|____\\\__
\                   /`
	}

	return art, colors
}

func (f *Boat) GetPos() Pos {
	return f.Pos
}

func (f *Boat) Tick(r *Renderer) {
	if f.ticks > 20 {
		f.ticks = 0
		f.Pos.X--
	} else {
		f.ticks++
	}
}
