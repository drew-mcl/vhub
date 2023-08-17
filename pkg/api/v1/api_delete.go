package api

import (
	"net/http"

	"vhub/pkg/data"

	"github.com/gorilla/mux"
)

// DeleteRegion handles the DELETE request to delete a region.
func DeleteRegion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName := vars["region"]

	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	if _, exists := data.GlobalData.Regions[regionName]; !exists {
		RespondWithError(w, http.StatusNotFound, "Region not found")
		return
	}

	delete(data.GlobalData.Regions, regionName)

	RespondWithJSON(w, http.StatusOK, "Region deleted successfully")
}

// DeleteEnvironment handles the DELETE request to delete an environment within a region.
func DeleteEnvironment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName := vars["region"]
	environmentName := vars["environment"]

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

	delete(region.Environments, environmentName)
	data.GlobalData.Regions[regionName] = region

	RespondWithJSON(w, http.StatusOK, "Environment deleted successfully")
}

// DeleteApp handles the DELETE request to delete an app within an environment.
func DeleteApp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName := vars["region"]
	environmentName := vars["environment"]
	appName := vars["app"]

	data.Mutex.Lock()
	defer data.Mutex.Unlock()

	region, ok1 := data.GlobalData.Regions[regionName]
	environment, ok2 := region.Environments[environmentName]
	if !ok1 || !ok2 {
		RespondWithError(w, http.StatusNotFound, "Region or environment not found")
		return
	}

	if _, exists := environment.Apps[appName]; !exists {
		RespondWithError(w, http.StatusNotFound, "App not found")
		return
	}

	delete(environment.Apps, appName)
	region.Environments[environmentName] = environment
	data.GlobalData.Regions[regionName] = region

	RespondWithJSON(w, http.StatusOK, "App deleted successfully")
}
