package modules

import (
	"fmt"
	"log"
	"time"

	"github.com/Glowman554/OpenOSC/oscmod/chatbox"
	"github.com/hypebeast/go-osc/osc"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type SysInfoModuleContainer struct {
	currentCpu    int
	currentMemory int
}

type SysInfoModule struct {
	container *SysInfoModuleContainer
}

func NewSysInfoModule() SysInfoModule {
	return SysInfoModule{
		container: &SysInfoModuleContainer{},
	}
}

func (m SysInfoModule) Name() string {
	return "System information"
}

func (m SysInfoModule) Init(client *osc.Client, dispatcher *osc.StandardDispatcher) error {
	return nil
}

func (m SysInfoModule) Tick(client *osc.Client, chatbox *chatbox.ChatBoxBuilder) error {
	m.triggerMeasure()

	chatbox.Placeholder("sysinfo.cpu", fmt.Sprintf("%d%%", m.container.currentCpu))
	chatbox.Placeholder("sysinfo.memory", fmt.Sprintf("%d%%", m.container.currentMemory))

	return nil
}

func (m SysInfoModule) triggerMeasure() {
	go func() {
		percent, err := cpu.Percent(time.Second, false)
		if err != nil {
			log.Printf("Failed to read cpu percentage: %v", err)
			return
		}
		m.container.currentCpu = int(percent[0])

		vm, err := mem.VirtualMemory()
		if err != nil {
			log.Printf("Failed to read memory percentage: %v", err)
			return
		}
		m.container.currentMemory = int(vm.UsedPercent)
	}()
}
