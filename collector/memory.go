package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
	"github.com/Konk32/go-pulse/snapshot"
)

func CollectMemory(ctx context.Context) (snapshot.Memory, error) {
	v, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return snapshot.Memory{}, fmt.Errorf("collect memory: %w", err)
	}
	return snapshot.Memory{
		TotalMB:      v.Total / 1024 / 1024,
		UsedMB:       v.Used / 1024 / 1024,
		UsagePercent: v.UsedPercent,
	}, nil
}