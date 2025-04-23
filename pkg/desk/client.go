package desk

import (
	"github.com/ready4god2513/desksdkgo/client"
)

// Client wraps the SDK client
type Client struct {
	*client.Client
}

// NewClient returns a new Teamwork Desk API client
func NewClient(baseURL, apiKey string) *Client {
	c := client.NewClient(baseURL, client.WithAPIKey(apiKey))
	return &Client{Client: c}
}
