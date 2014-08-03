package redisgreen

import (
	"fmt"
	"strconv"
)

// A Server is a single RedisGreen server.
type Server struct {
	// RedisGreen id
	Id string `json:"id"`

	// User-visible name of the server
	Name string `json:"name"`

	// Full Redis connection URL, including password, host, and port
	URL string `json:"url"`

	// Ids for all RedisGreen slaves of this server
	Slaves []string `json:"slaves"`
}

func (c *Client) ListServers() ([]Server, error) {
	servers := []Server{}
	err := c.get("/servers", &servers)
	return servers, err
}

// Allowed plans: minidev, dev, basic, starter, plus, large, xlarge, 2xlarge,
// 4xlarge, 15xlarge, 30xlarge, 60xlarge
//
// Allowed regions: us-east-1 (default), us-west-1, us-west-2, eu-west-1,
// ap-northeast-1, ap-southeast-1, ap-southeast-2, sa-east-1
func (c *Client) CreateServer(name, plan, region string, slaves int) (Server, error) {
	payload := map[string]string{
		"name":        name,
		"plan":        plan,
		"region":      region,
		"slave_count": strconv.Itoa(slaves),
	}

	server := Server{}
	err := c.post("/servers", payload, &server)
	return server, err
}

func (c *Client) GetServer(id string) (Server, error) {
	server := Server{}
	err := c.get(fmt.Sprintf("/servers/%s", id), &server)
	return server, err
}

func (c *Client) DeleteServer(id string) error {
	return c.del(fmt.Sprintf("/servers/%s", id))
}
