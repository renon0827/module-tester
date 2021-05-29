package parameters

import (
	"module-tester/repository"
	"net/http"
)

// HTTPRequest : HTTPリクエストのオプション
type HTTPRequest struct {
	SubPath string      `json:"SubPath"`
	Method  string      `json:"Method"`
	Header  http.Header `json:"Header"`
	Data    interface{} `json:"Data"`
}

// HTTPResponse : HTTPレスポンスのオプション
type HTTPResponse struct {
	StatusCode int         `json:"StatusCode"`
	Header     http.Header `json:"Header"`
	Data       interface{} `json:"Data"`
}

func (h HTTPRequest) DataToString() string {
	return dataToString(h.Data)
}

func (h *HTTPRequest) DataSetFromRepository(repo repository.Repository) {
	h.Data = dataSetFromRepository(h.Data, repo)
}

func (h HTTPResponse) DataToString() string {
	return dataToString(h.Data)
}

func (h *HTTPResponse) DataSetFromRepository(repo repository.Repository) {
	h.Data = dataSetFromRepository(h.Data, repo)
}
