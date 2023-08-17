package data

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var GlobalData = Data{Regions: make(map[string]Region)}
var Mutex = &sync.RWMutex{}
var Log = logrus.New()
var DataFilePath string
var BackupFilePath string

func LoadData() error {
	if err := loadDataFromFile(DataFilePath); err != nil {
		Log.WithField("filePath", DataFilePath).Warn("Failed to load data from primary file. Attempting to load from backup.")

		// If the primary load fails, try loading from the backup file
		if err := loadDataFromFile(BackupFilePath); err != nil {
			Log.WithField("filePath", BackupFilePath).Error("Failed to load data from backup file")
			return err
		}
	}

	Log.Info("Successfully loaded data from file")
	return nil
}

func CreateFileIfNotExists(filePath string) error {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			Log.WithField("filePath", filePath).Error("Failed to create file")
			return err
		}
		defer file.Close()

		initialData := Data{
			Regions: make(map[string]Region),
			// You can add more initialization here if needed
		}

		jsonData, err := json.Marshal(initialData)
		if err != nil {
			Log.WithField("filePath", filePath).Error("Failed to marshal initial data to JSON")
			return err
		}

		if _, err := file.Write(jsonData); err != nil {
			Log.WithField("filePath", filePath).Error("Failed to initialize JSON file")
			return err
		}
		Log.WithField("filePath", filePath).Info("File created")
	}
	return nil
}

func loadDataFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &GlobalData); err != nil {
		return err
	}
	return nil
}

func SaveData(filePath string) error {
	data, err := json.Marshal(GlobalData)
	if err != nil {
		return err
	}

	// Write to the backup file first
	if err = os.WriteFile(BackupFilePath, data, 0644); err != nil {
		return err
	}

	// Write to the primary file
	if err = os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	Log.WithField("filePath", filePath).Debug("Successfully saved data to file")
	return nil
}
