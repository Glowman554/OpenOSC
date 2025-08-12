package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Chatbox     []string `json:"chatbox"`
	SendIP      string   `json:"sendIP"`
	SendPort    int      `json:"sendPort"`
	ReceivePort int      `json:"receivePort"`
}

func LoadConfig(filename string) (*Config, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		defaultConfig := Config{
			Chatbox: []string{
				"🎵 {media.title} - {media.artist}",
				"{media.progress}",
				"CPU: {sysinfo.cpu}, Memory: {sysinfo.memory}",
				"Controlling Player: {media.control.player}",
				"Media Player: {media.player}",
			},
			SendIP:      "127.0.0.1",
			SendPort:    9000,
			ReceivePort: 9001,
		}

		data, err := json.MarshalIndent(defaultConfig, "", "  ")
		if err != nil {
			log.Printf("failed to marshal default config: %v", err)
			return nil, err
		}
		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			fmt.Printf("failed to create config file: %v", err)
			return nil, err
		}
		return &defaultConfig, nil
	} else if err != nil {
		log.Printf("failed to stat config file: %v", err)
		return nil, err
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("failed to read config file: %v", err)
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Printf("failed to unmarshal config: %v", err)
		return nil, err
	}

	return &config, nil
}
