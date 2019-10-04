package irccip

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// KeyCode is the type used to specify an IRCC Code
type KeyCode string

// The default set of KeyCodes documented in the Sony Bravia IRCC Codes page:
// https://pro-bravia.sony.net/develop/integrate/ircc-ip/ircc-codes/index.html
const (
	KeyPower       KeyCode = "AAAAAQAAAAEAAAAVAw=="
	KeyInput       KeyCode = "AAAAAQAAAAEAAAAlAw=="
	KeySyncMenu    KeyCode = "AAAAAgAAABoAAABYAw=="
	KeyHdmi1       KeyCode = "AAAAAgAAABoAAABaAw=="
	KeyHdmi2       KeyCode = "AAAAAgAAABoAAABbAw=="
	KeyHdmi3       KeyCode = "AAAAAgAAABoAAABcAw=="
	KeyHdmi4       KeyCode = "AAAAAgAAABoAAABdAw=="
	KeyNum1        KeyCode = "AAAAAQAAAAEAAAAAAw=="
	KeyNum2        KeyCode = "AAAAAQAAAAEAAAABAw=="
	KeyNum3        KeyCode = "AAAAAQAAAAEAAAACAw=="
	KeyNum4        KeyCode = "AAAAAQAAAAEAAAADAw=="
	KeyNum5        KeyCode = "AAAAAQAAAAEAAAAEAw=="
	KeyNum6        KeyCode = "AAAAAQAAAAEAAAAFAw=="
	KeyNum7        KeyCode = "AAAAAQAAAAEAAAAGAw=="
	KeyNum8        KeyCode = "AAAAAQAAAAEAAAAHAw=="
	KeyNum9        KeyCode = "AAAAAQAAAAEAAAAIAw=="
	KeyNum0        KeyCode = "AAAAAQAAAAEAAAAJAw=="
	KeyDot         KeyCode = "AAAAAgAAAJcAAAAdAw=="
	KeyCC          KeyCode = "AAAAAgAAAJcAAAAoAw=="
	KeyRed         KeyCode = "AAAAAgAAAJcAAAAlAw=="
	KeyGreen       KeyCode = "AAAAAgAAAJcAAAAmAw=="
	KeyYellow      KeyCode = "AAAAAgAAAJcAAAAnAw=="
	KeyBlue        KeyCode = "AAAAAgAAAJcAAAAkAw=="
	KeyUp          KeyCode = "AAAAAQAAAAEAAAB0Aw=="
	KeyDown        KeyCode = "AAAAAQAAAAEAAAB1Aw=="
	KeyRight       KeyCode = "AAAAAQAAAAEAAAAzAw=="
	KeyLeft        KeyCode = "AAAAAQAAAAEAAAA0Aw=="
	KeyConfirm     KeyCode = "AAAAAQAAAAEAAABlAw=="
	KeyHelp        KeyCode = "AAAAAgAAAMQAAABNAw=="
	KeyDisplay     KeyCode = "AAAAAQAAAAEAAAA6Aw=="
	KeyOptions     KeyCode = "AAAAAgAAAJcAAAA2Aw=="
	KeyBack        KeyCode = "AAAAAgAAAJcAAAAjAw=="
	KeyHome        KeyCode = "AAAAAQAAAAEAAABgAw=="
	KeyVolumeUp    KeyCode = "AAAAAQAAAAEAAAASAw=="
	KeyVolumeDown  KeyCode = "AAAAAQAAAAEAAAATAw=="
	KeyMute        KeyCode = "AAAAAQAAAAEAAAAUAw=="
	KeyAudio       KeyCode = "AAAAAQAAAAEAAAAXAw=="
	KeyChannelUp   KeyCode = "AAAAAQAAAAEAAAAQAw=="
	KeyChannelDown KeyCode = "AAAAAQAAAAEAAAARAw=="
	KeyPlay        KeyCode = "AAAAAgAAAJcAAAAaAw=="
	KeyPause       KeyCode = "AAAAAgAAAJcAAAAZAw=="
	KeyStop        KeyCode = "AAAAAgAAAJcAAAAYAw=="
	KeyFlashPlus   KeyCode = "AAAAAgAAAJcAAAB4Aw=="
	KeyFlashMinus  KeyCode = "AAAAAgAAAJcAAAB5Aw=="
	KeyPrev        KeyCode = "AAAAAgAAAJcAAAA8Aw=="
	KeyNext        KeyCode = "AAAAAgAAAJcAAAA9Aw=="
)

const (
	authHeaderName        = "X-Auth-PSK"
	soapActionHeaderName  = "SOAPACTION"
	soapActionHeaderValue = `"urn:schemas-sony-com:service:IRCC:1#X_SendIRCC"`
)

// sendKeyCodeHeaders returns the headers for a SendKeyCode request
func (c *Client) sendKeyCodeHeaders() http.Header {
	return http.Header{
		authHeaderName:       []string{c.preSharedKey},
		soapActionHeaderName: []string{soapActionHeaderValue},
		"Content-Type":       []string{"text/xml; charset=UTF-8"},
	}
}

// keyCodeBody returns a buffer containing the SOAP request to send a key code.
// The given key can be one of the constants defined in this package, or a custom one.
func keyCodeBody(key KeyCode) *bytes.Buffer {
	return bytes.NewBufferString(`<?xml version="1.0" encoding="utf-8"?>
<s:Envelope
	xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"
	s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <s:Body>
        <u:X_SendIRCC xmlns:u="urn:schemas-sony-com:service:IRCC:1">
            <IRCCCode>` + string(key) + `</IRCCCode>
        </u:X_SendIRCC>
    </s:Body>
</s:Envelope>`)
}

// HACK: pull out the description - SOAP XML parsing is too painful to be worth it just for this...
var keyCodeErrorRegexp = regexp.MustCompile(`.*<errorCode>(.+)</errorCode>\n.*<errorDescription>(.+)</errorDescription>.*`)

// SendKeyCode sends the given key code using the given client.
// key may be a KeyCode constant defined in the package, or your own key code.
func (c *Client) SendKeyCode(key KeyCode) (err error) {
	req, err := http.NewRequest(http.MethodPost, c.url, keyCodeBody(key))
	if err != nil {
		return err
	}

	req.Header = c.sendKeyCodeHeaders()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	// Sony helpfully return a generic 500 HTTP status for most errors...
	// So we'll need to pull the specific error out of the SOAP fault in the response body.
	if resp.StatusCode == http.StatusInternalServerError {
		defer resp.Body.Close()
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		submatch := keyCodeErrorRegexp.FindSubmatch(respBody)
		if len(submatch) == 3 {
			return fmt.Errorf("SOAP fault: %v (%v)", string(submatch[2]), string(submatch[1]))
		}
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP error: %v", resp.Status)
	}

	return nil
}
