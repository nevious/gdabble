package world

import (
	"cmp"
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/character"
	"main/grid"
	"main/utils"
	"slices"
)

/* -------------------------------------------------------------------------------------- */

type World struct {
	tileSize      float32
	player        character.Entity
	availableMaps *[]GameMapInterface
	currentMap    GameMapInterface
	mapDataDir    string
	renderList    []RenderListItem
	// A cell we know the player is allowed to be on, used to fall back on
	safeCell *rl.Vector2
}

type RenderListItem struct {
	Src     rl.Rectangle
	Dst     rl.Rectangle
	Texture rl.Texture2D
	zindex  int // We manage zindex in world, nowhere else
}

/* Public --------------------------------------------------------------------------------------- */

// Update world state
func (w *World) Update() {
	if !w.mapsAreLoaded() {
		w.loadMaps()
	}

	w.updateCurrentMap()
	w.updatePlayer()
	w.updateEnemies()
	w.updateRenderList()
}

// Return current player position
func (w *World) GetPlayerPosition() *rl.Vector2 {
	return w.player.GetCurrentPosition()
}

// Set the player position to the given `position` in
// Cell-Space
func (w *World) SetPlayerPosition(position *rl.Vector2) {
	w.player.SetTargetPosition(position)
	utils.LogPlayerTransition(rl.LogDebug, *w.player.GetCurrentPosition(), *w.player.GetTargetPosition())
}

// Get world Tilesize set by config
func (w *World) GetTileSize() float32 {
	return w.tileSize
}

// Get the currently selected map
func (w *World) GetMap() GameMapInterface {
	return w.currentMap
}

// Retrieve tiles in render order
func (w *World) GetRenderListItems() []RenderListItem {
	return w.renderList
}

/* Private -------------------------------------------------------------------------------------- */

func (w *World) updateRenderList() {
	// BUG:
	// Some trees still look weird when the player is on the same cell.
	mapSize := w.currentMap.GetSize()
	result := []RenderListItem{}

	// Calculate the tiles
	for x := float32(0); x < mapSize.X; x++ {
		for y := float32(0); y < mapSize.Y; y++ {
			destPosition := rl.NewVector2(x*w.tileSize, y*w.tileSize)
			textureRectangles := w.currentMap.GetTileAt(int(x), int(y))

			for _, renderItem := range textureRectangles {
				dstX, dstY := destPosition.X, destPosition.Y
				if renderItem.Scale > 1 {
					// multicell textures are rendered with an offset
					// that is equal to half the texture width. So for a 3-tile
					// texture, it renders a 1.5 cell offset, which is why trees are placed between cells.
					dstX = destPosition.X - (w.tileSize * renderItem.Scale / 2)
					dstY = destPosition.Y - (w.tileSize*renderItem.Scale - w.tileSize)
				}

				destRect := rl.NewRectangle(dstX, dstY, w.tileSize*renderItem.Scale, w.tileSize*renderItem.Scale)
				result = append(result,
					RenderListItem{
						Src:     renderItem.TexRect,
						Dst:     destRect,
						zindex:  renderItem.Zindex,
						Texture: *w.currentMap.GetTexture(),
					})
			}
		}
	}

	// Calculate entities - for now it's just the player so we won't overcomplicate it
	// it's okay to be messy for now and will be cleaner once we utilizie Entity as an interface
	playerTex, playerTexRect := w.player.GetSprite()
	playerPosition := w.player.GetCurrentPosition()
	playerGameRect := rl.NewRectangle(
		playerPosition.X-float32(w.tileSize/2),
		playerPosition.Y-float32(w.tileSize/2),
		w.tileSize, w.tileSize,
	)

	result = append(result, RenderListItem{
		Src:     *playerTexRect,
		Dst:     playerGameRect,
		zindex:  2,
		Texture: *playerTex,
	})

	// Sort the renderlist based on zIndex and Y-based depth
	slices.SortFunc(result, func(a, b RenderListItem) int {
		// If the Height is larger than the tilesize, it means we're dealing with a texture that's
		// bigger than a single tile. We adjust for this by shifting the rendering position
		// and we need to adjust for this again here
		var offsetA float32 = 0
		var offsetB float32 = 0

		if a.Dst.Height > w.tileSize {
			offsetA = a.Dst.Height - w.tileSize
		}

		if b.Dst.Height > w.tileSize {
			offsetB = b.Dst.Height - w.tileSize
		}

		if a.zindex > b.zindex {
			return 1
		}
		if a.zindex < b.zindex {
			return -1
		}

		return cmp.Compare(a.Dst.Y+offsetA, b.Dst.Y+offsetB)
	})

	w.renderList = result
}

// Set current map to players map
func (w *World) updateCurrentMap() {
	playerMapId := w.player.GetCurrentMap()
	if worldMap := w.currentMap; worldMap == nil || worldMap.GetId() != playerMapId {
		w.setWorldMap(w.player.GetCurrentMap())
	}
}

// Update Player Position and verify access to cell
func (w *World) updatePlayer() {
	// Very primitive collision detection
	if w.safeCell == nil {
		w.safeCell = grid.GetCellFromPixelPosition(w.player.GetCurrentPosition(), w.tileSize)
	}

	if cell := grid.GetCellFromPixelPosition(w.player.GetCurrentPosition(), w.tileSize); cell != w.safeCell {
		if !w.currentMap.isCellWalkable(*cell) {
			utils.LogDebug("Cell %+v is not walkable, moving to: %+v", cell, w.safeCell)
			w.player.SetTargetPosition(grid.GetCenterCellCoordinates(w.safeCell, w.tileSize))
			w.player.StopMoving()
			return // Abort all consecutive evaluations. Player must move
		}
		w.safeCell = cell
	}

	if w.player.IsMoving() {
		if w.player.ReachDestination() {
			w.player.StopMoving()
		}
		w.player.MoveForFrame()
	}
}

// Update Enemeies Position
func (w *World) updateEnemies() {}

// Check if maps have been loaded yet
func (w *World) mapsAreLoaded() bool {
	if w.availableMaps != nil {
		return true
	}

	return false
}

// Set the active map
func (w *World) setWorldMap(id int) {
	for _, m := range *w.availableMaps {
		if m.GetId() == id {
			w.currentMap = m
			return
		}
	}
	utils.LogError("Didn't set a map when called to")
}

// Load All maps from w.mapDataDir
// does not set the currently chosen map
func (w *World) loadMaps() {
	files := utils.GatherMapMetaDataFiles(w.mapDataDir)
	var result []GameMapInterface

	for _, mapJson := range files {
		gameMap := NewGameMapFromJson(mapJson)
		result = append(result, gameMap)
	}

	w.availableMaps = &result
}

/* Constructor Method --------------------------------------------------------------------------- */

// Create a new world
func NewWorld(tileSize float32, player character.Entity, mapDataDir string) *World {
	return &World{
		player:     player,
		tileSize:   tileSize,
		mapDataDir: mapDataDir,
	}
}
