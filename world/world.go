package world

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"main/character"
)

type World struct {
	tileSize float32
	player   character.Entity
	gameMap  GameMapInterface
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
	return w.gameMap
}

func (w *World) LoadMap()   {}
func (w *World) UnloadMap() {}

func NewWorld(tileSize float32, player character.Entity) *World {
	m := GameMap{
		id:             1,
		dimensions:     rl.Vector2{X: 30, Y: 15},
		transitions:    nil,
		terrainMapFile: "",
	}

	return &World{
		player:   player,
		tileSize: tileSize,
		gameMap:  &m,
	}
}
