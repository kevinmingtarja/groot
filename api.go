package main

type LogsResponse struct {
	Logs []ErrorLog `json:"logs"`
}

type LogResponse struct {
	Log ErrorLog `json:"log"`
}
