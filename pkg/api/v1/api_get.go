package api

import (
	"net/http"
	"vhub/pkg/data"

	"github.com/gorilla/mux"
)

// GetRegion handles the GET request to retrieve a specific region.
func GetRegion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName := vars["region"]

	data.Mutex.RLock()
	defer data.Mutex.RUnlock()

	region, exists := data.GlobalData.Regions[regionName]
	if !exists {
		RespondWithError(w, http.StatusNotFound, "Region not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, region)
}

// GetEnvironment handles the GET request to retrieve a specific environment within a region.
func GetEnvironment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName := vars["region"]
	environmentName := vars["environment"]

	data.Mutex.RLock()
	defer data.Mutex.RUnlock()

	region, ok1 := data.GlobalData.Regions[regionName]
	environment, ok2 := region.Environments[environmentName]
	if !ok1 || !ok2 {
		RespondWithError(w, http.StatusNotFound, "Region or environment not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, environment)
}

// GetApp handles the GET request to retrieve a specific app within an environment.
func GetApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName := vars["region"]
	environmentName := vars["environment"]
	appName := vars["app"]

	data.Mutex.RLock()
	defer data.Mutex.RUnlock()

	region, ok1 := data.GlobalData.Regions[regionName]
	environment, ok2 := region.Environments[environmentName]
	app, ok3 := environment.Apps[appName]
	if !ok1 || !ok2 || !ok3 {
		RespondWithError(w, http.StatusNotFound, "Region, environment, or app not found")
		return
	}

	RespondWithJSON(w, http.StatusOK, app)
}
