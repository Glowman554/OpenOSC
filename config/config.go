package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type LeashConfig struct {
	WalkDeadzone       float64 `json:"walkDeadzone"`
	RunDeadzone        float64 `json:"runDeadzone"`
	StrengthMultiplier float64 `json:"strengthMultiplier"`
	UpDownDeadzone     float64 `json:"upDownDeadzone"`
	UpDownCompensation float64 `json:"upDownCompensation"`
	TurningDeadzone    float64 `json:"turningDeadzone"`
	TurningMultiplier  float64 `json:"turningMultiplier"`
	TurningGoal        float64 `json:"turningGoal"`
	LeashDirection     string  `json:"leashDirection"`
}

type Config struct {
	Chatbox        []string    `json:"chatbox"`
	SendIP         string      `json:"sendIP"`
	SendPort       int         `json:"sendPort"`
	ReceivePort    int         `json:"receivePort"`
	OpenShockToken string      `json:"openShockToken"`
	ActiveModules  []string    `json:"activeModules"`
	LeashConfig    LeashConfig `json:"leashConfig"`
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
			SendIP:         "127.0.0.1",
			SendPort:       9000,
			ReceivePort:    9001,
			OpenShockToken: "",
			ActiveModules: []string{
				"media_chatbox",
				"media_control",
				"leash",
				"sysinfo",
			},
			LeashConfig: LeashConfig{
				WalkDeadzone:       0.15,
				RunDeadzone:        0.70,
				StrengthMultiplier: 1.2,
				UpDownDeadzone:     0.5,
				UpDownCompensation: 0.5,
				TurningDeadzone:    0.15,
				TurningMultiplier:  0.8,
				TurningGoal:        90.0,
				LeashDirection:     "north",
			},
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
