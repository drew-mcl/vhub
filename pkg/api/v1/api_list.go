package api

import (
	"net/http"

	"vhub/pkg/data"

	"github.com/gorilla/mux"
)

// ListRegions handles the GET request for listing all regions.
func ListRegions(w http.ResponseWriter, r *http.Request) {
	data.Mutex.RLock()
	defer data.Mutex.RUnlock()

	regions := make([]data.Region, 0, len(data.GlobalData.Regions))
	for _, region := range data.GlobalData.Regions {
		regions = append(regions, region)
	}

	RespondWithJSON(w, http.StatusOK, regions)
}

// ListEnvironments handles the GET request for listing all environments in a region.
func ListEnvironments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName, ok := vars["region"]

	data.Mutex.RLock()
	defer data.Mutex.RUnlock()

	region, ok := data.GlobalData.Regions[regionName]
	if !ok {
		RespondWithError(w, http.StatusNotFound, "Region not found")
		return
	}

	environments := make([]data.Environment, 0, len(region.Environments))
	for _, env := range region.Environments {
		environments = append(environments, env)
	}

	RespondWithJSON(w, http.StatusOK, environments)
}

// ListApps handles the GET request for listing all apps in an environment.
func ListApps(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	regionName, ok1 := vars["region"]
	environmentName, ok2 := vars["environment"]

	data.Mutex.RLock()
	defer data.Mutex.RUnlock()

	region, ok1 := data.GlobalData.Regions[regionName]
	environment, ok2 := region.Environments[environmentName]
	if !ok1 || !ok2 {
		RespondWithError(w, http.StatusNotFound, "Region or environment not found")
		return
	}

	apps := make([]data.App, 0, len(environment.Apps))
	for _, app := range environment.Apps {
		apps = append(apps, app)
	}

	RespondWithJSON(w, http.StatusOK, apps)
}
