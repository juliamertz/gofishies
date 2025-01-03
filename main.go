package main

import (
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

func main() {
	screen, err := tcell.NewScreen()
	check(err)
  err = screen.Init()
  check(err)
	width, height := screen.Size()

	renderer := Renderer{
		tickRate: 5,
		paused:   false,
    screen: screen,
		entities: []Renderable{
			&Goldfish{Pos: Pos{X: 0, Y: 15}},
			&Whale{Pos: Pos{X: width, Y: 20}},
			&Seaweed{Pos: Pos{X: 10, Y: height - 5}, length: 5},
			&Seaweed{Pos: Pos{X: 13, Y: height - 3}, length: 4},
			&Waves{Pos: Pos{X: 0, Y: 5}},
		}}

	defer renderer.screen.Fini()
	go inputHandler(&renderer)

	for  {
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
		switch b[0] {
		// ctrl-c
		case 3:
			r.screen.Fini()
			os.Exit(0)
			// space
		case 32:
			r.paused = !r.paused
			// j
		case 106:
			r.tickRate += 5
			// k
		case 107:
			if r.tickRate <= 5 {
				continue
			}
			r.tickRate -= 5
		}
	}
}
