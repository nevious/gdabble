package animation

import (
	"main/types"
)

type Animation interface {
	// Get the next frame
	NextFrame(flipY bool) types.RenderItem
}
