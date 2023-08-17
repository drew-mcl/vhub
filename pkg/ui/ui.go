package ui

import (
	"html/template"
	"log"
	"net/http"
	"vhub/pkg/checker" // Make sure to import your checker package
	"vhub/pkg/data"
)

type ViewData struct {
	Regions map[string]data.Region `json:"regions"`
	Health  []checker.HealthStatus `json:"health"` // Health status
}

func RenderTemplate(w http.ResponseWriter, regionData map[string]data.Region, healthData []checker.HealthStatus) {
	tmpl, err := template.ParseFiles("templates/template.html")
	if err != nil {
		log.Println("Template parse error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	viewData := ViewData{
		Regions: regionData,
		Health:  healthData,
	}

	err = tmpl.Execute(w, viewData)
	if err != nil {
		log.Println("Template execution error: ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
