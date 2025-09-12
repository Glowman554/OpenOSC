package oscmod

import (
	"log"

	"github.com/hypebeast/go-osc/osc"
)

type Player struct {
	client *osc.Client
}

func NewPlayer(client *osc.Client) *Player {
	return &Player{
		client: client,
	}
}

func (p *Player) Run() {
	p.send("/input/Run", true)
}

func (p *Player) StopRun() {
	p.send("/input/Run", false)
}

func (p *Player) MoveVertical(v float32) {
	p.send("/input/Vertical", v)
}

func (p *Player) MoveHorizontal(v float32) {
	p.send("/input/Horizontal", v)
}

func (p *Player) LookHorizontal(v float32) {
	p.send("/input/LookHorizontal", v)
}

func (p *Player) send(path string, value any) {
	msg := osc.NewMessage(path)
	msg.Append(value)
	err := p.client.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
