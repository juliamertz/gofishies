package main

import "github.com/gdamore/tcell/v2"

type Entity struct {
	Id           string
	kind         EntityKind
	pos          Pos
	Facing       Direction
	Tick         int
	defaultColor tcell.Color

	currentFrame int
	frames       *[]Frame
	shouldKill   bool

	width  int
	height int

	update func(*Entity, *Engine)
}

type EntityKind int

const (
	SmallFish EntityKind = iota
	LargeFish
	WaterLine
	Vehicle
	Other
)

func createEntity(
	id string,
	kind EntityKind,
	art []string,
	colorMap []string,
	defaultColor tcell.Color,
	pos Pos,
	facing Direction,
	update func(*Entity, *Engine),
) Entity {
	var frames []Frame
	// generate frames for entity
	for i, frame := range art {
		var cMap string
		if i < len(colorMap) {
			cMap = colorMap[i]
		} else {
			// fallback to first entry if frame doesn't have corresponding color map
			cMap = colorMap[0]
		}

		// assume all art is left facing
		if facing == Right {
			frame = mirrorAsciiArt(trimArt(frame))
			cMap = reverseArt(trimArt(cMap))
		}

		frames = append(frames, generateFrame(frame, cMap, defaultColor))
	}

	return Entity{
		Id:           id,
		pos:          pos,
		kind:         kind,
		Facing:       facing,
		defaultColor: defaultColor,
		frames:       &frames,
		update:       update,

		// assume all frames are of same height/width
		height: len(frames[0].cells),
		width:  len(frames[0].cells[0]),
	}
}

func (e *Entity) IsOffscreen(cols int, lines int) bool {
	if e.pos.X >= cols || e.pos.X+e.width <= 0 {
		return true
	}

	if e.pos.Y > lines || e.pos.Y+e.height <= 0 {
		return true
	}

	return false
}

func (e *Entity) Move(dir Direction) {
	switch dir {
	case Left:
		e.pos.X--
	case Right:
		e.pos.X++
	}
}

func (e *Entity) NextFrame() {
	e.currentFrame++
	if e.currentFrame >= len(*e.frames) {
		e.currentFrame = 0
	}
}

// Guess where the bubble should spawn
func (e *Entity) LikelyBubblePos() Pos {
	x := e.pos.X
	if e.Facing == Right {
		x += e.width
	}

	return Pos{X: x, Y: e.pos.Y + (e.height / 2)}
}

func (entity *Entity) GetCollisions(r *Engine) []Entity {
	var buff []Entity
	for _, e := range r.entities {
		if e.Id == entity.Id {
			continue
		}
		b1 := entity.BoundingBox()
		b2 := e.BoundingBox()
		if checkCollision(b1, b2) {
			buff = append(buff, e)
		}
	}
	return buff
}

type BoundingBox struct {
	xMin, yMin int
	xMax, yMax int
}

func checkCollision(a, b BoundingBox) bool {
	return a.xMin < b.xMax && a.xMax > b.xMin && a.yMin < b.yMax && a.yMax > b.yMin
}

// get corner positions of entity relative to it's parent frame
// left-top, right-top, left-bottom, right-bottom
func (e *Entity) BoundingBox() BoundingBox {
	frame := (*e.frames)[e.currentFrame]
	w := frame.width()
	h := frame.height()

	return BoundingBox{
		xMin: e.pos.X,
		xMax: e.pos.X + w,
		yMin: e.pos.Y,
		yMax: e.pos.Y + h,
	}
}
