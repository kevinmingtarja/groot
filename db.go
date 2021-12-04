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

func setupDatabase() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return db, err
	}

	fmt.Println("Successfully connected to database!")
	return db, nil
}

func (s *server) ErrorLog(id int) (ErrorLog, error) {
	var errorLog ErrorLog

	row := s.db.QueryRow("SELECT * FROM error_logs WHERE id = $1", id)
	if err := row.Scan(&errorLog.ID, &errorLog.Time, &errorLog.RequestURL, &errorLog.StackTrace, &errorLog.UserAgent, &errorLog.HTTPCode, &errorLog.AppName, &errorLog.FunctionName); err != nil {
		if err == sql.ErrNoRows {
			return errorLog, fmt.Errorf("no error log with given ID found")
		}
		return errorLog, fmt.Errorf(err.Error())
	}

	return errorLog, nil
}

func (s *server) ErrorLogByURL(url string) ([]ErrorLog, error) {
	logs := []ErrorLog{} // this is to prevent a nil slice which prevents a null response

	rows, err := s.db.Query("SELECT * FROM error_logs WHERE request_url = $1", url)
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

func (s *server) CreateErrorLog(e *ErrorLog) (int, error) {
	var id int

	query := "INSERT INTO error_logs (time, request_url, stack_trace, user_agent, http_code, app_name, function_name) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	err := s.db.QueryRow(query, time.Now(), e.RequestURL, e.StackTrace, e.UserAgent, e.HTTPCode, e.AppName, e.FunctionName).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (s *server) ChatID(appName string) (int, error) {
	var chatID int

	query := "SELECT chat_id FROM chat_ids WHERE app_name = $1"
	row := s.db.QueryRow(query, appName)
	if err := row.Scan(&chatID); err != nil {
		if err == sql.ErrNoRows {
			return chatID, fmt.Errorf("no such chat found")
		}
		return chatID, fmt.Errorf(err.Error())
	}

	return chatID, nil
}

func (s *server) SetChatID(c *Chat) error {
	if c.AppName == "" {
		return fmt.Errorf("app_name cannot be empty")
	}

	query := "INSERT INTO chat_ids (app_name, chat_id) VALUES ($1, $2) ON CONFLICT (app_name) DO UPDATE SET chat_id = $2"
	_, err := s.db.Exec(query, c.AppName, c.ChatID)
	if err != nil {
		return err
	}

	return nil
}
