package evaluators

// textEvaluator : 2つのテキストの値を比較するためのオブジェクト
type textEvaluator struct {
	Source    string
	Condition string
}

// Evaluation : テキストを比較し、結果を返却する
func (h textEvaluator) Evaluation() EvaluationResult {
	if h.Source != h.Condition {
		message := "--[ソース]--\n"
		message += h.Source + "\n"
		message += "--[条件]--\n"
		message += h.Condition
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

// CreateTextEvaluator : 2つのテキストの値を比較するオブジェクトを生成する
func CreateTextEvaluator(source string, condition string) Evaluator {
	return textEvaluator{
		Source:    source,
		Condition: condition,
	}
}
