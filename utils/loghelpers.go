package utils

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nevious/gdabble/entity"
)

func LogPlayerTransition(level rl.TraceLogLevel, ent entity.Entity, current, target rl.Vector2) {
	rl.TraceLog(
		level, "%s", fmt.Sprintf("Entity <%s> move from: %+v to: %+v", ent.GetEntityType(), current, target),
	)
}

func LogDebug(data string, args ...any) {
	rl.TraceLog(rl.LogDebug, "%s", fmt.Sprintf(data, args...))
}

func LogError(data string, args ...any) {
	rl.TraceLog(rl.LogError, "%s", fmt.Sprintf("%s", args...))
}
