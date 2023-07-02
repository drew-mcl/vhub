package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type AppVersions struct {
	Environments map[string]map[string]int
	sync.RWMutex
}

func NewAppVersions() *AppVersions {
	return &AppVersions{
		Environments: make(map[string]map[string]int),
	}
}

func (av *AppVersions) GetVersions(env string) map[string]int {
	av.RLock()
	defer av.RUnlock()
	return av.Environments[env]
}

func (av *AppVersions) SetVersion(env, app string, version int) {
	av.Lock()
	defer av.Unlock()
	if av.Environments[env] == nil {
		av.Environments[env] = make(map[string]int)
	}
	av.Environments[env][app] = version
}

// APIHandler handles the API endpoints
func (av *AppVersions) APIHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/version":
		if r.Method == http.MethodGet {
			av.handleGet(w, r)
		} else if r.Method == http.MethodPost {
			av.handlePost(w, r)
		} else {
			http.NotFound(w, r)
		}
	case "/":
		if r.Method == http.MethodGet {
			av.handleIndex(w, r)
		} else {
			http.NotFound(w, r)
		}
	default:
		http.NotFound(w, r)
	}
}

// handleGet handles GET requests to retrieve the version number of an app in a given environment
func (av *AppVersions) handleGet(w http.ResponseWriter, r *http.Request) {
	env := r.URL.Query().Get("env")
	app := r.URL.Query().Get("app")

	version, ok := av.GetVersions(env)[app]
	if !ok {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(version)
	if err != nil {
		log.Printf("Failed to encode version: %v", err)
		http.Error(w, "Failed to encode version", http.StatusInternalServerError)
	}
}

// handlePost handles POST requests
func (av *AppVersions) handlePost(w http.ResponseWriter, r *http.Request) {
	env := r.FormValue("env")
	app := r.FormValue("app")
	version := r.FormValue("version")

	var v int
	if _, err := fmt.Sscanf(version, "%d", &v); err != nil {
		http.Error(w, "Invalid version number", http.StatusBadRequest)
		return
	}

	av.SetVersion(env, app, v)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Version %d has been set for app %s in environment %s.\n", v, app, env)
}

// handleIndex handles GET requests to the root endpoint and returns a landing page
func (av *AppVersions) handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Welcome to WhatVersion!</h1><p>You can retrieve and update app version information here.</p>")
}

func main() {
	configFile := flag.String("config", "config.json", "Path to the configuration file")
	port := flag.Int("port", 8080, "Port number for the HTTP server")
	flag.Parse()

	appVersions := NewAppVersions()

	config, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configurations: %v", err)
	}

	for _, env := range config.Environments {
		for app, version := range env.Versions {
			appVersions.SetVersion(env.Name, app, version)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: http.TimeoutHandler(http.HandlerFunc(appVersions.APIHandler), 5*time.Second, "Request timeout"),
	}

	go func() {
		log.Printf("Server is starting on port %d", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	waitForTerminationSignal(ctx)

	log.Println("Shutting down server gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown server gracefully: %v", err)
	} else {
		log.Println("Server has been shutdown gracefully")
	}
}

type Environment struct {
	Name     string         `json:"name"`
	Versions map[string]int `json:"versions"`
}

type Config struct {
	Environments []*Environment `json:"environments"`
}

func LoadConfig(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer f.Close()

	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &config, nil
}

func waitForTerminationSignal(ctx context.Context) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signalCh:
		log.Println("Termination signal received")
	case <-ctx.Done():
		log.Println("Termination from context done")
	}
}
