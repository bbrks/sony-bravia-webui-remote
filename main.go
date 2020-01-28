package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"

	_ "net/http/pprof"

	irccip "github.com/bbrks/irccip-go"
	"github.com/julienschmidt/httprouter"
)

const (
	appName    = "sony-bravia-webui-remote"
	appVersion = "0.1"
)

func main() {
	envIP := os.Getenv("SONY_BRAVIA_IP")
	flagIP := flag.String("ip", "", "The IP address of the Sony Bravia display")

	envPSK := os.Getenv("SONY_BRAVIA_PSK")
	flagPSK := flag.String("psk", "", "The configured Pre-Shared-Key (PSK) of the Sony Bravia display")

	flagBindAddr := flag.String("http", ":8080", "The address and/or port to bind to the HTTP server")
	flagLogQuiet := flag.Bool("q", false, "Disables all logging except for errors")
	flagPprofAddr := flag.String("pprof", "localhost:6060", "The address and/or port to bind to the pprof HTTP server")
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

	if *flagPprofAddr != "" {
		l, err := net.Listen("tcp", *flagPprofAddr)
		if err != nil {
			s.Log(LevelError, ctx, "Error listening on %v: %v", *flagBindAddr, err)
			os.Exit(1)
		}

		s.Log(LevelDebug, ctx, "Serving pprof at http://%s/debug/pprof/...", l.Addr().String())
		go func() {
			if err = http.Serve(l, nil); err != nil {
				s.Log(LevelError, ctx, "Error from pprof HTTP server: %v", err)
				os.Exit(1)
			}
		}()
	}

	s.irccipClient = irccip.NewClient("http://"+ip, psk)
	s.routes()

	l, err := net.Listen("tcp", *flagBindAddr)
	if err != nil {
		s.Log(LevelError, ctx, "Error listening on %v: %v", *flagBindAddr, err)
		os.Exit(1)
	}

	s.Log(LevelInfo, ctx, "Serving at http://%s", l.Addr().String())
	if err := http.Serve(l, s.middleware(s.router)); err != nil {
		s.Log(LevelError, ctx, "Error from HTTP server: %v", err)
		os.Exit(1)
	}
}
