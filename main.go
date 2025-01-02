package main

import (
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/term"
)

// how many milliseconds per tick
const cycles = 400

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// cooked, err := term.MakeRaw(int(os.Stdin.Fd()))
	// check(err)
	// defer term.Restore(int(os.Stdin.Fd()), cooked)

	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	check(err)
	renderer := Renderer{
		tickRate: 50,
    paused: false,
		entities: []Renderable{
			&Goldfish{Pos: Pos{X: 0, Y: 15}},
			// &Block{Pos: Pos{X: width - 20, Y: 0}},
			&Whale{Pos: Pos{X: width, Y: 20}},
			&Seaweed{Pos: Pos{X: 10, Y: height - 6}, length: 3},
			&Waves{Pos: Pos{X: 0, Y: 5}},
		}}


	screen, err := tcell.NewScreen()
  renderer.screen = screen
	check(err)
	if err = screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()

	renderer.InitCells()


	go inputHandler(&renderer)

  renderer.screen.Clear()
	renderer.Draw(screen)

	i := 0

	for i < cycles {
		if renderer.paused {
			time.Sleep(time.Duration(renderer.tickRate) * time.Millisecond)
			continue
		}
		renderer.InitCells()
		renderer.Tick()
		time.Sleep(time.Duration(renderer.tickRate) * time.Millisecond)
    renderer.screen.Clear()
		renderer.Draw(screen)
	}

	time.Sleep(5 * time.Second)
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
