package checker

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type HealthCheckConfig struct {
	EnableHealthCheck bool           `json:"enableHealthCheck"`
	HealthChecks      []HealthStatus `json:"healthChecks"`
}

type HealthStatus struct {
	Region      string
	Environment string
	URL         string
	Status      string
	LastChecked time.Time
}

var (
	statusData []HealthStatus
	mu         sync.RWMutex
)

func StartHealthChecks(configFile string) {
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		return
	}

	if !config.EnableHealthCheck {
		return
	}

	statusData = make([]HealthStatus, len(config.HealthChecks))
	copy(statusData, config.HealthChecks)

	for i := range statusData {
		statusData[i].Status = "Unknown"
	}

	go func() {
		for {
			for i := range statusData {
				checkService(i)
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}

func loadConfig(file string) (HealthCheckConfig, error) {
	var config HealthCheckConfig

	configFile, err := os.Open(file)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

func checkService(i int) {
	healthURL := statusData[i].URL + "/healthcheck"

	resp, err := http.Get(healthURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		statusData[i].Status = "Fail"
		statusData[i].LastChecked = time.Now()
		return
	}

	defer resp.Body.Close()

	statusData[i].Status = "OK"
	statusData[i].LastChecked = time.Now()
}

func GetHealthStatus() []HealthStatus {
	mu.RLock()
	defer mu.RUnlock()

	return append([]HealthStatus(nil), statusData...)
}
