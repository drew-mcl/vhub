package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func NewRouter(filePath string) (*mux.Router, error) {
	// Load data file path to models
	DataFilePath = filePath
	if err := LoadData(); err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	// Handle the root path separately
	router.HandleFunc("/", ServeHTML).Methods("GET")
	router.HandleFunc("/cvs", ServeCSV).Methods("GET")

	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	ListRoutes := func(w http.ResponseWriter, r *http.Request) {
		router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			template, err := route.GetPathTemplate()
			if err == nil {
				methods, _ := route.GetMethods()
				fmt.Fprintln(w, "Path:", template)
				fmt.Fprintln(w, "Methods:", strings.Join(methods, ","))
			}
			return nil
		})
	}

	apiRouter.HandleFunc("/", ListRoutes).Methods("GET")
	apiRouter.HandleFunc("/regions/{region}", CreateRegion).Methods("POST")
	apiRouter.HandleFunc("/regions", ListRegions).Methods("GET")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}", CreateEnvironment).Methods("POST")
	apiRouter.HandleFunc("/regions/{region}/environments", ListEnvironments).Methods("GET")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}/apps/{app}", CreateApp).Methods("POST")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}/apps", ListApps).Methods("GET")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}/apps/{app}", GetApp).Methods("GET")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}/apps/{app}/version", GetAppVersion).Methods("GET")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}/apps/{app}/version", UpdateAppVersion).Methods("PUT")

	return router, nil
}
