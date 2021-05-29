package evaluators

import (
	"reflect"
	"strconv"
	"strings"

	"module-tester/constants"
	"module-tester/repository"
)

type difType int
type jsonValueType int

type jsonValue = interface{}
type jsonObject = map[string]interface{}
type jsonArray = []jsonValue

type jsonCompare struct {
	difType   difType
	Child     interface{}
	ChildType jsonValueType
}

const (
	Equal difType = iota
	NotExist
	DifferentType
	DifferentValue
)

const (
	Null jsonValueType = iota
	Object
	Array
)

// jsonEvaluator : 2つのJSONの値を評価するためのオブジェクト
type jsonEvaluator struct {
	Source     jsonValue
	Condition  jsonValue
	Repository repository.Repository
}

// Evaluation : JSONの比較を行い、結果を返却する
func (e jsonEvaluator) Evaluation() EvaluationResult {
	compare := e.CompareJSON(e.Condition, e.Source)
	if compare.difType != Equal {
		message := e.CreateCompareMessage("root", "", compare)
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

// CompareJSON : JSONのオブジェクトを比較する
func (e jsonEvaluator) CompareJSON(src, val jsonValue) jsonCompare {
	differentType := jsonCompare{
		difType:   DifferentType,
		Child:     nil,
		ChildType: Null,
	}

	notExist := jsonCompare{
		difType:   NotExist,
		Child:     nil,
		ChildType: Null,
	}

	// 配列
	if value, castable := src.(jsonArray); castable {
		if value2, castable := val.(jsonArray); castable {
			difType := Equal
			child := make([]jsonCompare, len(value))
			for index := range value {
				if index >= len(value2) {
					child[index] = notExist
					difType = DifferentValue
					continue
				}
				child[index] = e.CompareJSON(value[index], value2[index])
				if child[index].difType != Equal {
					difType = DifferentValue
				}
			}
			return jsonCompare{
				difType:   difType,
				Child:     child,
				ChildType: Array,
			}
		} else {
			return differentType
		}
	}

	// オブジェクト
	if value, castable := src.(jsonObject); castable {
		if value2, castable := val.(jsonObject); castable {
			difType := Equal
			child := map[string]jsonCompare{}
			for key := range value {
				if _, exist := value2[key]; !exist {
					child[key] = notExist
					difType = DifferentValue
					continue
				}
				child[key] = e.CompareJSON(value[key], value2[key])
				if child[key].difType != Equal {
					difType = DifferentValue
				}
			}
			return jsonCompare{
				difType:   difType,
				Child:     child,
				ChildType: Object,
			}
		} else {
			return differentType
		}
	}

	// 値
	if reflect.TypeOf(src) != reflect.TypeOf(val) {
		return differentType
	}

	if v, castable := src.(string); e.Repository != nil && castable &&
		strings.HasPrefix(v, constants.ValiableValuePrefix) && strings.HasSuffix(v, constants.ValiableValueSuffix) {
		v = v[len(constants.ValiableValuePrefix) : len(v)-len(constants.ValiableValueSuffix)]
		if len(v) > 0 {
			e.Repository.Set(v, val)
			val = src
		}
	}

	if src != val {
		return jsonCompare{
			difType:   DifferentValue,
			Child:     nil,
			ChildType: Null,
		}
	}

	return jsonCompare{
		difType:   Equal,
		Child:     nil,
		ChildType: Null,
	}
}

// CreateCompareMessage : 比べた結果のメッセージを生成する
func (e jsonEvaluator) CreateCompareMessage(name, head string, compare jsonCompare) string {
	message := name + " : "
	t := ""

	switch compare.difType {
	case Equal:
		t += "一致"
	case NotExist:
		t += "存在しません"
	case DifferentType:
		t += "型が異なります"
	case DifferentValue:
		t += "値が異なります"
	}

	if compare.ChildType != Null {
		message += t + "\n"
		if compare.ChildType == Object {
			for key, item := range compare.Child.(map[string]jsonCompare) {
				message += e.CreateCompareMessage(key, head+"\t", item)
			}
		} else if compare.ChildType == Array {
			for index, item := range compare.Child.([]jsonCompare) {
				message += e.CreateCompareMessage(strconv.Itoa(index), head+"\t", item)
			}
		}
	} else {
		message += t
	}
	return head + message + "\n"
}

// CreateJSONEvaluator : 2つのJSONの値を評価するオブジェクトを生成する
func CreateJSONEvaluator(source jsonValue, condition jsonValue, repository repository.Repository) Evaluator {
	return jsonEvaluator{
		Source:     source,
		Condition:  condition,
		Repository: repository,
	}
}
