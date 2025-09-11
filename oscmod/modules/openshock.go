package modules

import (
	"fmt"
	"log"

	"github.com/Glowman554/OpenOSC/openshock"
	"github.com/Glowman554/OpenOSC/oscmod/chatbox"
	"github.com/hypebeast/go-osc/osc"
)

type OpenShockGroup struct {
	shockerIDs       []string
	currentDuration  int
	currentIntensity int
}

type OpenShockModuleContainer struct {
	currentDefaultGroup string

	maxDuration  int
	maxIntensity int

	groups map[string]*OpenShockGroup

	api *openshock.OpenShockApi
}

type OpenShockModule struct {
	container *OpenShockModuleContainer
}

func NewOpenShockModule(token string) OpenShockModule {
	return OpenShockModule{
		container: &OpenShockModuleContainer{
			currentDefaultGroup: "0",
			maxDuration:         30000,
			maxIntensity:        100,
			groups:              map[string]*OpenShockGroup{},
			api:                 openshock.NewOpenShockApi(token),
		},
	}
}

func (m OpenShockModule) Name() string {
	return "OpenShock"
}

func (m OpenShockModule) Id() string {
	return "openshock"
}

func (m OpenShockModule) Init(client *osc.Client, dispatcher *osc.StandardDispatcher) error {
	shockers, err := m.container.api.LoadShockers()
	if err != nil {
		return err
	}

	shockerIDs := []string{}
	for _, i := range shockers {
		shockerIDs = append(shockerIDs, i.Id)
		log.Printf("Found shocker %s (%s, %d, %s)", i.Name, i.Id, i.RfId, i.Model)
	}

	// Group 0 should always contain every possible shocker - the Default group
	m.registerGroup("0", shockerIDs, client, dispatcher)

	// TODO: make groups configurable

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/PiShock/Group", func(msg *osc.Message) {
		if group, ok := msg.Arguments[0].(int32); ok {
			if _, ok := m.container.groups[fmt.Sprint(group)]; ok {
				m.container.currentDefaultGroup = fmt.Sprint(group)
				log.Printf("Setting group to %d", group)
			} else {
				log.Printf("Invalid group %d", group)
			}
		}
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/PiShock/Duration", func(msg *osc.Message) {
		m.container.groups[m.container.currentDefaultGroup].handleDuration(msg, m)
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/PiShock/Intensity", func(msg *osc.Message) {
		m.container.groups[m.container.currentDefaultGroup].handleIntensity(msg, m)
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/PiShock/Shock", func(msg *osc.Message) {
		m.container.groups[m.container.currentDefaultGroup].handleShock(msg, client, m)
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/PiShock/Vibrate", func(msg *osc.Message) {
		m.container.groups[m.container.currentDefaultGroup].handleVibrate(msg, client, m)
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/PiShock/Beep", func(msg *osc.Message) {
		m.container.groups[m.container.currentDefaultGroup].handleBeep(msg, client, m)
	})
	if err != nil {
		return err
	}

	return nil
}

func (m OpenShockModule) Tick(client *osc.Client, chatbox *chatbox.ChatBoxBuilder) error {
	// chatbox.Placeholder("openshock.duration", fmt.Sprint(m.container.groups[m.container.currentDefaultGroup].currentDuration/1000)+"S")
	// chatbox.Placeholder("openshock.intensity", fmt.Sprint(m.container.groups[m.container.currentDefaultGroup].currentIntensity)+"%")
	// chatbox.Placeholder("openshock.group", fmt.Sprint(m.container.currentDefaultGroup))

	return nil
}

func (m OpenShockModule) registerGroup(groupID string, shockerIDs []string, client *osc.Client, dispatcher *osc.StandardDispatcher) error {
	group := &OpenShockGroup{
		shockerIDs:       shockerIDs,
		currentDuration:  0,
		currentIntensity: 0,
	}
	m.container.groups[groupID] = group

	err := dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/PiShock/Duration/"+fmt.Sprint(groupID), func(msg *osc.Message) {
		group.handleDuration(msg, m)
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/PiShock/Intensity/"+groupID, func(msg *osc.Message) {
		group.handleIntensity(msg, m)
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/PiShock/Shock/"+groupID, func(msg *osc.Message) {
		group.handleShock(msg, client, m)
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/PiShock/Vibrate/"+groupID, func(msg *osc.Message) {
		group.handleVibrate(msg, client, m)
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler("/avatar/parameters/VRCOSC/PiShock/Beep/"+groupID, func(msg *osc.Message) {
		group.handleBeep(msg, client, m)
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *OpenShockGroup) handleDuration(msg *osc.Message, m OpenShockModule) {
	if duration, ok := msg.Arguments[0].(float32); ok {
		g.currentDuration = int(float32(m.container.maxDuration) * duration)
		// log.Printf("Duration: %dms", g.currentDuration)
	}
}

func (g *OpenShockGroup) handleIntensity(msg *osc.Message, m OpenShockModule) {
	if intensity, ok := msg.Arguments[0].(float32); ok {
		g.currentIntensity = int(float32(m.container.maxIntensity) * intensity)
		// log.Printf("Intensity: %d%%", g.currentIntensity)
	}
}

func (g *OpenShockGroup) handleShock(msg *osc.Message, client *osc.Client, m OpenShockModule) {
	if shock, ok := msg.Arguments[0].(bool); ok && shock {
		m.container.api.SendCommand(g.currentIntensity, g.currentDuration, openshock.Shock, g.shockerIDs)
		g.sendSuccess(client)
	}
}

func (g *OpenShockGroup) handleVibrate(msg *osc.Message, client *osc.Client, m OpenShockModule) {
	if vibrate, ok := msg.Arguments[0].(bool); ok && vibrate {
		m.container.api.SendCommand(g.currentIntensity, g.currentDuration, openshock.Vibrate, g.shockerIDs)
		g.sendSuccess(client)
	}
}

func (g *OpenShockGroup) handleBeep(msg *osc.Message, client *osc.Client, m OpenShockModule) {
	if beep, ok := msg.Arguments[0].(bool); ok && beep {
		// Should BEEP but i don't want it too
		m.container.api.SendCommand(g.currentIntensity, g.currentDuration, openshock.Vibrate, g.shockerIDs)
		g.sendSuccess(client)
	}
}

func (g *OpenShockGroup) sendSuccess(client *osc.Client) {
	go func() {
		msg := osc.NewMessage("/avatar/parameters/VRCOSC/PiShock/Success")
		msg.Append(true)
		err := client.Send(msg)
		if err != nil {
			log.Println(err)
		}

		msg = osc.NewMessage("/avatar/parameters/VRCOSC/PiShock/Success")
		msg.Append(false)
		err = client.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}()
}
