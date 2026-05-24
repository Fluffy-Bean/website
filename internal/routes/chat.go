package routes

import (
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"git.leggy.dev/Fluffy/Website/internal/events"
	"git.leggy.dev/Fluffy/Website/internal/sse"
	"git.leggy.dev/Fluffy/Website/internal/web"
)

const maxMessageSize = 512 // chars

func RegisterChatRoutes(h *web.Handler, r *chi.Mux) {
	r.Get("/chat", chatGet(h))
	r.Post("/chat/send", chatSendPost(h))
	r.Get("/chat/connect", chatConnectGet(h))
}

func chatGet(h *web.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Template(w, r, "templates/pages/chat.html", web.Data{
			"MaxMessageSize": maxMessageSize,
		})
	}
}

func chatSendPost(h *web.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := h.GetToken(r)
		if err != nil {
			slog.Error("verify token", "error", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		conn := h.SSE.GetConnection(claims.Sub)
		if conn == nil {
			slog.Error("no connection found for sub", "sub", claims.Sub)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)

			return
		}

		var input struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			slog.Error("decode chat message", "error", err)
			http.Error(w, "Malformed input", http.StatusBadRequest)

			return
		}

		if input.Message == "" || len(input.Message) > maxMessageSize {
			http.Error(w, "Malformed input", http.StatusBadRequest)

			return
		}

		go h.SSE.Broadcast(sse.Message{
			Name:  conn.Name,
			Value: input.Message,
		})

		w.WriteHeader(http.StatusAccepted)
	}
}

func chatConnectGet(h *web.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rc := http.NewResponseController(w)

		conn := sse.NewConnection(randomName())

		claims, _ := h.GetToken(r)
		if claims != nil {
			conn.ID = claims.Sub
			conn.Name = claims.Name
		}

		h.SSE.Subscribe(conn)
		defer h.SSE.Unsubscribe(conn)

		// They didn't already have a claim (or it was expired) so they need a new one
		if claims == nil {
			if err := h.SetToken(w, conn.ID, conn.Name); err != nil {
				slog.Error("set token", "error", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)

				return
			}
		}

		w.Header().Set("content-type", "text/event-stream")
		w.Header().Set("cache-control", "no-cache")
		w.Header().Set("connection", "keep-alive")

		w.WriteHeader(http.StatusOK)

		if err := rc.Flush(); err != nil {
			slog.Error("flush headers", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)

			return
		}

		go h.Events.BroadcastEvent(events.UserJoined{
			ID:   conn.ID,
			Name: conn.Name,
		})
		defer func() {
			go h.Events.BroadcastEvent(events.UserLeft{
				ID:   conn.ID,
				Name: conn.Name,
			})
		}()

		rc.SetWriteDeadline(time.Now().Add(h.SSE.Heartbeat * 2))

		for {
			select {
			case _, ok := <-conn.Heartbeat:
				if !ok {
					return
				}

				rc.SetWriteDeadline(time.Now().Add(h.SSE.Heartbeat * 2))
				w.Write([]byte(":\n\n"))

				if err := rc.Flush(); err != nil {
					slog.Error("flush heartbeat data", "error", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)

					return
				}

			case message, ok := <-conn.Messages:
				if !ok {
					return
				}

				usernameBase64 := base64Encode([]byte(message.Name))
				messageBase64 := base64Encode([]byte(message.Value))

				data := usernameBase64 + "." + messageBase64

				w.Write([]byte("event: message\n"))
				w.Write([]byte("data: " + data + "\n"))
				w.Write([]byte("\n\n"))

				if err := rc.Flush(); err != nil {
					slog.Error("flush event data", "error", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)

					return
				}

			case <-ctx.Done():
				return
			}
		}
	}
}

func base64Encode(src []byte) string {
	return base64.RawStdEncoding.EncodeToString(src)
}

func randomName() string {
	animals := []string{
		"Albatross", "Alligator", "Alpaca", "Ant", "Antelope", "Ape", "Armadillo",
		"Donkey",
		"Baboon", "Badger", "Bat", "Bear", "Beaver", "Bee", "Beetle", "Binturong", "Bison", "Boar", "Bobcat", "Bull", "Butterfly",
		"Camel", "Cardinal", "Cat", "Cattle", "Cheetah", "Chicken", "Chinchilla", "Chough", "Cobra", "Cockroach", "Cod", "Cormorant", "Cow", "Crab", "Crane", "Crocodile", "Crow",
		"Deer", "Dog", "Squalidae", "Dogfish", "Dolphin", "Donkey", "Dove", "Duck", "Dunlin",
		"Eagle", "Echidna", "Eel", "Elephant", "Elk", "Emu",
		"Falcon", "Ferret", "Finch", "Flamingo", "Fly", "Fox", "Frog",
		"Gecko", "Gerbil", "Giant panda", "Giraffe", "Gnat", "Wildebeest", "Gnu", "Goat", "Goldfinch", "Goosander", "Goose", "Gorilla", "Goshawk", "Grasshopper", "Grouse", "Guanaco", "Guinea fowl", "Guinea pig",
		"Hamster", "Hare", "Hawk", "Hedgehog", "Heron", "Herring", "Hippopotamus", "Hornet", "Horse", "Human", "Hyena",
		"Jaguar", "Jay", "Jellyfish", "Junglefowl",
		"Kangaroo", "Kingbird", "Kinkajou", "Koala",
		"Ladybug", "Lapwing", "Lark", "Lemur", "Leopard", "Lion", "Lizard", "Llama", "Lobster", "Locust", "Lyrebird",
		"Mallard", "Meerkat", "Mole", "Mongoose", "Monkey", "Moose", "Mosquito", "Moth", "Mouse",
		"Narwhal", "Newt", "Nightingale",
		"Octopus", "Opossum", "Otter", "Ox", "Owl", "Oyster",
		"Panda", "Partridge", "Peafowl", "Peccar", "Pelica", "Pengui", "Pheasant", "Pig", "Pigeon", "Platypus", "Polar bear", "Pony", "Porcupine", "Prairie dog", "Pug",
		"Quail",
		"Rabbit", "Raccoon", "Sheep", "Ram", "Rat", "Raven", "Rhinoceros", "Rook",
		"Salamander", "Salmon", "Sand dollar", "Sandpiper", "Sardine", "Seahorse", "Pinniped", "Seal", "Sea otter", "Shark", "Sheep", "Skunk", "Snail", "Snake", "Spider", "Spoonbill", "Squid", "Squirrel", "Starfish", "Starling", "Stingray", "Swan",
		"Termite", "Thrush", "Tiger", "Toad", "Toucan", "Turkey", "Turtle",
		"Wallaby", "Wasp", "Weasel", "Whale", "Wildebeest", "Wolf", "Wombat",
		"Zebra",
	}

	i := rand.IntN(len(animals))

	return animals[i]
}
