package stubs

import (
	"encoding/json"
	"module-tester/commons"
	"module-tester/evaluators"
	"module-tester/parameters"
	"module-tester/repository"
	"time"

	"github.com/Microsoft/go-winio"
)

// namedPipeStub : 名前付きパイプサーバーをホストするスタブ
type namedPipeStub struct {
	Response   parameters.NamedPipeServerOption
	Condition  interface{}
	Path       string
	Repository repository.Repository
}

// Listen : 名前付きパイプサーバーを立ち上げる
func (s namedPipeStub) Listen() ProcessResult {
	result := ProcessResult{
		Evaluator: evaluators.CreateFailedEvaluator(nil),
		Response:  s.Response,
		Begin:     time.Now(),
	}

	pipeConfig := winio.PipeConfig{
		SecurityDescriptor: "S:(ML;;NW;;;LW)D:(A;;0x12019f;;;WD)",
		InputBufferSize:    s.Response.InputBufferSize,
		OutputBufferSize:   s.Response.OutputBufferSize,
	}

	listener, err := winio.ListenPipe(`\\.\pipe\`+s.Path, &pipeConfig)
	if err != nil {
		return returnsFailedStubResult(&result, err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		return returnsFailedStubResult(&result, err)
	}
	defer conn.Close()

	data := make([]byte, s.Response.InputBufferSize)
	l, err := conn.Read(data)
	if err != nil {
		return returnsFailedStubResult(&result, err)
	}
	data = data[0:l]

	conn.Write([]byte(s.Response.DataToString()))
	result.End = time.Now()

	if s.Condition == nil {
		result.Evaluator = evaluators.CreateCorrectEvaluator()
	} else {
		if val, castable := s.Condition.(string); castable {
			result.Evaluator = evaluators.CreateTextEvaluator(
				string(data),
				val,
			)
			result.Request = val
		} else {
			var requestData interface{}
			conditionBytes, cErr := json.Marshal(s.Condition)
			if dErr := json.Unmarshal(data, &requestData); dErr == nil && cErr == nil {
				result.Evaluator = evaluators.CreateJSONEvaluator(
					requestData,
					s.Condition,
					s.Repository,
				)
				result.Request = requestData
			} else {
				result.Evaluator = evaluators.CreateTextEvaluator(
					string(data),
					string(conditionBytes),
				)
				result.Request = string(conditionBytes)
			}
		}
	}

	return result
}

// CreatenamedPipeStub : 名前付きパイプサーバーをホストするスタブを生成する
func CreatenamedPipeStub(option StubOption, repository repository.Repository) (Stub, error) {
	var response parameters.NamedPipeServerOption
	err := commons.ToStruct(option.Response, &response)
	if err != nil {
		return failedStub{}, err
	}
	response.DataSetFromRepository(repository)

	return namedPipeStub{
		Response:   response,
		Condition:  option.Condition,
		Path:       option.Path,
		Repository: repository,
	}, nil
}
