package redisgreen

import (
	"fmt"
)

type Monitor struct {
	// RedisMonitor id
	Id string `json:"id"`
	// User-visible name
	Name string `json:"name"`
	// Full Redis URL to be monitored
	Url string `json:"url"`
}

func (c *Client) ListMonitors() ([]Monitor, error) {
	monitors := []Monitor{}
	err := c.get("/monitors", &monitors)
	return monitors, err
}

func (c *Client) CreateMonitor(name, url string) (Monitor, error) {
	payload := map[string]string{
		"name": name,
		"url":  url,
	}

	monitor := Monitor{}
	err := c.post("/monitors", payload, &monitor)
	return monitor, err
}

func (c *Client) GetMonitor(id string) (Monitor, error) {
	monitor := Monitor{}
	err := c.get(fmt.Sprintf("/monitors/%s", id), &monitor)
	return monitor, err
}

func (c *Client) DeleteMonitor(id string) error {
	return c.del(fmt.Sprintf("/monitors/%s", id))
}
