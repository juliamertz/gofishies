package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

func Waves(pos Pos, width int) Entity {
	art := strings.Repeat("~", width) + "\n"
	colors := strings.Repeat("b", width) + "\n"
	waveHeight := 3
	for i := 0; i < waveHeight; i++ {
		for j := 0; j < width; j++ {
			r := RNG.IntN(4)
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

	return createEntity(
		Other,
		[]string{art},
		[]string{colors},
		tcell.ColorWhite,
		pos,
		Right, // irrelevant for waves
		func(e *Entity, r *Engine) {},
	)
}

func Bubble(pos Pos) Entity {
	return createEntity(
		Other,
		[]string{".", "o", "O"},
		[]string{"b"},
		tcell.ColorGray,
		pos,
		Left,
		func(e *Entity, r *Engine) {
			// kill bubble if it approaches the sealine
			if e.pos.Y < r.seaLevel+3 {
				e.shouldKill = true
				return
			}
			// kill bubble in case it collides with another entity
			collisions := e.GetCollisions(r)
			// also make sure the bubble isn't killed on spawn with a grace period
			if len(collisions) > 0 && e.Tick > 100 {
				e.shouldKill = true
			}

			if e.Tick%20 == 0 {
				e.pos.Y--
			}
			if e.Tick > 75 {
				e.currentFrame = 1
			} else if e.Tick > 150 {
				e.currentFrame = 2
			}
		},
	)
}

func Castle(facing Direction, pos Pos) Entity {
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

	return createEntity(
		Other,
		[]string{art},
		[]string{colors},
		tcell.ColorGray,
		pos,
		facing,
		func(e *Entity, r *Engine) {},
	)
}

func Boat(variation int, facing Direction, pos Pos) Entity {
	var art string
	var colors string

	switch variation {
	case 0:
		art = `
              __/___
        _____/______|
_______/_____\_______\_____
\              < < <       |`
		colors = `
              __/___
        _____/______|
_______/_____\_______\_____
\dddddddddddddd<d<d<ddddddd|`

	case 1:
		art = `
           |    |    |
         (_(  (_(  (_(
       /(___((___((___(
     //(_____(____(____(
__///____|____|____|_____
    \                   /`
		colors = `
           y    y    y
         www  www  www
       wwwwwwwwwwwwwwww
     wwwwwwwwwwwwwwwwwww
yywwwyyyyyyyyyyyyyyyyyyyy
    ydddddddddddddddddddy`
	}

	return createEntity(
		Vehicle,
		[]string{art},
		[]string{colors},
		tcell.ColorGray,
		pos,
		facing,
		func(e *Entity, r *Engine) {
			if e.Tick%20 == 0 {
				e.Move(e.Facing)
			}
		},
	)
}
