package drivers

import (
	"errors"
	"module-tester/evaluators"
	"module-tester/repository"
	"strings"
	"time"
)

type DriverOption struct {
	Type      string      `json:"Type"`
	Path      string      `json:"Path"`
	Request   interface{} `json:"Request"`
	Condition interface{} `json:"Condition"`
}

type DriveResult struct {
	Request   interface{}
	Response  interface{}
	Begin     time.Time
	End       time.Time
	Evaluator evaluators.Evaluator
}

// Driver : driver interface
type Driver interface {
	Drive() DriveResult
}

type failedDriver struct {
	Err    error
	Result DriveResult
}

func (d failedDriver) Drive() DriveResult {
	return returnsFailedDriveResult(&d.Result, d.Err)
}

// CreateDriver : DriverOptionからDriverを作成する
func CreateDriver(option DriverOption, repository repository.Repository) (Driver, error) {
	switch strings.ToUpper(option.Type) {
	case "HTTP", "":
		option.Type = "HTTP"
		return CreateHTTPDriver(option, repository)
	case "NAMEDPIPE":
		return CreateNamedPipeDriver(option, repository)
	}
	return failedDriver{}, errors.New("'" + option.Type + "' は存在しないDriverの種類です")
}

func returnsFailedDriveResult(result *DriveResult, err error) DriveResult {
	result.Evaluator = evaluators.CreateFailedEvaluator(err)
	return *result
}
