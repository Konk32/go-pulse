package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Konk32/go-pulse/snapshot"
)

type Sender struct {
	serverURL  string
	httpClient *http.Client
}

func New(serverURL string) *Sender {
	return &Sender{
		serverURL: serverURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *Sender) Send(ctx context.Context, snap snapshot.Snapshot) error {
	body, err := json.Marshal(snap)
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.serverURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post snapshot: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	return nil
}