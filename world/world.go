package world

import (
	"cmp"
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/entity"
	"main/grid"
	"main/utils"
	"math/rand"
	"slices"
)

/* -------------------------------------------------------------------------------------- */

type World struct {
	tileSize      float32
	availableMaps *[]GameMapInterface
	currentMap    GameMapInterface
	mapDataDir    string
	renderList    []RenderListItem
	// A cell we know the player is allowed to be on, used to fall back on
	safeCell *rl.Vector2
	entities []entity.Entity
}

type RenderListItem struct {
	Src     rl.Rectangle
	Dst     rl.Rectangle
	Texture rl.Texture2D
	Scale   float32
	zindex  int // We manage zindex in world, nowhere else
}

/* Public --------------------------------------------------------------------------------------- */

// Update world state
func (w *World) Update() {
	if !w.mapsAreLoaded() {
		w.loadMaps()
	}

	w.updateCurrentMap()
	w.updateEntities()
	w.updateRenderList()
}

// Return entitys `ent`'s position
func (w *World) GetEntityPosition(wEntity entity.Entity) *rl.Vector2 {
	return wEntity.GetCurrentPosition()
}

// Set the entity's position to the given `position` in
// Cell-Space
func (w *World) SetEntityPosition(wEntity entity.Entity, position *rl.Vector2) {
	wEntity.SetTargetPosition(position)
	utils.LogPlayerTransition(
		rl.LogDebug, wEntity, *wEntity.GetCurrentPosition(), *wEntity.GetTargetPosition(),
	)
}

// Get the currenty state of the given entity
func (w *World) GetEntityActionState(wEntity entity.Entity) entity.EntityMovementState {
	return wEntity.GetActionState()
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

// Return the player entity
func (w *World) GetPlayer() entity.Entity {
	for _, entityEl := range w.entities {
		if entityEl.GetEntityType() == entity.PlayerType {
			return entityEl
		}
	}
	panic("No player could be found")
}

/* Private -------------------------------------------------------------------------------------- */

func (w *World) updateRenderList() {
	mapSize := w.currentMap.GetSize()
	result := []RenderListItem{}

	// Calculate map render items
	for x := float32(0); x < mapSize.X; x++ {
		for y := float32(0); y < mapSize.Y; y++ {
			destPosition := rl.NewVector2(x*w.tileSize, y*w.tileSize)
			textureRectangles := w.currentMap.GetTileAt(int(x), int(y))

			for _, renderItem := range textureRectangles {
				dstX := destPosition.X - (w.tileSize * renderItem.Scale / 2)
				dstY := destPosition.Y - (w.tileSize*renderItem.Scale - w.tileSize)

				destRect := rl.NewRectangle(
					dstX, dstY, w.tileSize*renderItem.Scale, w.tileSize*renderItem.Scale,
				)

				result = append(result,
					RenderListItem{
						Src:     renderItem.TexRect,
						Dst:     destRect,
						zindex:  renderItem.Zindex,
						Scale:   renderItem.Scale,
						Texture: *renderItem.Texture,
					})
			}
		}
	}

	for _, ent := range w.entities {
		entityRenderItem := ent.GetSprite()
		entityPosition := ent.GetCurrentPosition()

		entityDestRect := rl.NewRectangle(
			entityPosition.X-float32(w.tileSize*entityRenderItem.Scale/2),
			entityPosition.Y-float32(w.tileSize*entityRenderItem.Scale/2),
			w.tileSize*entityRenderItem.Scale, w.tileSize*entityRenderItem.Scale,
		)

		result = append(result, RenderListItem{
			Src:     entityRenderItem.TexRect,
			Dst:     entityDestRect,
			zindex:  entityRenderItem.Zindex,
			Scale:   entityRenderItem.Scale,
			Texture: *entityRenderItem.Texture,
		})
	}

	// BUG:
	// Y-coordinate sorting is not doing what i would expect.
	// Sometimes player is in front of the middle cell of  atree
	// Sort the renderlist based on zIndex and Y-based depth
	slices.SortFunc(result, func(a, b RenderListItem) int {
		var offsetA float32 = 0
		var offsetB float32 = 0

		if a.Dst.Height > w.tileSize {
			offsetA = a.Dst.Height - (w.tileSize * a.Scale)
		}

		if b.Dst.Height > w.tileSize {
			offsetB = b.Dst.Height - (w.tileSize * a.Scale)
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
	player := w.GetPlayer()
	playerMapId := player.GetCurrentMap()
	if worldMap := w.currentMap; worldMap == nil || worldMap.GetId() != playerMapId {
		w.setWorldMap(player.GetCurrentMap())
	}
}

// Update Player Position and verify access to cell
func (w *World) updatePlayer() {
	// Very primitive collision detection
	player := w.GetPlayer()
	if w.safeCell == nil {
		w.safeCell = grid.GetCellFromPixelPosition(player.GetCurrentPosition(), w.tileSize)
	}

	if cell := grid.GetCellFromPixelPosition(player.GetCurrentPosition(), w.tileSize); cell != w.safeCell {
		if !w.currentMap.isCellWalkable(*cell) {
			utils.LogDebug("Cell %+v is not walkable, moving to: %+v", cell, w.safeCell)
			player.SetTargetPosition(grid.GetCenterCellCoordinates(w.safeCell, w.tileSize))
			player.StopMoving()
			return // Abort all consecutive evaluations. Player must move
		}
		w.safeCell = cell
	}

	player.Update()
}

// Update all entities
// Orchestration Function, dispatch to Type specific update function
func (w *World) updateEntities() {
	var enemyCount int = 0

	for _, ent := range w.entities {
		switch ent.GetEntityType() {
		case entity.PlayerType:
			w.updatePlayer()
		case entity.EnemeyType:
			w.updateEnemies()
			enemyCount++
		}
	}

	// Spawn an enemy if there is none
	// Double looping for now is fine
	if enemyCount < 10 {
		vec := rl.NewVector2(float32(rand.Intn(40)), float32(rand.Intn(20)))
		spawn := grid.GetCenterCellCoordinates(&vec, w.tileSize)
		w.entities = append(w.entities, entity.NewEnemy(*spawn))
	}
}

// Update Enemeies Position
func (w *World) updateEnemies() {
	for _, ent := range w.entities {
		if ent.GetEntityType() != entity.EnemeyType {
			continue
		}
		utils.LogDebug("Entity: %+v updating", ent)
	}
}

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
func NewWorld(tileSize float32, player entity.Entity, mapDataDir string) *World {
	return &World{
		tileSize:   tileSize,
		mapDataDir: mapDataDir,
		entities:   []entity.Entity{player},
	}
}
