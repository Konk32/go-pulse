package snapshot

import "time"

type CPU struct {
	UsagePercent float64 `json:"usage_percent"`
}

type Memory struct {
	TotalMB     uint64  `json:"total_mb"`
	UsedMB      uint64  `json:"used_mb"`
	UsagePercent float64 `json:"usage_percent"`
}

type Disk struct {
	Path         string  `json:"path"`
	TotalGB      uint64  `json:"total_gb"`
	UsedGB       uint64  `json:"used_gb"`
	UsagePercent float64 `json:"usage_percent"`
}

type Snapshot struct {
	AgentID   string    `json:"agent_id"`
	Timestamp time.Time `json:"timestamp"`
	CPU       CPU       `json:"cpu"`
	Memory    Memory    `json:"memory"`
	Disk      Disk      `json:"disk"`
}