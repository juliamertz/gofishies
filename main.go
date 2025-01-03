package main

import (
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// how many milliseconds per tick
const cycles = 400

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	screen, err := tcell.NewScreen()
	check(err)
	defer screen.Fini()

	width, height := screen.Size()
	renderer := Renderer{
		tickRate: 5,
		paused:   false,
		screen:   screen,
    // canvas: NewCanvas(width, height),
		entities: []Renderable{
			&Goldfish{Pos: Pos{X: 0, Y: 15}},
			&Whale{Pos: Pos{X: width, Y: 20}},
			&Seaweed{Pos: Pos{X: 10, Y: height - 5}, length: 5},
			&Seaweed{Pos: Pos{X: 13, Y: height - 3}, length: 4},
			&Waves{Pos: Pos{X: 0, Y: 5}},
		}}

  // seperate thread to handle ctrl-c and some other controls
	go inputHandler(&renderer)

  // main render loop
	for {
		if renderer.paused {
			time.Sleep(time.Duration(renderer.tickRate) * time.Millisecond)
			continue
		}

		renderer.Tick()
		time.Sleep(time.Duration(renderer.tickRate) * time.Millisecond)
		renderer.screen.Clear()
		renderer.Draw(screen)
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
		// fmt.Println("I got the byte", b, "("+string(b)+")")
	}

}
