package main

import (
	"net/http"

	irccip "github.com/bbrks/irccip-go"
	"github.com/julienschmidt/httprouter"
)

type server struct {
	Logger     // Embed Log()
	router     *httprouter.Router
	httpClient *http.Client

	irccipClient *irccip.Client
}

func write(w http.ResponseWriter, status int, body []byte) {
	w.WriteHeader(status)
	_, _ = w.Write(body)
	_, _ = w.Write([]byte("\n"))
}
