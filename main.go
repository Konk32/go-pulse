package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Konk32/go-pulse/collector"
	"github.com/Konk32/go-pulse/config"
	"github.com/Konk32/go-pulse/sender"
	"github.com/Konk32/go-pulse/snapshot"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load("config.yaml")
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger.Info("vigil-pulse starting",
		"agent_id", cfg.AgentID,
		"server_url", cfg.ServerURL,
		"interval", cfg.Interval,
	)

	s := sender.New(cfg.ServerURL)

	// Signal handling — SIGINT (Ctrl+C) and SIGTERM (systemd stop)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	// Run once immediately before the first tick
	collect(ctx, logger, cfg, s)

	for {
		select {
		case <-ticker.C:
			collect(ctx, logger, cfg, s)
		case <-ctx.Done():
			logger.Info("shutdown signal received, exiting cleanly")
			return
		}
	}
}

func collect(ctx context.Context, logger *slog.Logger, cfg *config.Config, s *sender.Sender) {
	// Each round gets its own deadline — 5s to collect + send
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var (
		wg   sync.WaitGroup
		mu   sync.Mutex
		snap = snapshot.Snapshot{
			AgentID:   cfg.AgentID,
			Timestamp: time.Now().UTC(),
		}
		collectionErr error
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		cpu, err := collector.CollectCPU(ctx)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			collectionErr = err
			return
		}
		snap.CPU = cpu
	}()

	go func() {
		defer wg.Done()
		mem, err := collector.CollectMemory(ctx)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			collectionErr = err
			return
		}
		snap.Memory = mem
	}()

	go func() {
		defer wg.Done()
		disk, err := collector.CollectDisk(ctx, "C:\\")
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			collectionErr = err
			return
		}
		snap.Disk = disk
	}()

	wg.Wait()

	if collectionErr != nil {
		logger.Error("collection failed", "error", collectionErr)
		return
	}

	logger.Info("snapshot collected",
		"cpu_percent", snap.CPU.UsagePercent,
		"mem_percent", snap.Memory.UsagePercent,
		"disk_percent", snap.Disk.UsagePercent,
	)

	if err := s.Send(ctx, snap); err != nil {
		logger.Warn("failed to send snapshot, will retry next tick", "error", err)
		return
	}

	logger.Info("snapshot sent successfully")
}