package api

import (
	"fmt"
	"net/http"
	"strings"
	"vhub/pkg/data" // Update with the actual import path to the data package

	"github.com/gorilla/mux"
)

func NewRouter() (*mux.Router, error) {
	// Load data
	if err := data.LoadData(); err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	// Handle the root path separately
	router.HandleFunc("/", ServeHTML).Methods("GET")
	router.HandleFunc("/healthcheck", HealthCheck).Methods("GET")

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

	// Regions
	apiRouter.HandleFunc("/regions", ListRegions).Methods("GET")
	apiRouter.HandleFunc("/regions", CreateRegion).Methods("POST")
	apiRouter.HandleFunc("/regions/{region}", GetRegion).Methods("GET")
	apiRouter.HandleFunc("/regions/{region}", UpdateRegion).Methods("PUT")
	apiRouter.HandleFunc("/regions/{region}", DeleteRegion).Methods("DELETE")

	// Environments
	apiRouter.HandleFunc("/regions/{region}/environments", ListEnvironments).Methods("GET")
	apiRouter.HandleFunc("/regions/{region}/environments", CreateEnvironment).Methods("POST")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}", GetEnvironment).Methods("GET")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}", UpdateEnvironment).Methods("PUT")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}", DeleteEnvironment).Methods("DELETE")

	// Apps
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}/apps", ListApps).Methods("GET")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}/apps", CreateApp).Methods("POST")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}/apps/{app}", GetApp).Methods("GET")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}/apps/{app}", UpdateApp).Methods("PUT")
	apiRouter.HandleFunc("/regions/{region}/environments/{environment}/apps/{app}", DeleteApp).Methods("DELETE")

	return router, nil
}
