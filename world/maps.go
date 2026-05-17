package world

/*
 * Texture to Terrain Mapping
 *
 * TODO: This is description is no longer accurate.
 * For each channel we have 8 values that can be encoded and there are no composits
 * between colors. The only "composite" is black and white. Each channel allowes for
 * encoding of 8 fields for a total of 24 values.
 * For each value, we must know which tilemap to load and which portion
 * of the tilemap to load for a given pixel.
 *
 * For edges, we need to know what's around the given pixel
 *
 * All information must be available from the json, otherwise we have a split brain situation.
 * Therefore each json must define:
 * - All textures to load with a reference
 * - Each Color used must map to both a texture AND a coordinate
 */

import (
	"encoding/json"
	rl "github.com/gen2brain/raylib-go/raylib"
	"io"
	"main/types"
	"main/utils"
	"os"
)

type GameMapInterface interface {
	GetSize() *rl.Vector2
	GetTileAt(x, y int) []types.RenderItem
	GetId() int
	isCellWalkable(position rl.Vector2) bool
}

type JsonMetaData struct {
	ID              int                 `json:"id"`
	TerrainMapFile  string              `json:"terrainMapFile"`
	Texture         string              `json:"texture"`
	TextureTileSize int                 `json:"textureTileSize"`
	Colors          map[string]ColorMap `json:"colors"`
}

type ColorMap struct {
	Locations  [][]int `json:"location"` // X, Y, Scale, Depth
	Accessible bool    `json:"accessible"`
}

type GameMap struct {
	id              int
	terrainMap      *rl.Image
	texture         rl.Texture2D
	colors          map[string]ColorMap
	textureTileSize int
}

func (m *GameMap) GetId() int {
	return m.id
}

func (m *GameMap) GetSize() *rl.Vector2 {
	return &rl.Vector2{
		X: float32(m.terrainMap.Width),
		Y: float32(m.terrainMap.Height),
	}
}

// Get the textures rectangle for given grid coordinates
// MUST be grid coordinates, not world or screen coordinates
func (m *GameMap) GetTileAt(x, y int) []types.RenderItem {
	ttSize := float32(m.textureTileSize)
	var result []types.RenderItem

	mapColor := rl.GetImageColor(*m.terrainMap, int32(x), int32(y))
	hexCol := utils.RaylibColorToHex(mapColor)

	for _, loc := range m.colors[hexCol].Locations {
		srcX, srcY, scale, zindex := float32(loc[0]), float32(loc[1]), float32(loc[2]), int(loc[3])
		multiplier := ttSize * scale
		rayRect := rl.NewRectangle(srcX*ttSize, srcY*ttSize, multiplier, multiplier)
		rect := types.NewRenderItem(rayRect, scale, zindex, &m.texture)
		result = append(result, *rect)
	}

	return result
}

/* Private -------------------------------------------------------------------------------------- */

func (m *GameMap) isCellWalkable(position rl.Vector2) bool {
	c := rl.GetImageColor(*m.terrainMap, int32(position.X), int32(position.Y))
	col := utils.RaylibColorToHex(c)
	return m.colors[col].Accessible
}

/* Constructor method --------------------------------------------------------------------------- */

// Create a new map object from a json file given by `path`
func NewGameMapFromJson(path string) GameMapInterface {
	pFD, err := os.Open(path)
	if err != nil {
		utils.LogError("Unable to read %s: %+v", path, err)
	}

	content, err := io.ReadAll(pFD)
	if err != nil {
		utils.LogError("Unable to read from file: %+v", err)
	}

	jsonMetaData := &JsonMetaData{}
	gameMap := &GameMap{}
	err = json.Unmarshal(content, jsonMetaData)
	if err != nil {
		utils.LogError("Unable to parse json: %+v", err)
	}

	utils.LogDebug("New JSON Loaded: %+v", jsonMetaData)
	terrainImage := rl.LoadImage(jsonMetaData.TerrainMapFile)
	gameMap.id = jsonMetaData.ID
	gameMap.terrainMap = terrainImage
	gameMap.texture = rl.LoadTexture(jsonMetaData.Texture)
	gameMap.colors = jsonMetaData.Colors
	gameMap.textureTileSize = jsonMetaData.TextureTileSize

	return gameMap
}
