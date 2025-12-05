package openstack

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"vm-monitor/internal/config"
)

type OpenStackClient struct {
	cfg config.Config
}

func NewOpenStackClient(cfg config.Config) *OpenStackClient {
	return &OpenStackClient{cfg: cfg}
}

func (c *OpenStackClient) ListServers() []Server {
	url := fmt.Sprintf("%s/servers/detail", c.cfg.OpenStackURL)
	log.Printf("[INFO] Calling OpenStack: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("[ERROR] building request:", err)
		return nil
	}
	req.Header.Set("X-Auth-Token", c.cfg.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("[ERROR] performing request:", err)
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("[ERROR] reading response body:", err)
		return nil
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("[ERROR] OpenStack returned %d %s\n", res.StatusCode, res.Status)
		log.Printf("[ERROR] Response body: %s\n", string(body))
		return nil
	}

	var data ServerListResponse
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("[ERROR] Decode error:", err)
		log.Println("[ERROR] Raw response body:", string(body))
		return nil
	}

	log.Printf("[INFO] Retrieved %d servers from OpenStack\n", len(data.Servers))
	return data.Servers
}