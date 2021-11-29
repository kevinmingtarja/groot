package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	host = "postgresql"
	port = 5432
)

func Initialize() (*Database, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	db := &Database{conn}
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}

	fmt.Println("Successfully connected to database!")
	return db, nil
}

func (db Database) ErrorLog(id int) (ErrorLog, error) {
	var log ErrorLog

	row := db.Conn.QueryRow("SELECT * FROM error_logs WHERE id = ?", id)
	if err := row.Scan(&log.ID, &log.Time, &log.StackTrace, &log.UserAgent, &log.HTTPCode, &log.AppName, &log.FunctionName); err != nil {
		if err == sql.ErrNoRows {
			return log, fmt.Errorf("no such log found")
		}
		return log, fmt.Errorf(err.Error())
	}

	return log, nil
}
