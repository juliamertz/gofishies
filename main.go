package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
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
		Boat(1, Right, Pos{X: 10, Y: 0}),
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

// get random entity kind that we can spawn
func (c *EntityCaps) GetKind() *EntityKind {
	// TODO: better way to prevent too many cycles
	for i := 0; i < 9; i++ {
		kind := EntityKind(RNG.IntN(3))
		switch kind {
		case SmallFish:
			if c.SmallFish >= 10 {
				continue
			}
		case LargeFish:
			if c.LargeFish >= 2 {
				continue
			}
		case Vehicle:
			if c.Vehicle >= 1 {
				continue
			}
		case WaterLine:
			if c.WaterLine >= 2 {
				continue
			}

		}
		return &kind
	}

	return nil
}

type Spawner struct {
	caps EntityCaps
}

func (e *Engine) spawnRandomEntity() {
	kind := e.spawner.caps.GetKind()
	if kind == nil {
		return
	}

	sWidth, sHeight := e.screen.Size()
	facing := Direction(RNG.IntN(2))
	var x int
	switch facing {
	case Left:
		x = sWidth - 5
	case Right:
		x = -5
	}
	pos := Pos{
		Y: e.seaLevel*2 + RNG.IntN(sHeight-e.seaLevel),
		X: x,
	}

	switch *kind {
	case SmallFish:
		f := Fish(RNG.IntN(6), facing, RandColor(), pos)
		f.Id = fmt.Sprintf("%s_%d", f.Id, len(e.entities))
		e.SpawnEntity(f)
	case LargeFish:
		e.SpawnEntity(Whale(facing, pos))

	case WaterLine:
		pos.Y = e.seaLevel
		e.SpawnEntity(Duck(facing, pos))
	}

	e.spawner.caps.increment(*kind)
}

func (r *Engine) spawnEntity(e Entity) {
	r.entities = append(r.entities, e)
}

func main() {
	screen, err := tcell.NewScreen()
	check(err)
	check(screen.Init())

	e := Engine{
		tickRate: 2,
		seaLevel: 5,
		paused:   false,
		screen:   screen,
		entities: mkSea(screen.Size()),
		spawner: Spawner{
			caps: EntityCaps{},
		},
	}

	defer e.screen.Fini()
	go eventHandler(&e)

	for {
		if e.paused {
			time.Sleep(time.Duration(e.tickRate) * time.Millisecond)
			continue
		}

		// TODO:
		if RandOneIn(100) {
			e.spawnRandomEntity()
		}

		drawCurrent(&e)
		e.Tick()
	}
}

func drawCurrent(e *Engine) {
	// update all entities
	e.screen.Clear()
	start := time.Now()

	// create empty canvas
	width, height := e.screen.Size()
	e.frame = makeFrame(width, height)

	// draw all entities to canvas
	err := e.Draw()
	check(err)

	if e.debug {
		elapsed := time.Now().Sub(start)
		fps := 1 / elapsed.Seconds()
		ser, err := json.MarshalIndent(e.spawner.caps, "", "  ")
		check(err)
		lines := []string{
			fmt.Sprintf("entities: %d", len(e.entities)),
			fmt.Sprintf("tickRate: %d", e.tickRate),
			fmt.Sprintf("fps: %d", int(fps)),
			fmt.Sprintf("entities: %s", string(ser)),
		}
		for i, line := range lines {
			e.DrawText(line, Pos{Y: i})
		}
	}

	e.screen.Show()
	time.Sleep(time.Duration(10/e.tickRate) * time.Millisecond)
}

func eventHandler(e *Engine) {
	for {
		ev := e.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			e.Reset(mkSea(e.screen.Size()))
		case *tcell.EventKey:
			// FIX: this doesn't work for some reason...
			// switch ev.Key() {
			// case tcell.KeyEscape | tcell.KeyCtrlC:
			// }
			if ev.Key() == tcell.Key(3) {
				e.screen.Fini()
				os.Exit(0)
			}

			switch ev.Rune() {
			case 's':
				e.spawnRandomEntity()
			case 'r':
				e.Reset(mkSea(e.screen.Size()))
			case 'q':
				e.screen.Fini()
				os.Exit(0)
			case ' ':
				e.paused = !e.paused
			case 'x':
				e.debug = !e.debug
				drawCurrent(e)
				e.screen.Show()
			case 'j':
				if e.tickRate <= 1 {
					continue
				}
				if e.tickRate < 10 {
					e.tickRate -= 1
				} else {
					e.tickRate -= 5
				}
			case 'k':
				if e.tickRate < 10 {
					e.tickRate += 1
				} else {
					e.tickRate += 5
				}
			}
		}
	}
}
