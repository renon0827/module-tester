package scenarios

import (
	"module-tester/drivers"
	"module-tester/stubs"
)

type ScenarioTest struct {
	Scenarios []Scenario `json:"Scenarios"`
}

type Scenario struct {
	Name        string          `json:"Name"`
	Description string          `json:"Description"`
	Sequence    []ModuleSetting `json:"Sequence"`
}

type ModuleSetting struct {
	Driver *drivers.DriverOption `json:"Driver"`
	Stubs  *[]stubs.StubOption   `json:"Stubs"`
}
