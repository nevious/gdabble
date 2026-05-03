package world

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Transition struct {
	position      *rl.Vector2
	destinationID int
	spawnPoint    *rl.Vector2 // Where on the destination to spawn
}

type GameMapInterface interface {
	GetSize() *rl.Vector2
	GetTexture() *rl.Texture2D
}

type GameMap struct {
	id             int
	dimensions     rl.Vector2
	transitions    []Transition
	terrainMapFile string
	tileSetFile    string
	texture        *rl.Texture2D
}

func (m *GameMap) GetSize() *rl.Vector2 {
	return &m.dimensions
}

func (m *GameMap) GetTexture() *rl.Texture2D {
	if m.texture == nil {
		texture := rl.LoadTexture(m.tileSetFile)
		m.texture = &texture
	}

	return m.texture
}

func NewGameMap(id int, dimensions rl.Vector2, tileSetFile string) GameMapInterface {
	return &GameMap{
		id:             id,
		dimensions:     dimensions,
		transitions:    nil,
		tileSetFile:    tileSetFile,
		terrainMapFile: "",
	}

}
