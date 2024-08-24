package s4t


import (
	"net/http"
	"time"
	"fmt"
	"s4t-sdk-module/pkg/read_conf"
)

type Client struct {
	HTTPClient *http.Client
	AuthToken string
	Port string
	Endpoint string
	Timeout time.Duration
}


type ClientOption func( *Client )


func NewClient(endpoint string, options ...ClientOption) *Client {
	c := &Client{
		HTTPClient: &http.Client{},
		Endpoint: endpoint,
		Timeout: time.Second * 30,
	}

	for _, option := range options {
		option(c)
	}

	c.HTTPClient.Timeout = c.Timeout

	return c
}


func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.Timeout = timeout
	}
}

func GetClientConnection() (*Client, error) {
	auth_req, err := read_config.ReadConfiguration()
	
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %v", err)
	}

	client := NewClient("http://" + auth_req.S4tAuthData.Ip)
	client.AuthToken = auth_req.S4tAuthData.Token	
	client.Port = auth_req.S4tAuthData.Port

	return client, nil



}


