package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	irccip "github.com/bbrks/irccip-go"
)

func (s *server) handleKeyPress() http.HandlerFunc {
	type request struct {
		KeyCode irccip.KeyCode `json:"key_code"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			write(w, http.StatusBadRequest, []byte(err.Error()))
		}
		defer r.Body.Close()

		var req request
		err = json.Unmarshal(b, &req)
		if err != nil {
			write(w, http.StatusBadRequest, []byte(err.Error()))
		}

		s.Log(LevelDebug, r.Context(), "sending key code: %s", req.KeyCode)
		err = s.irccipClient.SendKeyCode(req.KeyCode)
		if err != nil {
			s.Log(LevelError, r.Context(), "SendKeyCode error: %v", err)
			write(w, http.StatusBadRequest, []byte("SendKeyCode error: "+err.Error()))
		}
	}
}
