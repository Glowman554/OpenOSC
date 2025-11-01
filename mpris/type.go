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

type LoopType string

const (
	None     LoopType = "None"
	Track    LoopType = "Track"
	Playlist LoopType = "Playlist"
)

type CurrentlyPlaying struct {
	Title    string
	Artist   []string
	Album    string
	Duration time.Duration
	Position time.Duration
	Status   PlayStatus
	Shuffle  bool
	Loop     LoopType
}
