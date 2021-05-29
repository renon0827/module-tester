package drivers

import (
	"encoding/json"
	"io/ioutil"
	"module-tester/commons"
	"module-tester/evaluators"
	"module-tester/parameters"
	"module-tester/repository"
	"net/http"
	"strings"
	"time"
)

// httpDriver : HTTPリクエストを行うオブジェクト
type httpDriver struct {
	Request    parameters.HTTPRequest
	Condition  *parameters.HTTPResponse
	Path       string
	Repository repository.Repository
}

func httpRequest(method, path, data string, header http.Header) (*http.Response, error) {
	request, err := http.NewRequest(method, path, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	if header != nil {
		request.Header = header
	}

	client := http.Client{}
	return client.Do(request)
}

// Drive : HTTPリクエストを実行し、結果を返却する
func (d httpDriver) Drive() DriveResult {
	result := DriveResult{}
	if !strings.HasPrefix(d.Request.SubPath, "/") {
		d.Request.SubPath += "/"
	}

	result.Begin = time.Now()
	response, err := httpRequest(
		d.Request.Method,
		d.Path,
		d.Request.DataToString(),
		d.Request.Header,
	)
	if err != nil {
		return returnsFailedDriveResult(&result, err)
	}
	defer response.Body.Close()
	result.End = time.Now()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return returnsFailedDriveResult(&result, err)
	}

	var responseData interface{}
	if err := json.Unmarshal(data, &responseData); err != nil {
		responseData = string(data)
	}

	result.Request = d.Request
	result.Response = parameters.HTTPResponse{
		Header:     response.Header,
		StatusCode: response.StatusCode,
		Data:       responseData,
	}
	if d.Condition != nil {
		result.Evaluator = evaluators.CreateHTTPResponseEvaluator(
			result.Response.(parameters.HTTPResponse),
			*d.Condition,
			d.Repository,
		)
	} else {
		result.Evaluator = evaluators.CreateCorrectEvaluator()
	}

	return result
}

// CreateHTTPDriver : HTTPリクエストを行うオブジェクトを作成する
func CreateHTTPDriver(option DriverOption, repository repository.Repository) (Driver, error) {
	var request parameters.HTTPRequest
	err := commons.ToStruct(option.Request, &request)
	if err != nil {
		return failedDriver{
			Err: err,
		}, err
	}
	request.DataSetFromRepository(repository)

	var condition *parameters.HTTPResponse
	if option.Condition != nil {
		var c parameters.HTTPResponse
		condition = &c
		err = commons.ToStruct(option.Condition, condition)
		if err != nil {
			return failedDriver{
				Err: err,
			}, err
		}
	}

	if !strings.HasPrefix(request.SubPath, "/") {
		request.SubPath = "/" + request.SubPath
	}

	return httpDriver{
		Request:    request,
		Condition:  condition,
		Path:       option.Path + request.SubPath,
		Repository: repository,
	}, nil
}
