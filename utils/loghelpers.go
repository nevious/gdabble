package utils

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func LogPlayerTransition(level rl.TraceLogLevel, current, target rl.Vector2) {
	rl.TraceLog(
		level, "%s", fmt.Sprintf("Player move from: %+v to: %+v", current, target),
	)
}

func LogDebug(data string, args ...any) {
	rl.TraceLog(rl.LogDebug, "%s", fmt.Sprintf(data, args...))
}

func LogError(data string, args ...any) {
	rl.TraceLog(rl.LogError, "%s", fmt.Sprintf("%s", args...))
}
