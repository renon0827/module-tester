package drivers

import (
	"encoding/json"
	"module-tester/commons"
	"module-tester/evaluators"
	"module-tester/parameters"
	"module-tester/repository"
	"net"
	"time"

	"github.com/Microsoft/go-winio"
)

// namedPipeDriver : 名前付きパイプによるリクエストを行うオブジェクト
type namedPipeDriver struct {
	Request    parameters.NamedPipeClientOption
	Condition  interface{}
	Path       string
	Repository repository.Repository
}

// Drive : 名前付きパイプによるリクエストを行う
func (d namedPipeDriver) Drive() DriveResult {
	result := DriveResult{}
	var err error
	var conn net.Conn

	result.Begin = time.Now()
	var timeout *time.Duration
	if d.Request.TimeoutSeconds != 0 {
		t := time.Second * time.Duration(d.Request.TimeoutSeconds)
		timeout = &t
	}
	for {
		conn, err = winio.DialPipe(`\\.\pipe\`+d.Path, timeout)
		if err != nil {
			if timeout != nil && time.Now().Unix() < result.Begin.Add(*timeout).Unix() {
				continue
			}
			return returnsFailedDriveResult(&result, err)
		}
		break
	}
	result.Request = d.Request.Data

	conn.Write([]byte(d.Request.DataToString()))

	data := make([]byte, d.Request.BufferSize)
	l, err := conn.Read(data)
	if err != nil {
		return returnsFailedDriveResult(&result, err)
	}
	result.End = time.Now()

	data = data[0:l]
	result.Response = string(data)

	if val, castable := d.Condition.(string); castable {
		result.Evaluator = evaluators.CreateTextEvaluator(
			string(data),
			val,
		)
	} else {
		var responseData interface{}
		conditionBytes, cErr := json.Marshal(d.Condition)
		if dErr := json.Unmarshal(data, &responseData); dErr == nil && cErr == nil {
			result.Evaluator = evaluators.CreateJSONEvaluator(
				responseData,
				d.Condition,
				d.Repository,
			)
		} else {
			result.Evaluator = evaluators.CreateTextEvaluator(
				string(data),
				string(conditionBytes),
			)
		}
	}

	return result
}

// CreateNamedPipeDriver : 名前付きパイプによるリクエストを行うオブジェクトを生成する
func CreateNamedPipeDriver(option DriverOption, repository repository.Repository) (Driver, error) {
	var request parameters.NamedPipeClientOption
	err := commons.ToStruct(option.Request, &request)
	if err != nil {
		return failedDriver{
			Err: err,
		}, err
	}
	request.DataSetFromRepository(repository)

	return namedPipeDriver{
		Request:    request,
		Condition:  option.Condition,
		Path:       option.Path,
		Repository: repository,
	}, nil
}
