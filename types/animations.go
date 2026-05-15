package types

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type SpriteAnimation struct {
	SpriteSheet  *rl.Texture2D
	FrameCount   int
	FPS          int
	CurrentFrame int
	FrameSize    rl.Vector2
	timer        float32
}

type Animation interface {
	// Get the next frame
	NextFrame()
	// Load the Texture
	LoadSprite()
	// Unload the Texture
	UnloadSprite()
}

func NewAnimation(sprite string, fCount, fps, fInitial int, frameSize rl.Vector2, timer float32) {

}
