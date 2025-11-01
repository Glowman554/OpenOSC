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
	TurningEnabled     bool    `json:"turningEnabled"`
}

type OpenShockConfig struct {
	APIToken          string `json:"apiToken"`
	MaximumIntensity  int    `json:"maximumIntensity"`
	MaximumDurationMS int    `json:"maximumDurationMS"`
}

type OpenShockControlConfig struct {
	MaximumIntensity   int                 `json:"maximumIntensity"`
	MaximumDurationMS  int                 `json:"maximumDurationMS"`
	Mapping            map[string][]string `json:"mapping"`
	DurationParameter  string              `json:"durationParameter"`
	IntensityParameter string              `json:"intensityParameter"`
}

type Config struct {
	Chatbox                []string               `json:"chatbox"`
	ChatboxDebug           bool                   `json:"chatboxDebug"`
	SendIP                 string                 `json:"sendIP"`
	SendPort               int                    `json:"sendPort"`
	ReceivePort            int                    `json:"receivePort"`
	ActiveModules          []string               `json:"activeModules"`
	LeashConfig            LeashConfig            `json:"leashConfig"`
	OpenShockConfig        OpenShockConfig        `json:"openShockConfig"`
	OpenShockControlConfig OpenShockControlConfig `json:"openShockControlConfig"`
}

var defaultConfig = Config{
	Chatbox: []string{
		"🎵 {media.title} - {media.artist}",
		"{media.progress}",
		"CPU: {sysinfo.cpu}, Memory: {sysinfo.memory}",
		"{sysinfo.time.12h} / {sysinfo.time.24h}",
	},
	ChatboxDebug: false,
	SendIP:       "127.0.0.1",
	SendPort:     9000,
	ReceivePort:  9001,
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
		TurningEnabled:     false,
	},
	OpenShockConfig: OpenShockConfig{
		APIToken:          "",
		MaximumIntensity:  100,
		MaximumDurationMS: 30000,
	},
	OpenShockControlConfig: OpenShockControlConfig{
		MaximumIntensity:   100,
		MaximumDurationMS:  10000,
		Mapping:            map[string][]string{},
		DurationParameter:  "/avatar/parameters/Shock/Duration",
		IntensityParameter: "/avatar/parameters/Shock/Intensity",
	},
}

func LoadConfig(filename string) (*Config, error) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
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

	err = updateConfigFile(filename)
	if err != nil {
		log.Printf("failed to update config file: %v", err)
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

func updateConfigFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("failed to read config file: %v", err)
		return err
	}

	var config map[string]any
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Printf("failed to unmarshal config: %v", err)
		return err
	}

	changed := false

	// add leash config
	if _, ok := config["leashConfig"]; !ok {
		changed = true
		config["leashConfig"] = defaultConfig.LeashConfig
	}

	// add turningEnabled option
	leashConfig, _ := config["leashConfig"].(map[string]any)
	if _, ok := leashConfig["turningEnabled"]; !ok {
		changed = true
		leashConfig["turningEnabled"] = defaultConfig.LeashConfig.TurningEnabled
	}

	// update old openshock
	var openShockToken string = ""
	if val, ok := config["openShockToken"]; ok {
		changed = true
		openShockToken = val.(string)
		delete(config, "openShockToken")
	}

	if _, ok := config["openShockConfig"]; !ok {
		changed = true
		openShockConfig := defaultConfig.OpenShockConfig
		openShockConfig.APIToken = openShockToken
		config["openShockConfig"] = openShockConfig

	}

	// add openshock control config
	if _, ok := config["openShockControlConfig"]; !ok {
		changed = true
		config["openShockControlConfig"] = defaultConfig.OpenShockControlConfig
	}

	// add chatboxDebug
	if _, ok := config["chatboxDebug"]; !ok {
		changed = true
		config["chatboxDebug"] = defaultConfig.ChatboxDebug

	}

	if changed {
		log.Println("Writing updated config")
		configJson, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Printf("failed to marshal config: %v", err)
			return err
		}

		err = os.WriteFile(filename, configJson, 0)
		if err != nil {
			log.Printf("failed to write config file: %v", err)
			return err
		}
	}

	return nil
}
