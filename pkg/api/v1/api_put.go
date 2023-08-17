package api

import (
	"encoding/json"
	"net/http"
	"time"
	"vhub/pkg/data"

	"github.com/gorilla/mux"
)

// UpdateRegion handles the PUT request to update an existing region.
func UpdateRegion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName := vars["region"]

	var region data.Region

	if err := json.NewDecoder(r.Body).Decode(&region); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	if _, exists := data.GlobalData.Regions[regionName]; !exists {
		RespondWithError(w, http.StatusNotFound, "Region not found")
		return
	}

	data.GlobalData.Regions[regionName] = region

	RespondWithJSON(w, http.StatusOK, region)
}

// UpdateEnvironment handles the PUT request to update an existing environment.
func UpdateEnvironment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName := vars["region"]
	environmentName := vars["environment"]

	var environment data.Environment

	if err := json.NewDecoder(r.Body).Decode(&environment); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	region, exists := data.GlobalData.Regions[regionName]
	if !exists {
		RespondWithError(w, http.StatusNotFound, "Region not found")
		return
	}

	if _, exists := region.Environments[environmentName]; !exists {
		RespondWithError(w, http.StatusNotFound, "Environment not found")
		return
	}

	region.Environments[environmentName] = environment
	data.GlobalData.Regions[regionName] = region

	RespondWithJSON(w, http.StatusOK, environment)
}

// UpdateApp handles the PUT request to update an existing app or just update the version.
func UpdateApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName := vars["region"]
	environmentName := vars["environment"]
	appName := vars["app"]

	var app data.App

	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	region, ok1 := data.GlobalData.Regions[regionName]
	environment, ok2 := region.Environments[environmentName]
	oldApp, ok3 := environment.Apps[appName]
	if !ok1 || !ok2 || !ok3 {
		RespondWithError(w, http.StatusNotFound, "Region, environment, or app not found")
		return
	}

	// If only the version is updated, retain other fields and update the date
	if app.Name == "" {
		app.Name = oldApp.Name
		app.Route = oldApp.Route
		app.Date = time.Now().Format(time.RFC3339) // Update date
	}

	environment.Apps[appName] = app
	region.Environments[environmentName] = environment
	data.GlobalData.Regions[regionName] = region

	RespondWithJSON(w, http.StatusOK, app)
}
