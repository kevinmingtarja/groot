package main

import "time"

type ErrorLog struct {
	ID           int
	Time         time.Time
	StackTrace   string
	UserAgent    string
	HTTPCode     int
	AppName      string
	FunctionName string
}

type ErrorLogStore interface {
	ErrorLog(id int) (ErrorLog, error)
	ErrorLogByApp(appName string) ([]ErrorLog, error)
	ErrorLogByFunction(functionName string) ([]ErrorLog, error)
	ErrorLogs() ([]ErrorLog, error)
	CreateErrorLog(e *ErrorLog) error
}
