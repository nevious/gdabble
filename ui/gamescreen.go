package ui

/*
 * NOTE
 * I think we do _a_lot_ of game logic in the UI package, that should be handled in
 * the world. UI should just receive stuff to render, not neccessarily device where
 * to render what
 */

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nevious/gdabble/grid"
	"github.com/nevious/gdabble/utils"
	"github.com/nevious/gdabble/world"
)

type GameScreen struct {
	world     *world.World
	font      *rl.Font
	fontColor rl.Color
	parent    Screen
	camera    *rl.Camera2D
	highlight rl.Color
}

// Update tracking of the camera and zoomdd
// Positions are casted to int and then to float32 again. This enforces
// whole number values
func (s *GameScreen) updateCamera() {
	player := s.world.GetPlayer()
	playerPosition := s.world.GetEntityPosition(player)
	s.camera.Target.X = float32(int(playerPosition.X))
	s.camera.Target.Y = float32(int(playerPosition.Y))
	s.camera.Offset.X = float32(int(rl.GetScreenWidth()) / 2)
	s.camera.Offset.Y = float32(int(rl.GetScreenHeight()) / 2)

	// Camera Zoom
	if wheelMove := rl.GetMouseWheelMove(); wheelMove > 0 {
		s.camera.Zoom++
	} else if wheelMove < 0 && s.camera.Zoom >= 2 {
		s.camera.Zoom--
	}
}

// Set the players position via world to the clicked cell
func (s *GameScreen) handlePlayerInput() {
	player := s.world.GetPlayer()
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {

		vec := rl.GetScreenToWorld2D(rl.GetMousePosition(), *s.camera)
		clickedCell := grid.GetCellFromPixelPosition(&vec, s.world.GetTileSize())

		if grid.CellWithinMapBounds(clickedCell, s.world.GetMap().GetSize()) {
			s.world.SetEntityPosition(player, grid.GetCenterCellCoordinates(clickedCell, s.world.GetTileSize()))
		}
	}
}

// Render the overall interface
// NOTE: This function lives in screen space, not world or cell space
func (s *GameScreen) renderOverlayInterface() {
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

// Render the Cell Border under the mouse
func (s *GameScreen) renderNavigationCross() {
	tSize := s.world.GetTileSize()
	rl.BeginMode2D(*s.camera)
	_, gridHighlightColor := s.highlight, s.highlight
	gridHighlightColor.A = 255
	mousePosition := rl.GetScreenToWorld2D(rl.GetMousePosition(), *s.camera)
	mouseCell := *grid.GetCellFromPixelPosition(&mousePosition, s.world.GetTileSize())

	if grid.CellWithinMapBounds(&mouseCell, s.world.GetMap().GetSize()) {
		rect := rl.NewRectangle((mouseCell.X*tSize)-1, (mouseCell.Y*tSize)-1, tSize, tSize)
		rl.DrawRectangleLinesEx(rect, 1, gridHighlightColor)
	}
	rl.EndMode2D()
}

// Render the Map on the screeen
func (s *GameScreen) renderMap() {
	rotation := rl.NewVector2(0, 0)

	rl.BeginMode2D(*s.camera)
	for _, tile := range s.world.GetRenderListItems() {
		rl.DrawTexturePro(tile.Texture, tile.Src, tile.Dst, rotation, 0, rl.White)
	}
	rl.EndMode2D()
}

// Render UI elements of the character
func (s *GameScreen) renderCharacterElements() {
	rl.BeginMode2D(*s.camera)
	player := s.world.GetPlayer()
	playerPosition := *s.world.GetEntityPosition(player)

	cx, cy := int32(playerPosition.X), int32(playerPosition.Y)

	// some debugging information about player position
	rl.DrawText(fmt.Sprintf("Row: %d | %d", cx, cy), cx, cy+20, 10, rl.RayWhite)
	worldPos := rl.GetScreenToWorld2D(playerPosition, *s.camera)
	rl.DrawText(fmt.Sprintf("World: %0f | %0f", worldPos.X, worldPos.Y), cx, cy+30, 10, rl.RayWhite)
	screenPos := rl.GetWorldToScreen2D(worldPos, *s.camera)
	rl.DrawText(fmt.Sprintf("Screen: %0f | %0f", screenPos.X, screenPos.Y), cx, cy+40, 10, rl.RayWhite)
	mousePosition := rl.GetMousePosition()
	rl.DrawText(fmt.Sprintf("Mouse: %0f | %0f", mousePosition.X, mousePosition.Y), cx, cy+50, 10, rl.RayWhite)
	playerState := s.world.GetEntityActionState(player)
	rl.DrawText(fmt.Sprintf("State: %s", playerState), cx, cy+60, 10, rl.RayWhite)

	rl.EndMode2D()
}

// Function is called once per frame before drawing
// is enabled
func (s *GameScreen) Update() Screen {
	if rl.IsKeyPressed(rl.KeyBackspace) {
		return s.parent
	}

	s.world.Update() // Update world state and game logic
	s.handlePlayerInput()
	s.updateCamera()

	return s
}

// Set parent screen of this screen
func (s *GameScreen) SetParent(parent Screen) Screen {
	s.parent = parent
	return s
}

// Orchestration function with no parameters
// call to child-functions rendering specific things
func (s *GameScreen) Render() {
	s.renderMap()
	s.renderCharacterElements()
	s.renderNavigationCross()
	s.renderOverlayInterface()
}

// Create a new UI Screen handling the game world
func NewGameScreen(gameWorld *world.World, font *rl.Font, color, highlight rl.Color, cam *rl.Camera2D) Screen {
	return &GameScreen{
		world:     gameWorld,
		font:      font,
		fontColor: color,
		camera:    cam,
		highlight: highlight,
	}
}
