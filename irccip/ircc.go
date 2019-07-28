// Package irccip implements Sony's InfraRed Compatible Control over Internet Protocol (IRCC-IP) for Bravia displays.
package irccip

import (
	"net/http"
)

const irccPath = "/sony/ircc"

// NewClient returns a new IRCC-IP client.
// - clientAddr is the HTTP URL of your Sony Bravia display. (e.g. http://192.168.1.12)
// - preSharedKey is your pre-configured authentication code to control the display.
func NewClient(clientAddr string, preSharedKey string) *Client {
	return &Client{
		url:          clientAddr + irccPath,
		preSharedKey: preSharedKey,
		httpClient:   http.DefaultClient,
	}
}

// SetHTTPClient allows you to specify a custom HTTP Client for communicating with the display.
func (c *Client) SetHTTPClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

// Client can send IRCC-IP commands to a compatible Sony Bravia display
type Client struct {
	url          string
	preSharedKey string
	httpClient   *http.Client
}
