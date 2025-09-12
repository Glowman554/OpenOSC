package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Glowman554/OpenOSC/config"
	"github.com/Glowman554/OpenOSC/oscmod"
	"github.com/Glowman554/OpenOSC/oscmod/chatbox"
	"github.com/Glowman554/OpenOSC/oscmod/modules"
	"github.com/hypebeast/go-osc/osc"
	"github.com/mitchellh/go-ps"
)

func isVRChatRunning() bool {
	processes, err := ps.Processes()
	if err != nil {
		return false
	}

	for _, p := range processes {
		if strings.Contains(strings.ToLower(p.Executable()), "vrchat") {
			return true
		}
	}
	return false
}

func main() {
	configPath := flag.String("config", "config.json", "Path to the config file")
	flag.Parse()

	config, err := config.LoadConfig(*configPath)

	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
		panic(err)
	}

	// for true {
	// 	log.Println("Waiting for VRChat")

	// 	if isVRChatRunning() {
	// 		break
	// 	}

	// 	time.Sleep(2 * time.Second)
	// }

	modules := []oscmod.OSCModule{
		modules.NewMediaChatBoxModule(),
		modules.NewMediaControlModule(),
		modules.NewSysInfoModule(),
		modules.NewOpenShockModule(config.OpenShockConfig),
		modules.NewLeashModule(config.LeashConfig),
	}
	activeModules := []oscmod.OSCModule{}

	log.Println("Starting...")

	client := osc.NewClient(config.SendIP, config.SendPort)

	dispatcher := osc.NewStandardDispatcher()
	server := &osc.Server{Addr: fmt.Sprintf("0.0.0.0:%d", config.ReceivePort), Dispatcher: dispatcher}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
			panic(err)
		}
	}()

	for _, module := range modules {
		activate := false
		for _, activeModuleId := range config.ActiveModules {
			if module.Id() == activeModuleId {
				activate = true
			}
		}

		if activate {
			err := module.Init(client, dispatcher)
			if err != nil {
				log.Fatalf("Failed to initialize module: %s (%v)", module.Name(), err)
			} else {
				log.Printf("Initialized %s", module.Name())
				activeModules = append(activeModules, module)
			}
		}
	}

	chatbox := chatbox.NewChatBoxBuilder()
	for _, i := range config.Chatbox {
		chatbox.AddLine(i)
	}

	for true {
		chatbox.BeginTick()

		for _, module := range activeModules {
			err := module.Tick(client, chatbox)
			if err != nil {
				log.Printf("Failed to tick module: %s (%v)", module.Name(), err)
			}
		}

		chatbox.EndTick(client)

		time.Sleep(2 * time.Second)
	}
}
