package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Snapshot struct {
	AgentID   string    `json:"agent_id"`
	Timestamp time.Time `json:"timestamp"`
	CPU       struct {
		UsagePercent float64 `json:"usage_percent"`
	} `json:"cpu"`
	Memory struct {
		TotalMB      uint64  `json:"total_mb"`
		UsedMB       uint64  `json:"used_mb"`
		UsagePercent float64 `json:"usage_percent"`
	} `json:"memory"`
	Disk struct {
		Path         string  `json:"path"`
		TotalGB      uint64  `json:"total_gb"`
		UsedGB       uint64  `json:"used_gb"`
		UsagePercent float64 `json:"usage_percent"`
	} `json:"disk"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	http.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var snap Snapshot
		if err := json.NewDecoder(r.Body).Decode(&snap); err != nil {
			logger.Error("failed to decode snapshot", "error", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		logger.Info("heartbeat received",
			"agent_id", snap.AgentID,
			"timestamp", snap.Timestamp.Format(time.RFC3339),
			"cpu_percent", snap.CPU.UsagePercent,
			"mem_used_mb", snap.Memory.UsedMB,
			"mem_total_mb", snap.Memory.TotalMB,
			"disk_used_gb", snap.Disk.UsedGB,
			"disk_total_gb", snap.Disk.TotalGB,
		)

		w.WriteHeader(http.StatusOK)
	})

	logger.Info("mock server listening", "addr", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}