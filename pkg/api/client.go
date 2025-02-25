package s4t

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MIKE9708/s4t-sdk-go/pkg/read_conf"
	"net/http"
	"time"
)

const (
	IP        = "127.0.0.1"
	PORT      = "8812"
	AUTH_PORT = "5000"
)

type Client struct {
	HTTPClient *http.Client
	AuthToken  string
	Port       string
	AuthPort   string
	Endpoint   string
	Timeout    time.Duration
}

func GetClientConnection(authReqData read_config.AuthRequest_1) (*Client, error) {
	client := NewClient("http://" + IP)
	client.Port = PORT
	client.AuthPort = AUTH_PORT
	token, err := client.Authenticate(client, &authReqData)
	if err != nil {
		return nil, fmt.Errorf("Error authenticating: %v", err)
	}
	client.AuthToken = token
	return client, nil
}

func (c *Client) Authenticate(client *Client, auth_req *read_config.AuthRequest_1) (string, error) {
	var authetication_data struct {
		Auth read_config.AuthRequest_1 `json:"auth"`
	}
	authetication_data.Auth = *auth_req

	jsonBody, err := json.Marshal(&authetication_data)

	if err != nil {
		return "", fmt.Errorf("Error marshaling to JSON: %d\n", err)

	}
	req, err := http.NewRequest("POST", client.Endpoint+":"+client.AuthPort+"/v3/auth/tokens", bytes.NewBuffer(jsonBody))

	if err != nil {
		return "", fmt.Errorf("authentication request failed: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.HTTPClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("Request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("Authentication failed with status code: %d", resp.StatusCode)
	}

	token := resp.Header.Get("X-Subject-Token")

	if token == "" {
		return "", fmt.Errorf("No token found in the response")
	}

	return token, nil
}

type ClientOption func(*Client)

func NewClient(endpoint string, options ...ClientOption) *Client {
	c := &Client{
		HTTPClient: &http.Client{},
		Endpoint:   endpoint,
		Timeout:    time.Second * 90,
	}

	for _, option := range options {
		option(c)
	}

	c.HTTPClient.Timeout = c.Timeout

	return c
}
