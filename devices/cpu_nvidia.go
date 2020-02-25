package devices

import (
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"time"
)

func init() {
	if nvml.Init() == nil {
		shutdownFuncs = append(shutdownFuncs, func() { nvml.Shutdown() })
		deviceCounts = append(deviceCounts, func(b bool) (int, error) {
			r, e := nvml.GetDeviceCount()
			return int(r), e
		})
		devicePercents = append(devicePercents, func(t time.Duration, b bool) ([]float64, error) {
			return nil, nil
		})
	}
}
