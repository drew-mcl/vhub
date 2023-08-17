package api

import (
	"encoding/json"
	"net/http"
	"vhub/pkg/data"

	"github.com/gorilla/mux"
)

// CreateRegion handles the POST request to create a new region.
func CreateRegion(w http.ResponseWriter, r *http.Request) {
	var region data.Region

	if err := json.NewDecoder(r.Body).Decode(&region); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	if _, exists := data.GlobalData.Regions[region.Name]; exists {
		RespondWithError(w, http.StatusConflict, "Region already exists")
		return
	}

	// Initialize the Apps map to an empty map
	region.Environments = make(map[string]data.Environment)

	data.GlobalData.Regions[region.Name] = region

	if err := data.SaveData(data.DataFilePath); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to save data")
		return
	}

	RespondWithJSON(w, http.StatusCreated, region)
}

// CreateEnvironment handles the POST request to create a new environment within a region.
func CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName := vars["region"]

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

	if region.Environments == nil {
		region.Environments = make(map[string]data.Environment)
	}

	if _, exists := region.Environments[environment.Name]; exists {
		RespondWithError(w, http.StatusConflict, "Environment already exists in this region")
		return
	}

	// Initialize the Apps map to an empty map
	environment.Apps = make(map[string]data.App)

	region.Environments[environment.Name] = environment
	data.GlobalData.Regions[regionName] = region

	if err := data.SaveData(data.DataFilePath); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to save data")
		return
	}

	RespondWithJSON(w, http.StatusCreated, environment)
}

// CreateApp handles the POST request to create a new app within an environment.
func CreateApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName := vars["region"]
	environmentName := vars["environment"]

	var app data.App

	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	region, ok1 := data.GlobalData.Regions[regionName]
	environment, ok2 := region.Environments[environmentName]
	if !ok1 || !ok2 {
		RespondWithError(w, http.StatusNotFound, "Region or environment not found")
		return
	}

	// Initialize the Apps map if it's nil
	if environment.Apps == nil {
		environment.Apps = make(map[string]data.App)
	}

	if _, exists := environment.Apps[app.Name]; exists {
		RespondWithError(w, http.StatusConflict, "App already exists in this environment")
		return
	}

	environment.Apps[app.Name] = app
	region.Environments[environmentName] = environment
	data.GlobalData.Regions[regionName] = region

	if err := data.SaveData(data.DataFilePath); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to save data")
		return
	}

	RespondWithJSON(w, http.StatusCreated, app)
}
