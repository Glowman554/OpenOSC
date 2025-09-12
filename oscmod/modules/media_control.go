package modules

import (
	"log"

	"github.com/Glowman554/OpenOSC/mpris"
	"github.com/Glowman554/OpenOSC/oscmod/chatbox"
	"github.com/hypebeast/go-osc/osc"
)

type MediaControlModuleContainer struct {
	dbus           *mpris.DBUSInterface
	currentPlayer  *string
	seekToPosition float32
}

type MediaControlModule struct {
	container *MediaControlModuleContainer
}

func NewMediaControlModule() MediaControlModule {
	return MediaControlModule{
		container: &MediaControlModuleContainer{
			dbus:          &mpris.DBUSInterface{},
			currentPlayer: nil,
		},
	}
}

func (m MediaControlModule) Name() string {
	return "Meida control"
}

func (m MediaControlModule) Id() string {
	return "media_control"
}

func (m MediaControlModule) Init(client *osc.Client, dispatcher *osc.StandardDispatcher) error {
	err := m.container.dbus.Connect()
	if err != nil {
		return err
	}

	// VRCOSC/Media/Muted // idk?
	// VRCOSC/Media/Volume // maybe?

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/Media/Play", func(msg *osc.Message) {
		if play, ok := msg.Arguments[0].(bool); ok {
			if m.container.currentPlayer == nil {
				return
			}

			if play {
				m.container.dbus.Play(*m.container.currentPlayer)
			} else {
				m.container.dbus.Pause(*m.container.currentPlayer)
			}
		}
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/Media/Next", func(msg *osc.Message) {
		if next, ok := msg.Arguments[0].(bool); ok && next {
			if m.container.currentPlayer == nil {
				return
			}
			m.container.dbus.Next(*m.container.currentPlayer)
		}
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/Media/Previous", func(msg *osc.Message) {
		if previous, ok := msg.Arguments[0].(bool); ok && previous {
			if m.container.currentPlayer == nil {
				return
			}
			m.container.dbus.Previous(*m.container.currentPlayer)
		}
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/Media/Repeat", func(msg *osc.Message) {
		if mode, ok := msg.Arguments[0].(int32); ok {
			if m.container.currentPlayer == nil {
				return
			}

			switch mode {
			case 0: // No repeat
				m.container.dbus.Loop(*m.container.currentPlayer, mpris.None)
			case 1: // Repeat track
				m.container.dbus.Loop(*m.container.currentPlayer, mpris.Track)
			case 2: // Repeat playlist
				m.container.dbus.Loop(*m.container.currentPlayer, mpris.Playlist)
			}
		}
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/Media/Shuffle", func(msg *osc.Message) {
		if shuffle, ok := msg.Arguments[0].(bool); ok {
			if m.container.currentPlayer == nil {
				return
			}

			m.container.dbus.Shuffle(*m.container.currentPlayer, shuffle)
		}
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/Media/Seeking", func(msg *osc.Message) {
		if seek, ok := msg.Arguments[0].(bool); ok {
			if !seek {
				if m.container.currentPlayer == nil {
					return
				}

				m.container.dbus.Seek(*m.container.currentPlayer, m.container.seekToPosition)
			}
		}
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/Media/Position", func(msg *osc.Message) {
		if position, ok := msg.Arguments[0].(float32); ok {
			m.container.seekToPosition = position
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func (m MediaControlModule) Tick(client *osc.Client, chatbox *chatbox.ChatBoxBuilder) error {
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
			ratio := float32(playing.Position) / float32(playing.Duration)

			msg := osc.NewMessage("/avatar/parameters/VRCOSC/Media/Position")
			msg.Append(ratio)
			err := client.Send(msg)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
				return err
			}
		}

		m.container.currentPlayer = nil
		if playing.Status == mpris.Playing || playing.Status == mpris.Paused {
			m.container.currentPlayer = &player
			chatbox.Placeholder("media.control.player", player)

			msg := osc.NewMessage("/avatar/parameters/VRCOSC/Media/Play")
			if playing.Status == mpris.Playing {
				msg.Append(true)
			} else {
				msg.Append(false)
			}
			err := client.Send(msg)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
				return err
			}

			msg = osc.NewMessage("/avatar/parameters/VRCOSC/Media/Repeat")
			msg.Append(m.loopTypeToId(playing.Loop))
			err = client.Send(msg)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
				return err
			}

			msg = osc.NewMessage("/avatar/parameters/VRCOSC/Media/Shuffle")
			msg.Append(playing.Shuffle)
			err = client.Send(msg)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
				return err
			}

			break
		}

	}

	return nil
}

func (m MediaControlModule) loopTypeToId(loop mpris.LoopType) int32 {
	switch loop {
	case mpris.None:
		return 0
	case mpris.Track:
		return 1
	case mpris.Playlist:
		return 2
	default:
		return 0
	}
}
