package api

import (
	"github.com/gorilla/mux"
)

func NewRouter(filePath string) (*mux.Router, error) {
	// Load data file path to models
	DataFilePath = filePath
	if err := LoadData(); err != nil {
		return nil, err
	}

	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	router.HandleFunc("/regions/{region}", CreateRegion).Methods("POST")
	router.HandleFunc("/regions", ListRegions).Methods("GET")
	router.HandleFunc("/regions/{region}/environments/{environment}", CreateEnvironment).Methods("POST")
	router.HandleFunc("/regions/{region}/environments", ListEnvironments).Methods("GET")
	router.HandleFunc("/regions/{region}/environments/{environment}/apps/{app}", CreateApp).Methods("POST")
	router.HandleFunc("/regions/{region}/environments/{environment}/apps", ListApps).Methods("GET")
	router.HandleFunc("/regions/{region}/environments/{environment}/apps/{app}", GetApp).Methods("GET")

	return router, nil
}
