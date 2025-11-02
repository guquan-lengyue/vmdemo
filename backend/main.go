package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func proxyWS(w http.ResponseWriter, r *http.Request, target string, verbose bool) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}
	defer ws.Close()

	if verbose {
		log.Printf("New WS client: %s -> %s", r.RemoteAddr, target)
	}

	tcpConn, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("dial target %s error: %v", target, err)
		return
	}
	defer tcpConn.Close()

	// WS -> TCP
	errc := make(chan error, 2)

	go func() {
		for {
			mt, msg, err := ws.ReadMessage()
			if err != nil {
				errc <- err
				return
			}
			// Only send the payload bytes to TCP
			if mt == websocket.TextMessage || mt == websocket.BinaryMessage {
				_, werr := tcpConn.Write(msg)
				if werr != nil {
					errc <- werr
					return
				}
			}
		}
	}()

	// TCP -> WS
	go func() {
		buf := make([]byte, 8192)
		for {
			n, rerr := tcpConn.Read(buf)
			if rerr != nil {
				errc <- rerr
				return
			}
			// Send as binary message
			werr := ws.WriteMessage(websocket.BinaryMessage, buf[:n])
			if werr != nil {
				errc <- werr
				return
			}
		}
	}()

	// Wait for first error
	e := <-errc
	if verbose {
		log.Printf("closing proxy for %s: %v", r.RemoteAddr, e)
	}
}

func main() {
	var webDir string
	var certFile string
	var keyFile string
	var verbose bool

	flag.StringVar(&webDir, "web", "", "Serve static files from this directory")
	flag.StringVar(&certFile, "cert", "", "TLS certificate file (optional)")
	flag.StringVar(&keyFile, "key", "", "TLS key file (optional, defaults to cert)")
	flag.BoolVar(&verbose, "v", false, "Verbose logging")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [--web dir] [--cert cert.pem [--key key.pem]] [source_addr:]source_port target_addr:target_port\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}

	sourceArg := args[0]
	targetArg := args[1]

	var listenAddr string
	// parse sourceArg
	if strings.Contains(sourceArg, ":") {
		listenAddr = sourceArg
	} else {
		listenAddr = ":" + sourceArg
	}

	// parse targetArg: must contain host:port
	if !strings.Contains(targetArg, ":") {
		fmt.Fprintln(os.Stderr, "target must be host:port")
		os.Exit(2)
	}
	target := targetArg

	if keyFile == "" {
		keyFile = certFile
	}

	if verbose {
		log.Printf("Proxying %s -> %s", listenAddr, target)
		if webDir != "" {
			log.Printf("Serving web directory: %s", webDir)
		}
		if certFile != "" {
			log.Printf("TLS enabled using cert=%s key=%s", certFile, keyFile)
		}
	}

	mux := http.NewServeMux()

	// If webDir provided, serve static files
	if webDir != "" {
		fs := http.FileServer(http.Dir(webDir))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// serve static files for normal HTTP requests
			fs.ServeHTTP(w, r)
		})
	}

	// WebSocket endpoint will accept any path and proxy to target
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If this is a websocket upgrade, handle proxy
		if strings.EqualFold(r.Header.Get("Connection"), "upgrade") || websocket.IsWebSocketUpgrade(r) {
			proxyWS(w, r, target, verbose)
			return
		}
		// Otherwise serve static if configured, else 403
		if webDir != "" {
			// file server already handles root, but this keeps behavior simple
			http.FileServer(http.Dir(webDir)).ServeHTTP(w, r)
			return
		}
		http.Error(w, "403 Permission Denied", http.StatusForbidden)
	})

	srv := &http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}

	// Start server
	if certFile != "" {
		// load cert to ensure readability
		_, cerr := tls.LoadX509KeyPair(certFile, keyFile)
		if cerr != nil {
			log.Fatalf("Failed to load cert/key: %v", cerr)
		}
		if verbose {
			log.Printf("Starting TLS server on %s", listenAddr)
		}
		log.Fatal(srv.ListenAndServeTLS(certFile, keyFile))
	} else {
		if verbose {
			log.Printf("Starting HTTP server on %s", listenAddr)
		}
		log.Fatal(srv.ListenAndServe())
	}
}
