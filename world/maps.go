package world

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Transition struct {
	position      *rl.Vector2
	destinationID int
	spawnPoint    *rl.Vector2
}

type GameMapInterface interface {
	GetSize() *rl.Vector2
}

type GameMap struct {
	id             int
	dimensions     rl.Vector2
	transitions    []Transition
	terrainMapFile string
}

func (m *GameMap) GetSize() *rl.Vector2 {
	return &m.dimensions
}
