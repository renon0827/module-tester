package scenarios

import "time"

type ScenarioTestResult struct {
	ScenarioName string             `json:"ScenarioName"`
	Results      []ModuleTestResult `json:"Results"`
}

type ModuleTestResult struct {
	Driver *TestResult   `json:"Driver"`
	Stubs  *[]TestResult `json:"Stubs"`
}

type TestResult struct {
	Correct bool              `json:"Correct"`
	Message string            `json:"Message"`
	Detail  *TestResultDetail `jsou:"Detail"`
}

type TestResultDetail struct {
	Score    float64     `json:"Score"`
	Request  interface{} `json:"Request"`
	Response interface{} `json:"Response"`
	Begin    time.Time   `json:"Begin"`
	End      time.Time   `json:"End"`
}
