package world

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/character"
	"main/utils"
)

type World struct {
	tileSize      float32
	player        character.Entity
	availableMaps *[]GameMapInterface
	currentMap    GameMapInterface
	mapDataDir    string
}

func (w *World) Update() {
	if w.player.IsMoving() {
		// if player has arrived, stop the move
		if w.player.HasArrived() {
			w.player.StopMoving()
		}

		w.player.UpdatePosition(rl.GetFrameTime())
	}

}

func (w *World) GetPlayer() character.Entity {
	return w.player
}

func (w *World) GetTileSize() float32 {
	return w.tileSize
}

func (w *World) GetMap() GameMapInterface {
	return w.currentMap
}

// Check if maps have been loaded yet
func (w *World) MapsAreLoaded() bool {
	if w.availableMaps != nil {
		return true
	}

	return false
}

func (w *World) SetPlayerMapPosition(id int) {
	for _, _map := range *w.availableMaps {
		if _map.GetId() == id {
			w.currentMap = _map
			return
		}
	}
	utils.LogDebug("Didn't set a map when called to")
}

// Load All maps from w.mapDataDir
// does not set the currently chosen map
func (w *World) LoadMaps() {
	files := utils.GatherMapMetaDataFiles(w.mapDataDir)
	var result []GameMapInterface

	for _, mapJson := range files {
		gameMap := NewGameMapFromJson(mapJson)
		result = append(result, gameMap)
	}

	w.availableMaps = &result
}

func NewWorld(tileSize float32, player character.Entity, mapDataDir string) *World {
	//gameMap := NewGameMap(1, rl.NewVector2(40, 20), "./assets/TinySwordsFreePack/Terrain/Tileset/Tilemap_color1.png", 64)

	newWorld := &World{
		player:     player,
		tileSize:   tileSize,
		mapDataDir: mapDataDir,
	}

	//newWorld.LoadMap(gameMap)
	return newWorld
}
