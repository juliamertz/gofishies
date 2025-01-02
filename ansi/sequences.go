package ansi

import (
	"fmt"
	"io"
)

type Pos struct {
	X int
	Y int
}

const Esc = "\x1b"
const Clear = Esc + "[2J"
const MoveHome = Esc + "[H"

func MoveTo(line int, column int) string {
	return Esc + fmt.Sprintf("[%d;%dH", line+1, column+1) + Esc + fmt.Sprintf("[%d;%df", line+1, column+1)
}

func Print(seq string) {
	fmt.Print(seq)
}

func PrintLines(pos Pos, lines []string) {
	for i := 0; i < len(lines); i++ {
		Print(MoveTo(pos.Y+i, pos.X))
		Print(lines[i])
	}
}

type AnsiColor interface {
	code(c *AnsiColor) int
}

func Color(code int) string {
	return Esc + fmt.Sprintf("[1;%dm", code)
}

type ForegroundColor int
const (
	Black  ForegroundColor = 30
	Red    ForegroundColor = 31
	Green  ForegroundColor = 32
	Yellow ForegroundColor = 33
	White  ForegroundColor = 37
)

func ColorFromByte(b byte) *int {
	var result int
	switch b {
	case 'r':
		result = 31
	case 'g':
		result = 32
	case 'y':
		result = 33
	case 'b':
		result = 34
	case 'm':
		result = 35
	case 'c':
		result = 36
	case 'w':
		result = 37
	case 'd':
		result = 39
	}
	return &result
}

type Cell struct {
	Content byte
	Fg     ForegroundColor 
	Bg      int
}

// Print cell at current cursor
func (c *Cell) Print(w io.Writer) {
  if c.Content == 0 {
    w.Write([]byte{' '})
    return
  }

	if c.Bg != 0 {
    w.Write([]byte(Color(c.Bg)))
	}

    w.Write([]byte(Color(int(c.Fg)) + string(c.Content)))
}

// func ColorizeArt(art string, colors string, base int) [][]Cell {
//   lines := strings.Split(art, "\n")
//   buff := make([][]Cell, )
// 	for i := range r.cells {
// 		buff[i] = make([]ansi.Cell, width, width)
// 	}
// }

// func Black(str string) string {
//   return color(30) + str + color(40)
// }

// Red 	31 	41
// Green 	32 	42
// Yellow 	33 	43
// Blue 	34 	44
// Magenta 	35 	45
// Cyan 	36 	46
// White 	37 	47
// Default 	39 	49
// Reset 	0 	0
