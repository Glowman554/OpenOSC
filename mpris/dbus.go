package mpris

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
)

type DBUSInterface struct {
	session *dbus.Conn
}

func (d *DBUSInterface) Connect() error {
	con, err := dbus.SessionBus()
	if err != nil {
		log.Printf("Failed to connect to session bus: %v", err)
		return err
	}

	d.session = con
	return nil
}

func (d *DBUSInterface) LoadPlayers() ([]string, error) {

	var names []string
	err := d.session.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names)
	if err != nil {
		log.Printf("Failed to list D-Bus names: %v", err)
		return nil, err
	}

	players := []string{}
	for _, name := range names {
		if strings.HasPrefix(name, "org.mpris.MediaPlayer2.") {
			players = append(players, name)
			// log.Printf("Found player %s", name)
		}
	}

	return players, nil
}

func (d *DBUSInterface) LoadCurrentlyPlaying(player string) (*CurrentlyPlaying, error) {
	obj := d.session.Object(player, "/org/mpris/MediaPlayer2")
	metaVariant, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
	if err != nil {
		log.Printf("Failed to get Metadata: %v", err)
		return nil, err
	}
	metadata := metaVariant.Value().(map[string]dbus.Variant)

	title := d.getString(metadata, "xesam:title")
	artist := d.getStringList(metadata, "xesam:artist")
	album := d.getString(metadata, "xesam:album")
	duration := time.Duration(d.getInt64(metadata, "mpris:length")) * time.Microsecond

	posVariant, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Position")
	if err != nil {
		log.Printf("Failed to get Position: %v", err)
		return nil, err
	}

	position := time.Duration(posVariant.Value().(int64)) * time.Microsecond

	statusVariant, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
	if err != nil {
		log.Printf("Failed to get PlaybackStatus: %v", err)
		return nil, err
	}

	status := d.stringToStatus(statusVariant.Value().(string))

	shuffleVariant, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Shuffle")
	shuffle := false
	if err != nil {
		// log.Printf("Failed to get Shuffle: %v", err)
	} else {
		shuffle = shuffleVariant.Value().(bool)
	}

	loopVariant, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.LoopStatus")
	loopStatus := None
	if err != nil {
		// log.Printf("Failed to get LoopStatus: %v", err)
	} else {
		loopStatus = d.stringToLoopType(loopVariant.Value().(string))
	}

	return &CurrentlyPlaying{
		Title:    title,
		Artist:   artist,
		Album:    album,
		Duration: duration,
		Position: position,
		Status:   status,
		Shuffle:  shuffle,
		Loop:     loopStatus,
	}, nil
}

func (d *DBUSInterface) commonCall(player string, command string) error {
	obj := d.session.Object(player, "/org/mpris/MediaPlayer2")

	call := obj.Call(command, 0)
	if call.Err != nil {
		log.Printf("Failed to call %s: %v", command, call.Err)
		return call.Err
	}

	return nil
}

func (d *DBUSInterface) Play(player string) error {
	return d.commonCall(player, "org.mpris.MediaPlayer2.Player.Play")
}

func (d *DBUSInterface) Pause(player string) error {
	return d.commonCall(player, "org.mpris.MediaPlayer2.Player.Pause")
}

func (d *DBUSInterface) PlayPause(player string) error {
	return d.commonCall(player, "org.mpris.MediaPlayer2.Player.PlayPause")
}

func (d *DBUSInterface) Stop(player string) error {
	return d.commonCall(player, "org.mpris.MediaPlayer2.Player.Stop")
}

func (d *DBUSInterface) Next(player string) error {
	return d.commonCall(player, "org.mpris.MediaPlayer2.Player.Next")
}

func (d *DBUSInterface) Previous(player string) error {
	return d.commonCall(player, "org.mpris.MediaPlayer2.Player.Previous")
}

func (d *DBUSInterface) Shuffle(player string, enabled bool) error {
	obj := d.session.Object(player, "/org/mpris/MediaPlayer2")

	call := obj.Call("org.freedesktop.DBus.Properties.Set", 0, "org.mpris.MediaPlayer2.Player", "Shuffle", dbus.MakeVariant(enabled))
	if call.Err != nil {
		log.Printf("Failed to set shuffle: %v", call.Err)
		return call.Err
	}

	return nil
}

func (d *DBUSInterface) Loop(player string, status LoopType) error {
	obj := d.session.Object(player, "/org/mpris/MediaPlayer2")

	call := obj.Call("org.freedesktop.DBus.Properties.Set", 0, "org.mpris.MediaPlayer2.Player", "LoopStatus", dbus.MakeVariant(status))
	if call.Err != nil {
		log.Printf("Failed to set loop: %v", call.Err)
		return call.Err
	}

	return nil
}

func (d *DBUSInterface) Seek(player string, position float32) error {
	obj := d.session.Object(player, "/org/mpris/MediaPlayer2")

	variant, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
	if err != nil {
		log.Printf("Failed to get Metadata: %v", err)
		return err
	}

	metadata := variant.Value().(map[string]dbus.Variant)

	durationMicros := d.getInt64(metadata, "mpris:length")

	raw := metadata["mpris:trackid"].Value()
	var trackId dbus.ObjectPath

	switch v := raw.(type) {
	case dbus.ObjectPath:
		trackId = v
	case string:
		trackId = dbus.ObjectPath(v)
	default:
		log.Printf("Unexpected type for mpris:trackid: %T", v)
		return fmt.Errorf("invalid type for trackid: %T", v)
	}

	targetMicros := int64(position * float32(durationMicros))

	call := obj.Call("org.mpris.MediaPlayer2.Player.SetPosition", 0, trackId, targetMicros)
	if call.Err != nil {
		log.Printf("Failed to seek: %v", call.Err)
		return err
	}

	return nil
}

func (d *DBUSInterface) getString(m map[string]dbus.Variant, key string) string {
	if v, ok := m[key]; ok {
		return v.Value().(string)
	}
	return ""
}

func (d *DBUSInterface) getStringList(m map[string]dbus.Variant, key string) []string {
	if v, ok := m[key]; ok {
		return v.Value().([]string)
	}
	return []string{}
}

func (d *DBUSInterface) getInt64(m map[string]dbus.Variant, key string) int64 {
	if v, ok := m[key]; ok {
		i64, ok := v.Value().(int64)
		if !ok {
			return int64(v.Value().(uint64))
		}
		return i64
	}
	return 0
}

func (d *DBUSInterface) stringToStatus(status string) PlayStatus {
	switch status {
	case "Playing":
		return Playing
	case "Paused":
		return Paused
	case "Stopped":
		return Stopped
	default:
		return Unknown
	}
}

func (d *DBUSInterface) stringToLoopType(status string) LoopType {
	switch status {
	case "None":
		return None
	case "Track":
		return Track
	case "Playlist":
		return Playlist
	default:
		return None
	}
}
