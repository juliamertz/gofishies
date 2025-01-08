package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func mkSea(width int, height int) []Entity {
	entities := []Entity{
		Waves(Pos{Y: 5}, width),
		Castle(Left, Pos{X: width - 35, Y: height - 14}),
		Fish(1, Right, tcell.ColorBlue, Pos{X: 2, Y: 15}),
		Fish(2, Right, tcell.ColorYellow, Pos{X: 40, Y: 24}),
		Boat(1, Right, Pos{X: 10, Y: 0}),
		Duck(Right, Pos{X: 5, Y: 5}),
		Duck(Right, Pos{X: 15, Y: 5}),
		Duck(Right, Pos{X: 25, Y: 5}),
		Whale(Right, Pos{X: 5, Y: 15}),
	}

	// spawn some seaweed
	maxInitialHeight := 8
	plantsAmount := width / 10
	var plants []Entity
	for i := 0; i < plantsAmount; i++ {
		x := RNG.IntN(width)
	inner:
		for {
			if slices.ContainsFunc(plants, func(n Entity) bool {
				return n.pos.X == x
			}) {
				x = RNG.IntN(width)
			} else {
				break inner
			}
		}
		pHeight := RNG.IntN(maxInitialHeight)
		pos := Pos{X: x, Y: height - pHeight}
		plants = append(plants, Seaweed(pHeight, pos))
	}

	return append(plants, entities...)
}

type EntityCap struct {
	smallFish int
	largeFish int
	boat      int
}

func main() {
	screen, err := tcell.NewScreen()
	check(err)
	err = screen.Init()
	check(err)
	r := Renderer{
		tickRate: 2,
		seaLevel: 5,
		paused:   false,
		screen:   screen,
		entities: mkSea(screen.Size()),
	}

	defer r.screen.Fini()
	go eventHandler(&r)

	for {
		if r.paused {
			time.Sleep(time.Duration(r.tickRate) * time.Millisecond)
			continue
		}

		drawCurrent(&r)
		r.Tick()
	}
}

func drawCurrent(r *Renderer) {
	// update all entities
	r.screen.Clear()

	start := time.Now()

	// create empty canvas
	width, height := r.screen.Size()
	r.frame = makeFrame(width, height)

	// draw all entities to canvas
	err := r.Draw()
	check(err)

	elapsed := time.Now().Sub(start)

	if r.debug {
		ser, err := json.MarshalIndent(r.entities, "", "  ")
		check(err)
		lines := []string{
			fmt.Sprintf("entities: %d", len(r.entities)),
			fmt.Sprintf("tickRate: %d", r.tickRate),
			fmt.Sprintf("lastRender: %v", elapsed),
			fmt.Sprintf("entities: %s", string(ser)),
		}
		for i, line := range lines {
			r.DrawText(line, Pos{Y: i})
		}
	}

	r.screen.Show()
	time.Sleep(time.Duration(10/r.tickRate) * time.Millisecond)
}

func join(lines []string) string {
	return strings.Join(lines, "\n")
}

func eventHandler(r *Renderer) {
	for {
		ev := r.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			r.entities = mkSea(r.screen.Size())
		case *tcell.EventKey:
			// FIX: this doesn't work for some reason...
			// switch ev.Key() {
			// case tcell.KeyEscape | tcell.KeyCtrlC:
			// }
			if ev.Key() == tcell.Key(3) {
				r.screen.Fini()
				os.Exit(0)
			}

			switch ev.Rune() {
			case 's':
				r.spawnRandomEntity()
			case 'r':
				r.entities = mkSea(r.screen.Size())
			case 'q':
				r.screen.Fini()
				os.Exit(0)
			case ' ':
				r.paused = !r.paused
			case 'x':
				r.debug = !r.debug
				drawCurrent(r)
				r.screen.Show()
			case 'j':
				if r.tickRate <= 1 {
					continue
				}
				if r.tickRate < 10 {
					r.tickRate -= 1
				} else {
					r.tickRate -= 5
				}
			case 'k':
				if r.tickRate < 10 {
					r.tickRate += 1
				} else {
					r.tickRate += 5
				}
			}
		}
	}
}
