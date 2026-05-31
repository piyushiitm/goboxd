package models

type RunRequest struct {
	Language string `json:"language"`
	Source   string `json:"source"`
}

type RunResponse struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}
