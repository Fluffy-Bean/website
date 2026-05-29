package handlers

import (
	"fmt"

	"git.leggy.dev/Fluffy/Website/internal/events"
	"git.leggy.dev/Fluffy/Website/internal/sse"
	"git.leggy.dev/Fluffy/Website/internal/web"
)

func RegisterNewSongHandler(h *web.Handler) func(data events.NewSong) {
	return func(data events.NewSong) {
		h.SSE.Broadcast(sse.Message{
			Name:  "System",
			Value: fmt.Sprintf("Fluffy is now listening to %s by %s", data.Title, data.Artist),
		})
	}
}
