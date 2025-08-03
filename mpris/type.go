package mpris

import (
	"time"
)

type PlayStatus int

const (
	Playing PlayStatus = iota
	Paused
	Stopped
	Unknown
)

type LoopType int

const (
	None LoopType = iota
	Track
	Playlist
)

type CurrentlyPlaying struct {
	Title    string
	Artist   []string
	Album    string
	Duration time.Duration
	Position time.Duration
	Status   PlayStatus
}
