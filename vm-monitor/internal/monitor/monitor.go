package monitor

import (
	"fmt"
	"log"
	"time"

	"vm-monitor/internal/config"
	"vm-monitor/internal/notifier"
	"vm-monitor/internal/openstack"
)

type Monitor struct {
	osClient   *openstack.OpenStackClient
	notify     notifier.Notifier
	cfg        config.Config
	lastAlerts map[string]string
}

func NewMonitor(osClient *openstack.OpenStackClient, notify notifier.Notifier, cfg config.Config) *Monitor {
	return &Monitor{
		osClient:   osClient,
		notify:     notify,
		cfg:        cfg,
		lastAlerts: make(map[string]string),
	}
}

func (m *Monitor) Run() {
	interval := time.Duration(m.cfg.IntervalSec)
	if interval <= 0 {
		interval = 20
	}

	log.Printf("[INFO] Starting VM monitor; interval=%ds\n", interval)

	ticker := time.NewTicker(interval * time.Second)
	defer ticker.Stop()

	for {
		m.pollOnce()
		<-ticker.C
	}
}

func (m *Monitor) pollOnce() {
	log.Println("[INFO] Checking OpenStack instances...")

	servers := m.osClient.ListServers()
	if servers == nil {
		log.Println("[WARN] OpenStack returned no servers (nil slice or decode error)")
		return
	}

	log.Printf("[INFO] Retrieved %d servers from OpenStack\n", len(servers))
	for _, s := range servers {
		m.checkServer(&s)
	}
}

func (m *Monitor) checkServer(s *openstack.Server) {
	alertKey := ""
	var msg string

	switch {
	case s.Status == "ERROR":
		alertKey = "ERROR"
		msg = fmt.Sprintf(
			":x: *VM %s* (%s) is in *ERROR* state\nvm_state=%s, task_state=%s, power_state=%d, progress=%d%%",
			s.Name, s.ID, s.VMState, s.TaskState, s.PowerState, s.Progress,
		)

	case s.Status == "BUILD" && s.Progress < 100:
		alertKey = "BUILD_STUCK"
		msg = fmt.Sprintf(
			":warning: *VM %s* (%s) appears stuck in BUILD — progress=%d%%",
			s.Name, s.ID, s.Progress,
		)

	case s.Status == "ACTIVE" && s.PowerState != 1:
		alertKey = "ACTIVE_BAD_POWER"
		msg = fmt.Sprintf(
			":warning: *VM %s* (%s) is ACTIVE but power_state=%d",
			s.Name, s.ID, s.PowerState,
		)

	default:
		if prev, ok := m.lastAlerts[s.ID]; ok && prev != "" {
			log.Printf("[INFO] VM %s (%s) recovered from previous state: %s\n", s.Name, s.ID, prev)
		}
		m.lastAlerts[s.ID] = ""
		return
	}

	if prev, ok := m.lastAlerts[s.ID]; ok && prev == alertKey {
		return
	}

	m.lastAlerts[s.ID] = alertKey
	log.Printf("[INFO] Alerting for VM %s (%s): %s\n", s.Name, s.ID, alertKey)
	m.notify.Send(msg)
}
