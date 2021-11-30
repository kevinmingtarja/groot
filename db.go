package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

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

	db := &Database{conn}
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}

	fmt.Println("Successfully connected to database!")
	return db, nil
}

func (db *Database) ErrorLog(id int) (ErrorLog, error) {
	var errorLog ErrorLog

	row := db.Conn.QueryRow("SELECT * FROM error_logs WHERE id = $1", id)
	if err := row.Scan(&errorLog.ID, &errorLog.Time, &errorLog.RequestURL, &errorLog.StackTrace, &errorLog.UserAgent, &errorLog.HTTPCode, &errorLog.AppName, &errorLog.FunctionName); err != nil {
		if err == sql.ErrNoRows {
			return errorLog, fmt.Errorf("no such errorLog found")
		}
		return errorLog, fmt.Errorf(err.Error())
	}

	return errorLog, nil
}

func (db *Database) ErrorLogByURL(url string) ([]ErrorLog, error) {
	logs := []ErrorLog{} // this is to prevent a nil slice which prevents a null response

	rows, err := db.Conn.Query("SELECT * FROM error_logs WHERE request_url = $1", url)
	if err != nil {
		return nil, fmt.Errorf("ErrorLogByURL - Query Error: %s", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var errorLog ErrorLog
		if err := rows.Scan(&errorLog.ID, &errorLog.Time, &errorLog.RequestURL, &errorLog.StackTrace, &errorLog.UserAgent, &errorLog.HTTPCode, &errorLog.AppName, &errorLog.FunctionName); err != nil {
			return nil, fmt.Errorf("ErrorLogByURL - Scan Error: %s", err.Error())
		}
		logs = append(logs, errorLog)
	}

	return logs, nil
}

func (db *Database) CreateErrorLog(e *ErrorLog) error {
	query := "INSERT INTO error_logs (time, request_url, stack_trace, user_agent, http_code, app_name, function_name) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)"
	res, err := db.Conn.Exec(query, time.Now(), e.RequestURL, e.StackTrace, e.UserAgent, e.HTTPCode, e.AppName, e.FunctionName)
	if err != nil {
		return err
	}
	fmt.Println(res)

	return nil
}
