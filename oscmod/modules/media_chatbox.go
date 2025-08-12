package modules

import (
	"fmt"
	"strings"
	"time"

	"github.com/Glowman554/OpenOSC/mpris"
	"github.com/Glowman554/OpenOSC/oscmod/chatbox"
	"github.com/hypebeast/go-osc/osc"
)

type MediaChatBoxModuleContainer struct {
	dbus *mpris.DBUSInterface
}

type MediaChatBoxModule struct {
	container *MediaChatBoxModuleContainer
}

func NewMediaChatBoxModule() MediaChatBoxModule {
	return MediaChatBoxModule{
		container: &MediaChatBoxModuleContainer{
			dbus: &mpris.DBUSInterface{},
		},
	}
}

func (m MediaChatBoxModule) Name() string {
	return "Meida ChatBox"
}

func (m MediaChatBoxModule) Init(client *osc.Client, dispatcher *osc.StandardDispatcher) error {
	err := m.container.dbus.Connect()
	return err
}

func (m MediaChatBoxModule) Tick(client *osc.Client, chatbox *chatbox.ChatBoxBuilder) error {
	players, err := m.container.dbus.LoadPlayers()
	if err != nil {
		return err
	}

	for _, player := range players {
		playing, err := m.container.dbus.LoadCurrentlyPlaying(player)
		if err != nil {
			return err
		}

		if playing.Status == mpris.Playing {
			chatbox.Placeholder("media.title", playing.Title)
			chatbox.Placeholder("media.album", playing.Album)
			chatbox.Placeholder("media.artist", strings.Join(playing.Artist, ", "))
			chatbox.Placeholder("media.progress", m.makeProgressBar(playing.Position, playing.Duration, 15))
			chatbox.Placeholder("media.player", player)

			break
		}

	}

	return nil
}

func (m MediaChatBoxModule) formatDuration(dur time.Duration) string {
	minutes := int(dur.Minutes())
	seconds := int(dur.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (m MediaChatBoxModule) makeProgressBar(position time.Duration, duration time.Duration, width int) string {
	ratio := float64(position) / float64(duration)
	filled := int(ratio * float64(width))

	if filled > width {
		filled = width
	}

	bar := fmt.Sprintf("|%s%s| %s / %s",
		strings.Repeat("x", m.max(filled, 0)),
		strings.Repeat("-", m.max(width-filled, 0)),
		m.formatDuration(position),
		m.formatDuration(duration),
	)

	return bar
}

func (m MediaChatBoxModule) max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
