package main

import (
	"log"
)

func main() {
	trackerconfig := GetConfig()
	ds := NewDataSource(trackerconfig.DatabaseConfig)

	trackingData := GetTrackingData(trackerconfig.Paths.TrackerPath)
	PutTrackingData(ds, trackingData)

	log.Println("Added your tracking data")
}
