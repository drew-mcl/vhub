package api

import (
	"encoding/csv"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func sendClientResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message)) // Send the message to the client
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func handleLogError(statusCode int, fields logrus.Fields, msg string) {
	Log.WithFields(fields).Warn(msg)
}

func handleLogSuccess(statusCode int, fields logrus.Fields, msg string) {
	Log.WithFields(fields).Info(msg)
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

	Mutex.RLock()
	defer Mutex.RUnlock()

	tmpl, err := template.ParseFiles("templates/template.html")
	if err != nil {
		handleLogError(http.StatusInternalServerError, logrus.Fields{"error": err}, "Failed to parse template")
		return
	}
	if err := tmpl.Execute(w, Regions); err != nil {
		handleLogError(http.StatusInternalServerError, logrus.Fields{"error": err}, "Failed to execute template")
		return
	}
	Log.WithFields(logrus.Fields{}).Debug("Served HTML data")
}

func ServeCSV(w http.ResponseWriter, r *http.Request) {

	Mutex.RLock()
	defer Mutex.RUnlock()

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
			handleLogError(http.StatusInternalServerError, logrus.Fields{"error": err}, "Failed to write CSV")
			return
		}
	}
	writer.Flush()
	Log.WithFields(logrus.Fields{}).Debug("Served CSV data")
}

func ListRegions(w http.ResponseWriter, r *http.Request) {

	Mutex.RLock()
	defer Mutex.RUnlock()

	if err := sendJSONResponse(w, http.StatusOK, Regions); err != nil {
		http.Error(w, "Failed to encode data", http.StatusInternalServerError)
	}
}

func ListEnvironments(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]

	Mutex.RLock()
	defer Mutex.RUnlock()

	if regionExists(region) {
		if err := sendJSONResponse(w, http.StatusOK, Regions[region]); err != nil {
			http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		}
		return
	}
	msg := "Incorrect input"
	handleLogError(http.StatusNotFound, logrus.Fields{"region": region}, msg)
	sendClientResponse(w, http.StatusNotFound, msg)
}

func ListApps(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	environment := mux.Vars(r)["environment"]

	Mutex.RLock()
	defer Mutex.RUnlock()

	if environmentExists(region, environment) {
		if err := sendJSONResponse(w, http.StatusOK, Regions[region][environment].Apps); err != nil {
			http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		}
		return
	}
	msg := "Incorrect input"
	handleLogError(http.StatusNotFound, logrus.Fields{"region": region, "environment": environment}, msg)
	sendClientResponse(w, http.StatusNotFound, msg)
}

func GetApp(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	environment := mux.Vars(r)["environment"]
	app := mux.Vars(r)["app"]

	Mutex.RLock()
	defer Mutex.RUnlock()

	if appExists(region, environment, app) {
		if err := sendJSONResponse(w, http.StatusOK, Regions[region][environment].Apps[app]); err != nil {
			http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		}
		return
	}
	msg := "Incorrect input"
	handleLogError(http.StatusNotFound, logrus.Fields{"region": region, "environment": environment}, msg)
	sendClientResponse(w, http.StatusNotFound, msg)
}

func GetAppVersion(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	environment := mux.Vars(r)["environment"]
	app := mux.Vars(r)["app"]

	Mutex.RLock()
	defer Mutex.RUnlock()

	if appExists(region, environment, app) {
		version := Regions[region][environment].Apps[app].Version
		if err := sendJSONResponse(w, http.StatusOK, map[string]string{"version": version}); err != nil {
			http.Error(w, "Failed to encode data", http.StatusInternalServerError)
		}
		return

	}
	msg := "Incorrect input"
	handleLogError(http.StatusNotFound, logrus.Fields{"region": region, "environment": environment}, msg)
	sendClientResponse(w, http.StatusNotFound, msg)
}

