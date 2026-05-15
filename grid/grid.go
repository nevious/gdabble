package grid

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

func GetCellFromPixelPosition(position *rl.Vector2, cellWidth float32) *rl.Vector2 {
	x := math.Floor(float64(position.X / float32(cellWidth)))
	y := math.Floor(float64(position.Y / float32(cellWidth)))

	return &rl.Vector2{
		X: float32(x), Y: float32(y),
	}
}

// Return World Coordinates
func GetCenterCellCoordinates(pos *rl.Vector2, cellWidth float32) *rl.Vector2 {
	x := pos.X*float32(cellWidth) + float32(cellWidth/2)
	y := pos.Y*float32(cellWidth) + float32(cellWidth/2)

	return &rl.Vector2{
		X: x, Y: y,
	}
}

// Check if a given `pos` is within the `bounds`
// All in world coordinates, not pixel or screen
func CellWithinMapBounds(cell, size *rl.Vector2) bool {
	if cell.X >= 0 && cell.X < size.X && cell.Y >= 0 && cell.Y < size.Y {
		return true
	}
	return false
}

// Calculate the linerar interpolation for a value
func Lerp(start, target, time float32) float32 {
	return start + (target-start)*time
}

// For X, if we move left, the diff will be negative, otherwise positive
// We never want to move more than the calculated step size which is
// player.speed * rl.FrameTime, to result in a pixel/s "step"
// returning current +/- speed means the position is never exceeding this step size
// while the min/max is a clamp not to go beyond the target
func ApproachPoint(current, target, speed float32) float32 {
	diff := target - current
	if diff < 0 {
		return max(current-speed, target)
	}
	return min(current+speed, target)
}
