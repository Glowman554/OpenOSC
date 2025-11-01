package modules

import (
	"fmt"
	"log"

	"github.com/Glowman554/OpenOSC/config"
	"github.com/Glowman554/OpenOSC/openshock"
	"github.com/Glowman554/OpenOSC/oscmod/chatbox"
	"github.com/hypebeast/go-osc/osc"
)

type OpenShockControlModuleContainer struct {
	config config.OpenShockControlConfig
	api    *openshock.OpenShockApi

	currentDuration  int
	currentIntensity int
}

type OpenShockControlModule struct {
	container *OpenShockControlModuleContainer
}

func NewOpenShockControlModule(configOpenShock config.OpenShockConfig, config config.OpenShockControlConfig) OpenShockControlModule {
	return OpenShockControlModule{
		container: &OpenShockControlModuleContainer{
			config:           config,
			api:              openshock.NewOpenShockApi(configOpenShock.APIToken),
			currentDuration:  0,
			currentIntensity: 0,
		},
	}
}

func (m OpenShockControlModule) Name() string {
	return "OpenShock control"
}

func (m OpenShockControlModule) Id() string {
	return "openshock_control"
}

func (m OpenShockControlModule) Init(client *osc.Client, dispatcher *osc.StandardDispatcher) error {
	shockers, err := m.container.api.LoadShockersShared()
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler(m.container.config.DurationParameter, func(msg *osc.Message) {
		if duration, ok := msg.Arguments[0].(float32); ok {
			m.container.currentDuration = int(float32(m.container.config.MaximumDurationMS) * duration)
			// log.Printf("Duration: %dms", m.container.currentDuration)
		}
	})
	if err != nil {
		return err
	}

	err = dispatcher.AddMsgHandler(m.container.config.IntensityParameter, func(msg *osc.Message) {
		if intensity, ok := msg.Arguments[0].(float32); ok {
			m.container.currentIntensity = int(float32(m.container.config.MaximumIntensity) * intensity)
			// log.Printf("Intensity: %d%%", m.container.currentIntensity)
		}
	})
	if err != nil {
		return err
	}

	for key, s := range m.container.config.Mapping {

		shockerIDs := []string{}
		for _, i := range s {
			if shocker, ok := shockers[i]; ok {
				shockerIDs = append(shockerIDs, shocker.Id)
			} else {
				return fmt.Errorf("failed to find %s", i)
			}
		}

		log.Printf("Registering handler %s for %d shockers", key, len(shockerIDs))

		err = dispatcher.AddMsgHandler(key, func(msg *osc.Message) {
			if trigger, ok := msg.Arguments[0].(bool); ok && trigger {
				err := m.container.api.SendCommand(m.container.currentIntensity, m.container.currentDuration, openshock.Shock, shockerIDs)
				if err != nil {
					log.Printf("Failed to send command %v", err)
				}
				// log.Printf("Shock: %s", msg.Address)
			}
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (m OpenShockControlModule) Tick(client *osc.Client, chatbox *chatbox.ChatBoxBuilder) error {
	return nil
}
