package irccip

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// testServerPSK is the PSK for the server used in the following tests.
const testServerPSK = "0000"

func TestSendKeyCode(t *testing.T) {
	tests := []struct {
		name                 string
		key                  KeyCode
		clientPSK            string
		errorContains        string
		mockedResponseStatus int
		mockedResponse       []byte
	}{
		{
			name:                 "valid key",
			key:                  KeyLeft,
			clientPSK:            testServerPSK,
			errorContains:        "",
			mockedResponseStatus: http.StatusOK,
			mockedResponse: []byte(`<?xml version="1.0"?>
<s:Envelope
    xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"
    s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>
    <u:X_SendIRCCResponse xmlns:u="urn:schemas-sony-com:service:IRCC:1">
    </u:X_SendIRCCResponse>
  </s:Body>
</s:Envelope>`),
		},
		{
			name:                 "invalid auth",
			key:                  KeyLeft,
			clientPSK:            "0001",
			errorContains:        "HTTP error: 403 Forbidden",
			mockedResponseStatus: http.StatusForbidden,
			mockedResponse:       []byte(``),
		},
		{
			name:                 "invalid key",
			key:                  KeyLeft + "a",
			clientPSK:            testServerPSK,
			errorContains:        "SOAP fault: Cannot accept the IRCC Code (800)",
			mockedResponseStatus: http.StatusInternalServerError,
			mockedResponse: []byte(`<?xml version="1.0"?>
<s:Envelope
    xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"
    s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
  <s:Body>
    <s:Fault>
      <faultcode>s:Client</faultcode>
      <faultstring>UPnPError</faultstring>
      <detail>
        <UPnPError xmlns="urn:schemas-upnp-org:control-1-0">
          <errorCode>800</errorCode>
          <errorDescription>Cannot accept the IRCC Code</errorDescription>
        </UPnPError>
      </detail>
    </s:Fault>
  </s:Body>
</s:Envelope>`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, server := NewTestClientServer(tt.clientPSK, testServerPSK, tt.mockedResponseStatus, tt.mockedResponse)
			defer server.Close()

			err := client.SendKeyCode(tt.key)
			if err == nil {
				if tt.errorContains != "" {
					t.Fatalf("expected error to contain %q, but got no error", err.Error())
				}
			} else if !strings.Contains(err.Error(), tt.errorContains) {
				t.Fatalf("error does not contain %q: %v", tt.errorContains, err)
			}
		})

	}
}

// NewTestClientServer returns a client and a test server for the given mocks.
func NewTestClientServer(clientPSK, serverPSK string, mockedResponseStatus int, mockedResponse []byte) (client *Client, server *httptest.Server) {
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the given path has the irccPath suffix
		if !strings.HasSuffix(r.URL.Path, irccPath) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(""))
			return
		}

		// Check auth key
		reqPSK := r.Header.Get(authHeaderName)
		if reqPSK != serverPSK {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(""))
			return
		}

		// Return the mocked response
		w.WriteHeader(mockedResponseStatus)
		w.Write(mockedResponse)
	}))

	client = NewClient(server.URL, clientPSK)
	client.SetHTTPClient(server.Client())

	return client, server
}

func ExampleClient_SendKeyCode() {
	// Create a client to send commands to a remote display
	c := NewClient("http://192.168.1.12", "0000")

	// Send the 'Home' key code
	if err := c.SendKeyCode(KeyHome); err != nil {
		log.Fatalf("SendKeyCode error: %v", err)
	}
}
