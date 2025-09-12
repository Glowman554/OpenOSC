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

func (m SysInfoModule) Id() string {
	return "sysinfo"
}

func (m SysInfoModule) Init(client *osc.Client, dispatcher *osc.StandardDispatcher) error {
	return nil
}

func (m SysInfoModule) Tick(client *osc.Client, chatbox *chatbox.ChatBoxBuilder) error {
	m.triggerMeasure()
	time24h, time12h := m.getCurrentTime()

	chatbox.Placeholder("sysinfo.cpu", fmt.Sprintf("%d%%", m.container.currentCpu))
	chatbox.Placeholder("sysinfo.memory", fmt.Sprintf("%d%%", m.container.currentMemory))
	chatbox.Placeholder("sysinfo.time.12h", time12h)
	chatbox.Placeholder("sysinfo.time.24h", time24h)

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

func (m SysInfoModule) getCurrentTime() (string, string) {
	now := time.Now()
	time24h := now.Format("15:04:05")
	time12h := now.Format("03:04:05 PM")
	return time24h, time12h
}
