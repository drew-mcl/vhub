package api

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type Environment struct {
	Name string         `json:"name"`
	Apps map[string]App `json:"apps"`
}

type App struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var Regions = make(map[string]map[string]Environment)
var Mutex = &sync.Mutex{}
var Log = logrus.New()
var DataFilePath string

func LoadData() error {
	if _, err := os.Stat(DataFilePath); os.IsNotExist(err) {
		defaultRegions := []string{"amer", "emea", "apac"}
		defaultEnvs := []string{"production", "qa", "uat", "dev", "dr"}

		Mutex.Lock()
		defer Mutex.Unlock()
		for _, region := range defaultRegions {
			Regions[region] = make(map[string]Environment)
			for _, env := range defaultEnvs {
				Regions[region][env] = Environment{Name: env, Apps: make(map[string]App)}
			}
		}

		Log.Info("Created initial regions and environments")
		if err := SaveData(); err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	data, err := os.ReadFile(DataFilePath)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"filePath": DataFilePath,
		}).Error("Failed to read data file")
		return err
	}

	if err := json.Unmarshal(data, &Regions); err != nil {
		Log.WithFields(logrus.Fields{
			"filePath": DataFilePath,
		}).Error("Failed to parse data file")
		return err
	}

	Log.Info("Successfully loaded data from file")
	return nil
}

func SaveData() error {
	data, err := json.Marshal(Regions)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"filePath": DataFilePath,
		}).Error("Failed to serialize data")
		return err
	}

	err = os.WriteFile(DataFilePath, data, 0644)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"filePath": DataFilePath,
		}).Error("Failed to write data to file")
		return err
	}

	Log.Info("Successfully saved data to file")
	return nil
}
