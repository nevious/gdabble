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

// Update camera Target and Offset
func (s *gameScreen) handleCameraInput() {
	// Player position
	p := s.world.GetPlayer()
	s.camera.Target.X = float32(int(p.GetCurrentPosition().X))
	s.camera.Target.Y = float32(int(p.GetCurrentPosition().Y))
	// Casting to int before, so the value is clamped to a near number
	s.camera.Offset.X = float32(int(rl.GetScreenWidth()) / 2)
	s.camera.Offset.Y = float32(int(rl.GetScreenHeight()) / 2)

	// Camera Zoom
	if wheelMove := rl.GetMouseWheelMove(); wheelMove > 0 {
		s.camera.Zoom++
	} else if wheelMove < 0 && s.camera.Zoom >= 2 {
		s.camera.Zoom--
	}
}

func (s *gameScreen) handlePlayerInput() {
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		p := s.world.GetPlayer()

		vec := rl.GetScreenToWorld2D(rl.GetMousePosition(), *s.camera)
		clickedCell := grid.GetCellFromPixelPosition(&vec, s.world.GetTileSize())
		worldMapSize := s.world.GetMap().GetSize()

		// TODO
		// Whether a cell is inside the world grid is something
		// that should behandled in the grid package
		if clickedCell.X >= 0 && clickedCell.X < worldMapSize.X &&
			clickedCell.Y >= 0 && clickedCell.Y < worldMapSize.Y {
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
}

func (s *gameScreen) drawOverlayInterface() {
	detail := fmt.Sprintf("%dx%d@%d FPS", rl.GetScreenWidth(), rl.GetScreenHeight(), rl.GetFPS())
	worldMap := s.world.GetMap().GetSize()

	rl.DrawText(detail, 10, int32(rl.GetScreenHeight())-30, s.font.BaseSize+10, s.fontColor)
	rl.DrawCircle(0, 0, 15, rl.Red)
	rl.DrawCircle(int32(rl.GetScreenWidth()), 0, 15, rl.Blue)
	rl.DrawCircle(0, int32(rl.GetScreenHeight()), 15, rl.Purple)
	rl.DrawCircle(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), 15, rl.Orange)

	t := fmt.Sprintf("%0.f x %0.f", worldMap.X, worldMap.Y)
	rl.DrawText(
		t, utils.GetCenterForText(t, 20, *s.font), 25, 20, s.fontColor,
	)

}

func (s *gameScreen) drawGridLines() {
	tSize := s.world.GetTileSize()
	rl.BeginMode2D(*s.camera)
	gridColor, gridHighlightColor := s.highlight, s.highlight
	gridColor.A = 100
	gridHighlightColor.A = 255
	mousePosition := rl.GetScreenToWorld2D(rl.GetMousePosition(), *s.camera)
	mouseCell := *grid.GetCellFromPixelPosition(&mousePosition, s.world.GetTileSize())
	rect := rl.NewRectangle((mouseCell.X*tSize)-1, (mouseCell.Y*tSize)-1, tSize, tSize)
	rl.DrawRectangleLinesEx(rect, 1, gridHighlightColor)
	rl.EndMode2D()
}

func (s *gameScreen) drawTileMap() {
	// Dummy function, draw a single tile across everything
	// not what we're gonna have
	worldMap := s.world.GetMap().GetSize()
	tSize := s.world.GetTileSize()
	texture := *s.world.GetMap().GetTexture()

	rl.BeginMode2D(*s.camera)

	for x := float32(0); x < worldMap.X; x++ {
		for y := float32(0); y < worldMap.Y; y++ {
			destPosition := rl.NewVector2(x*tSize, y*tSize)
			destRect := rl.NewRectangle(destPosition.X, destPosition.Y, tSize, tSize)
			textureRect := s.world.GetMap().GetTileAt(int(x), int(y))
			rl.DrawTexturePro(texture, *textureRect, destRect, rl.NewVector2(0, 0), 0, rl.White)
		}
	}

	rl.EndMode2D()
}

func (s *gameScreen) drawPlayer() {
	rl.BeginMode2D(*s.camera)
	p := s.world.GetPlayer()
	if p.GetCurrentPosition() != nil {
		pPos := p.GetCurrentPosition()
		cx, cy := int32(pPos.X), int32(pPos.Y)
		rl.DrawCircle(cx, cy, 13, rl.DarkGray) // Acting as a shaodw for now

		// some debugging information about player position
		rl.DrawText(fmt.Sprintf("%d | %d - Raw", cx, cy), cx, cy+20, 10, rl.SkyBlue)
		worldPos := rl.GetScreenToWorld2D(*pPos, *s.camera)
		rl.DrawText(fmt.Sprintf("%0f | %0f - World", worldPos.X, worldPos.Y), cx, cy+30, 10, rl.SkyBlue)
		screenPos := rl.GetWorldToScreen2D(worldPos, *s.camera)
		rl.DrawText(fmt.Sprintf("%0f | %0f - Back2Screen", screenPos.X, screenPos.Y), cx, cy+40, 10, rl.SkyBlue)

		// p.Draw()
		playerSprite := p.GetCharacterSprite()
		rl.DrawTexture(
			*playerSprite,
			int32(pPos.X)-playerSprite.Width/2,
			int32(pPos.Y)-playerSprite.Height,
			rl.White,
		)
	}
	rl.EndMode2D()
}

/*
 * Interface Functions
 */

func (s *gameScreen) SetParent(parent Screen) Screen {
	s.parent = parent
	return s
}

// Function is called once per frame before drawing
func (s *gameScreen) HandleInput() Screen {
	if rl.IsKeyPressed(rl.KeyBackspace) {
		return s.parent
	}

	s.world.Update() // Update world state and game logic
	s.handlePlayerInput()
	s.handleCameraInput()

	return s
}

// Function is called once the Drawing canvas has started
// it draws the current state
func (s *gameScreen) Draw() {
	// Dispatch drawing, each function should call their own BeginMode2D
	// if needed and provide their own state data like dimensions
	s.drawTileMap()
	s.drawGridLines()
	s.drawPlayer()
	s.drawOverlayInterface()
}

func NewGameScreen(gameWorld *world.World, font *rl.Font, color, highlight rl.Color, cam *rl.Camera2D) Screen {
	return &gameScreen{
		world:     gameWorld,
		font:      font,
		fontColor: color,
		camera:    cam,
		highlight: highlight,
	}
}
