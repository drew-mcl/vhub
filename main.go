package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	api "vhub/pkg/api/v1"

	"github.com/sirupsen/logrus"
)

var InitData bool

func main() {
	// Command-line flags
	host := flag.String("host", "localhost", "Define host of the server")
	port := flag.String("port", "8080", "Define port of the server")

	filePath := flag.String("filePath", "/var/vhub/data.json", "Define path of the data file")
	flag.Parse()

	// Initialize and check the router
	router, err := api.NewRouter(*filePath)
	if err != nil {
		logrus.Fatalf("Failed to initialize router: %v", err)
	}

	// Create a new server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", *host, *port),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Server failed to start: %v", err)
		}
	}()

	logrus.Infof("Server is listening on %s", server.Addr)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	// Create a deadline for the current context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("Server shutdown failed: %v", err)
	}

	logrus.Println("Server exited properly")
}
