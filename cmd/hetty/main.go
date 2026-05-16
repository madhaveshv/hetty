package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const (
	// defaultAddr is the default address the HTTP proxy and admin interface listen on.
	defaultAddr = ":8080"
	// defaultCertsDir is the default directory for storing CA certificate and key.
	defaultCertsDir = ".hetty"
)

// version is set at build time via ldflags.
var version = "dev"

func main() {
	// CLI flags
	addr := flag.String("addr", defaultAddr, "Address to listen on (e.g. \":8080\")")
	certsDir := flag.String("certs", "", "Directory for CA certificate and key (default: $HOME/.hetty)")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	printVersion := flag.Bool("version", false, "Print version and exit")

	flag.Parse()

	if *printVersion {
		fmt.Printf("hetty %s\n", version)
		os.Exit(0)
	}

	// Resolve certs directory
	if *certsDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("[FATAL] Could not determine home directory: %v", err)
		}
		*certsDir = filepath.Join(homeDir, defaultCertsDir)
	}

	if err := os.MkdirAll(*certsDir, 0700); err != nil {
		log.Fatalf("[FATAL] Could not create certs directory %q: %v", *certsDir, err)
	}

	if *verbose {
		log.Printf("[INFO] Certs directory: %s", *certsDir)
		log.Printf("[INFO] Listening on %s", *addr)
	}

	// Set up TCP listener
	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("[FATAL] Failed to listen on %s: %v", *addr, err)
	}

	// Placeholder mux — proxy and admin UI handlers will be registered here
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hetty %s — proxy running\n", version)
	})

	srv := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	// Graceful shutdown on SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[INFO] hetty %s listening on %s", version, listener.Addr())
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] Server error: %v", err)
		}
	}()

	<-quit
	log.Println("[INFO] Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("[FATAL] Server forced to shutdown: %v", err)
	}

	log.Println("[INFO] Server stopped.")
}
