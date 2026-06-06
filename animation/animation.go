package animation

import (
	"github.com/nevious/gdabble/types"
)

type Animation interface {
	// Get the next frame
	NextFrame(flipY bool) types.RenderItem
}
