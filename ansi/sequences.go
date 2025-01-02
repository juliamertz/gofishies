package ansi

import (
	"bytes"
	"fmt"
	"strings"
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

func color(code int) string {
	return Esc + fmt.Sprintf("[1;%dm", code)
}

var Black = color(30)
var Red = color(31)
var Green = color(32)
var Yellow = color(33)
var White = color(37)
var Reset = color(0)

func ColorFromByte(b byte) *int {
  var result int;
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

func ColorizeArt(art string, colors string, base int) string {
	var lastColor int
	lines := strings.Split(art, "\n")
	colorLines := strings.Split(colors, "\n")
	buff := bytes.NewBufferString(color(base))

	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[y]); x++ {
			var newColor int = base

			if x < len(colorLines[y]) {
        code := ColorFromByte(colorLines[y][x])
        if code == nil {
          newColor = base
        } else {
          newColor = *code
        }
			}

			if newColor != lastColor {
				buff.WriteString(color(newColor))
			}
			lastColor = newColor
			buff.WriteByte(lines[y][x])
		}
		buff.WriteByte('\n')
	}

	// buff.WriteString(Reset)
	return buff.String()
}

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
