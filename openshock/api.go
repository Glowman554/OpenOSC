package openshock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type OpenShockApi struct {
	token string
}

func NewOpenShockApi(token string) *OpenShockApi {
	return &OpenShockApi{
		token: token,
	}
}

func (o *OpenShockApi) LoadShockers() ([]ShockerEntry, error) {
	req, err := http.NewRequest("GET", "https://api.openshock.app/1/shockers/own", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Open-Shock-Token", o.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", string(body))
	}

	message := LoadShockersMessage{}
	err = json.Unmarshal(body, &message)
	if err != nil {
		return nil, err
	}

	shockers := []ShockerEntry{}
	for _, i := range message.Data {
		shockers = append(shockers, i.Shockers...)
	}

	return shockers, nil
}

func (o *OpenShockApi) LoadShockersShared() (map[string]ShockerEntry, error) {
	req, err := http.NewRequest("GET", "https://api.openshock.app/1/shockers/shared", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Open-Shock-Token", o.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", string(body))
	}

	message := LoadShockersSharedMessage{}
	err = json.Unmarshal(body, &message)
	if err != nil {
		return nil, err
	}

	shockers := make(map[string]ShockerEntry)
	for _, i := range message.Data {
		for _, j := range i.Devices {
			for _, k := range j.Shockers {
				shockers[j.Name+":"+k.Name] = k
			}
		}
	}

	return shockers, nil
}

func (o *OpenShockApi) SendCommand(intensity int, duration int, command ShockType, shockerIDs []string) error {
	commands := ShockerControlMessage{
		Shocks:     []ShockControl{},
		CustomName: "OpenOSC",
	}

	for _, id := range shockerIDs {
		command := ShockControl{
			Id:        id,
			Type:      string(command),
			Intensity: intensity,
			Duration:  duration,
			Exclusive: true,
		}
		commands.Shocks = append(commands.Shocks, command)
	}

	jsonData, err := json.Marshal(commands)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.openshock.app/2/shockers/control", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Open-Shock-Token", o.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		log.Printf("Status Code: %d", resp.StatusCode)
		log.Printf("Response Body: %s", string(body))
	}

	return nil
}
