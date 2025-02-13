package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

var RNG = rand.New(rand.NewPCG(uint64(0), uint64(time.Now().UnixNano())))

func SeededRandIntN(n int, seed uint64) int {
	src := rand.NewPCG(uint64(seed), uint64(seed))
	rng := rand.New(src)
	return rng.IntN(n)
}

func RandOneIn(n int) bool {
	src := rand.NewPCG(uint64(n), uint64(time.Now().UnixNano()))
	rng := rand.New(src)
	j := rng.IntN(n)
	return j == 1
}

func join(lines []string) string {
	return strings.Join(lines, "\n")
}

func findArtWidth(art string) int {
	lines := strings.Split(art, "\n")
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	return width
}

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// remove leading/trailing newlines from art string
func trimArt(val string) string {
	var buff []string
	for _, line := range strings.Split(val, "\n") {
		if strings.TrimSpace(line) != "" {
			buff = append(buff, line)
		}
	}
	return join(buff)
}

// reverse and padd whitespace
func reverseArt(s string) string {
	width := findArtWidth(s)
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = reverseString(line)
		if len(lines[i]) < width {
			lines[i] = strings.Repeat(" ", width-len(lines[i])) + lines[i]
		}
	}
	return strings.Join(lines, "\n")
}

// reverse string and replace mapped ascii symbols
func mirrorAsciiArt(art string) string {
	// TODO: Preserve text
	symbolMap := map[rune]rune{
		'<':  '>',
		'>':  '<',
		'(':  ')',
		')':  '(',
		'[':  ']',
		']':  '[',
		'{':  '}',
		'}':  '{',
		'/':  '\\',
		'\\': '/',
	}

	reversed := reverseArt(art)
	lines := strings.Split(reversed, "\n")
	for i := range lines {
		for j, ch := range lines[i] {
			if value, exists := symbolMap[ch]; exists {
				runes := []rune(lines[i])
				runes[j] = value
				lines[i] = string(runes)
			}
		}
	}
	return strings.Join(lines, "\n")
}

func assertEq(a any, b any, msg string) {
	if a != b {
		fmt.Println("Assertion failed! ", msg)
		fmt.Printf("%v == %v\n", a, b)
		os.Exit(1)
	}
}

func ColorFromRune(r rune) tcell.Color {
	switch r {
	case 'r':
		return tcell.ColorRed
	case 'g':
		return tcell.ColorGreen
	case 'y':
		return tcell.ColorYellow
	case 'b':
		return tcell.ColorBlue
	case 'p':
		return tcell.ColorPurple
	case 'c':
		return tcell.ColorLightCyan
	case 'w':
		return tcell.ColorWhite
	case 'd':
		return tcell.ColorDefault
	}
	return tcell.ColorNone
}

// Get random fill color
func RandColor() tcell.Color {
	switch RNG.IntN(14) {
	case 0:
		return tcell.ColorMaroon
	case 1:
		return tcell.ColorGreen
	case 2:
		return tcell.ColorOlive
	case 3:
		return tcell.ColorNavy
	case 4:
		return tcell.ColorPurple
	case 5:
		return tcell.ColorTeal
	case 6:
		return tcell.ColorSilver
	case 7:
		return tcell.ColorRed
	case 8:
		return tcell.ColorLime
	case 9:
		return tcell.ColorYellow
	case 10:
		return tcell.ColorBlue
	case 11:
		return tcell.ColorFuchsia
	case 12:
		return tcell.ColorAqua
	case 13:
		return tcell.ColorWhite
	}

	panic("unreachable")
}
