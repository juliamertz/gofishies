package main

import (
	"gofishies/ansi"
)

// Whale

type Whale struct {
	Pos   ansi.Pos
	cycle int
}

func (f *Whale) DefaultColor() *int {
	return ansi.ColorFromByte('g')
}

func (f *Whale) Render(r *Renderer) (string, string) {
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
 o  //  wgg       ))))      ===/__((
    ))           )))))))        __((
    \\     \)     ))))    __===\ _((
     \\_______________====      \_((
                                 \((`

	return art, colors
}

func (f *Whale) GetPos() ansi.Pos {
	return f.Pos
}

func (f *Whale) Tick(r *Renderer) {
	if f.cycle == 6 {
		f.cycle = 0
		f.Pos.X--
	} else {
		if f.cycle == 3 {
			// r.fleeting = append(r.fleeting, &Bubble{Pos: f.Pos})
		}
		f.cycle++
	}
}

// Bubble

type Bubble struct {
	Pos   ansi.Pos
	cycle int
}

func (f *Bubble) DefaultColor() *int {
	return ansi.ColorFromByte('w')
}

func (f *Bubble) Render(r *Renderer) (string, string) {
	return "o", "w"
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

func (f *Goldfish) DefaultColor() *int {
	return ansi.ColorFromByte('c')
}

func (g *Goldfish) Render(r *Renderer) (string, string) {
	art := []string{
		"  _ ",
		"><_>",
	}

	colors := []string{
		"  c ",
		"yycc",
	}
	return join(art), join(colors)
}

func (g *Goldfish) GetPos() ansi.Pos {
	return g.Pos
}

func (g *Goldfish) Tick(r *Renderer) {
	if g.cycle == 10 {
		g.cycle = 0
		g.Pos.X++
	} else {
		g.cycle++
	}
}

//
// // TESTING: block
//
// type Block struct {
// 	Pos ansi.Pos
// }
//
// func (f *Block) Render() *string {
// 	art := []string{
// 		"rrbnggy",
// 		"rbrnbyb",
// 	}
//
// 	colors := []string{
// 		"rrb ggy",
// 		"rbr byb",
// 	}
//
// 	return cutAndColorize(join(art), join(colors), 'g', f.Pos)
// }
//
// func (f *Block) GetPos() ansi.Pos {
// 	return f.Pos
// }
//
// func (f *Block) Tick(r *Renderer) {
// 	f.Pos.X--
// }
