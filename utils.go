package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	toml "github.com/pelletier/go-toml"
)

func PutTrackingData(ds DataSource, trackingData TrackingData) {
	for key, val := range trackingData.TimedHabits {
		tracker := Tracker{
			Date:   trackingData.Date,
			Action: key,
			Val:    val,
		}

		ds.AddTracker(tracker)
	}

	for key, val := range trackingData.QuantitativeHabits {
		tracker := Tracker{
			Date:   trackingData.Date,
			Action: key,
			Val:    val,
		}

		ds.AddTracker(tracker)
	}
}

func GetConfig() TrackerConfig {
	var trackerConfig TrackerConfig
	doc, err := os.ReadFile("./config/config.toml")
	if err != nil {
		log.Println("error while reading the config,", err)
	}

	err = toml.Unmarshal(doc, &trackerConfig)
	if err != nil {
		log.Println("error while unmarshalling the data in config,", err)
	}

	return trackerConfig
}

func GetTrackingData(filePath string) TrackingData {
	var trackingData TrackingData
	doc, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Error while reading TOML file:", err)
	}

	err = toml.Unmarshal(doc, &trackingData)
	if err != nil {
		log.Println("Error unmarshalling TOML file:", err)
	}

	return trackingData
}

func AddDummyData(ds DataSource) {
	years := [...]int{2024}
	months := [...]int{1, 2, 3, 4}
	habits := [...]string{"coding", "workout"}

	for _, year := range years {
		for _, month := range months {
			for day := 1; day < 31; day++ {
				for _, habit := range habits {
					tracker := Tracker{
						Date:   fmt.Sprintf("%d-%d-%d", year, month, day),
						Action: habit,
						Val:    rand.Intn(60-10) + 10,
					}
					ds.AddTracker(tracker)
				}
			}
		}
	}
}