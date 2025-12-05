package config

import "os"

type Config struct {
	OpenStackURL string
	ProjectID    string
	Token        string
	SlackWebhook string
	IntervalSec  int
}

func LoadConfig() Config {
	return Config{
		OpenStackURL: os.Getenv("OS_URL"),
		ProjectID:    os.Getenv("OS_PROJECT_ID"),
		Token:        os.Getenv("OS_TOKEN"),
		SlackWebhook: os.Getenv("SLACK_WEBHOOK"),
		IntervalSec:  20,
	}
}
