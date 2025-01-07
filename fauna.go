package main

import (
	"github.com/gdamore/tcell/v2"
)

func Fish(variation int, facing Direction, pos Pos) Entity {
	var art string
	var colors string

	switch variation {
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
 /ddddddd'\  /
<d'dd)rgydd=<
 \d\dddddd/  \
  ''\'"'"'`
	case 3:
		art = `
    /
 /--\ /
<o)  =<
 \__/ \
  \
`
		colors = `
    /
 /--\ /
<wbdd=<
 \__/ \
  \
  `
	case 4:
		art = join([]string{
			"  ,/..   ",
			"<')   `=<",
			" ``\\``   ",
		})
		colors = join([]string{
			"  yyrr   ",
			"yw)ddd`=<",
			" yyyyy   ",
		})
	case 5:
		art = join([]string{
" /,  ",
"<')=<",
" \\`  ",
		})
		colors = join([]string{
" y,  ",
"<w)=<",
" y`  ",
		})
  case 6:
    art = `
 __
/o \/
\__/\
    `
    colors = `
 __
/w \/
\__/\
    `
	}

	return createEntity(
		"fish",
		[]string{art},
		[]string{colors},
		tcell.ColorOrange,
		pos,
		facing,
		func(e *Entity, r *Renderer) {
			if e.Tick%10 == 0 {
				e.Move(e.Facing)
				if RandOneIn(30) {
					r.SpawnEntity(Bubble(e.LikelyBubblePos()))
				}
			}
		},
	)
}

func Whale(facing Direction, pos Pos) Entity {
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
		"wwdddddwwdddddwwwwdddd__===\\d_(w",
		" wwwwwwwwwwwwwwwwwwwww      \\_(w",
		"                             \\((",
	})

	return createEntity(
		"whale",
		[]string{art},
		[]string{colors},
		tcell.ColorGreen,
		pos,
		facing,
		func(e *Entity, r *Renderer) {
			if e.Tick%20 == 0 {
				e.Move(e.Facing)
			}

			if RandOneIn(200) {
				r.SpawnEntity(Bubble(e.LikelyBubblePos()))
			}
		},
	)
}

func Duck(facing Direction, pos Pos) Entity {
	art := []string{
		`
  _
=(')____,
 (' =~~/`,
		`
  _
>(')____,
 (' =~~/`,
	}

	colors :=
		`
  _
y(')____,
 ('d=~~/`

	return createEntity(
		"duck",
		art,
		[]string{colors},
		tcell.ColorWhite,
		pos,
		facing,
		func(e *Entity, r *Renderer) {
			if e.Tick%30 == 0 {
				e.Move(e.Facing)
			}

			if e.Tick%20 == 0 {
				e.NextFrame()
			}
		},
	)
}
