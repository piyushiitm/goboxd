package models

type RunRequest struct {
	Language         string       `json:"language"`
	Source           string       `json:"source"`
	SourceFilename   string       `json:"source_filename,omitempty"`
	ArtifactFilename string       `json:"artifact_filename,omitempty"`
	Build            *PhaseConfig `json:"build,omitempty"`
	Run              *PhaseConfig `json:"run,omitempty"`
	Tests            []TestCase   `json:"tests"`
}

type PhaseConfig struct {
	Limits *Limits  `json:"limits,omitempty"`
	Flags  []string `json:"flags,omitempty"`
}

type Limits struct {
	WallTimeS    int `json:"wall_time_s,omitempty"`
	MemoryKB     int `json:"memory_kb,omitempty"`
	MaxProcesses int `json:"max_processes,omitempty"`
}

type TestCase struct {
	Stdin          string `json:"stdin"`
	ExpectedStdout string `json:"expected_stdout"`
}

type RunResponse struct {
	Status string       `json:"status"`
	Build  *BuildResult `json:"build,omitempty"`
	Tests  []TestResult `json:"tests,omitempty"`
}

type BuildResult struct {
	Status     string `json:"status"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	DurationMS int64  `json:"duration_ms"`
}

type TestResult struct {
	Status       string `json:"status"`
	Stdout       string `json:"stdout"`
	Stderr       string `json:"stderr"`
	DurationMS   int64  `json:"duration_ms"`
	MemoryPeakKB int64  `json:"memory_peak_kb"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
