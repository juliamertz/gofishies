package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	return []Entity{
		&Castle{Pos: Pos{X: width - 34, Y: height - 14}},
		&Fish{variation: 0, direction: Right, Pos: Pos{X: width / 3, Y: 15}},
		&Fish{variation: 1, direction: Right, Pos: Pos{X: 20, Y: 20}},
		&Fish{variation: 2, direction: Right, Pos: Pos{X: 2, Y: 10}},
		&Whale{direction: Left, Pos: Pos{X: width - 5, Y: 20}},
		&Seaweed{Pos: Pos{X: 10, Y: height - 5}, length: 5},
		&Seaweed{Pos: Pos{X: 13, Y: height - 3}, length: 4},
		&Waves{Pos: Pos{X: 0, Y: 5}},
		&Duck{direction: Right, Pos: Pos{X: 5, Y: 2}},
		&Boat{variation: 1, Pos: Pos{X: width - 30, Y: 0}},
	}
}

type Spawnable interface {
	Spawn(*Renderer)
	Clone() Spawnable
}

type Spawner struct {
	renderer *Renderer
	pool     []Spawnable
}

func (s *Spawner) spawnRandom() {
	i := RNG.IntN(len(s.pool))
	if i > len(s.pool) {
		panic("random number is bigger than pool")
	}

	s.pool[i].Clone().Spawn(s.renderer)
}

func main() {
	screen, err := tcell.NewScreen()
	check(err)
	err = screen.Init()
	check(err)

	r := Renderer{
		tickRate: 1,
		seaLevel: 5,
		paused:   false,
		screen:   screen,
		entities: mkSea(screen.Size()),
	}

	s := Spawner{
		renderer: &r,
		pool:     []Spawnable{&Fish{}, &Whale{}},
	}

	defer r.screen.Fini()
	go eventHandler(&r, &s)

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

	// create empty canvas
	width, height := r.screen.Size()
	r.canvas = NewCanvas(width, height)

	// draw all entities to canvas
	err := r.Draw()
	check(err)

	if r.debug {
		ser, err := json.MarshalIndent(r.entities, "", "  ")
		check(err)
		lines := []string{
			fmt.Sprintf("entities: %d", len(r.entities)),
			fmt.Sprintf("tickRate: %d", r.tickRate),
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

func eventHandler(r *Renderer, s *Spawner) {
	for {
		ev := r.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			r.entities = mkSea(r.screen.Size())
		case *tcell.EventKey:
      if ev.Key() == tcell.Key(3) {
				r.screen.Fini()
				os.Exit(0)
      }
      // this doesn't work for some reason...
			// switch ev.Key() {
			// case tcell.KeyEscape | tcell.KeyCtrlC:
			// }

			switch ev.Rune() {
			case 's':
				s.spawnRandom()
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
