package events

type UserJoined struct {
	ID   int64
	Name string
}

type UserLeft struct {
	ID   int64
	Name string
}

type NewSong struct {
	Title  string
	Artist string
}
