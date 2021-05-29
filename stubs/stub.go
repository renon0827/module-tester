package stubs

import (
	"errors"
	"module-tester/evaluators"
	"module-tester/repository"
	"strings"
	"time"
)

type StubOption struct {
	Type           string      `json:"Type"`
	Path           string      `json:"Path"`
	TimeoutSeconds int         `json:"TimeoutSeconds"`
	Response       interface{} `json:"Response"`
	Condition      interface{} `json:"Condition"`
}

type ProcessResult struct {
	Request   interface{}
	Response  interface{}
	Begin     time.Time
	End       time.Time
	Evaluator evaluators.Evaluator
}

// Stub : スタブのインターフェース
type Stub interface {
	Listen() ProcessResult
}

// failedStub : 必ず失敗するスタブ
type failedStub struct {
	Err    error
	Result ProcessResult
}

// GetEvaluator : FailedEvaluatorを返す
func (d failedStub) GetEvaluator() evaluators.Evaluator {
	return evaluators.CreateFailedEvaluator(d.Err)
}

func (d failedStub) Listen() ProcessResult {
	return returnsFailedStubResult(&d.Result, d.Err)
}

// CreateStub : create stub
func CreateStub(option StubOption, repository repository.Repository) (Stub, error) {
	switch strings.ToUpper(option.Type) {
	case "HTTP", "":
		option.Type = "HTTP"
		return CreateHTTPStub(option, repository)
	case "NAMEDPIPE":
		return CreatenamedPipeStub(option, repository)
	}
	return failedStub{}, errors.New("'" + option.Type + "' は存在しないStubの種類です")
}

func returnsFailedStubResult(result *ProcessResult, err error) ProcessResult {
	result.Evaluator = evaluators.CreateFailedEvaluator(err)
	return *result
}
