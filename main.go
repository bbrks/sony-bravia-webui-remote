package main

import (
	"context"
	"flag"
	"net/http"
	"os"

	irccip "github.com/bbrks/irccip-go"
	"github.com/julienschmidt/httprouter"
)

const (
	appName    = "sony-bravia-webui-remote"
	appVersion = "0.1"
)

func main() {
	envIP := os.Getenv("SONY_BRAVIA_IP")
	flagIP := flag.String("ip", "", "Your Unique Property Reference Number (UPRN)")

	envPSK := os.Getenv("SONY_BRAVIA_PSK")
	flagPSK := flag.String("psk", "", "Your Unique Property Reference Number (UPRN)")

	flagBindAddr := flag.String("http", ":8080", "The address and/or port to bind to the HTTP server")
	flagLogQuiet := flag.Bool("q", false, "Disables all logging except for errors")
	flag.Parse()

	s := server{
		Logger:     newDefaultLogger(*flagLogQuiet),
		router:     httprouter.New(),
		httpClient: http.DefaultClient,
	}
	ctx := context.Background()

	ip := *flagIP
	if ip == "" {
		ip = envIP
	}
	if ip == "" {
		s.Log(LevelError, ctx, "Display IP was not set. Set with \"SONY_BRAVIA_IP\" env var or -ip=\"1.2.3.4\" flag")
		os.Exit(1)
	}

	psk := *flagPSK
	if psk == "" {
		psk = envPSK
	}
	if psk == "" {
		s.Log(LevelError, ctx, "Pre-shared key was not set. Set with \"SONY_BRAVIA_PSK\" env var or -psk=\"1234\" flag")
		os.Exit(1)
	}

	s.irccipClient = irccip.NewClient("http://"+ip, psk)

	s.routes()

	s.Log(LevelInfo, ctx, "Starting server")
	if err := http.ListenAndServe(*flagBindAddr, s.router); err != nil {
		s.Log(LevelError, ctx, err.Error())
		os.Exit(1)
	}
}
