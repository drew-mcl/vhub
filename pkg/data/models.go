package data

type Data struct {
	Regions map[string]Region `json:"regions"`
}

type Region struct {
	Name         string                 `json:"name"`
	Environments map[string]Environment `json:"environments"`
}

type Environment struct {
	Name string         `json:"name"`
	Apps map[string]App `json:"apps"`
}

type App struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Route   string `json:"route"`
	Date    string `json:"date"`
}
