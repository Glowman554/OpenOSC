package gpuinfo

type GPUUsage struct {
	Index         int
	Name          string
	Utilization   int
	MemoryUsedMB  int
	MemoryTotalMB int
	Vendor        string
}

type Provider interface {
	Read() ([]GPUUsage, error)
}
