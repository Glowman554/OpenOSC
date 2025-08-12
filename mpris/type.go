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
	Track             = "Track"
	Playlist          = "Playlist"
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
