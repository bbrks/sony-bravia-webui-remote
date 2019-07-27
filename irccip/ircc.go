// Package irccip implements Sony's InfraRed Compatible Control over Internet Protocol (IRCC-IP) for Bravia displays
package irccip

import (
	"net/http"
)

const irccPath = "/sony/ircc"

// NewClient returns a new IRCC-IP client.
// url is a HTTP url to the IP address of your Sony Bravia display.
// preSharedKey is the configured authentication code to control the display.
// httpClient may be set if you require a non-default HTTP client.
func NewClient(clientAddr string, preSharedKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		url:          clientAddr + irccPath,
		preSharedKey: preSharedKey,
		httpClient:   httpClient,
	}
}

// Client can send IRCC-IP commands to a compatible Sony Bravia display
type Client struct {
	url          string
	preSharedKey string
	httpClient   *http.Client
}