// ------------------ STATE MOD ----------------------------
func CreateRegion(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]

	Mutex.Lock()
	defer Mutex.Unlock()

	if !regionExists(region) {
		msg := "Region already exists"
		handleLogError(http.StatusNotFound, logrus.Fields{"region": region}, msg)
		sendClientResponse(w, http.StatusNotFound, msg)
		return
	}
	Regions[region] = make(map[string]Environment)

	if err := SaveData(DataFilePath); err != nil {
		Log.Error("Failed to save data:", err)
	}

	msg := "Created new environment"
	Log.WithFields(logrus.Fields{"region": region}).Info(msg)
	sendClientResponse(w, http.StatusCreated, msg)
}

func CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	environment := mux.Vars(r)["environment"]

	Mutex.Lock()
	defer Mutex.Unlock()

	if !regionExists(region) {
		msg := "Region not found"
		handleLogError(http.StatusNotFound, logrus.Fields{"region": region}, msg)
		sendClientResponse(w, http.StatusNotFound, msg)
		return
	}

	if environmentExists(region, environment) {
		msg := "Environment already exists"
		handleLogError(http.StatusConflict, logrus.Fields{"region": region, "environment": environment}, msg)
		sendClientResponse(w, http.StatusConflict, msg)
		return
	}

	Regions[region][environment] = Environment{Name: environment, Apps: make(map[string]App)}

	if err := SaveData(DataFilePath); err != nil {
		Log.Error("Failed to save data:", err)
	}

	msg := "Created new environment"
	Log.WithFields(logrus.Fields{"region": region, "environment": environment}).Info(msg)
	sendClientResponse(w, http.StatusCreated, msg)
}

func CreateApp(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	environment := mux.Vars(r)["environment"]
	appName := mux.Vars(r)["app"]
	var requestBody struct {
		Version string `json:"version,omitempty"`
	}
	_ = json.NewDecoder(r.Body).Decode(&requestBody)

	if requestBody.Version == "" {
		requestBody.Version = "undefined"
	}

	Mutex.Lock()
	defer Mutex.Unlock()

	if !environmentExists(region, environment) {
		msg := "Environment not found"
		handleLogError(http.StatusNotFound, logrus.Fields{"region": region, "environment": environment}, msg)
		sendClientResponse(w, http.StatusNotFound, msg)
		return
	}

	if appExists(region, environment, appName) {
		msg := "App already exists"
		handleLogError(http.StatusConflict, logrus.Fields{"region": region, "environment": environment, "app": appName}, msg)
		sendClientResponse(w, http.StatusConflict, msg)
		return
	}

	newApp := App{Name: appName, Version: requestBody.Version}
	Regions[region][environment].Apps[appName] = newApp

	if err := SaveData(DataFilePath); err != nil {
		Log.Error("Failed to save data:", err)
	}

	msg := "Created new app"
	Log.WithFields(logrus.Fields{"region": region, "environment": environment, "app": appName}).Info(msg)
	sendClientResponse(w, http.StatusCreated, msg)
}

func UpdateAppVersion(w http.ResponseWriter, r *http.Request) {
	region := mux.Vars(r)["region"]
	environment := mux.Vars(r)["environment"]
	appName := mux.Vars(r)["app"]

	var requestData map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		msg := "Failed to decode request body"
		handleLogError(http.StatusBadRequest, logrus.Fields{"error": err}, msg)
		sendClientResponse(w, http.StatusBadRequest, msg)
		return
	}
	version, ok := requestData["version"]
	if !ok {
		msg := "Version not provided"
		handleLogError(http.StatusBadRequest, logrus.Fields{}, msg)
		sendClientResponse(w, http.StatusBadRequest, msg)
		return
	}

	Mutex.Lock()
	defer Mutex.Unlock()

	if !appExists(region, environment, appName) {
		msg := "App not found"
		handleLogError(http.StatusNotFound, logrus.Fields{"region": region, "environment": environment, "app": appName}, msg)
		sendClientResponse(w, http.StatusNotFound, msg)
		return
	}

	app := Regions[region][environment].Apps[appName]
	app.Version = version
	Regions[region][environment].Apps[appName] = app

	if err := SaveData(DataFilePath); err != nil {
		Log.Error("Failed to save data:", err)
	}

	msg := "Updated app version"
	Log.WithFields(logrus.Fields{"region": region, "environment": environment, "app": appName}).Info(msg)
	sendClientResponse(w, http.StatusOK, msg)
}
