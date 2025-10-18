package modules

import (
	"log"
	"strconv"

	"github.com/Glowman554/OpenOSC/gpuinfo"
	"github.com/Glowman554/OpenOSC/oscmod/chatbox"
	"github.com/hypebeast/go-osc/osc"
)

type GpuInfoModuleContainer struct {
	usageAMD    []gpuinfo.GPUUsage
	usageNVIDIA []gpuinfo.GPUUsage

	providerAMD    *gpuinfo.AMDProvider
	providerNVIDIA *gpuinfo.NvidiaProvider
}

type GpuInfoModule struct {
	container *GpuInfoModuleContainer
}

func NewGpuInfoModule() GpuInfoModule {
	return GpuInfoModule{
		container: &GpuInfoModuleContainer{
			usageAMD:       []gpuinfo.GPUUsage{},
			usageNVIDIA:    []gpuinfo.GPUUsage{},
			providerAMD:    nil,
			providerNVIDIA: nil,
		},
	}
}

func (m GpuInfoModule) Name() string {
	return "GPU information"
}

func (m GpuInfoModule) Id() string {
	return "gpuinfo"
}

func (m GpuInfoModule) Init(client *osc.Client, dispatcher *osc.StandardDispatcher) error {
	if gpuinfo.CanUseAMDProvider() {
		m.container.providerAMD = gpuinfo.NewAMDProvider()
		log.Print("Enabled AMD provider")
	}

	if gpuinfo.CanUseNvidiaProvider() {
		m.container.providerNVIDIA = gpuinfo.NewNvidiaProvider()
		log.Print("Enabled NVIDIA provider")
	}

	return nil
}

func (m GpuInfoModule) Tick(client *osc.Client, chatbox *chatbox.ChatBoxBuilder) error {
	m.triggerMeasure()

	for _, info := range m.container.usageAMD {
		m.register(chatbox, "amd", info)
	}

	for _, info := range m.container.usageNVIDIA {
		m.register(chatbox, "nvidia", info)
	}

	return nil
}

func (m GpuInfoModule) register(chatbox *chatbox.ChatBoxBuilder, prefix string, info gpuinfo.GPUUsage) {
	chatbox.Placeholder("gpuinfo."+prefix+strconv.Itoa(info.Index)+".name", info.Name)
	chatbox.Placeholder("gpuinfo."+prefix+strconv.Itoa(info.Index)+".vendor", info.Vendor)
	chatbox.Placeholder("gpuinfo."+prefix+strconv.Itoa(info.Index)+".usage", strconv.Itoa(info.Utilization)+"%")
	chatbox.Placeholder(
		"gpuinfo."+prefix+strconv.Itoa(info.Index)+".memory",
		strconv.Itoa(int(float64(info.MemoryUsedMB)/float64(info.MemoryTotalMB)*100))+"%",
	)
	chatbox.Placeholder("gpuinfo."+prefix+strconv.Itoa(info.Index)+".memory.total", strconv.Itoa(info.MemoryTotalMB))
	chatbox.Placeholder("gpuinfo."+prefix+strconv.Itoa(info.Index)+".memory.used", strconv.Itoa(info.MemoryUsedMB))
}

func (m GpuInfoModule) triggerMeasure() {
	go func() {
		if m.container.providerAMD != nil {
			amd, err := m.container.providerAMD.Read()
			if err != nil {
				log.Print("%v", err)
			}
			m.container.usageAMD = amd
		}

		if m.container.providerNVIDIA != nil {
			nvidia, err := m.container.providerNVIDIA.Read()
			if err != nil {
				log.Print("%v", err)
			}
			m.container.usageNVIDIA = nvidia
		}
	}()
}
