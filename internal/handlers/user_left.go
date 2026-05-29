package handlers

import (
	"fmt"

	"git.leggy.dev/Fluffy/Website/internal/events"
	"git.leggy.dev/Fluffy/Website/internal/sse"
	"git.leggy.dev/Fluffy/Website/internal/web"
)

func RegisterUserLeftHandler(h *web.Handler) func(data events.UserLeft) {
	return func(data events.UserLeft) {
		h.SSE.Broadcast(sse.Message{
			Name:  "System",
			Value: fmt.Sprintf("%s has left the room.", data.Name),
		})
	}
}
