package main

import (
	"math/rand/v2"
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

// func mkSea(width int, height int) []Renderable {
// 	return []Renderable{
// 		&Goldfish{Pos: Pos{X: 0, Y: 15}},
// 		&Whale{Pos: Pos{X: width - 5, Y: 20}},
// 		&Seaweed{Pos: Pos{X: 10, Y: height - 5}, length: 5},
// 		&Seaweed{Pos: Pos{X: 13, Y: height - 3}, length: 4},
// 		&Waves{Pos: Pos{X: 0, Y: 5}},
// 	}
// }

type Spawnable interface {
	Spawn(r *Renderer)
}

type Spawner struct {
	renderer *Renderer
	pool     []Spawnable
}

func (s *Spawner) spawnRandom(r *Renderer) {
  i := rand.IntN(len(s.pool)-1)
  s.pool[i].Spawn(r)
}

func main() {
	screen, err := tcell.NewScreen()
	check(err)
	err = screen.Init()
	check(err)
	// width, height := screen.Size()

	renderer := Renderer{
		tickRate: 10,
		seaLevel: 5,
		paused:   false,
		screen:   screen,
		// entities: mkSea(width, height),
	}

	defer renderer.screen.Fini()
	go inputHandler(&renderer)

	spawner := Spawner{
		renderer: &renderer,
		pool:     []Spawnable{&Goldfish{}, &Whale{}},
	}

	spawner.spawnRandom(&renderer)
	spawner.spawnRandom(&renderer)
	spawner.spawnRandom(&renderer)

	for {
		if renderer.paused {
			time.Sleep(time.Duration(renderer.tickRate) * time.Millisecond)
			continue
		}
		renderer.Tick()
		renderer.screen.Clear()
		err := renderer.Draw(screen)
		check(err)
		time.Sleep(time.Duration(renderer.tickRate) * time.Millisecond)
	}
}

func join(lines []string) string {
	return strings.Join(lines, "\n")
}

func inputHandler(r *Renderer) {
	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		if len(b) < 1 {
			continue
		}
		r.stdin = append(r.stdin, b[0])
		switch b[0] {
		// ctrl-c
		case 3:
			r.screen.Fini()
			os.Exit(0)
			// r
		// case 114:
		// 	r.entities = mkSea(r.screen.Size())
			// space
		case 32:
			r.paused = !r.paused
			// j
		case 106:
			if r.tickRate < 10 {
				r.tickRate += 1
			} else {
				r.tickRate += 5
			}
			// k
		case 107:
			if r.tickRate <= 2 {
				continue
			}
			if r.tickRate < 10 {
				r.tickRate -= 1
			} else {
				r.tickRate -= 5
			}
		}
	}
}
