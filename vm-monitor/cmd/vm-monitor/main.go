package main

import (
	"log"
	"vm-monitor/internal/config"
	"vm-monitor/internal/monitor"
	"vm-monitor/internal/notifier"
	"vm-monitor/internal/openstack"
)

func main() {
	cfg := config.LoadConfig()
	if cfg.OpenStackURL == "" || cfg.ProjectID == "" || cfg.Token == "" || cfg.SlackWebhook == "" {
		log.Fatal("Missing required env vars: OS_URL, OS_PROJECT_ID, OS_TOKEN, SLACK_WEBHOOK")
	}

	osClient := openstack.NewOpenStackClient(cfg)
	slack := notifier.NewSlackNotifier(cfg.SlackWebhook)
	m := monitor.NewMonitor(osClient, slack, cfg)

	m.Run()
}
