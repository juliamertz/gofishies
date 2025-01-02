package main

import (
	"fmt"
	"gofishies/ansi"
	"os"
	"slices"
	"strings"

	"golang.org/x/term"
)

type Renderable interface {
	Render() *string
	Tick(*Renderer)
	GetPos() ansi.Pos
	// setPos(Pos)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Renderer struct {
  paused bool

	entities []Renderable
	// Trails or bubbles that fish can leave behind
	fleeting []Renderable
}

func (r *Renderer) Tick() {
	for _, e := range r.entities {
		e.Tick(r)
	}
}

func (r *Renderer) Draw() {
	for _, e := range r.entities {
    rendered := e.Render();
    if rendered == nil {
      continue
    }
		lines := strings.Split(*rendered, "\n")
		ansi.PrintLines(e.GetPos(), lines)
	}
}

func cutVisible(content string, pos ansi.Pos) *string {
	_, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic("unable to get terminal size")
	}

	lines := strings.Split(content, "\n")

	if pos.Y < 0 {
		skip := pos.Y * -1
		lines = slices.Delete(lines, 0, skip)
	}

	lastIdx := len(lines) - 1

	if (pos.Y + lastIdx) > height {
		skip := (pos.Y + lastIdx) - height
		if skip > lastIdx {
			return nil
		}
		start := lastIdx - skip
		lines = slices.Delete(lines, start, lastIdx)
	}

  if pos.X < 0 {
    start := pos.X * -1
    if start > len(lines[0]) -1 {
      return nil
    }
    for idx, line := range lines {
      fmt.Print("idx: %d", idx)
      lines[idx] = line[start:]
    }
  }

	// for i := 0; i < len(lines); i++ {
	//    line := lines[i]
	//    nextPos := pos
	// 	if (pos.Y + i) < 0 || (pos.Y + i) > height-1 {
	// 		continue
	// 	}
	//    if pos.X < 0 {
	//      // continue
	//      if pos.X + len(line) < 0 {
	//        continue
	//      }
	//      // fmt.Printf("line: %v", line)
	//      line = line[pos.X * -1:]
	//      nextPos.X = 0
	//    }
	// 	Print(MoveTo(nextPos.Y+i, nextPos.X))
	// 	Print(line)
	// }

	res := strings.Join(lines, "\n")
	return &res
}
