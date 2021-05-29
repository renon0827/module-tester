package stubs

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"module-tester/commons"
	"module-tester/evaluators"
	"module-tester/parameters"
	"module-tester/repository"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// httpStub : HTTPサーバーをホストするスタブ
type httpStub struct {
	Condition  *parameters.HTTPRequest
	Response   parameters.HTTPResponse
	Path       string
	Timeout    int
	Repository repository.Repository
}

type httpServerHandler struct {
	Handler func(w http.ResponseWriter, r *http.Request)
}

func (h *httpServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Handler(w, r)
}

// Listen : HTTPサーバーを立ち上げる
func (s httpStub) Listen() ProcessResult {
	resultChan := make(chan ProcessResult, 1)
	closeSignal := make(chan bool, 1)
	timeoutFlag := true
	srv := &http.Server{
		Addr: s.Path,
		Handler: &httpServerHandler{
			Handler: func(w http.ResponseWriter, r *http.Request) {
				result := ProcessResult{
					Evaluator: nil,
					Response:  s.Response,
					Begin:     time.Now(),
				}
				content, err := ioutil.ReadAll(r.Body)
				if err != nil {
					result.Evaluator = evaluators.CreateFailedEvaluator(err)
					resultChan <- result
					return
				}

				var requestData interface{}
				if err := json.Unmarshal(content, &requestData); err != nil {
					requestData = string(content)
				}
				result.Request = parameters.HTTPRequest{
					Data:    requestData,
					Header:  r.Header,
					Method:  r.Method,
					SubPath: r.URL.Path,
				}

				for key, header := range s.Response.Header {
					for _, item := range header {
						w.Header().Add(key, item)
					}
				}
				if s.Response.StatusCode != 0 {
					w.WriteHeader(s.Response.StatusCode)
				}
				w.Write([]byte(s.Response.DataToString()))

				result.End = time.Now()
				resultChan <- result
				closeSignal <- true
				timeoutFlag = false
			},
		},
	}

	go func() {
		time.Sleep(time.Second * time.Duration(s.Timeout))
		if timeoutFlag && s.Timeout != 0 {
			closeSignal <- true
			resultChan <- ProcessResult{
				Evaluator: evaluators.CreateFailedEvaluator(errors.New("受付がタイムアウトしました")),
			}
		}
	}()

	go func() {
		<-closeSignal
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Infoln(err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Infoln(err)
	}

	result := <-resultChan
	if result.Evaluator != nil {
		return result
	}

	if s.Condition == nil {
		result.Evaluator = evaluators.CreateCorrectEvaluator()
	} else if request, castable := result.Request.(parameters.HTTPRequest); castable {
		result.Evaluator = evaluators.CreateHTTPRequestEvaluator(
			request,
			*s.Condition,
			s.Repository,
		)
	} else {
		return result
	}
	return result
}

// CreateHTTPStub : HTTPサーバーをホストするスタブを生成する
func CreateHTTPStub(option StubOption, repository repository.Repository) (Stub, error) {
	var response parameters.HTTPResponse
	err := commons.ToStruct(option.Response, &response)
	if err != nil {
		return failedStub{
			Err: err,
		}, err
	}
	response.DataSetFromRepository(repository)

	var condition *parameters.HTTPRequest
	if option.Condition != nil {
		var c parameters.HTTPRequest
		condition = &c
		err = commons.ToStruct(option.Condition, condition)
		if err != nil {
			return failedStub{}, err
		}
	}

	return httpStub{
		Condition:  condition,
		Response:   response,
		Path:       option.Path,
		Timeout:    option.TimeoutSeconds,
		Repository: repository,
	}, nil
}
