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
		// Fish(1, Right, tcell.ColorBlue, Pos{X: 2, Y: 15}),
		// Fish(2, Right, tcell.ColorYellow, Pos{X: 40, Y: 24}),
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

type EntityKind int

const (
	SmallFish EntityKind = iota
	LargeFish
	Vehicle
)

type EntityCaps struct {
	smallFish int
	largeFish int
	vehicle   int
}

// get random entity kind that we can spawn
func (c *EntityCaps) GetKind() *EntityKind {
  // TODO: better way to prevent too many cycles
	for i := 0; i < 9; i++ {
		kind := EntityKind(RNG.IntN(2))
		switch kind {
		case SmallFish:
			if c.smallFish >= 10 {
				continue
			}
		case LargeFish:
			if c.largeFish >= 2 {
				continue
			}
		case Vehicle:
			if c.vehicle >= 1 {
				continue
			}

		}
		return &kind
	}

	// if c.smallFish < 10 {
	//   return SmallFish
	// }
	//
  return nil
}

type Spawner struct {
	caps     EntityCaps
	renderer *Renderer
}

func (s *Spawner) spawnRandomEntity() {
	sWidth, sHeight := s.renderer.screen.Size()

	// assume small fish for now
	facing := Direction(RNG.IntN(2))
	var x int
	switch facing {
	case Left:
		x = sWidth - 5
	case Right:
		x = -5
	}

	f := Fish(RNG.IntN(6), facing, RandColor(), Pos{
		Y: s.renderer.seaLevel*2 + RNG.IntN(sHeight-s.renderer.seaLevel),
		X: x,
	})

	f.Id = fmt.Sprintf("%s_%d", f.Id, len(s.renderer.entities))
	s.renderer.SpawnEntity(f)
}

func (s *Spawner) spawnEntity(e Entity) {
	s.renderer.entities = append(s.renderer.entities, e)
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

	s := Spawner{
		renderer: &r,
		caps:     EntityCaps{},
	}

	defer r.screen.Fini()
	go eventHandler(&r, &s)

	for {
		if r.paused {
			time.Sleep(time.Duration(r.tickRate) * time.Millisecond)
			continue
		}

		// TODO:
		if RandOneIn(100) {
			s.spawnRandomEntity()
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

	if r.debug {
		elapsed := time.Now().Sub(start)
		fps := 1 / elapsed.Seconds()
		ser, err := json.MarshalIndent(r.entities, "", "  ")
		check(err)
		lines := []string{
			fmt.Sprintf("entities: %d", len(r.entities)),
			fmt.Sprintf("tickRate: %d", r.tickRate),
			fmt.Sprintf("fps: %d", int(fps)),
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
				s.spawnRandomEntity()
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
