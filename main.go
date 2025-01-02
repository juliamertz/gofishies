package main

import (
	"fmt"
	"gofishies/ansi"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

// how many milliseconds per tick
const cycles = 400

func main() {
	cooked, err := term.MakeRaw(int(os.Stdin.Fd()))
	check(err)
	defer term.Restore(int(os.Stdin.Fd()), cooked)

	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	check(err)
	renderer := Renderer{entities: []Renderable{
		&Goldfish{Pos: ansi.Pos{X: 0, Y: 0}},
		// &Block{Pos: ansi.Pos{X: width - 20, Y: 0}},
		&Whale{Pos: ansi.Pos{X: width, Y: 20}},
    &Seaweed{Pos:ansi.Pos{X: 10, Y: height-6}, length: 3},
	}}

	tickRate := 50

	fmt.Println(ansi.Clear)
	renderer.Draw()

	i := 0

	go func() {
		var b []byte = make([]byte, 1)
		for {
			os.Stdin.Read(b)

			if len(b) < 1 {
				continue
			}
			switch b[0] {
			// ctrl-c
			case 3:
				os.Exit(0)
				// space
			case 32:
				renderer.paused = !renderer.paused
				// j
			case 106:
				tickRate += 5
				// k
			case 107:
        if tickRate <= 5 {
          continue
        }
				tickRate -= 5
			}
			// fmt.Println("I got the byte", b, "("+string(b)+")")
		}
	}()

	for i < cycles {
		if renderer.paused {
			time.Sleep(time.Duration(tickRate) * time.Millisecond)
			continue
		}
		renderer.Tick()
		time.Sleep(time.Duration(tickRate) * time.Millisecond)
		fmt.Println(ansi.Clear)
		renderer.Draw()
	}

	time.Sleep(5 * time.Second)
}

func join(lines []string) string {
	return strings.Join(lines, "\n")
}
