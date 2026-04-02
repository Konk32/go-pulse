package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/Konk32/go-pulse/snapshot"
)

func CollectCPU(ctx context.Context) (snapshot.CPU, error) {
	percents, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		return snapshot.CPU{}, fmt.Errorf("collect cpu: %w", err)
	}
	return snapshot.CPU{UsagePercent: percents[0]}, nil
}