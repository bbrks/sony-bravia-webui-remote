package irccip

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type KeyCode string

// The default set of KeyCodes reported from a KDL-43W809C TV
const (
	KeyNum1                           KeyCode = "AAAAAQAAAAEAAAAAAw=="
	KeyNum2                           KeyCode = "AAAAAQAAAAEAAAABAw=="
	KeyNum3                           KeyCode = "AAAAAQAAAAEAAAACAw=="
	KeyNum4                           KeyCode = "AAAAAQAAAAEAAAADAw=="
	KeyNum5                           KeyCode = "AAAAAQAAAAEAAAAEAw=="
	KeyNum6                           KeyCode = "AAAAAQAAAAEAAAAFAw=="
	KeyNum7                           KeyCode = "AAAAAQAAAAEAAAAGAw=="
	KeyNum8                           KeyCode = "AAAAAQAAAAEAAAAHAw=="
	KeyNum9                           KeyCode = "AAAAAQAAAAEAAAAIAw=="
	KeyNum0                           KeyCode = "AAAAAQAAAAEAAAAJAw=="
	KeyNum11                          KeyCode = "AAAAAQAAAAEAAAAKAw=="
	KeyNum12                          KeyCode = "AAAAAQAAAAEAAAALAw=="
	KeyEnter                          KeyCode = "AAAAAQAAAAEAAAALAw=="
	KeyGGuide                         KeyCode = "AAAAAQAAAAEAAAAOAw=="
	KeyChannelUp                      KeyCode = "AAAAAQAAAAEAAAAQAw=="
	KeyChannelDown                    KeyCode = "AAAAAQAAAAEAAAARAw=="
	KeyVolumeUp                       KeyCode = "AAAAAQAAAAEAAAASAw=="
	KeyVolumeDown                     KeyCode = "AAAAAQAAAAEAAAATAw=="
	KeyMute                           KeyCode = "AAAAAQAAAAEAAAAUAw=="
	KeyTvPower                        KeyCode = "AAAAAQAAAAEAAAAVAw=="
	KeyAudio                          KeyCode = "AAAAAQAAAAEAAAAXAw=="
	KeyMediaAudioTrack                KeyCode = "AAAAAQAAAAEAAAAXAw=="
	KeyTv                             KeyCode = "AAAAAQAAAAEAAAAkAw=="
	KeyInput                          KeyCode = "AAAAAQAAAAEAAAAlAw=="
	KeyTvInput                        KeyCode = "AAAAAQAAAAEAAAAlAw=="
	KeyTvAntennaCable                 KeyCode = "AAAAAQAAAAEAAAAqAw=="
	KeyWakeUp                         KeyCode = "AAAAAQAAAAEAAAAuAw=="
	KeyPowerOff                       KeyCode = "AAAAAQAAAAEAAAAvAw=="
	KeySleep                          KeyCode = "AAAAAQAAAAEAAAAvAw=="
	KeyRight                          KeyCode = "AAAAAQAAAAEAAAAzAw=="
	KeyLeft                           KeyCode = "AAAAAQAAAAEAAAA0Aw=="
	KeySleepTimer                     KeyCode = "AAAAAQAAAAEAAAA2Aw=="
	KeyAnalog2                        KeyCode = "AAAAAQAAAAEAAAA4Aw=="
	KeyTvAnalog                       KeyCode = "AAAAAQAAAAEAAAA4Aw=="
	KeyDisplay                        KeyCode = "AAAAAQAAAAEAAAA6Aw=="
	KeyJump                           KeyCode = "AAAAAQAAAAEAAAA7Aw=="
	KeyPicOff                         KeyCode = "AAAAAQAAAAEAAAA+Aw=="
	KeyPictureOff                     KeyCode = "AAAAAQAAAAEAAAA+Aw=="
	KeyTeletext                       KeyCode = "AAAAAQAAAAEAAAA/Aw=="
	KeyVideo1                         KeyCode = "AAAAAQAAAAEAAABAAw=="
	KeyVideo2                         KeyCode = "AAAAAQAAAAEAAABBAw=="
	KeyAnalogRgb1                     KeyCode = "AAAAAQAAAAEAAABDAw=="
	KeyHome                           KeyCode = "AAAAAQAAAAEAAABgAw=="
	KeyExit                           KeyCode = "AAAAAQAAAAEAAABjAw=="
	KeyPictureMode                    KeyCode = "AAAAAQAAAAEAAABkAw=="
	KeyConfirm                        KeyCode = "AAAAAQAAAAEAAABlAw=="
	KeyUp                             KeyCode = "AAAAAQAAAAEAAAB0Aw=="
	KeyDown                           KeyCode = "AAAAAQAAAAEAAAB1Aw=="
	KeyClosedCaption                  KeyCode = "AAAAAgAAAKQAAAAQAw=="
	KeyComponent1                     KeyCode = "AAAAAgAAAKQAAAA2Aw=="
	KeyComponent2                     KeyCode = "AAAAAgAAAKQAAAA3Aw=="
	KeyWide                           KeyCode = "AAAAAgAAAKQAAAA9Aw=="
	KeyEPG                            KeyCode = "AAAAAgAAAKQAAABbAw=="
	KeyPAP                            KeyCode = "AAAAAgAAAKQAAAB3Aw=="
	KeyTenKey                         KeyCode = "AAAAAgAAAJcAAAAMAw=="
	KeyBSCS                           KeyCode = "AAAAAgAAAJcAAAAQAw=="
	KeyDdata                          KeyCode = "AAAAAgAAAJcAAAAVAw=="
	KeyStop                           KeyCode = "AAAAAgAAAJcAAAAYAw=="
	KeyPause                          KeyCode = "AAAAAgAAAJcAAAAZAw=="
	KeyPlay                           KeyCode = "AAAAAgAAAJcAAAAaAw=="
	KeyRewind                         KeyCode = "AAAAAgAAAJcAAAAbAw=="
	KeyForward                        KeyCode = "AAAAAgAAAJcAAAAcAw=="
	KeyDot                            KeyCode = "AAAAAgAAAJcAAAAdAw=="
	KeyRec                            KeyCode = "AAAAAgAAAJcAAAAgAw=="
	KeyReturn                         KeyCode = "AAAAAgAAAJcAAAAjAw=="
	KeyBlue                           KeyCode = "AAAAAgAAAJcAAAAkAw=="
	KeyRed                            KeyCode = "AAAAAgAAAJcAAAAlAw=="
	KeyGreen                          KeyCode = "AAAAAgAAAJcAAAAmAw=="
	KeyYellow                         KeyCode = "AAAAAgAAAJcAAAAnAw=="
	KeySubTitle                       KeyCode = "AAAAAgAAAJcAAAAoAw=="
	KeyCS                             KeyCode = "AAAAAgAAAJcAAAArAw=="
	KeyBS                             KeyCode = "AAAAAgAAAJcAAAAsAw=="
	KeyDigital                        KeyCode = "AAAAAgAAAJcAAAAyAw=="
	KeyOptions                        KeyCode = "AAAAAgAAAJcAAAA2Aw=="
	KeyMedia                          KeyCode = "AAAAAgAAAJcAAAA4Aw=="
	KeyPrev                           KeyCode = "AAAAAgAAAJcAAAA8Aw=="
	KeyNext                           KeyCode = "AAAAAgAAAJcAAAA9Aw=="
	KeyDpadCenter                     KeyCode = "AAAAAgAAAJcAAABKAw=="
	KeyCursorUp                       KeyCode = "AAAAAgAAAJcAAABPAw=="
	KeyCursorDown                     KeyCode = "AAAAAgAAAJcAAABQAw=="
	KeyCursorLeft                     KeyCode = "AAAAAgAAAJcAAABNAw=="
	KeyCursorRight                    KeyCode = "AAAAAgAAAJcAAABOAw=="
	KeyShopRemoteControlForcedDynamic KeyCode = "AAAAAgAAAJcAAABqAw=="
	KeyFlashPlus                      KeyCode = "AAAAAgAAAJcAAAB4Aw=="
	KeyFlashMinus                     KeyCode = "AAAAAgAAAJcAAAB5Aw=="
	KeyDemoMode                       KeyCode = "AAAAAgAAAJcAAAB8Aw=="
	KeyAnalog                         KeyCode = "AAAAAgAAAHcAAAANAw=="
	KeyMode3D                         KeyCode = "AAAAAgAAAHcAAABNAw=="
	KeyDigitalToggle                  KeyCode = "AAAAAgAAAHcAAABSAw=="
	KeyDemoSurround                   KeyCode = "AAAAAgAAAHcAAAB7Aw=="
	KeyStar                           KeyCode = "AAAAAgAAABoAAAA7Aw=="
	KeyAudioMixUp                     KeyCode = "AAAAAgAAABoAAAA8Aw=="
	KeyAudioMixDown                   KeyCode = "AAAAAgAAABoAAAA9Aw=="
	KeyPhotoFrame                     KeyCode = "AAAAAgAAABoAAABVAw=="
	KeyTv_Radio                       KeyCode = "AAAAAgAAABoAAABXAw=="
	KeySyncMenu                       KeyCode = "AAAAAgAAABoAAABYAw=="
	KeyHdmi1                          KeyCode = "AAAAAgAAABoAAABaAw=="
	KeyHdmi2                          KeyCode = "AAAAAgAAABoAAABbAw=="
	KeyHdmi3                          KeyCode = "AAAAAgAAABoAAABcAw=="
	KeyHdmi4                          KeyCode = "AAAAAgAAABoAAABdAw=="
	KeyTopMenu                        KeyCode = "AAAAAgAAABoAAABgAw=="
	KeyPopUpMenu                      KeyCode = "AAAAAgAAABoAAABhAw=="
	KeyOneTouchTimeRec                KeyCode = "AAAAAgAAABoAAABkAw=="
	KeyOneTouchView                   KeyCode = "AAAAAgAAABoAAABlAw=="
	KeyDUX                            KeyCode = "AAAAAgAAABoAAABzAw=="
	KeyFootballMode                   KeyCode = "AAAAAgAAABoAAAB2Aw=="
	KeyiManual                        KeyCode = "AAAAAgAAABoAAAB7Aw=="
	KeyNetflix                        KeyCode = "AAAAAgAAABoAAAB8Aw=="
	KeyAssists                        KeyCode = "AAAAAgAAAMQAAAA7Aw=="
	KeyFeaturedApp                    KeyCode = "AAAAAgAAAMQAAABEAw=="
	KeyFeaturedAppVOD                 KeyCode = "AAAAAgAAAMQAAABFAw=="
	KeyGooglePlay                     KeyCode = "AAAAAgAAAMQAAABGAw=="
	KeyActionMenu                     KeyCode = "AAAAAgAAAMQAAABLAw=="
	KeyHelp                           KeyCode = "AAAAAgAAAMQAAABNAw=="
	KeyTvSatellite                    KeyCode = "AAAAAgAAAMQAAABOAw=="
	KeyWirelessSubwoofer              KeyCode = "AAAAAgAAAMQAAAB+Aw=="
	KeyAndroidMenu                    KeyCode = "AAAAAgAAAMQAAABPAw=="
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

// Hack to pull out the description - SOAP XML parsing is too painful to be worth it just for this...
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
