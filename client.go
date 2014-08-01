// Package redisgreen implements a basic client for the RedisGreen API.
package redisgreen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
)

var (
	// Root URL for all API requests to RedisGreen. Omit the trailing slash.
	ApiUrl = "https://dashboard.redisgreen.net"

	// User agent string for all API requests to RedisGreen.
	UserAgent = fmt.Sprintf("redisgreen-go/0.0.1 (%s; %s)", runtime.GOOS, runtime.GOARCH)
)

// A Client is a RedisGreen API client that acts as a single user.
type Client struct {
	Token string
}

// ListServers returns a slice of all RedisGreen servers owned by this client's
// account.
func (c *Client) ListServers() ([]Server, error) {
	servers := []Server{}
	err := c.get("/servers", &servers)
	return servers, err
}

// CreateServer builds a new RedisGreen server from scratch and returns the
// created server information. Keep in mind that some servers may take up to 10
// minutes to become available.
//
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

// GetServer returns the Server record for the RedisGreen server with the given
// ID.
func (c *Client) GetServer(id string) (Server, error) {
	server := Server{}
	err := c.get(fmt.Sprintf("/servers/%s", id), &server)
	return server, err
}

// DeleteServer tears down a running RedisGreen server. It does not shut down
// any of the server's slaves.
func (c *Client) DeleteServer(id string) error {
	return c.del(fmt.Sprintf("/servers/%s", id))
}

func (c *Client) get(path string, val interface{}) error {
	body, err := c.do("GET", path, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, val)
}

func (c *Client) post(path string, payload, val interface{}) error {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	body, err := c.do("POST", path, bytes.NewReader(payloadJson))
	if err != nil {
		return err
	}
	return json.Unmarshal(body, val)
}

func (c *Client) del(path string) error {
	_, err := c.do("DELETE", path, nil)
	return err
}

func (c *Client) do(method, path string, reqBody io.Reader) (respBody []byte, err error) {
	req, err := http.NewRequest(method, ApiUrl+path, reqBody)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Accept", "application/vnd.redisgreen.1+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("X-API-Token", c.Token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err = readBody(resp.Body)
	if resp.StatusCode/100 == 2 || err != nil {
		return respBody, err
	}

	// Convert JSON-formatted error to Go error
	respErr := &Error{}
	err = json.Unmarshal(respBody, respErr)
	if err == nil {
		err = respErr
	}

	return nil, err
}

func readBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return ioutil.ReadAll(body)
}
