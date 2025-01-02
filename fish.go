package main

import (
	"gofishies/ansi"
)


// returns nil if the rendered string is completely out of view
func cutAndColorize(art string, colors string, base byte, pos ansi.Pos) *string {
	cutArt := cutVisible(art, pos)
	cutColors := cutVisible(colors, pos)
	if cutArt == nil || cutColors == nil {
		return nil
	}
	colored := ansi.ColorizeArt(*cutArt, *cutColors, *ansi.ColorFromByte(base))
	return &colored
}

// Whale

type Whale struct {
	Pos   ansi.Pos
	cycle int
}

func (f *Whale) Render() *string {
	art :=
		` o
o      ______/~/~/~/__           /((
  o  // __            ====__    /_((
 o  //  @))       ))))      ===/__((
    ))           )))))))        __((
    \\     \)     ))))    __===\ _((
     \\_______________====      \_((
                                 \((`

	colors :=
		` o
o      ______/~/~/~/__           /((
  o  // __            ====__    /_((
 o  //  @))       ))))      ===/__((
    ))           )))))))        __((
    \\     \)     ))))    __===\ _((
     \\_______________====      \_((
                                 \((`

	return cutAndColorize(art, colors, 'g', f.Pos)
}

func (f *Whale) GetPos() ansi.Pos {
	return f.Pos
}

func (f *Whale) Tick(r *Renderer) {
	if f.cycle == 6 {
		f.cycle = 0
		f.Pos.X--
	} else {
		f.cycle++
	}
}

// Bubble

type Bubble struct {
	Pos   ansi.Pos
	cycle int
}

func (f *Bubble) Render() string {
	return ansi.ColorizeArt("o", " ", *ansi.ColorFromByte('b'))
}

func (f *Bubble) GetPos() ansi.Pos {
	return f.Pos
}

func (f *Bubble) Tick(r *Renderer) {
	if f.cycle == 6 {
		f.cycle = 0
		f.Pos.Y++
	} else {
		f.cycle++
	}
}

// Goldfish

type Goldfish struct {
	Pos   ansi.Pos
	cycle int
}

func (f *Goldfish) Render() *string {
	art := []string{
		"  _ ",
		"><_>",
	}

	colors := []string{
		"    ",
		"    ",
	}

	return cutAndColorize(join(art), join(colors), 'c', f.Pos)
}

func (f *Goldfish) GetPos() ansi.Pos {
	return f.Pos
}

func (f *Goldfish) Tick(r *Renderer) {
	if f.cycle == 10 {
		f.cycle = 0
		f.Pos.X++
	} else {
		f.cycle++
	}
}

// TESTING: block

type Block struct {
	Pos ansi.Pos
}

func (f *Block) Render() *string {
	art := []string{
		"rrbnggy",
		"rbrnbyb",
	}

	colors := []string{
		"rrb ggy",
		"rbr byb",
	}

	return cutAndColorize(join(art), join(colors), 'g', f.Pos)
}

func (f *Block) GetPos() ansi.Pos {
	return f.Pos
}

func (f *Block) Tick(r *Renderer) {
	f.Pos.X--
}
