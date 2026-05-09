package utils

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

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
