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
)

var (
	// Root URL for all API requests to RedisGreen. Omit the trailing slash.
	ApiUrl = "https://dashboard.redisgreen.net"

	// User agent string for all API requests to RedisGreen.
	UserAgent = fmt.Sprintf("redisgreen-go/0.0.1 (%s; %s)", runtime.GOOS, runtime.GOARCH)
)

// A Client is a RedisGreen API client that acts as the user with the given API
// token.
type Client struct {
	Token string
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
