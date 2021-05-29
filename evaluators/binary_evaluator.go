package evaluators

import (
	"reflect"
)

// binaryEvaluator : 2つのバイナリデータを比較するためのオブジェクト
type binaryEvaluator struct {
	Source    []byte
	Condition []byte
}

// Evaluation : 2つのバイナリデータを比較し、結果を返却する
func (h binaryEvaluator) Evaluation() EvaluationResult {
	if !reflect.DeepEqual(h.Source, h.Condition) {
		message := "バイナリデータが異なります"
		return EvaluationResult{
			Correct: false,
			Message: message,
		}
	}
	return EvaluationResult{
		Correct: true,
		Message: CorrectMessage,
	}
}

// CreateBinaryEvaluator : 2つのバイナリの値を比較するオブジェクトを生成する
func CreateBinaryEvaluator(source []byte, condition []byte) Evaluator {
	return binaryEvaluator{
		Source:    source,
		Condition: condition,
	}
}
