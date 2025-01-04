package main

import (
	"math/rand/v2"
	"strings"
	"time"
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

// reverse and padd whitespace 
func reverseArt(s string) string {
	width := findArtWidth(s)
  lines := strings.Split(s, "\n")
  for i, line := range lines {
    lines[i] = reverseString(line)
    if len(lines[i]) < width {
      lines[i] = strings.Repeat(" ", width - len(lines[i])) + lines[i] 
    }
  }
  return strings.Join(lines, "\n")
}

func flipAsciiArt(art string) string {
  // TODO: expand symbol map
	symbolMap := map[rune]rune{
		'<':  '>',
		'>':  '<',
		'(':  ')',
		')':  '(',
		'[':  ']',
		']':  '[',
		'/':  '\\',
		'\\': '/',
	}

  reversed := reverseArt(art)
	lines := strings.Split(reversed, "\n")
	for i, _ := range lines {
		// lines[i] = reverseArt(line)

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
