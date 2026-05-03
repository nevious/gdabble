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

func (w *World) LoadMap(gameMap GameMapInterface) {
	w.gameMap = gameMap
}

func NewWorld(tileSize float32, player character.Entity) *World {
	gameMap := NewGameMap(1, rl.NewVector2(40, 20), "./assets/terrain/basicSand.png")

	newWorld := &World{
		player:   player,
		tileSize: tileSize,
	}

	newWorld.LoadMap(gameMap)
	return newWorld
}
