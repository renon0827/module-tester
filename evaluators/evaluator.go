package evaluators

import (
	"errors"
)

const CorrectMessage = "テストは正しく終了しました"

type EvaluationResult struct {
	Correct bool
	Message string
}

type Evaluator interface {
	Evaluation() EvaluationResult
}

// failedEvaluator : 必ずNGになるEvaluator
type failedEvaluator struct {
	Err error
}

// Evaluation : 必ずNGになる
func (e failedEvaluator) Evaluation() EvaluationResult {
	if e.Err == nil {
		e.Err = errors.New("FailedEvaluatorのEvaluationが呼ばれました")
	}
	return EvaluationResult{
		Correct: false,
		Message: e.Err.Error(),
	}
}

// correctEvaluator : 必ずOKになるEvaluator
type correctEvaluator struct {
}

// Evaluation : 必ずOKになる
func (e correctEvaluator) Evaluation() EvaluationResult {
	return EvaluationResult{
		Correct: true,
		Message: CorrectMessage,
	}
}

// CreateFailedEvaluator : 必ずNGになるEvaluatorを生成する
func CreateFailedEvaluator(err error) Evaluator {
	return failedEvaluator{
		Err: err,
	}
}

// CreateCorrectEvaluator : 必ずOKになるEvaluatorを生成する
func CreateCorrectEvaluator() Evaluator {
	return correctEvaluator{}
}
