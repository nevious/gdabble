package types

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type RenderItem struct {
	TexRect rl.Rectangle
	Scale   float32
	Zindex  int
	Texture *rl.Texture2D
}

// Return a new RenderItem containing a
// TexRect: depicting the source of the render in the texture
// Scale: Identifying what scaling needs to be applied on the destination
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
