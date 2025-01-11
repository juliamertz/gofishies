package main

import (
	"github.com/gdamore/tcell/v2"
)

func Seaweed(length int, pos Pos) Entity {
	var parts []string
	var colorParts []string

	for i := 0; i < length; i++ {
		if i%2 == 0 {
			parts = append(parts, "( ")
			colorParts = append(colorParts, "g ")
		} else {
			parts = append(parts, " )")
			colorParts = append(colorParts, " g")
		}
	}

	art := []string{join(parts), mirrorAsciiArt(join(parts))}

	return createEntity(
		"seaweed",
		Other,
		art,
		[]string{join(colorParts)},
		tcell.ColorGreen,
		pos,

		Right, // n.a.
		func(e *Entity, r *Engine) {
			if e.Tick%40 == 0 {
				e.NextFrame()
			}
		},
	)
}
