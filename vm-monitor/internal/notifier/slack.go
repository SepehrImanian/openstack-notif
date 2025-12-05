package notifier

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type SlackNotifier struct {
	Webhook string
	client  *http.Client
}

func NewSlackNotifier(url string) *SlackNotifier {
	return &SlackNotifier{
		Webhook: url,
		client:  &http.Client{},
	}
}

func (s *SlackNotifier) Send(text string) {
	if s.Webhook == "" {
		log.Println("[WARN] Slack webhook URL is empty; skipping send")
		return
	}

	payload, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		log.Println("[ERROR] Slack payload marshal failed:", err)
		return
	}

	resp, err := s.client.Post(s.Webhook, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("[ERROR] Slack send failed:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[ERROR] Slack responded with status %d: %s\n", resp.StatusCode, string(body))
	}
}
