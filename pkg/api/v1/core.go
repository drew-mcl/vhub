package api

import (
	"net/http"
	"vhub/pkg/data"

	"vhub/pkg/checker"

	"vhub/pkg/ui"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, map[string]string{"status": "OK"})
}

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	// Get region data and health data
	data.Mutex.RLock()
	defer data.Mutex.RUnlock()

	regionData := data.GlobalData.Regions
	healthData := checker.GetHealthStatus()

	// Render the template
	ui.RenderTemplate(w, regionData, healthData)
}
