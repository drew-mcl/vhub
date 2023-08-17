package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"vhub/pkg/api/v1"
	"vhub/pkg/checker"
	"vhub/pkg/data"

	"github.com/sirupsen/logrus"
)

func StartBackupSavingInterval(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)

			backupFilePath := data.BackupFilePath
			if err := data.SaveData(backupFilePath); err != nil {
				data.Log.WithFields(logrus.Fields{
					"error": err,
				}).Error("Failed to save backup data at interval")
			} else {
				data.Log.Debug("Backup data saved successfully at interval")
			}
		}
	}()
}

func main() {
	// Command-line flags
	host := flag.String("host", "localhost", "Define host of the server")
	port := flag.String("port", "8080", "Define port of the server")
	filePath := flag.String("filePath", "", "Define path of the data file")
	enableHealthCheck := flag.Bool("checker", false, "Enable health check")
	checkerConfig := flag.String("checker-config", "config/checker.json", "supply config for checker")
	flag.Parse()

	if *filePath == "" {
		logrus.Warn("No -filePath flag provided. Defaulting to $(pwd)/data.json")
		executable, err := os.Executable()
		if err != nil {
			panic(err)
		}
		executableDir := filepath.Dir(executable)
		data.DataFilePath = filepath.Join(executableDir, "data.json")
	} else {
		data.DataFilePath = *filePath
	}

	// Set the backup file path based on the primary data file path
	data.BackupFilePath = data.DataFilePath + ".backup"

	// Create file if not exists
	if err := data.CreateFileIfNotExists(data.DataFilePath); err != nil {
		log.Fatal(err)
	}

	if err := data.CreateFileIfNotExists(data.BackupFilePath); err != nil {
		log.Fatal(err)
	}

	if err := data.LoadData(); err != nil {
		logrus.Fatalf("Failed to load data: %v", err)
	}

	// Initialize and check the router
	router, err := api.NewRouter() // Update this according to your new routing setup
	if err != nil {
		logrus.Fatalf("Failed to initialize router: %v", err)
	}

	// Create a new server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", *host, *port),
		Handler: router,
	}

	// Start the backup saving interval
	StartBackupSavingInterval(5 * time.Minute)

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Server failed to start: %v", err)
		}
	}()

	logrus.Infof("Server is listening on %s", server.Addr)

	if *enableHealthCheck {
		checker.StartHealthChecks(*checkerConfig)
	}

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
