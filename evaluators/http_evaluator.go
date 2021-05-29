package evaluators

import (
	"encoding/json"
	"module-tester/parameters"
	"module-tester/repository"
	"strings"
)

// httpResponseEvaluator : HTTPレスポンスを比較するためのオブジェクト
type httpResponseEvaluator struct {
	Source     parameters.HTTPResponse
	Condition  parameters.HTTPResponse
	Repository repository.Repository
}

// httpRequestEvaluator : HTTPリクエストを比較するためのオブジェクト
type httpRequestEvaluator struct {
	Source     parameters.HTTPRequest
	Condition  parameters.HTTPRequest
	Repository repository.Repository
}

// Evaluation : HTTPレスポンスを比較し、結果を返却する
func (h httpResponseEvaluator) Evaluation() EvaluationResult {
	if h.Condition.StatusCode != 0 && h.Condition.StatusCode != h.Source.StatusCode {
		return EvaluationResult{
			Correct: false,
			Message: "ステータスコードが異なります",
		}
	}

	if h.Condition.Header != nil {
		for key, value := range h.Condition.Header {
			if _, exist := h.Source.Header[key]; exist {
				for _, item := range value {
					if !strings.Contains(h.Source.Header.Get(key), item) {
						return EvaluationResult{
							Correct: false,
							Message: "条件ヘッダーが含まれていません : '" + key + "' の '" + item + "'",
						}
					}
				}
			} else {
				return EvaluationResult{
					Correct: false,
					Message: "ヘッダーが不十分です : '" + key + "'",
				}
			}
		}
	}

	result := EvaluationResult{
		Correct: true,
		Message: CorrectMessage,
	}
	if h.Condition.Data != nil {
		evaluated := false
		if _, err := json.Marshal(h.Condition.Data); err == nil {
			j := jsonEvaluator{
				Source:     h.Source.Data,
				Condition:  h.Condition.Data,
				Repository: h.Repository,
			}.Evaluation()
			result.Correct = j.Correct
			result.Message = j.Message
			evaluated = true
		} else {
			if source, castable := h.Source.Data.(string); castable {
				if condition, castable := h.Condition.Data.(string); castable {
					t := textEvaluator{
						Source:    source,
						Condition: condition,
					}.Evaluation()

					result.Correct = t.Correct
					result.Message = t.Message
					evaluated = true
				}
			}
		}

		if !evaluated {
			result.Correct = false
			result.Message = "比較できないHTTPレスポンスデータです"
		}
	}

	return result
}

// Evaluation : HTTPリクエストを比較し、結果を返却する
func (h httpRequestEvaluator) Evaluation() EvaluationResult {
	if h.Condition.Method != "" && h.Condition.Method != h.Source.Method {
		return EvaluationResult{
			Correct: false,
			Message: "HTTPメソッドが異なります",
		}
	}

	if h.Condition.SubPath != "" && h.Condition.SubPath != h.Source.SubPath {
		return EvaluationResult{
			Correct: false,
			Message: "パスが異なります",
		}
	}

	if h.Condition.Header != nil {
		for key, value := range h.Condition.Header {
			if _, exist := h.Source.Header[key]; exist {
				for _, item := range value {
					if !strings.Contains(h.Source.Header.Get(key), item) {
						return EvaluationResult{
							Correct: false,
							Message: "条件ヘッダーが含まれていません : '" + key + "' の '" + item + "'",
						}
					}
				}
			} else {
				return EvaluationResult{
					Correct: false,
					Message: "ヘッダーが不十分です : '" + key + "'",
				}
			}
		}
	}

	result := EvaluationResult{
		Correct: true,
		Message: CorrectMessage,
	}
	if h.Condition.Data != nil {
		evaluated := false
		if _, err := json.Marshal(h.Condition.Data); err == nil {
			j := jsonEvaluator{
				Source:     h.Source.Data,
				Condition:  h.Condition.Data,
				Repository: h.Repository,
			}.Evaluation()
			result.Correct = j.Correct
			result.Message = j.Message
			evaluated = true
		} else {
			if source, castable := h.Source.Data.(string); castable {
				if condition, castable := h.Condition.Data.(string); castable {
					t := textEvaluator{
						Source:    source,
						Condition: condition,
					}.Evaluation()

					result.Correct = t.Correct
					result.Message = t.Message
					evaluated = true
				}
			}
		}

		if !evaluated {
			result.Correct = false
			result.Message = "比較できないHTTPリクエストデータです"
		}
	}

	return result
}

// CreateHTTPResponseEvaluator : HTTPレスポンスを比較するためのオブジェクトを生成する
func CreateHTTPResponseEvaluator(source parameters.HTTPResponse, condition parameters.HTTPResponse, repository repository.Repository) Evaluator {
	return httpResponseEvaluator{
		Source:     source,
		Condition:  condition,
		Repository: repository,
	}
}

// CreateHTTPRequestEvaluator : HTTPレスポンスを比較するためのオブジェクトを生成する
func CreateHTTPRequestEvaluator(source parameters.HTTPRequest, condition parameters.HTTPRequest, repository repository.Repository) Evaluator {
	return httpRequestEvaluator{
		Source:     source,
		Condition:  condition,
		Repository: repository,
	}
}
