package collector

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/disk"
	"github.com/Konk32/go-pulse/snapshot"
)

func CollectDisk(ctx context.Context, path string) (snapshot.Disk, error) {
	u, err := disk.UsageWithContext(ctx, path)
	if err != nil {
		return snapshot.Disk{}, fmt.Errorf("collect disk: %w", err)
	}
	return snapshot.Disk{
		Path:         path,
		TotalGB:      u.Total / 1024 / 1024 / 1024,
		UsedGB:       u.Used / 1024 / 1024 / 1024,
		UsagePercent: u.UsedPercent,
	}, nil
}