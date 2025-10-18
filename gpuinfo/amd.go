package gpuinfo

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
)

type AMDProvider struct {
}

func NewAMDProvider() *AMDProvider {
	return &AMDProvider{}
}

func (p *AMDProvider) Read() ([]GPUUsage, error) {
	cmd := exec.Command("rocm-smi", "--showuse", "--showmeminfo", "vram", "--showproductname", "--json")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute rocm-smi: %w", err)
	}

	var data map[string]map[string]string
	if err := json.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	bracketRe := regexp.MustCompile(`\[(.*?)\]`)
	var gpus []GPUUsage

	for card, vals := range data {
		var gpu GPUUsage
		fmt.Sscanf(card, "card%d", &gpu.Index)

		rawName := vals["Card Series"]
		if match := bracketRe.FindStringSubmatch(rawName); len(match) == 2 {
			gpu.Name = match[1]
		} else if rawName != "" {
			gpu.Name = rawName
		} else {
			gpu.Name = fmt.Sprintf("AMD GPU %d", gpu.Index)
		}

		gpu.Vendor = vals["Card Vendor"]
		if gpu.Vendor == "" {
			gpu.Vendor = "AMD"
		}

		fmt.Sscanf(vals["GPU use (%)"], "%d", &gpu.Utilization)

		var usedB, totalB float64
		fmt.Sscanf(vals["VRAM Total Used Memory (B)"], "%f", &usedB)
		fmt.Sscanf(vals["VRAM Total Memory (B)"], "%f", &totalB)
		gpu.MemoryUsedMB = int(usedB / 1024 / 1024)
		gpu.MemoryTotalMB = int(totalB / 1024 / 1024)

		gpus = append(gpus, gpu)
	}

	return gpus, nil
}

func CanUseAMDProvider() bool {
	if _, err := exec.LookPath("rocm-smi"); err == nil {
		return true
	}
	return false
}
