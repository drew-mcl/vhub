package api

import (
	"encoding/csv"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func handleLogError(w http.ResponseWriter, statusCode int, fields logrus.Fields, msg string) {
	Log.WithFields(fields).Warn(msg)
	w.WriteHeader(statusCode)
}

func regionExists(region string) bool {
	_, ok := Regions[region]
	return ok
}

func environmentExists(region, environment string) bool {
	if !regionExists(region) {
		return false
	}
	_, ok := Regions[region][environment]
	return ok
}

func appExists(region, environment, app string) bool {
	if !environmentExists(region, environment) {
		return false
	}
	_, ok := Regions[region][environment].Apps[app]
	return ok
}

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	Mutex.Lock()
	defer Mutex.Unlock()

	tmpl, err := template.ParseFiles("templates/template.html")
	if err != nil {
		handleLogError(w, http.StatusInternalServerError, logrus.Fields{"error": err}, "Failed to parse template")
		return
	}
	if err := tmpl.Execute(w, Regions); err != nil {
		handleLogError(w, http.StatusInternalServerError, logrus.Fields{"error": err}, "Failed to execute template")
		return
	}
	Log.WithFields(logrus.Fields{}).Debug("Served HTML data")
}

func ServeCSV(w http.ResponseWriter, r *http.Request) {
	Mutex.Lock()
	defer Mutex.Unlock()

	// Prepare CSV data
	records := [][]string{
		{"Region", "Environment", "Name", "App Name", "Version"},
	}
	for region, environments := range Regions {
		for env, details := range environments {
			for _, app := range details.Apps {
				records = append(records, []string{region, env, details.Name, app.Name, app.Version})
			}
		}
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=data.csv")
	writer := csv.NewWriter(w)
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			handleLogError(w, http.StatusInternalServerError, logrus.Fields{"error": err}, "Failed to write CSV")
			return
		}
	}
	writer.Flush()
	Log.WithFields(logrus.Fields{}).Debug("Served CSV data")
}

func ListRegions(w http.ResponseWriter, r *http.Request) {
	Mutex.Lock()
	defer Mutex.Unlock()
	if err := sendJSONResponse(w, http.StatusOK, Regions); err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
	}
}

func CreateRegion(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	Mutex.Lock()
	defer Mutex.Unlock()
	if regionExists(region) {
		handleLogError(w, http.StatusConflict, logrus.Fields{"region": region}, "Region already exists")
		return
	}
	Regions[region] = make(map[string]Environment)
	handleLogError(w, http.StatusCreated, logrus.Fields{"region": region}, "Created new region")
}

func ListEnvironments(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	Mutex.Lock()
	defer Mutex.Unlock()
	if regionExists(region) {
		if err := sendJSONResponse(w, http.StatusOK, Regions[region]); err != nil {
			http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		}
		return
	}
	handleLogError(w, http.StatusNotFound, logrus.Fields{"region": region}, "Region not found")
}

func CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	environment := mux.Vars(r)["environment"]
	Mutex.Lock()
	defer Mutex.Unlock()
	if regionExists(region) {
		if environmentExists(region, environment) {
			handleLogError(w, http.StatusConflict, logrus.Fields{"region": region, "environment": environment}, "Environment already exists")
			return
		}
		Regions[region][environment] = Environment{Name: environment, Apps: make(map[string]App)}
		handleLogError(w, http.StatusCreated, logrus.Fields{"region": region, "environment": environment}, "Created new environment")
		return
	}
	handleLogError(w, http.StatusNotFound, logrus.Fields{"region": region}, "Region not found")
}

func ListApps(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	environment := mux.Vars(r)["environment"]
	Mutex.Lock()
	defer Mutex.Unlock()
	if environmentExists(region, environment) {
		if err := sendJSONResponse(w, http.StatusOK, Regions[region][environment].Apps); err != nil {
			http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		}
		return
	}
	handleLogError(w, http.StatusNotFound, logrus.Fields{"region": region, "environment": environment}, "Environment not found")
}

func CreateApp(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	environment := mux.Vars(r)["environment"]
	appName := mux.Vars(r)["app"]
	var newApp App
	_ = json.NewDecoder(r.Body).Decode(&newApp)
	newApp.Name = appName
	Mutex.Lock()
	defer Mutex.Unlock()
	if environmentExists(region, environment) {
		if appExists(region, environment, appName) {
			handleLogError(w, http.StatusConflict, logrus.Fields{"region": region, "environment": environment, "app": appName}, "App already exists")
			return
		}
		Regions[region][environment].Apps[appName] = newApp
		handleLogError(w, http.StatusCreated, logrus.Fields{"region": region, "environment": environment, "app": appName}, "Created new app")
		return
	}
	handleLogError(w, http.StatusNotFound, logrus.Fields{"region": region, "environment": environment}, "Environment not found")
}

func GetApp(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	environment := mux.Vars(r)["environment"]
	app := mux.Vars(r)["app"]
	Mutex.Lock()
	defer Mutex.Unlock()
	if appExists(region, environment, app) {
		if err := sendJSONResponse(w, http.StatusOK, Regions[region][environment].Apps[app]); err != nil {
			http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		}
		return
	}
	handleLogError(w, http.StatusNotFound, logrus.Fields{"region": region, "environment": environment, "app": app}, "App not found")
}
