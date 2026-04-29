package ui

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/grid"
	"main/utils"
	"main/world"
)

type gameScreen struct {
	world     *world.World
	cellWidth int
	font      *rl.Font
	fontColor rl.Color
	parent    Screen
}

func (s *gameScreen) SetParent(parent Screen) Screen {
	s.parent = parent
	return s
}

// Function is called once per frame at the beginning
// of the game loop before raylib is ready to be drawed() on
// it manipulates state, not the ui
func (s *gameScreen) Update() Screen {
	s.world.Update() // Update world state, do not draw anything

	if rl.IsKeyPressed(rl.KeyBackspace) {
		return s.parent
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		p := s.world.GetPlayer()
		vec := rl.GetMousePosition()
		clickedCell := grid.GetCellFromPixelPosition(&vec, s.cellWidth)
		p.SetTargetPosition(
			grid.GetCenterCellCoordinates(clickedCell, s.cellWidth),
		)
		utils.LogPlayerTransition(
			rl.LogDebug,
			*p.GetCurrentPosition(),
			*p.GetTargetPosition(),
		)
	}
	return s
}

// Function is called once the Drawing canvas has started
// it craws the current state
func (s *gameScreen) Draw() {
	// Some debugging data
	detail := fmt.Sprintf(
		"%dx%d@%d FPS",
		rl.GetScreenWidth(),
		rl.GetScreenHeight(),
		rl.GetFPS(),
	)
	rl.DrawText(
		detail,
		10,
		int32(rl.GetScreenHeight())-30,
		s.font.BaseSize+10,
		s.fontColor,
	)

	rl.DrawCircle(0, 0, 15, rl.Red)
	rl.DrawCircle(
		int32(rl.GetScreenWidth()), 0, 15, rl.Blue,
	)
	rl.DrawCircle(
		0, int32(rl.GetScreenHeight()), 15, rl.Purple,
	)
	rl.DrawCircle(
		int32(rl.GetScreenWidth()),
		int32(rl.GetScreenHeight()),
		15,
		rl.Orange,
	)

	vec := *grid.GetCoordinateSystem(s.cellWidth)
	t := fmt.Sprintf("%0.f x %0.f", vec.X, vec.Y)
	rl.DrawText(
		t, utils.GetCenterForText(t, 20, *s.font), 150, 20, s.fontColor,
	)

	for x := 0; x < int(vec.X); x++ {
		from_x, to_x := int32(x*s.cellWidth-1), int32(x*s.cellWidth)
		rl.DrawLine(from_x, 0, to_x, int32(rl.GetScreenHeight()), rl.Pink)
	}

	for y := 0; y < int(vec.Y); y++ {
		from_y, to_y := int32(y*s.cellWidth-1), int32(y*s.cellWidth)
		rl.DrawLine(0, from_y, int32(rl.GetScreenWidth()), to_y, rl.Pink)
	}

	p := s.world.GetPlayer()
	if p.GetCurrentPosition() != nil {
		pPos := p.GetCurrentPosition()
		cx, cy := int32(pPos.X), int32(pPos.Y)
		rl.DrawCircle(cx, cy, 15, rl.White)
	}

}

func NewGameScreen(world *world.World, cellWidth int, font *rl.Font, color, highlight rl.Color) Screen {
	// TODO: highlight ommited!
	return &gameScreen{
		world:     world,
		cellWidth: cellWidth,
		font:      font,
		fontColor: color,
	}
}
