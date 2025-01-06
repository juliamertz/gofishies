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
	return tcell.ColorGray
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
 |_______|__|_|_|_|__|_______|`

	colors := `
                rr           
                             
              yyy            
             y   y           
            y     y          
           y       y         
                             
                             
                             
              yyy            
             yy yy           
            y y y y          
            yyyyyyy           `

	return art, colors
}

func (c *Castle) GetPos() Pos {
	return c.Pos
}

func (c *Castle) Tick(r *Renderer) {}

// Boat

type Boat struct {
	Pos       Pos
	ticks     int
	variation int
	direction Direction
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
        ___\__                                                                                           
       |______\_____                                                                                     
 _____/_______/_____\_______                                                                             
|       > > >              /`
		colors = `
        ___\__                                                                                           
       |______\_____                                                                                     
 _____/_______/_____\_______                                                                             
|ddddddd>d>d>dddddddddddddd/`

	case 1:
		art = `
   |    |    |                 
   )_)  )_)  )_)              
  )___))___))___)\            
 )____)____)_____)\\
_____|____|____|____\\\__
\                   /`
		colors = `
   w    w    w                 
   www  www  www              
  wwwwwwwwwwwwwww\            
 wwwwwwwwwwwwwwwww\\
_____|____|____|____\\\__
\ddddddddddddddddddd/`
	}

	if b.direction == Left {
		return flipArt(art, colors)
	}
	return art, colors
}

func (b *Boat) GetPos() Pos {
	return b.Pos
}

func (b *Boat) Tick(r *Renderer) {
	if b.ticks > 20 {
		b.ticks = 0
		if b.direction == Left {
			b.Pos.X--
		} else {
			b.Pos.X++
		}
	} else {
		b.ticks++
	}
}
