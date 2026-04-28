package utils

import rl "github.com/gen2brain/raylib-go/raylib"

// Return the point on X-axis at which a given text with size is centered
func GetCenterForText(data string, size int32, font rl.Font) int32 {
	vec := rl.MeasureTextEx(font, data, float32(size), 0)
	return int32(
		(float32(rl.GetScreenWidth()) / 2) - (vec.X / 2),
	)
}
