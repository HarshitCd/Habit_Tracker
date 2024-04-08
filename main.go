package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/go-sql-driver/mysql"
	toml "github.com/pelletier/go-toml"
)

const (
	user     = "root"
	password = "root"
	filePath = "./template/tracker.toml"
)

type Tracker struct {
	Date   string
	Action string
	Val    int
}

func (t *Tracker) isEmpty() bool {
	return t.Date == "" || t.Action == ""
}

type TrackingData struct {
	Date               string         `toml:"date"`
	TimedHabits        map[string]int `toml:"timed-habits"`
	QuantitativeHabits map[string]int `toml:"quantitative-habits"`
}

func main() {
	cfg := mysql.Config{
		User:   user,
		Passwd: password,
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "trackers",
	}

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	log.Println("Connected!")

	trackingData := GetTrackingData()
	PutTrackingData(db, trackingData)
}

func PutTrackingData(db *sql.DB, trackingData TrackingData) {
	for key, val := range trackingData.TimedHabits {
		tracker := Tracker{
			Date:   trackingData.Date,
			Action: key,
			Val:    val,
		}

		AddTracker(db, tracker)
	}

	for key, val := range trackingData.QuantitativeHabits {
		tracker := Tracker{
			Date:   trackingData.Date,
			Action: key,
			Val:    val,
		}

		AddTracker(db, tracker)
	}
}

func GetTrackingData() TrackingData {
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

func AddDummyData(db *sql.DB) {
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
					AddTracker(db, tracker)
				}
			}
		}
	}
}

func AddTracker(db *sql.DB, t Tracker) {
	if t.isEmpty() {
		log.Println("no tracking data provided")
		return
	}

	tracker, err := SelectTracker(db, t)
	if err != nil {
		log.Println(err.Error())
	}

	if tracker.isEmpty() {
		_, err = db.Exec("INSERT INTO habit_tracker (date, action, val) VALUES (?, ?, ?)", t.Date, t.Action, t.Val)
		if err != nil {
			log.Println("could not insert into the DB", t)
		}
	} else {
		_, err = db.Exec("UPDATE habit_tracker SET val = ? WHERE date = ? AND action = ?", t.Val, t.Date, t.Action)
		if err != nil {
			log.Println("could not update into the DB", t)
		}
	}

}

func SelectTracker(db *sql.DB, t Tracker) (Tracker, error) {
	var tracker Tracker
	rows, err := db.Query("SELECT * FROM habit_tracker WHERE date = ? AND action = ?", t.Date, t.Action)
	if err != nil {
		return tracker, fmt.Errorf("failed fetch data")
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&tracker.Date, &tracker.Action, &tracker.Val); err != nil {
			return tracker, fmt.Errorf("could not to scan the data")
		}
	}

	return tracker, nil
}

func DeleteTracker(db *sql.DB, t Tracker) {
	_, err := db.Exec("DELETE FROM habit_tracker WHERE date = ? AND habit = ?", t.Date, t.Action)
	if err != nil {
		log.Println("could not delete tracker", t)
	}
}
