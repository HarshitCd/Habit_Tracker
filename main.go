package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/go-sql-driver/mysql"
)

const (
	user     = "root"
	password = "root"
)

type Tracker struct {
	date     string
	habit    string
	duration int
}

func (t *Tracker) isEmpty() bool {
	return t.date == "" || t.habit == ""
}

func main() {
	cfg := mysql.Config{
		User:   user,
		Passwd: password,
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "habit_tracker",
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

	fmt.Println("Connected!")
	AddDummyData(db)
}

func AddDummyData(db *sql.DB) {
	years := [...]int{2024}
	months := [...]int{1, 2, 3}

	for _, year := range years {
		for _, month := range months {
			for day := 1; day < 31; day++ {
				tracker := Tracker{
					date:     fmt.Sprintf("%d-%d-%d", year, month, day),
					habit:    "workout",
					duration: rand.Intn(60-10) + 10,
				}
				AddTracker(db, tracker)
			}
		}
	}
}

func AddTracker(db *sql.DB, t Tracker) {
	tracker, err := SelectTracker(db, t)
	if err != nil {
		log.Println(err.Error())
	}

	if tracker.isEmpty() {
		_, err = db.Exec("INSERT INTO tracker (date, habit, duration) VALUES (?, ?, ?)", t.date, t.habit, t.duration)
		if err != nil {
			log.Println("could not insert into the DB", t)
		}
	} else {
		_, err = db.Exec("UPDATE tracker SET duration = ? WHERE date = ? AND habit = ?", t.duration, t.date, t.habit)
		if err != nil {
			log.Println("could not update into the DB", t)
		}
	}

}

func SelectTracker(db *sql.DB, t Tracker) (Tracker, error) {
	var tracker Tracker
	rows, err := db.Query("SELECT * FROM tracker WHERE date = ? AND habit = ?", t.date, t.habit)
	if err != nil {
		return tracker, fmt.Errorf("failed fetch data")
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&tracker.date, &tracker.habit, &tracker.duration); err != nil {
			return tracker, fmt.Errorf("could not to scan the data")
		}
	}

	return tracker, nil
}

func DeleteTracker(db *sql.DB, t Tracker) {
	_, err := db.Exec("DELETE FROM tracker WHERE date = ? AND habit = ?", t.date, t.habit)
	if err != nil {
		log.Println("could not delete tracker", t)
	}
}
