package handlers

import (
	"fmt"

	"git.leggy.dev/Fluffy/Website/internal/events"
	"git.leggy.dev/Fluffy/Website/internal/sse"
	"git.leggy.dev/Fluffy/Website/internal/web"
)

func RegisterUserJoinedHandler(h *web.Handler) func(data events.UserJoined) {
	return func(data events.UserJoined) {
		h.SSE.Broadcast(sse.Message{
			Name:  "System",
			Value: fmt.Sprintf("Welcome to the chatroom %s!", data.Name),
		})
	}
}
