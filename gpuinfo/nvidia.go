package gpuinfo

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type NvidiaProvider struct {
}

func NewNvidiaProvider() *NvidiaProvider {
	return &NvidiaProvider{}
}

func (p *NvidiaProvider) Read() ([]GPUUsage, error) {
	cmd := exec.Command("nvidia-smi",
		"--query-gpu=name,utilization.gpu,memory.used,memory.total",
		"--format=csv,noheader,nounits",
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute nvidia-smi: %w", err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	var gpus []GPUUsage
	index := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 4 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		util, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		memUsed, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
		memTotal, _ := strconv.Atoi(strings.TrimSpace(parts[3]))

		gpus = append(gpus, GPUUsage{
			Index:         index,
			Name:          name,
			Vendor:        "NVIDIA",
			Utilization:   util,
			MemoryUsedMB:  memUsed,
			MemoryTotalMB: memTotal,
		})
		index++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan output: %w", err)
	}

	return gpus, nil
}

func CanUseNvidiaProvider() bool {
	if _, err := exec.LookPath("nvidia-smi"); err == nil {
		return true
	}
	return false
}
