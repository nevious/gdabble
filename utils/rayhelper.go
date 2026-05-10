package utils

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RenderItem struct {
	TexRect rl.Rectangle
	Scale   float32
	Zindex  int
	Texture *rl.Texture2D
}

// Return the point on X-axis at which a given text with size is centered
func GetCenterForText(data string, size int32, font rl.Font) int32 {
	vec := rl.MeasureTextEx(font, data, float32(size), 0)
	return int32(
		(float32(rl.GetScreenWidth()) / 2) - (vec.X / 2),
	)
}

// Return a numerical string representation of rl.Color  without Alpa channel
func RaylibColorToHex(color rl.Color) string {
	value := fmt.Sprintf("%02x%02x%02x", color.R, color.G, color.B)
	return value
}

// Return a new RenderItem containing a
// Rect: depicting the source of the render
// Scale: depting what sort of scaling needs to be applied on the destination
// Zindex: Indicating the ordinal position of the thing to be rendered
// Texture: Referencing what texture to puill from
func NewRenderItem(texRect rl.Rectangle, scale float32, zindex int, texture *rl.Texture2D) *RenderItem {
	return &RenderItem{
		TexRect: texRect,
		Scale:   scale,
		Zindex:  zindex,
		Texture: texture,
	}
}
