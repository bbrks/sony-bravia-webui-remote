package irccip

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testPSK = "0000"

func TestSendKeyCode(t *testing.T) {
	tests := []struct {
		name                 string
		key                  KeyCode
		psk                  string
		errorContains        string
		mockedResponseStatus int
		mockedResponse       []byte
	}{
		{
			name:                 "valid key",
			key:                  KeyLeft,
			psk:                  testPSK,
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
			psk:                  "0001",
			errorContains:        "HTTP error: 403 Forbidden",
			mockedResponseStatus: http.StatusForbidden,
			mockedResponse:       []byte(``),
		},
		{
			name:                 "invalid key",
			key:                  KeyLeft + "a",
			psk:                  testPSK,
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

			backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				reqPSK := r.Header.Get(authHeaderName)
				if reqPSK != testPSK {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte(""))
					return
				}

				w.WriteHeader(tt.mockedResponseStatus)
				w.Write(tt.mockedResponse)
			}))
			c := NewClient(backend.URL, testPSK)
			c.SetHTTPClient(backend.Client())

			err := c.SendKeyCode(tt.key)
			if err == nil {
				if tt.errorContains != "" {
					t.Fatalf("expected error to contain %q, but got no error", err.Error())
				}
			} else if !strings.Contains(err.Error(), tt.errorContains) {
				t.Fatalf("error does not contain %q: %v", tt.errorContains, err)
			}

			backend.Close()
		})

	}
}

func ExampleSendKeyCode() {
	// Create a client to send commands to a remote display
	c := NewClient("http://192.168.1.12", "0000")

	// Send the 'Home' key code
	if err := c.SendKeyCode(KeyHome); err != nil {
		log.Fatalf("SendKeyCode error: %v", err)
	}
}
