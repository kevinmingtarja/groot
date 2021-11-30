package main

import "time"

type ErrorLog struct {
	ID           int       `json:"id"`
	Time         time.Time `json:"time"`
	RequestURL   string    `json:"request_url"`
	StackTrace   *string   `json:"stack_trace"`
	UserAgent    *string   `json:"user_agent"`
	HTTPCode     int       `json:"http_code"`
	AppName      string    `json:"app_name"`
	FunctionName *string   `json:"function_name"`
}

type ErrorLogStore interface {
	ErrorLog(id int) (ErrorLog, error)
	ErrorLogByURL(url string) ([]ErrorLog, error)
	ErrorLogByApp(appName string) ([]ErrorLog, error)
	ErrorLogByFunction(functionName string) ([]ErrorLog, error)
	ErrorLogs() ([]ErrorLog, error)
	CreateErrorLog(e *ErrorLog) error
}

type Chat struct {
	AppName string `json:"app_name"`
	ChatID  int    `json:"chat_id"`
}

type ChatStore interface {
	ChatID(appName string) (int, error)
	SetChatID(appName string, chatID int) error
}
