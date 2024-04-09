package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
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

type DBConfig struct {
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	Net       string `toml:"net"`
	Address   string `toml:"address"`
	DBName    string `toml:"db-name"`
	TableName string `toml:"table-name"`
}

type PathsConfig struct {
	TrackerPath string `toml:"tracker-path"`
}

type TrackerConfig struct {
	DatabaseConfig DBConfig    `toml:"db-config"`
	Paths          PathsConfig `toml:"paths"`
}

type DataSource struct {
	DbConfig DBConfig
	DB       *sql.DB
}

func NewDataSource(dbConfig DBConfig) DataSource {
	var ds DataSource
	cfg := mysql.Config{
		User:   dbConfig.Username,
		Passwd: dbConfig.Password,
		Net:    dbConfig.Net,
		Addr:   dbConfig.Address,
		DBName: dbConfig.DBName,
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

	ds.DbConfig = dbConfig
	ds.DB = db

	return ds
}

func (ds *DataSource) AddTracker(t Tracker) {
	if t.isEmpty() {
		log.Println("no tracking data provided")
		return
	}

	tracker, err := ds.SelectTracker(t)
	if err != nil {
		log.Println(err.Error())
	}

	if tracker.isEmpty() {
		_, err = ds.DB.Exec("INSERT INTO habit_tracker (date, action, val) VALUES (?, ?, ?)", t.Date, t.Action, t.Val)
		if err != nil {
			log.Println("could not insert into the DB", t)
		}
	} else {
		_, err = ds.DB.Exec("UPDATE habit_tracker SET val = ? WHERE date = ? AND action = ?", t.Val, t.Date, t.Action)
		if err != nil {
			log.Println("could not update into the DB", t)
		}
	}

}

func (ds *DataSource) SelectTracker(t Tracker) (Tracker, error) {
	var tracker Tracker
	rows, err := ds.DB.Query("SELECT * FROM habit_tracker WHERE date = ? AND action = ?", t.Date, t.Action)
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

func (ds *DataSource) DeleteTracker(t Tracker) {
	_, err := ds.DB.Exec("DELETE FROM habit_tracker WHERE date = ? AND habit = ?", t.Date, t.Action)
	if err != nil {
		log.Println("could not delete tracker", t)
	}
}