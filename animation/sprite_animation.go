package animation

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nevious/gdabble/types"
)

type SpriteAnimation struct {
	SpriteSheet  rl.Texture2D
	FrameCount   int
	CurrentFrame int
	FrameSize    float32
	FPS          float32
	Scale        float32
	Zindex       int
	timer        float32
}

/* Interface Implementation --------------------------------------------------------------------- */

// Get the next frame if the timer has passed
func (sa *SpriteAnimation) NextFrame(reverseDirection bool) types.RenderItem {
	sa.timer += rl.GetFrameTime()
	frameTime := 1 / sa.FPS

	if sa.timer >= frameTime {
		sa.CurrentFrame = (sa.CurrentFrame + 1) % sa.FrameCount
		sa.timer -= frameTime
	}

	var sizeX float32
	if reverseDirection {
		sizeX = -sa.FrameSize
	} else {
		sizeX = sa.FrameSize
	}

	// If rimer has passed return rectrangle with currentFrame+1 and reset currentFrame
	return *types.NewRenderItem(
		rl.NewRectangle(float32(sa.CurrentFrame)*sa.FrameSize, 0, sizeX, sa.FrameSize),
		sa.Scale, sa.Zindex, &sa.SpriteSheet,
	)
}

// Create a new animation
// sprite: File path, loaded in constructor
// fCount, fSize, fps, fInit -> frame Count, frame Size, fps, frame init (set current)
func NewSpriteAnimation(sprite string, fCount, fInit int, fSize, fps, scale float32, zindex int) Animation {
	return &SpriteAnimation{
		SpriteSheet:  rl.LoadTexture(sprite),
		FrameCount:   fCount,
		FrameSize:    fSize,
		FPS:          fps,
		CurrentFrame: fInit,
		Scale:        scale,
		Zindex:       zindex,
	}
}
