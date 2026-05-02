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
	font      *rl.Font
	fontColor rl.Color
	parent    Screen
	camera    *rl.Camera2D
	highlight rl.Color
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
	s.camera.Target = (*s.world.GetPlayer().GetCurrentPosition())

	if rl.IsKeyPressed(rl.KeyBackspace) {
		return s.parent
	}

	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		p := s.world.GetPlayer()

		vec := rl.GetScreenToWorld2D(rl.GetMousePosition(), *s.camera)
		clickedCell := grid.GetCellFromPixelPosition(&vec, s.world.GetTileSize())
		worldMapSize := s.world.GetMap().GetSize()

		if clickedCell.X >= 0 && clickedCell.X <= worldMapSize.X &&
			clickedCell.Y >= 0 && clickedCell.Y <= worldMapSize.Y {
			p.SetTargetPosition(
				grid.GetCenterCellCoordinates(clickedCell, s.world.GetTileSize()),
			)
		}
		utils.LogPlayerTransition(
			rl.LogDebug,
			*p.GetCurrentPosition(),
			*p.GetTargetPosition(),
		)
	}
	return s
}

func (s *gameScreen) DrawInterface() {
	detail := fmt.Sprintf("%dx%d@%d FPS", rl.GetScreenWidth(), rl.GetScreenHeight(), rl.GetFPS())
	worldMap := s.world.GetMap().GetSize()
	tSize := s.world.GetTileSize()

	rl.DrawText(detail, 10, int32(rl.GetScreenHeight())-30, s.font.BaseSize+10, s.fontColor)
	rl.DrawCircle(0, 0, 15, rl.Red)
	rl.DrawCircle(int32(rl.GetScreenWidth()), 0, 15, rl.Blue)
	rl.DrawCircle(0, int32(rl.GetScreenHeight()), 15, rl.Purple)
	rl.DrawCircle(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), 15, rl.Orange)

	vec := s.world.GetMap().GetSize()
	t := fmt.Sprintf("%0.f x %0.f", vec.X, vec.Y)
	rl.DrawText(
		t, utils.GetCenterForText(t, 20, *s.font), 25, 20, s.fontColor,
	)
	rl.BeginMode2D(*s.camera)
	for x := 0; x <= int(worldMap.X); x++ {
		fromX, toX := int32(x*int(tSize)-1), int32(x*int(tSize))
		rl.DrawText(fmt.Sprintf("%d", x), fromX, 0, 10, rl.White)
		rl.DrawLine(fromX, 0, toX, int32(worldMap.Y*tSize), rl.Pink)
	}

	for y := 0; y <= int(worldMap.Y); y++ {
		fromY, toY := int32(y*int(tSize)-1), int32(y*int(tSize))
		rl.DrawText(fmt.Sprintf("%d", y), 0, fromY, 10, rl.White)
		rl.DrawLine(0, fromY, int32(worldMap.X*tSize), toY, rl.Pink)
	}
	rl.EndMode2D()
}

// Function is called once the Drawing canvas has started
// it craws the current state
func (s *gameScreen) Draw() {
	s.DrawInterface()

	rl.BeginMode2D(*s.camera)
	p := s.world.GetPlayer()
	if p.GetCurrentPosition() != nil {
		pPos := p.GetCurrentPosition()
		cx, cy := int32(pPos.X), int32(pPos.Y)
		rl.DrawCircle(cx, cy, 13, rl.DarkGray)
		rl.DrawText(fmt.Sprintf("%d | %d - Dunno", cx, cy), cx, cy+20, 10, rl.SkyBlue)
		worldPos := rl.GetScreenToWorld2D(*pPos, *s.camera)
		rl.DrawText(fmt.Sprintf("%0f | %0f - World", worldPos.X, worldPos.Y), cx, cy+30, 10, rl.SkyBlue)
		screenPos := rl.GetWorldToScreen2D(worldPos, *s.camera)
		rl.DrawText(fmt.Sprintf("%0f | %0f - Back2Screen", screenPos.X, screenPos.Y), cx, cy+40, 10, rl.SkyBlue)
		p.Draw()
	}
	rl.EndMode2D()
}

func NewGameScreen(world *world.World, font *rl.Font, color, highlight rl.Color, cam *rl.Camera2D) Screen {
	return &gameScreen{
		world:     world,
		font:      font,
		fontColor: color,
		camera:    cam,
		highlight: highlight,
	}
}
