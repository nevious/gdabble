package world

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type tileType int

const (
	TL tileType = iota
	TN
	TR
	LN
	RN
	BL
	BN
	BR
	C
)

type Transition struct {
	position      *rl.Vector2
	destinationID int
	spawnPoint    *rl.Vector2 // Where on the destination to spawn
}

type GameMapInterface interface {
	GetSize() *rl.Vector2
	GetTexture() *rl.Texture2D
	GetTileAt(x, y int) *rl.Rectangle
}

type GameMap struct {
	id              int
	dimensions      rl.Vector2
	transitions     []Transition
	terrainData     [][]tileType
	tileSetFile     string
	texture         *rl.Texture2D
	textureTileSize int
}

func (m *GameMap) GetSize() *rl.Vector2 {
	return &m.dimensions
}

// Get the texture rectangle for given grid coordinates
// MUST be grid coordinates, not world or screen coordinates
func (m *GameMap) GetTileAt(x, y int) *rl.Rectangle {
	// First, look up the type of tile we need
	var srcX, srcY float32
	ttSize := float32(m.textureTileSize)

	switch tileType(m.terrainData[y][x]) {
	case TL:
		srcX, srcY = 0, 0
	case TN:
		srcX, srcY = 1, 0
	case TR:
		srcX, srcY = 2, 0
	case LN:
		srcX, srcY = 0, 1
	case RN:
		srcX, srcY = 2, 1
	case BL:
		srcX, srcY = 0, 2
	case BN:
		srcX, srcY = 1, 2
	case BR:
		srcX, srcY = 2, 2
	default:
		srcX, srcY = 1, 1
	}

	srcRect := rl.NewRectangle(srcX*ttSize, srcY*ttSize, ttSize, ttSize)
	return &srcRect
}

func (m *GameMap) GetTexture() *rl.Texture2D {
	if m.texture == nil {
		texture := rl.LoadTexture(m.tileSetFile)
		m.texture = &texture
	}

	return m.texture
}

func NewGameMap(id int, dimensions rl.Vector2, tileSetFile string, textureTileSize int) GameMapInterface {
	terrain := [][]tileType{
		{TL, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TN, TR},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{LN, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, C, RN},
		{BL, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BN, BR},
	}

	return &GameMap{
		id:              id,
		dimensions:      dimensions,
		transitions:     nil,
		tileSetFile:     tileSetFile,
		terrainData:     terrain,
		textureTileSize: textureTileSize,
	}

}
