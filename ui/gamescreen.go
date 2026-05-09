package ui

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/character"
	"main/grid"
	"main/utils"
	"main/world"
)

type GameScreen struct {
	world     *world.World
	worldMap  world.GameMapInterface
	font      *rl.Font
	fontColor rl.Color
	parent    Screen
	camera    *rl.Camera2D
	highlight rl.Color
	player    character.Entity
}

// TODO: Bruh... wth is a camerainput....
// Update camera Target and Offset
func (s *GameScreen) handleCameraInput() {
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

func (s *GameScreen) handlePlayerInput() {
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

func (s *GameScreen) renderNavigationCross() {
	tSize := s.world.GetTileSize()
	rl.BeginMode2D(*s.camera)
	_, gridHighlightColor := s.highlight, s.highlight
	gridHighlightColor.A = 255
	mousePosition := rl.GetScreenToWorld2D(rl.GetMousePosition(), *s.camera)
	mouseCell := *grid.GetCellFromPixelPosition(&mousePosition, s.world.GetTileSize())
	rect := rl.NewRectangle((mouseCell.X*tSize)-1, (mouseCell.Y*tSize)-1, tSize, tSize)
	rl.DrawRectangleLinesEx(rect, 1, gridHighlightColor)
	rl.EndMode2D()
}

func (s *GameScreen) loadPlayerIntoMap() {
	s.player = s.world.GetPlayer()
	playerMapId := s.player.GetCurrentMap()

	if s.worldMap = s.world.GetMap(); s.worldMap == nil || s.worldMap.GetId() != playerMapId {
		s.world.SetPlayerMapPosition(s.player.GetCurrentMap())
		s.worldMap = s.world.GetMap()
	}
}

// yStart, yStop in map coordinates
func (s *GameScreen) renderMap() {
	tSize := s.world.GetTileSize()
	texture := *s.worldMap.GetTexture()
	mapSize := s.worldMap.GetSize()

	rl.BeginMode2D(*s.camera)

	for x := float32(0); x < mapSize.X; x++ {
		for y := float32(0); y < mapSize.Y; y++ {
			destPosition := rl.NewVector2(x*tSize, y*tSize)
			textureRectangles := s.world.GetMap().GetTileAt(int(x), int(y))
			for _, rect := range textureRectangles {
				X, Y := destPosition.X, destPosition.Y
				if rect.Scale > 1 {
					X, Y = destPosition.X-((tSize*rect.Scale)/2), destPosition.Y-(tSize*rect.Scale-tSize)
				}
				destRect := rl.NewRectangle(X, Y, tSize*rect.Scale, tSize*rect.Scale)
				rl.DrawTexturePro(texture, rect.Rect, destRect, rl.NewVector2(0, 0), 0, rl.White)
			}
		}
	}

	rl.EndMode2D()
}

func (s *GameScreen) renderPlayer() {
	rl.BeginMode2D(*s.camera)
	p := s.world.GetPlayer()
	pPos := p.GetCurrentPosition()
	tSize := s.world.GetTileSize()

	if p.GetCurrentPosition() != nil {
		cx, cy := int32(pPos.X), int32(pPos.Y)

		// some debugging information about player position
		rl.DrawText(fmt.Sprintf("%d | %d - Raw", cx, cy), cx, cy+20, 10, rl.SkyBlue)
		worldPos := rl.GetScreenToWorld2D(*pPos, *s.camera)
		rl.DrawText(fmt.Sprintf("%0f | %0f - World", worldPos.X, worldPos.Y), cx, cy+30, 10, rl.SkyBlue)
		screenPos := rl.GetWorldToScreen2D(worldPos, *s.camera)
		rl.DrawText(fmt.Sprintf("%0f | %0f - Back2Screen", screenPos.X, screenPos.Y), cx, cy+40, 10, rl.SkyBlue)

		playerTexture, playerTextureLocation := p.GetCharacterSprite()
		playerRectangle := rl.NewRectangle(
			pPos.X-float32(tSize/2),
			pPos.Y-float32(tSize/2),
			tSize,
			tSize,
		)
		rl.DrawTexturePro(
			*playerTexture,
			*playerTextureLocation,
			playerRectangle,
			rl.NewVector2(0, 0),
			0,
			rl.White,
		)
	}
	rl.EndMode2D()
}

// Set parent screen to this screen
func (s *GameScreen) SetParent(parent Screen) Screen {
	s.parent = parent
	return s
}

// Function is called once per frame before drawing
func (s *GameScreen) HandleInput() Screen {
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
func (s *GameScreen) Draw() {
	if !s.world.MapsAreLoaded() {
		s.world.LoadMaps()
	}
	s.loadPlayerIntoMap()

	s.renderMap()
	s.renderPlayer()
	s.renderNavigationCross()
	s.renderOverlayInterface()
}

func NewGameScreen(gameWorld *world.World, font *rl.Font, color, highlight rl.Color, cam *rl.Camera2D) Screen {
	return &GameScreen{
		world:     gameWorld,
		font:      font,
		fontColor: color,
		camera:    cam,
		highlight: highlight,
	}
}
