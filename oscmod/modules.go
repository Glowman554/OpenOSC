package oscmod

import (
	"github.com/Glowman554/OpenOSC/oscmod/chatbox"
	"github.com/Glowman554/OpenOSC/oscmod/modules"
	"github.com/hypebeast/go-osc/osc"
)

type OSCModule interface {
	Name() string
	Init(client *osc.Client, dispatcher *osc.StandardDispatcher) error
	Tick(client *osc.Client, chatbox *chatbox.ChatBoxBuilder) error
}

var Modules = []OSCModule{
	modules.NewMediaChatBoxModule(),
	modules.NewMediaControlModule(),
	modules.NewSysInfoModule(),
}
