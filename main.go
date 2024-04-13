package main

import (
	"log"
	"os"
	"os/exec"
)

// const configPath string = "/Users/harshitcd/.config/habit_tracker/config.toml"

func main() {
	if len(os.Args) <= 2 {
		log.Println("doesn't have enough command line arguments")
		return
	}
	configPath := os.Args[1]
	commandLineArg := os.Args[2]
	switch commandLineArg {
	case "push":
		pushTrackersInfo(configPath)
	case "update":
		updateTrackersInfo(configPath)
	default:
		log.Println("invalid command line argument")
	}
}

func pushTrackersInfo(configPath string) {
	trackerconfig := GetConfig(configPath)
	ds := NewDataSource(trackerconfig.DatabaseConfig)

	trackingData := GetTrackingData(trackerconfig.Paths.TrackerPath)
	PutTrackingData(ds, trackingData)

	log.Println("added your tracking data")
}

func updateTrackersInfo(configPath string) {
	trackerconfig := GetConfig(configPath)
	cmdStruct := exec.Command(trackerconfig.Editor, trackerconfig.Paths.TrackerPath)
	_, err := cmdStruct.Output()
	if err != nil {
		log.Println(err)
	}
	log.Println("opened the tracker file")
}
